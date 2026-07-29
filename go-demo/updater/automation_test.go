package updater

import "testing"

func TestAutomatedVersionUpgrade103To104(t *testing.T) {
	if newer, err := IsNewer("1.0.3", "1.0.4"); err != nil || !newer {
		t.Fatalf("1.0.3 -> 1.0.4 must be an update, newer=%v err=%v", newer, err)
	}
}

func TestAutomatedRollbackAfterPending104(t *testing.T) {
	state := PersistentState{Current: "1.0.3", LastKnownGood: "1.0.3"}
	if err := StageUpdate(&state, "1.0.4"); err != nil {
		t.Fatal(err)
	}
	if state.Pending != "1.0.4" || state.UpdatedFrom != "1.0.3" {
		t.Fatalf("bad pending state: %+v", state)
	}
	if err := Rollback(&state, 1); err != nil {
		t.Fatal(err)
	}
	if state.Current != "1.0.3" || state.Status != StateRolledBack || state.RollbackAttempts != 1 {
		t.Fatalf("bad rollback state: %+v", state)
	}
}
