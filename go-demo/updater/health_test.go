package updater

import "testing"

func TestEvaluateStartupHealthy(t *testing.T) {
	s := PersistentState{Current: "1.0.0", Pending: "1.0.1", LastKnownGood: "1.0.0"}
	outcome, err := EvaluateStartup(&s, "1.0.1", 1)
	if err != nil {
		t.Fatalf("khong mong doi loi: %v", err)
	}
	if outcome != OutcomeHealthy {
		t.Fatalf("outcome = %s", outcome)
	}
	if s.LastKnownGood != "1.0.1" || s.Current != "1.0.1" || s.Pending != "" {
		t.Fatalf("state sau health-check sai: %+v", s)
	}
}

func TestEvaluateStartupRollback(t *testing.T) {
	s := PersistentState{Current: "1.0.1", Pending: "1.0.1", LastKnownGood: "1.0.0"}
	outcome, err := EvaluateStartup(&s, "1.0.0", 1)
	if err != nil {
		t.Fatalf("khong mong doi loi: %v", err)
	}
	if outcome != OutcomeRollback {
		t.Fatalf("outcome = %s", outcome)
	}
	if s.Current != "1.0.0" || s.RollbackAttempts != 1 || s.Status != StateRolledBack {
		t.Fatalf("state sau rollback sai: %+v", s)
	}
}

func TestEvaluateStartupKhongRollbackLapVoHan(t *testing.T) {
	s := PersistentState{Current: "1.0.1", Pending: "1.0.1", LastKnownGood: "1.0.0", RollbackAttempts: 1}
	if _, err := EvaluateStartup(&s, "1.0.0", 1); err == nil {
		t.Fatal("phai chan khi vuot maxRollbackAttempts")
	}
}

func TestEvaluateStartupNormal(t *testing.T) {
	s := PersistentState{}
	outcome, err := EvaluateStartup(&s, "1.0.0", 1)
	if err != nil {
		t.Fatalf("khong mong doi loi: %v", err)
	}
	if outcome != OutcomeNormal {
		t.Fatalf("outcome = %s", outcome)
	}
	if s.Current != "1.0.0" || s.LastKnownGood != "1.0.0" {
		t.Fatalf("state sai: %+v", s)
	}
}
