// Package updater chứa logic kiểm tra, download, verify, install và rollback update.
package updater

// UpdatePolicy đại diện cho cấu hình update policy local.
// Mapping sang JSON keys camelCase bằng struct tags.
type UpdatePolicy struct {
	Channel                   string `json:"channel"`
	AutoCheckOnStartup        bool   `json:"autoCheckOnStartup"`
	AutoCheckIntervalMinutes  int    `json:"autoCheckIntervalMinutes"`
	AutoDownload              bool   `json:"autoDownload"`
	AutoInstall               bool   `json:"autoInstall"`
	RequireUserConfirmation   bool   `json:"requireUserConfirmation"`
	RestartAutomatically      bool   `json:"restartAutomatically"`
	AllowDowngrade            bool   `json:"allowDowngrade"`
	RollbackOnStartupFailure  bool   `json:"rollbackOnStartupFailure"`
	HealthCheckTimeoutSeconds int    `json:"healthCheckTimeoutSeconds"`
	MaxRollbackAttempts       int    `json:"maxRollbackAttempts"`
	MinSupportedVersion       string `json:"minSupportedVersion"`
}

// DefaultPolicy trả về policy mặc định theo §7 của prompt.
func DefaultPolicy() UpdatePolicy {
	return UpdatePolicy{
		Channel:                   "stable",
		AutoCheckOnStartup:        true,
		AutoCheckIntervalMinutes:  60,
		AutoDownload:              true,
		AutoInstall:               true,
		RequireUserConfirmation:   false,
		RestartAutomatically:      true,
		AllowDowngrade:            false,
		RollbackOnStartupFailure:  true,
		HealthCheckTimeoutSeconds: 30,
		MaxRollbackAttempts:       1,
		MinSupportedVersion:       "1.0.0",
	}
}
