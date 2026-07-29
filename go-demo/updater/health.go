package updater

// StartupOutcome là kết quả health-check lúc khởi động.
type StartupOutcome string

const (
	// OutcomeNormal: khởi động bình thường, không có update đang chờ.
	OutcomeNormal StartupOutcome = "normal"
	// OutcomeHealthy: bản mới khởi động thành công -> ghi last-known-good.
	OutcomeHealthy StartupOutcome = "healthy"
	// OutcomeRollback: version đang chạy không phải bản pending -> đã rollback.
	OutcomeRollback StartupOutcome = "rolled-back"
)

// EvaluateStartup so sánh version đang chạy với state trên đĩa để quyết định:
//
//   - không có pending          -> normal (đồng bộ current/last-known-good)
//   - pending == đang chạy      -> healthy (commit last-known-good)
//   - pending != đang chạy      -> bản mới không lên được, rollback
//
// maxRollback chặn rollback lặp vô hạn (mặc định 1 lần).
func EvaluateStartup(s *PersistentState, runningVersion string, maxRollback int) (StartupOutcome, error) {
	if s.Pending == "" {
		if s.Current != runningVersion {
			s.Current = runningVersion
			s.Status = StateUpToDate
		}
		if s.LastKnownGood == "" {
			s.LastKnownGood = runningVersion
		}
		return OutcomeNormal, nil
	}
	if s.Pending == runningVersion {
		if err := MarkHealthy(s); err != nil {
			return OutcomeNormal, err
		}
		return OutcomeHealthy, nil
	}
	if err := Rollback(s, maxRollback); err != nil {
		return OutcomeNormal, err
	}
	return OutcomeRollback, nil
}
