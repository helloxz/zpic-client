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
		Title:  "Zpic Client",
		Width:  1200,
		Height: 750,
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		BackgroundColour: &options.RGBA{R: 27, G: 38, B: 54, A: 1},
		OnStartup:        app.startup,
		Bind: []interface{}{
			app,
			new_app,
		},
	})

	if err != nil {
		println("Error:", err.Error())
	}
}
