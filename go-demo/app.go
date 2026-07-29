package main

import (
	"context"
	"fmt"
	"runtime"
	"time"

	"go-demo/internal/version"
	"go-demo/updater"
)

// App là struct được bind vào frontend Wails.
type App struct {
	ctx context.Context
	svc *updater.Service
}

// NewApp tạo App với updater service đã khởi tạo.
func NewApp(svc *updater.Service) *App {
	return &App{svc: svc}
}

// startup được Wails gọi khi app khởi động. Tự kiểm tra update ở background.
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
	if !a.svc.Policy.AutoCheckOnStartup {
		return
	}
	go func() {
		checkCtx, cancel := context.WithTimeout(ctx, 10*time.Minute)
		defer cancel()
		_, _ = a.svc.CheckAndUpdate(checkCtx)
	}()
}

// GetVersion trả về version của bản build đang chạy.
func (a *App) GetVersion() string {
	return version.Version
}

// GetInfo trả về app name, version, OS, arch và public key đang pin.
func (a *App) GetInfo() map[string]string {
	return map[string]string{
		"app":       "go-demo",
		"version":   version.Version,
		"os":        runtime.GOOS,
		"arch":      runtime.GOARCH,
		"keyId":     updater.PinnedKeyID,
		"endpoint":  updater.ManifestURL,
	}
}

// GetUpdateState trả về trạng thái updater hiện tại.
// checking | up-to-date | update-available | downloading | verifying |
// installing | restarting | failed | rolled-back
func (a *App) GetUpdateState() string {
	return a.svc.State()
}

// GetLogs trả về log updater để hiển thị trong UI.
func (a *App) GetLogs() []string {
	return a.svc.Logs()
}

// CheckForUpdate cho user bấm kiểm tra thủ công.
func (a *App) CheckForUpdate() string {
	ctx, cancel := context.WithTimeout(a.ctx, 5*time.Minute)
	defer cancel()
	_, _ = a.svc.CheckAndUpdate(ctx)
	return a.svc.State()
}

// Greet giữ lại từ template Wails.
func (a *App) Greet(name string) string {
	return fmt.Sprintf("Hello %s, It's show time!", name)
}
