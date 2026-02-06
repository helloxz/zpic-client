package main

import (
	"context"
	"fmt"
	"zpic-client/model"
	"zpic-client/pkg"

	"zpic-client/core"

	"github.com/robfig/cron/v3"
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
	// 初始化数据库链接
	model.InitDB()
	// 启动定时任务
	go func() {
		crontab()
	}()

}

// Greet returns a greeting for the given name
func (a *App) Greet(name string) string {
	return fmt.Sprintf("Hello %s, It's show time!", name)
}

func crontab() {
	// 创建一个新的cron管理器
	c := cron.New(cron.WithSeconds())
	c.Start()
	// 这里写各种定时任务
	// 每1分钟执行一次URL上传任务
	c.AddFunc("@every 1m", func() {
		core.UploadTaskList()
	})
	// 每隔35s执行一次BatchUpload
	c.AddFunc("@every 35s", func() {
		core.BatchUpload()
	})
	// 每隔20s更新一次上传状态
	c.AddFunc("@every 20s", func() {
		core.UpdateOneTaskBatch(2)
	})
	// 这里使用select{}可以使主程序持续运行，否则主程序可能会结束，导致定时任务未能执行
	// select {}
}

// GetMessage 返回一个用于前端调用的测试字符串。
func (a *App) GetMessage() string {
	return "Hello from Go!"
}
