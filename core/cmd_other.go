//go:build !windows

package core

import "os/exec"

func hideWindow(cmd *exec.Cmd) {
	// 非 Windows 系统不需要处理
}
