package updater

import (
	"fmt"
	"runtime"
	"strings"

	"github.com/creativeprojects/go-selfupdate"
)

// UpdateResult chứa kết quả kiểm tra update.
type UpdateResult struct {
	HasUpdate    bool
	LatestVersion string
	ReleaseNotes string
	URL          string
}

// CheckForUpdate kiểm tra GitHub Releases xem có bản mới không.
// repoURL dạng "owner/repo", ví dụ "o25160526-pip/self-update-go-rust".
func CheckForUpdate(currentVersion, repoURL string) (*UpdateResult, error) {
	source, err := selfupdate.NewGitHubSource(nil)
	if err != nil {
		return nil, fmt.Errorf("create github source: %w", err)
	}

	updater, err := selfupdate.NewUpdater(
		selfupdate.WithSource(source),
	)
	if err != nil {
		return nil, fmt.Errorf("create updater: %w", err)
	}

	// Detect platform
	asset := detectAsset(repoURL)

	latest, found, err := updater.DetectVersion(asset, currentVersion)
	if err != nil {
		return nil, fmt.Errorf("detect version: %w", err)
	}
	if !found {
		return &UpdateResult{HasUpdate: false}, nil
	}

	// Compare semver
	isNewer, err := IsNewer(currentVersion, latest.ReleaseVersion)
	if err != nil {
		return nil, fmt.Errorf("compare version: %w", err)
	}

	return &UpdateResult{
		HasUpdate:     isNewer,
		LatestVersion: latest.ReleaseVersion,
		ReleaseNotes:  latest.ReleaseNotes,
		URL:           latest.URL,
	}, nil
}

// detectAsset trả về asset name pattern cho platform hiện tại.
// Go demo dùng suffix "-go" để không ghi đè với Rust.
func detectAsset(repoURL string) string {
	// Format: go-demo-windows-{arch}
	arch := runtime.GOARCH
	if arch == "amd64" {
		arch = "x64"
	}
	return fmt.Sprintf("go-demo-windows-%s", arch)
}

// PlatformString trả về platform identifier chuẩn.
// Ví dụ: "windows-x86_64"
func PlatformString() string {
	os := runtime.GOOS
	arch := runtime.GOARCH
	if arch == "amd64" {
		arch = "x86_64"
	}
	return fmt.Sprintf("%s-%s", os, arch)
}

// StateString trả về state name cho UI.
type UpdateState string

const (
	StateChecking       UpdateState = "checking"
	StateUpToDate       UpdateState = "up-to-date"
	StateUpdateAvailable UpdateState = "update-available"
	StateDownloading    UpdateState = "downloading"
	StateVerifying      UpdateState = "verifying"
	StateInstalling     UpdateState = "installing"
	StateRestarting     UpdateState = "restarting"
	StateFailed         UpdateState = "failed"
	StateRolledBack     UpdateState = "rolled-back"
)

// IsValidState kiểm tra state string có hợp lệ không.
func IsValidState(s string) bool {
	switch UpdateState(s) {
	case StateChecking, StateUpToDate, StateUpdateAvailable,
		StateDownloading, StateVerifying, StateInstalling,
		StateRestarting, StateFailed, StateRolledBack:
		return true
	}
	return false
}

// MockCheckForUpdate mô phỏng kiểm tra update cho --offline-test mode.
func MockCheckForUpdate(currentVersion string, mockHasUpdate bool) *UpdateResult {
	if mockHasUpdate {
		return &UpdateResult{
			HasUpdate:     true,
			LatestVersion: "1.0.1",
			ReleaseNotes: "Mock: Fix update flow and restart",
			URL:          "mock://artifact-url",
		}
	}
	return &UpdateResult{HasUpdate: false}
}

// ParseReleaseNotes trả về release notes ngắn (giới hạn 200 chars).
func ParseReleaseNotes(notes string) string {
	notes = strings.TrimSpace(notes)
	if len(notes) > 200 {
		return notes[:197] + "..."
	}
	return notes
}
