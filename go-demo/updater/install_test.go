package updater

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestInstallOverSelfAndRestore(t *testing.T) {
	dir := t.TempDir()
	exe := filepath.Join(dir, "app.exe")
	if err := os.WriteFile(exe, []byte("v1"), 0o700); err != nil {
		t.Fatal(err)
	}
	newFile := filepath.Join(dir, "staged", "app.exe")
	if err := os.MkdirAll(filepath.Dir(newFile), 0o700); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(newFile, []byte("v2"), 0o700); err != nil {
		t.Fatal(err)
	}

	backup, err := InstallOverSelf(exe, newFile, "1.0.0")
	if err != nil {
		t.Fatalf("install: %v", err)
	}
	if got, _ := os.ReadFile(exe); string(got) != "v2" {
		t.Fatalf("exe chua duoc thay the: %q", got)
	}
	if got, _ := os.ReadFile(backup); string(got) != "v1" {
		t.Fatalf("backup sai noi dung: %q", got)
	}

	if err := RestoreBackup(exe, backup); err != nil {
		t.Fatalf("restore: %v", err)
	}
	if got, _ := os.ReadFile(exe); string(got) != "v1" {
		t.Fatalf("rollback khong dua lai ban cu: %q", got)
	}
}

func TestInstallOverSelfMissingArtifact(t *testing.T) {
	dir := t.TempDir()
	exe := filepath.Join(dir, "app.exe")
	if err := os.WriteFile(exe, []byte("v1"), 0o700); err != nil {
		t.Fatal(err)
	}
	if _, err := InstallOverSelf(exe, filepath.Join(dir, "khong-ton-tai"), "1.0.0"); err == nil {
		t.Fatal("phai loi khi artifact khong ton tai")
	}
	if got, _ := os.ReadFile(exe); string(got) != "v1" {
		t.Fatalf("exe bi hong sau khi install fail: %q", got)
	}
}

func TestHealthMarkerLifecycle(t *testing.T) {
	dir := filepath.Join(t.TempDir(), "health")
	if HealthMarkerExists(dir, "1.0.1") {
		t.Fatal("marker khong nen ton tai luc dau")
	}
	if err := WriteHealthMarker(dir, "1.0.1"); err != nil {
		t.Fatal(err)
	}
	if !WaitForHealthMarker(dir, "1.0.1", time.Second, 10*time.Millisecond) {
		t.Fatal("phai thay marker")
	}
	if WaitForHealthMarker(dir, "1.0.2", 50*time.Millisecond, 10*time.Millisecond) {
		t.Fatal("khong duoc thay marker cua version khac")
	}
	if err := ClearHealthMarker(dir, "1.0.1"); err != nil {
		t.Fatal(err)
	}
	if HealthMarkerExists(dir, "1.0.1") {
		t.Fatal("marker phai bi xoa")
	}
	if err := ClearHealthMarker(dir, "1.0.1"); err != nil {
		t.Fatalf("xoa marker khong ton tai khong duoc loi: %v", err)
	}
}
