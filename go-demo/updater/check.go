package updater

import (
	"context"
	"fmt"
	"net/http"
	"runtime"
	"strings"
	"time"
)

type UpdateResult struct {
	HasUpdate     bool      `json:"hasUpdate"`
	LatestVersion string    `json:"latestVersion,omitempty"`
	ReleaseNotes  string    `json:"releaseNotes,omitempty"`
	URL           string    `json:"url,omitempty"`
	Manifest      *Manifest `json:"manifest,omitempty"`
}

type Checker struct {
	Client *http.Client
}

func NewChecker() *Checker {
	return &Checker{Client: &http.Client{Timeout: 20 * time.Second}}
}

func (c *Checker) Check(ctx context.Context, manifestURL, currentVersion string, policy UpdatePolicy) (*UpdateResult, error) {
	m, err := FetchManifest(ctx, c.Client, manifestURL)
	if err != nil {
		return nil, err
	}
	if err := m.Validate(); err != nil {
		return nil, err
	}
	if m.Channel != policy.Channel {
		return &UpdateResult{}, nil
	}
	if m.Platform != PlatformString() {
		return &UpdateResult{}, nil
	}
	newer, err := IsNewer(currentVersion, m.Version)
	if err != nil {
		return nil, fmt.Errorf("compare version: %w", err)
	}
	if !newer && !policy.AllowDowngrade {
		return &UpdateResult{}, nil
	}
	return &UpdateResult{HasUpdate: true, LatestVersion: m.Version, ReleaseNotes: ParseReleaseNotes(m.ReleaseNotes), URL: m.URL, Manifest: m}, nil
}

func PlatformString() string {
	arch := runtime.GOARCH
	if arch == "amd64" {
		arch = "x86_64"
	}
	return fmt.Sprintf("%s-%s", runtime.GOOS, arch)
}

type UpdateState string

const (
	StateChecking        UpdateState = "checking"
	StateUpToDate        UpdateState = "up-to-date"
	StateUpdateAvailable UpdateState = "update-available"
	StateDownloading     UpdateState = "downloading"
	StateVerifying       UpdateState = "verifying"
	StateInstalling      UpdateState = "installing"
	StateRestarting      UpdateState = "restarting"
	StateFailed          UpdateState = "failed"
	StateRolledBack      UpdateState = "rolled-back"
)

func IsValidState(s string) bool {
	switch UpdateState(s) {
	case StateChecking, StateUpToDate, StateUpdateAvailable, StateDownloading, StateVerifying, StateInstalling, StateRestarting, StateFailed, StateRolledBack:
		return true
	}
	return false
}
func MockCheckForUpdate(currentVersion string, hasUpdate bool) *UpdateResult {
	if !hasUpdate {
		return &UpdateResult{}
	}
	return &UpdateResult{HasUpdate: true, LatestVersion: "1.0.1", ReleaseNotes: "Mock: Fix update flow and restart", URL: "mock://artifact-url"}
}
func ParseReleaseNotes(notes string) string {
	notes = strings.TrimSpace(notes)
	if len(notes) > 200 {
		return notes[:197] + "..."
	}
	return notes
}
