package main

import (
	"embed"
	"flag"
	"fmt"
	"os"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"

	"go-demo/internal/version"
)

//go:embed all:frontend/dist
var assets embed.FS

func main() {
	// Parse CLI flags
	versionFlag := flag.Bool("version", false, "Print version and exit")
	printStateFlag := flag.Bool("print-update-state", false, "Print update state and exit")
	offlineTestFlag := flag.Bool("offline-test", false, "Use mock provider for offline testing")
	flag.Parse()

	// Handle --version
	if *versionFlag {
		fmt.Println(version.Version)
		os.Exit(0)
	}

	// Handle --print-update-state
	if *printStateFlag {
		// P2: will integrate with real updater state
		fmt.Println("checking")
		os.Exit(0)
	}

	// --offline-test flag is stored for updater to use (P2)
	if *offlineTestFlag {
		fmt.Fprintln(os.Stderr, "[info] offline-test mode enabled")
	}

	// Create an instance of the app structure
	app := NewApp()

	// Create application with options
	err := wails.Run(&options.App{
		Title:  "Go Demo - Auto Update",
		Width:  800,
		Height: 600,
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		BackgroundColour: &options.RGBA{R: 27, G: 38, B: 54, A: 1},
		OnStartup:        app.startup,
		Bind: []interface{}{
			app,
		},
	})

	if err != nil {
		println("Error:", err.Error())
	}
}
