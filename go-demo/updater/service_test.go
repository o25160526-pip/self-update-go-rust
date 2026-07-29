package updater

import (
	"context"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"aead.dev/minisign"
)

// demoPrivateKey đọc private key demo trong repo (keys/) và giải mã.
func demoPrivateKey(t *testing.T) minisign.PrivateKey {
	t.Helper()
	keyBytes, err := os.ReadFile(filepath.Join("..", "..", "keys", "demo-signing.key"))
	if err != nil {
		t.Fatalf("doc demo private key: %v", err)
	}
	pwBytes, err := os.ReadFile(filepath.Join("..", "..", "keys", "demo-signing.password"))
	if err != nil {
		t.Fatalf("doc demo password: %v", err)
	}
	key, err := minisign.DecryptKey(strings.TrimSpace(string(pwBytes)), []byte(strings.TrimSpace(string(keyBytes))))
	if err != nil {
		t.Fatalf("giai ma demo private key: %v", err)
	}
	return key
}

func TestDemoKeyKhopPinnedPublicKey(t *testing.T) {
	key := demoPrivateKey(t)
	payload := []byte("go-demo artifact")
	sig := base64.StdEncoding.EncodeToString(minisign.Sign(key, payload))
	if err := VerifyPinnedSignature(payload, sig); err != nil {
		t.Fatalf("demo key khong khop public key pin trong keys.go: %v", err)
	}
	if err := VerifyPinnedSignature([]byte("artifact bi sua"), sig); err == nil {
		t.Fatal("data bi sua ma signature van pass")
	}
}

// setupService dựng 1 HTTPS server phục vụ manifest + artifact đã ký thật.
func setupService(t *testing.T, artifact []byte, newVersion string) (*Service, string) {
	t.Helper()
	dir := t.TempDir()
	exe := filepath.Join(dir, "app.exe")
	if err := os.WriteFile(exe, []byte("OLD BUILD"), 0o700); err != nil {
		t.Fatal(err)
	}
	sig := base64.StdEncoding.EncodeToString(minisign.Sign(demoPrivateKey(t), artifact))
	sum := sha256.Sum256(artifact)

	var srv *httptest.Server
	mux := http.NewServeMux()
	mux.HandleFunc("/artifact", func(w http.ResponseWriter, _ *http.Request) {
		_, _ = w.Write(artifact)
	})
	mux.HandleFunc("/manifest", func(w http.ResponseWriter, _ *http.Request) {
		m := Manifest{
			Version:             newVersion,
			Channel:             "stable",
			PublishedAt:         time.Now().UTC().Format(time.RFC3339),
			ReleaseNotes:        "ban test",
			MinSupportedVersion: "1.0.0",
			Platform:            PlatformString(),
			URL:                 srv.URL + "/artifact",
			SHA256:              hex.EncodeToString(sum[:]),
			Signature:           sig,
			Size:                int64(len(artifact)),
		}
		_ = json.NewEncoder(w).Encode(m)
	})
	srv = httptest.NewTLSServer(mux)
	t.Cleanup(srv.Close)

	svc := &Service{
		Version:     "1.0.0",
		ExePath:     exe,
		Dir:         dir,
		ManifestURL: srv.URL + "/manifest",
		Policy:      DefaultPolicy(),
		Client:      srv.Client(),
		Spawn:       func(string, ...string) error { return nil },
		Exit:        func(int) {},
	}
	svc.Policy.HealthCheckTimeoutSeconds = 1
	return svc, exe
}

func TestServiceUpdateThanhCong(t *testing.T) {
	artifact := []byte("NEW BUILD 1.0.1")
	svc, exe := setupService(t, artifact, "1.0.1")
	exitCode := -1
	svc.Exit = func(c int) { exitCode = c }
	svc.Spawn = func(string, ...string) error {
		// giả lập bản mới khởi động thành công
		return WriteHealthMarker(HealthDir(svc.Dir), "1.0.1")
	}

	res, err := svc.CheckAndUpdate(context.Background())
	if err != nil {
		t.Fatalf("update loi: %v", err)
	}
	if res == nil || !res.HasUpdate || res.LatestVersion != "1.0.1" {
		t.Fatalf("khong phat hien ban moi: %+v", res)
	}
	if got, _ := os.ReadFile(exe); string(got) != string(artifact) {
		t.Fatalf("exe chua duoc thay the: %q", got)
	}
	if exitCode != 0 {
		t.Fatalf("mong doi thoat 0 sau khi ban moi healthy, nhan %d", exitCode)
	}
	st, err := LoadState(StatePath(svc.Dir))
	if err != nil {
		t.Fatal(err)
	}
	if st.Pending != "1.0.1" || st.UpdatedFrom != "1.0.0" {
		t.Fatalf("state pending sai: %+v", st)
	}
	if _, err := os.Stat(BackupPath(exe, "1.0.0")); err != nil {
		t.Fatalf("thieu backup de rollback: %v", err)
	}

	// Tiến trình mới khởi động: health-check phải commit last-known-good.
	newProc := &Service{
		Version: "1.0.1",
		ExePath: exe,
		Dir:     svc.Dir,
		Policy:  DefaultPolicy(),
		Spawn:   func(string, ...string) error { return nil },
		Exit:    func(int) {},
	}
	if err := newProc.Startup(); err != nil {
		t.Fatalf("startup ban moi: %v", err)
	}
	st, err = LoadState(StatePath(svc.Dir))
	if err != nil {
		t.Fatal(err)
	}
	if st.LastKnownGood != "1.0.1" || st.Pending != "" || st.Current != "1.0.1" {
		t.Fatalf("health-check chua commit last-known-good: %+v", st)
	}
}

func TestServiceRollbackKhiBanMoiKhongLen(t *testing.T) {
	artifact := []byte("BROKEN BUILD 1.0.1")
	svc, exe := setupService(t, artifact, "1.0.1")
	exitCode := -1
	svc.Exit = func(c int) { exitCode = c }
	spawns := 0
	svc.Spawn = func(string, ...string) error { spawns++; return nil } // khong bao giờ healthy

	if _, err := svc.CheckAndUpdate(context.Background()); err != nil {
		t.Fatalf("khong mong doi loi: %v", err)
	}
	if got, _ := os.ReadFile(exe); string(got) != "OLD BUILD" {
		t.Fatalf("chua rollback ve ban cu: %q", got)
	}
	if svc.State() != string(StateRolledBack) {
		t.Fatalf("state = %s, mong doi rolled-back", svc.State())
	}
	if spawns != 2 {
		t.Fatalf("mong doi spawn 2 lan (ban moi + ban cu), nhan %d", spawns)
	}
	if exitCode != 1 {
		t.Fatalf("exit code = %d", exitCode)
	}
	st, err := LoadState(StatePath(svc.Dir))
	if err != nil {
		t.Fatal(err)
	}
	if st.Current != "1.0.0" || st.RollbackAttempts != 1 {
		t.Fatalf("state sau rollback sai: %+v", st)
	}
}

func TestServiceUpToDate(t *testing.T) {
	svc, exe := setupService(t, []byte("SAME BUILD"), "1.0.0")
	res, err := svc.CheckAndUpdate(context.Background())
	if err != nil {
		t.Fatalf("check loi: %v", err)
	}
	if res.HasUpdate {
		t.Fatalf("cung version ma bao co update: %+v", res)
	}
	if svc.State() != string(StateUpToDate) {
		t.Fatalf("state = %s", svc.State())
	}
	if got, _ := os.ReadFile(exe); string(got) != "OLD BUILD" {
		t.Fatalf("khong duoc thay doi exe: %q", got)
	}
}

func TestServiceOfflineTest(t *testing.T) {
	svc, _ := setupService(t, []byte("x"), "1.0.1")
	svc.Offline = true
	res, err := svc.CheckAndUpdate(context.Background())
	if err != nil {
		t.Fatalf("offline check loi: %v", err)
	}
	if res.HasUpdate {
		t.Fatal("offline mode khong duoc bao co update")
	}
	if svc.State() != string(StateUpToDate) {
		t.Fatalf("state = %s", svc.State())
	}
}
