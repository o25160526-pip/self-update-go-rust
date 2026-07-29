package updater

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// RollbackPrevious restores the newest local backup created by the updater.
// It never downloads a binary, so manual rollback remains local-only.
func (s *Service) RollbackPrevious() error {
	backups, err := filepath.Glob(s.ExePath + ".backup-*")
	if err != nil { return err }
	if len(backups) == 0 { return fmt.Errorf("no previous local version is available") }
	sort.Slice(backups, func(i, j int) bool {
		a, _ := os.Stat(backups[i]); b, _ := os.Stat(backups[j])
		return a.ModTime().After(b.ModTime())
	})
	backup := backups[0]
	previous := strings.TrimPrefix(backup, s.ExePath+".backup-")
	if previous == "" { return fmt.Errorf("previous version metadata is missing") }
	if err := RestoreBackup(s.ExePath, backup); err != nil { return err }
	st, err := LoadState(StatePath(s.Dir)); if err != nil { return err }
	st.Current = previous; st.LastKnownGood = previous; st.Pending = ""; st.Status = StateRolledBack
	if err := SaveState(StatePath(s.Dir), st); err != nil { return err }
	s.setState(StateRolledBack)
	s.logf("manual rollback requested: restored %s", previous)
	if err := s.Spawn(s.ExePath, "--rolled-back"); err != nil { return err }
	s.Exit(0)
	return nil
}
