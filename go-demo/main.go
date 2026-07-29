package main

import (
	"context"
	"embed"
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"

	"go-demo/internal/version"
	"go-demo/updater"
)

//go:embed all:frontend/dist
var assets embed.FS

func main() {
	versionFlag := flag.Bool("version", false, "In version roi thoat")
	printStateFlag := flag.Bool("print-update-state", false, "In trang thai updater roi thoat")
	offlineTestFlag := flag.Bool("offline-test", false, "Dung mock provider, khong goi mang")
	checkOnlyFlag := flag.Bool("check-only", false, "Chay luong update o che do CLI roi thoat")
	updatedFrom := flag.String("updated-from", "", "Version truoc khi update (do tien trinh cu truyen vao)")
	rolledBack := flag.Bool("rolled-back", false, "Tien trinh nay chay sau khi rollback")
	flag.Parse()

	if *versionFlag {
		fmt.Println(version.Version)
		return
	}

	svc, err := updater.NewService(version.Version, *offlineTestFlag)
	if err != nil {
		fmt.Fprintln(os.Stderr, "[error] khoi tao updater:", err)
		os.Exit(1)
	}

	if *printStateFlag {
		fmt.Println(svc.PersistedStatus())
		return
	}

	if *updatedFrom != "" {
		svc.LogUpdatedFrom(*updatedFrom)
	}
	if *rolledBack {
		svc.LogRolledBack()
	}
	if err := svc.Startup(); err != nil {
		fmt.Fprintln(os.Stderr, "[warn] startup health-check:", err)
	}

	if *checkOnlyFlag {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
		defer cancel()
		if _, err := svc.CheckAndUpdate(ctx); err != nil {
			fmt.Fprintln(os.Stderr, "[error]", err)
		}
		for _, line := range svc.Logs() {
			fmt.Println(line)
		}
		fmt.Println(svc.State())
		return
	}

	app := NewApp(svc)
	if err := wails.Run(&options.App{
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
	}); err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err)
	}
}
