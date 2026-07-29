package main

import (
	"context"
	"fmt"
	"go-demo/internal/version"
	"runtime"
)

// App struct
type App struct {
	ctx context.Context
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
}

// GetVersion returns the current app version.
func (a *App) GetVersion() string {
	return version.Version
}

// GetInfo returns app name, version, OS, and architecture.
func (a *App) GetInfo() map[string]string {
	return map[string]string{
		"app":     "go-demo",
		"version": version.Version,
		"os":      runtime.GOOS,
		"arch":    runtime.GOARCH,
	}
}

// GetUpdateState returns the current updater state.
// States: checking, up-to-date, update-available, downloading, verifying,
// installing, restarting, failed, rolled-back
func (a *App) GetUpdateState() string {
	// P2: will integrate with real updater
	return "checking"
}

// Greet returns a greeting for the given name
func (a *App) Greet(name string) string {
	return fmt.Sprintf("Hello %s, It's show time!", name)
}
