package main

import (
	"embed"

	"zpic-client/core"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
)

//go:embed all:frontend/dist
var assets embed.FS

func main() {
	// Create an instance of the app structure
	app := NewApp()
	new_app := core.NewAppCore()

	// Create application with options
	err := wails.Run(&options.App{
		Title:     "Zpic Client",
		Width:     1200,
		Height:    750,
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
