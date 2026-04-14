package core

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"zpic-client/helper"
)

func (ac *AppCore) ClearLogs(ctx context.Context) {
	SplitLog()
}

var logFilePath = helper.GetUserConfigDir() + "/data/logs/error.log"

// 打开日志目录
func (ac *AppCore) OpenLogDirectory() error {
	var cmd *exec.Cmd
	path := helper.GetUserConfigDir() + "/data/logs"

	switch runtime.GOOS {
	case "windows":
		// Windows: explorer [路径]
		cmd = exec.Command("explorer", path)
	case "darwin":
		// macOS: open [路径]
		cmd = exec.Command("open", path)
	case "linux":
		// Linux: xdg-open [路径]
		cmd = exec.Command("xdg-open", path)
	default:
		return fmt.Errorf("不支持的平台: %s", runtime.GOOS)
	}

	// 使用 Start() 而不是 Run()，这样不会阻塞前端界面
	return cmd.Start()
}

// 切割日志
func SplitLog() bool {
	// 只保留日志最后1万行，然后保存

	const maxLines = 10000

	file, err := os.Open(logFilePath)
	if err != nil {
		return false
	}
	defer file.Close()

	lines := make([]string, 0, maxLines)
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		if len(lines) == maxLines {
			copy(lines, lines[1:])
			lines[maxLines-1] = scanner.Text()
		} else {
			lines = append(lines, scanner.Text())
		}
	}
	if err := scanner.Err(); err != nil {
		return false
	}

	tmpFile, err := os.Create(logFilePath)
	if err != nil {
		return false
	}
	defer tmpFile.Close()

	writer := bufio.NewWriter(tmpFile)
	for i, line := range lines {
		if i > 0 {
			if _, err := writer.WriteString("\n"); err != nil {
				return false
			}
		}
		if _, err := writer.WriteString(line); err != nil {
			return false
		}
	}
	if err := writer.Flush(); err != nil {
		return false
	}

	return true
}

// 读取最近的500行日志，并返回给前端
func (ac *AppCore) GetRecentLogs() ResData {
	// 启动一个新的goroutine来切割日志，避免阻塞主线程
	go func() {
		SplitLog()
	}()
	// 读取最后100行日志然后返回
	const maxLines = 100

	file, err := os.Open(logFilePath)
	if err != nil {
		return ResData{
			Status: false,
			Msg:    "打开日志文件失败：" + err.Error(),
			Data:   nil,
		}
	}
	defer file.Close()

	lines := make([]string, 0, maxLines)
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		if len(lines) == maxLines {
			copy(lines, lines[1:])
			lines[maxLines-1] = scanner.Text()
		} else {
			lines = append(lines, scanner.Text())
		}
	}
	if err := scanner.Err(); err != nil {
		return ResData{
			Status: false,
			Msg:    "读取日志文件失败：" + err.Error(),
			Data:   nil,
		}
	}

	return ResData{
		Status: true,
		Msg:    "获取日志成功",
		Data:   lines,
	}
}
