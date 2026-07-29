package updater

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
)

func DownloadAndVerify(ctx context.Context, client *http.Client, m Manifest, destination string, verifySignature func([]byte, string) error) error {
	if err := m.Validate(); err != nil {
		return err
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, m.URL, nil)
	if err != nil {
		return err
	}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("download artifact: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("download artifact: HTTP %d", resp.StatusCode)
	}
	if err := os.MkdirAll(filepath.Dir(destination), 0700); err != nil {
		return err
	}
	tmp := destination + ".part"
	f, err := os.OpenFile(tmp, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0600)
	if err != nil {
		return err
	}
	n, copyErr := io.Copy(f, io.LimitReader(resp.Body, m.Size+1))
	closeErr := f.Close()
	if copyErr != nil {
		return copyErr
	}
	if closeErr != nil {
		return closeErr
	}
	defer os.Remove(tmp)
	if n != m.Size {
		return fmt.Errorf("artifact size mismatch: got %d want %d", n, m.Size)
	}
	data, err := os.ReadFile(tmp)
	if err != nil {
		return err
	}
	if !VerifySHA256(data, m.SHA256) {
		return fmt.Errorf("artifact SHA-256 mismatch")
	}
	if verifySignature == nil {
		return fmt.Errorf("signature verifier is required")
	}
	if err := verifySignature(data, m.Signature); err != nil {
		return fmt.Errorf("signature verification: %w", err)
	}
	return os.Rename(tmp, destination)
}
