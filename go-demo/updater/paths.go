package updater

import (
	"os"
	"path/filepath"
)

// ManifestURL là endpoint manifest của bản Go trên GitHub Releases.
// `releases/latest/download/...` luôn trỏ tới release mới nhất (non-prerelease).
const ManifestURL = "https://github.com/o25160526-pip/self-update-go-rust/releases/latest/download/manifest-go.json"

// DataDir trả về thư mục dữ liệu runtime của app.
// Windows: %LOCALAPPDATA%\go-demo · Linux: ~/.cache/go-demo (dùng khi test).
func DataDir() (string, error) {
	base, err := os.UserCacheDir()
	if err != nil {
		return os.MkdirTemp("", "go-demo")
	}
	return filepath.Join(base, "go-demo"), nil
}

// StatePath là đường dẫn state.json (current / pending / last-known-good).
func StatePath(dir string) string { return filepath.Join(dir, "state.json") }

// DownloadDir là thư mục tạm riêng để tải artifact về.
func DownloadDir(dir string) string { return filepath.Join(dir, "update-tmp") }

// HealthDir chứa health marker do bản mới ghi ra khi khởi động thành công.
func HealthDir(dir string) string { return filepath.Join(dir, "health") }
