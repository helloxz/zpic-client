package main

import (
	"embed"
	"runtime"

	"zpic-client/core"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
)

//go:embed all:frontend/dist
var assets embed.FS

//go:embed bin/*
var embeddedBin embed.FS

func main() {
	// Create an instance of the app structure
	app := NewApp()
	new_app := core.NewAppCore()

	width := 1200
	height := 750
	// Windows 补偿标题栏和边框
	if runtime.GOOS == "windows" {
		height += 40 // 补偿约 30px
	}

	// Create application with options
	err := wails.Run(&options.App{
		Title:     "Zpic Client",
		Width:     width,
		Height:    height,
		MinWidth:  1200,
		MinHeight: 750,
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		// 启用单实例锁
		SingleInstanceLock: &options.SingleInstanceLock{
			UniqueId: "zpic-client", // 唯一标识符
		},
		// ============ 无边框窗口（可选） ============
		Frameless:        false,
		DisableResize:    false,
		BackgroundColour: &options.RGBA{R: 27, G: 38, B: 54, A: 1},
		OnStartup:        app.startup,
		OnShutdown:       new_app.ClearLogs,
		Bind: []interface{}{
			app,
			new_app,
		},
	})

	if err != nil {
		println("Error:", err.Error())
	}
}
