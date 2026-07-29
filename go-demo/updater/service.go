package updater

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"sync"
	"time"
)

// Service gom toàn bộ luồng update của app Go: check -> download -> verify ->
// install -> restart -> health-check -> rollback.
type Service struct {
	Version     string
	ExePath     string
	Dir         string
	ManifestURL string
	Policy      UpdatePolicy
	Offline     bool
	Client      *http.Client

	// Hook để test thay thế (mặc định: Relaunch / os.Exit).
	Spawn func(exePath string, args ...string) error
	Exit  func(code int)

	mu    sync.Mutex
	state UpdateState
	logs  []string
}

// NewService khởi tạo service với đường dẫn thật của tiến trình đang chạy.
func NewService(version string, offline bool) (*Service, error) {
	dir, err := DataDir()
	if err != nil {
		return nil, err
	}
	exe, err := os.Executable()
	if err != nil {
		return nil, err
	}
	return &Service{
		Version:     version,
		ExePath:     exe,
		Dir:         dir,
		ManifestURL: ManifestURL,
		Policy:      DefaultPolicy(),
		Offline:     offline,
		Client:      &http.Client{Timeout: 30 * time.Second},
		Spawn:       Relaunch,
		Exit:        os.Exit,
		state:       StateChecking,
	}, nil
}

// State trả về trạng thái updater hiện tại (một trong 9 state ở §6).
func (s *Service) State() string {
	s.mu.Lock()
	defer s.mu.Unlock()
	return string(s.state)
}

// Logs trả về log updater (tối đa 200 dòng gần nhất).
func (s *Service) Logs() []string {
	s.mu.Lock()
	defer s.mu.Unlock()
	out := make([]string, len(s.logs))
	copy(out, s.logs)
	return out
}

func (s *Service) setState(st UpdateState) {
	s.mu.Lock()
	s.state = st
	s.mu.Unlock()
	s.logf("state=%s", st)
}

func (s *Service) logf(format string, args ...any) {
	line := fmt.Sprintf("[%s] %s", time.Now().UTC().Format("15:04:05"), fmt.Sprintf(format, args...))
	s.mu.Lock()
	s.logs = append(s.logs, line)
	if len(s.logs) > 200 {
		s.logs = s.logs[len(s.logs)-200:]
	}
	s.mu.Unlock()
}

// LogUpdatedFrom ghi lại metadata bản trước (flag --updated-from).
func (s *Service) LogUpdatedFrom(previous string) {
	s.logf("khoi dong sau update, updated-from=%s", previous)
}

// LogRolledBack ghi lại việc tiến trình này chạy sau khi rollback.
func (s *Service) LogRolledBack() {
	s.logf("khoi dong sau rollback ve %s", s.Version)
	s.setState(StateRolledBack)
}

// PersistedStatus đọc status đã lưu trên đĩa (dùng cho --print-update-state).
func (s *Service) PersistedStatus() string {
	st, err := LoadState(StatePath(s.Dir))
	if err != nil || st.Status == "" {
		return string(StateChecking)
	}
	return string(st.Status)
}

// Startup chạy health-check lúc app khởi động: commit last-known-good nếu bản
// mới lên được, hoặc rollback nếu đang chạy version không phải bản pending.
func (s *Service) Startup() error {
	path := StatePath(s.Dir)
	st, err := LoadState(path)
	if err != nil {
		s.logf("doc state loi, dung state rong: %v", err)
		st = PersistentState{}
	}
	outcome, evalErr := EvaluateStartup(&st, s.Version, s.Policy.MaxRollbackAttempts)
	if evalErr != nil {
		s.logf("health-check: %v", evalErr)
	}
	switch outcome {
	case OutcomeHealthy:
		s.logf("health-check OK, last-known-good=%s", st.LastKnownGood)
		s.setState(StateUpToDate)
	case OutcomeRollback:
		s.logf("bản pending khong khoi dong duoc, da rollback ve %s", st.LastKnownGood)
		s.setState(StateRolledBack)
	default:
		s.setState(StateUpToDate)
	}
	if err := WriteHealthMarker(HealthDir(s.Dir), s.Version); err != nil {
		s.logf("ghi health marker loi: %v", err)
	}
	return SaveState(path, st)
}

// CheckAndUpdate thực hiện §8 bước 3-14: kiểm tra, và nếu policy cho phép thì
// tải, verify, cài, restart luôn (không cần user bấm gì).
func (s *Service) CheckAndUpdate(ctx context.Context) (*UpdateResult, error) {
	s.setState(StateChecking)
	if s.Offline {
		s.logf("offline-test: dung mock provider, khong goi mang")
		s.setState(StateUpToDate)
		return MockCheckForUpdate(s.Version, false), nil
	}
	res, err := (&Checker{Client: s.Client}).Check(ctx, s.ManifestURL, s.Version, s.Policy)
	if err != nil {
		s.setState(StateFailed)
		s.logf("kiem tra update loi: %v", err)
		return nil, err
	}
	if !res.HasUpdate {
		s.logf("da la ban moi nhat (%s)", s.Version)
		s.setState(StateUpToDate)
		return res, nil
	}
	s.setState(StateUpdateAvailable)
	s.logf("co ban moi: %s | %s", res.LatestVersion, res.ReleaseNotes)
	if !s.Policy.AutoDownload {
		s.logf("policy autoDownload=false, dung lai")
		return res, nil
	}
	if err := s.applyUpdate(ctx, res); err != nil {
		s.setState(StateFailed)
		s.logf("cap nhat that bai: %v", err)
		return res, err
	}
	return res, nil
}

func (s *Service) applyUpdate(ctx context.Context, res *UpdateResult) error {
	if res.Manifest == nil {
		return fmt.Errorf("manifest rong")
	}
	s.setState(StateDownloading)
	dest := filepath.Join(DownloadDir(s.Dir), res.LatestVersion, filepath.Base(s.ExePath))
	if err := DownloadAndVerify(ctx, s.Client, *res.Manifest, dest, VerifyPinnedSignature); err != nil {
		return err
	}
	s.setState(StateVerifying)
	s.logf("size + SHA-256 + Minisign OK (key %s)", PinnedKeyID)
	if !s.Policy.AutoInstall {
		s.logf("policy autoInstall=false, dung lai sau khi verify")
		return nil
	}

	s.setState(StateInstalling)
	statePath := StatePath(s.Dir)
	st, err := LoadState(statePath)
	if err != nil {
		return err
	}
	if st.Current == "" {
		st.Current = s.Version
	}
	if st.LastKnownGood == "" {
		st.LastKnownGood = s.Version
	}
	if err := StageUpdate(&st, res.LatestVersion); err != nil {
		return err
	}
	if err := SaveState(statePath, st); err != nil {
		return err
	}
	backup, err := InstallOverSelf(s.ExePath, dest, s.Version)
	if err != nil {
		return err
	}
	s.logf("da cai %s, backup=%s", res.LatestVersion, filepath.Base(backup))
	if !s.Policy.RestartAutomatically {
		return nil
	}

	s.setState(StateRestarting)
	_ = ClearHealthMarker(HealthDir(s.Dir), res.LatestVersion)
	if err := s.Spawn(s.ExePath, "--updated-from", s.Version); err != nil {
		_ = RestoreBackup(s.ExePath, backup)
		return fmt.Errorf("khong khoi dong duoc ban moi: %w", err)
	}
	timeout := time.Duration(s.Policy.HealthCheckTimeoutSeconds) * time.Second
	if WaitForHealthMarker(HealthDir(s.Dir), res.LatestVersion, timeout, 300*time.Millisecond) {
		s.logf("ban %s khoi dong OK, thoat tien trinh cu", res.LatestVersion)
		s.Exit(0)
		return nil
	}

	s.logf("ban moi khong bao healthy trong %s -> rollback", timeout)
	if err := RestoreBackup(s.ExePath, backup); err != nil {
		return err
	}
	if st, err := LoadState(statePath); err == nil {
		if rbErr := Rollback(&st, s.Policy.MaxRollbackAttempts); rbErr != nil {
			s.logf("rollback state: %v", rbErr)
		}
		if saveErr := SaveState(statePath, st); saveErr != nil {
			s.logf("luu state sau rollback: %v", saveErr)
		}
	}
	s.setState(StateRolledBack)
	if err := s.Spawn(s.ExePath, "--rolled-back"); err != nil {
		return fmt.Errorf("khong khoi dong lai ban cu: %w", err)
	}
	s.Exit(1)
	return nil
}
