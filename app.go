package main

import (
	"context"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	goruntime "runtime"
	"strings"
	"zpic-client/model"
	"zpic-client/pkg"

	"zpic-client/core"
	"zpic-client/helper"

	"github.com/robfig/cron/v3"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

// App struct 应用结构体
type App struct {
	ctx context.Context
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
	// 让core包也能获取startup上下文，之所以这样用是因为runtime应用启动回调，包需要上下文参数，不然其它binds调用会很麻烦
	core.SetCtx(ctx)
	// 初始化配置文件
	pkg.LoadConfig()
	// 拷贝嵌入的bin文件到运行目录
	error := CopyEmbeddedBin()
	if error != nil {
		fmt.Printf("复制嵌入的bin文件失败: %v\n", error)
		os.Exit(1)
	}
	// 初始化数据库链接
	model.InitDB()
	// 启动定时任务
	go func() {
		crontab(a)
	}()

}

// Greet returns a greeting for the given name
func (a *App) Greet(name string) string {
	return fmt.Sprintf("Hello %s, It's show time!", name)
}

func crontab(a *App) {
	// 创建一个新的cron管理器
	c := cron.New(cron.WithSeconds())
	c.Start()
	// 这里写各种定时任务
	// 每1分钟执行一次URL上传任务
	c.AddFunc("@every 40s", func() {
		core.UploadTaskList()
		runtime.EventsEmit(a.ctx, "refresh_url_upload_list")
	})
	// 每隔35s执行一次BatchUpload
	c.AddFunc("@every 35s", func() {
		core.BatchUpload()
	})
	// 每隔20s更新一次上传状态
	c.AddFunc("@every 20s", func() {
		core.UpdateOneTaskBatch(2)
		// 推送事件
		runtime.EventsEmit(a.ctx, "refresh_task")
	})
	// 这里使用select{}可以使主程序持续运行，否则主程序可能会结束，导致定时任务未能执行
	// select {}
}

// GetMessage 返回一个用于前端调用的测试字符串。
func (a *App) GetMessage() string {
	return "Hello from Go!"
}

func (a *App) GetRunDir(ctx context.Context) string {
	env := runtime.Environment(ctx)

	// env.BuildType 的值可能为:
	// "dev"        -> 执行 wails dev 时
	// "production" -> 执行 wails build 时
	// "debug"      -> 执行 wails build -debug 时

	if env.BuildType == "dev" {
		runDir, _ := os.Getwd()
		return runDir
	} else {
		runDir, _ := os.Executable()
		return runDir
	}
}

func CopyEmbeddedBin() error {
	runDir := helper.GetRunDir()
	targetBinDir := filepath.Join(runDir, "bin")

	var srcDir string
	switch goruntime.GOOS {
	case "windows":
		srcDir = "bin/windows"
	case "linux":
		srcDir = "bin/linux"
	case "darwin":
		srcDir = "bin/darwin"
	default:
		return nil
	}

	destSubDir := filepath.Join(targetBinDir, goruntime.GOOS)
	if _, err := os.Stat(destSubDir); err == nil {
		return nil
	}

	return fs.WalkDir(embeddedBin, srcDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		relPath, _ := strings.CutPrefix(path, srcDir)
		targetPath := filepath.Join(destSubDir, relPath)

		if d.IsDir() {
			return os.MkdirAll(targetPath, 0755)
		}

		srcFile, err := embeddedBin.Open(path)
		if err != nil {
			return err
		}
		defer srcFile.Close()

		destFile, err := os.OpenFile(targetPath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0755)
		if err != nil {
			return err
		}
		defer destFile.Close()

		_, err = io.Copy(destFile, srcFile)
		return err
	})
}
