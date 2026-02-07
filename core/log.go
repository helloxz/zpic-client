package core

import (
	"bufio"
	"os"
)

// 切割日志
func SplitLog() bool {
	// 日志文件路径
	logFilePath := "data/logs/error.log"
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
	// 日志文件路径
	logFilePath := "data/logs/error.log"
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
