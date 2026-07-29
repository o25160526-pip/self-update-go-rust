package updater

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

// BackupPath là vị trí lưu bản exe cũ để rollback.
func BackupPath(exePath, version string) string {
	safe := strings.NewReplacer(string(os.PathSeparator), "_", "/", "_", "\\", "_").Replace(version)
	return exePath + ".backup-" + safe
}

// CopyFile copy file (tạo thư mục đích nếu cần) với quyền thực thi.
func CopyFile(src, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()
	if err := os.MkdirAll(filepath.Dir(dst), 0o755); err != nil {
		return err
	}
	out, err := os.OpenFile(dst, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0o755)
	if err != nil {
		return err
	}
	if _, err := io.Copy(out, in); err != nil {
		out.Close()
		return err
	}
	return out.Close()
}

// InstallOverSelf cài bản mới vào đúng vị trí exe đang chạy.
//
// Trên Windows không xoá được file .exe đang chạy, nhưng ĐỔI TÊN thì được.
// Vì vậy: rename exe hiện tại thành backup (tiến trình đang chạy vẫn sống),
// rồi copy artifact mới vào đường dẫn gốc. Nếu copy lỗi thì rollback ngay.
func InstallOverSelf(exePath, newFile, currentVersion string) (string, error) {
	if _, err := os.Stat(newFile); err != nil {
		return "", fmt.Errorf("artifact moi khong ton tai: %w", err)
	}
	backup := BackupPath(exePath, currentVersion)
	_ = os.Remove(backup)
	if err := os.Rename(exePath, backup); err != nil {
		return "", fmt.Errorf("doi ten exe hien tai: %w", err)
	}
	if err := CopyFile(newFile, exePath); err != nil {
		_ = os.Rename(backup, exePath)
		return "", fmt.Errorf("copy artifact moi: %w", err)
	}
	return backup, nil
}

// RestoreBackup đưa bản backup trở lại vị trí exe (dùng khi rollback).
func RestoreBackup(exePath, backup string) error {
	if _, err := os.Stat(backup); err != nil {
		return fmt.Errorf("khong tim thay backup: %w", err)
	}
	broken := exePath + ".failed"
	_ = os.Remove(broken)
	if err := os.Rename(exePath, broken); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("doi ten exe loi: %w", err)
	}
	if err := os.Rename(backup, exePath); err != nil {
		return fmt.Errorf("phuc hoi backup: %w", err)
	}
	return nil
}

// Relaunch khởi động lại app (không chờ tiến trình con).
func Relaunch(exePath string, args ...string) error {
	cmd := exec.Command(exePath, args...)
	cmd.Dir = filepath.Dir(exePath)
	return cmd.Start()
}

// WriteHealthMarker được bản MỚI gọi khi khởi động thành công.
func WriteHealthMarker(dir, version string) error {
	if err := os.MkdirAll(dir, 0o700); err != nil {
		return err
	}
	stamp := time.Now().UTC().Format(time.RFC3339)
	return os.WriteFile(filepath.Join(dir, version+".ok"), []byte(stamp), 0o600)
}

// HealthMarkerExists kiểm tra bản version đã báo healthy chưa.
func HealthMarkerExists(dir, version string) bool {
	_, err := os.Stat(filepath.Join(dir, version+".ok"))
	return err == nil
}

// WaitForHealthMarker chờ bản mới báo healthy, trả về false nếu quá timeout.
func WaitForHealthMarker(dir, version string, timeout, interval time.Duration) bool {
	deadline := time.Now().Add(timeout)
	for {
		if HealthMarkerExists(dir, version) {
			return true
		}
		if time.Now().After(deadline) {
			return false
		}
		time.Sleep(interval)
	}
}

// ClearHealthMarker xoá marker cũ trước khi khởi động bản mới.
func ClearHealthMarker(dir, version string) error {
	err := os.Remove(filepath.Join(dir, version+".ok"))
	if os.IsNotExist(err) {
		return nil
	}
	return err
}
