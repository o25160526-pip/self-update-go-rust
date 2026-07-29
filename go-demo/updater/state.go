package updater

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

type PersistentState struct {
	Current          string      `json:"current"`
	Pending          string      `json:"pending,omitempty"`
	LastKnownGood    string      `json:"lastKnownGood"`
	UpdatedFrom      string      `json:"updatedFrom,omitempty"`
	RollbackAttempts int         `json:"rollbackAttempts"`
	Status           UpdateState `json:"status"`
	UpdatedAt        time.Time   `json:"updatedAt"`
}

func LoadState(path string) (PersistentState, error) {
	b, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		return PersistentState{Status: StateUpToDate}, nil
	}
	if err != nil {
		return PersistentState{}, err
	}
	var s PersistentState
	if err := json.Unmarshal(b, &s); err != nil {
		return s, err
	}
	return s, nil
}
func SaveState(path string, s PersistentState) error {
	s.UpdatedAt = time.Now().UTC()
	b, err := json.MarshalIndent(s, "", "  ")
	if err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(path), 0700); err != nil {
		return err
	}
	tmp := path + ".tmp"
	if err := os.WriteFile(tmp, b, 0600); err != nil {
		return err
	}
	return os.Rename(tmp, path)
}
func StageUpdate(s *PersistentState, version string) error {
	if version == "" {
		return fmt.Errorf("pending version is empty")
	}
	s.Pending = version
	s.UpdatedFrom = s.Current
	s.Status = StateInstalling
	return nil
}
func MarkHealthy(s *PersistentState) error {
	if s.Pending == "" {
		return fmt.Errorf("no pending version")
	}
	s.Current = s.Pending
	s.LastKnownGood = s.Pending
	s.Pending = ""
	s.RollbackAttempts = 0
	s.Status = StateUpToDate
	return nil
}
func Rollback(s *PersistentState, max int) error {
	if s.RollbackAttempts >= max {
		return fmt.Errorf("maximum rollback attempts reached")
	}
	if s.LastKnownGood == "" {
		return fmt.Errorf("last-known-good is empty")
	}
	s.Current = s.LastKnownGood
	s.Pending = ""
	s.RollbackAttempts++
	s.Status = StateRolledBack
	return nil
}
