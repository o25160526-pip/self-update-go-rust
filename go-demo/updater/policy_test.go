package updater

import "testing"

func TestDefaultPolicy(t *testing.T) {
	p := DefaultPolicy()
	if p.Channel != "stable" {
		t.Errorf("expected stable, got %s", p.Channel)
	}
	if !p.AutoCheckOnStartup {
		t.Error("expected AutoCheckOnStartup=true")
	}
	if p.RequireUserConfirmation {
		t.Error("expected RequireUserConfirmation=false")
	}
	if !p.RollbackOnStartupFailure {
		t.Error("expected RollbackOnStartupFailure=true")
	}
	if p.MaxRollbackAttempts != 1 {
		t.Errorf("expected MaxRollbackAttempts=1, got %d", p.MaxRollbackAttempts)
	}
}
