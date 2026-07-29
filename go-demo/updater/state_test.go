package updater

import "testing"

func TestStateLifecycle(t *testing.T) {
	s := PersistentState{Current: "1.0.0", LastKnownGood: "1.0.0"}
	if err := StageUpdate(&s, "1.0.1"); err != nil {
		t.Fatal(err)
	}
	if s.Pending != "1.0.1" {
		t.Fatal("not staged")
	}
	if err := MarkHealthy(&s); err != nil {
		t.Fatal(err)
	}
	if s.Current != "1.0.1" || s.LastKnownGood != "1.0.1" || s.Pending != "" {
		t.Fatalf("bad healthy state: %+v", s)
	}
}
func TestRollbackLimited(t *testing.T) {
	s := PersistentState{Current: "1.0.1", Pending: "1.0.1", LastKnownGood: "1.0.0"}
	if err := Rollback(&s, 1); err != nil {
		t.Fatal(err)
	}
	if s.Current != "1.0.0" || s.Status != StateRolledBack {
		t.Fatalf("bad rollback: %+v", s)
	}
	if err := Rollback(&s, 1); err == nil {
		t.Fatal("expected rollback limit")
	}
}
