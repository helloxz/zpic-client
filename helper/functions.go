package helper

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
)

// 记录日志到文件
func WriteLog(log_content string) {
	// 日志文件目录
	log_dir := "data/logs"
	// 检查日志目录是否存在，不存在则创建
	if _, err := os.Stat(log_dir); os.IsNotExist(err) {
		// 目录不存在，创建目录
		err := os.MkdirAll(log_dir, 0755) // 0755 权限允许所有者读写执行，组和其他用户只读执行
		if err != nil {
			// 如果创建目录过程中出错，可以进行错误处理
			panic(err) // 或者使用其他方式处理错误，比如记录日志等
		}
	}

	// 打开一个文件，如果文件不存在，则会创建一个新文件，如果文件存在，则会在文件末尾追加内容
	file, err := os.OpenFile("data/logs/error.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		fmt.Println(err)
	}
	defer file.Close()

	log.SetOutput(file)
	log.Print(log_content)
}

// 验证是否是文件夹
func V_dir(dir string) bool {
	dirinfo, err := os.Stat(dir)

	if err != nil {
		//fmt.Println(err)
		return false
	}

	if dirinfo.IsDir() {
		return true
	} else {
		return false
	}
}

// 递归创建目录，不存在会自动创建
func CreateDir(path string) error {
	// Stat returns file info. It will return an error if the directory does not exist.
	if _, err := os.Stat(path); os.IsNotExist(err) {
		// MkdirAll creates a directory named path along with any necessary parents.
		// The permission bits perm (before umask) are set for the directory.
		return os.MkdirAll(path, 0755) // 0755 commonly used permissions
	} else if err != nil {
		return err // Return other potential errors
	}
	return nil // If the directory already exists, do nothing
}

// 获取运行目录
func GetRunDir() string {
	runDir, _ := os.Executable()
	// 判断路径中是否包含：build，如果包含，则说明是开发环境，返回项目根目录
	if strings.Contains(runDir, "build") {
		runDir, _ := os.Getwd()
		// fmt.Println("开发环境")
		return runDir
	} else {
		// fmt.Println("生产环境")
		// 去掉可执行文件的文件名，返回目录
		return filepath.Dir(runDir)
	}
}

// 获取用户配置目录
func GetUserConfigDir() string {
	userConfigDir, _ := os.UserConfigDir()
	userConfigDir += "/zpic"
	return userConfigDir
}
