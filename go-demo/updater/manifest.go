package updater

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"
)

type Manifest struct {
	Version             string `json:"version"`
	Channel             string `json:"channel"`
	PublishedAt         string `json:"publishedAt"`
	ReleaseNotes        string `json:"releaseNotes"`
	MinSupportedVersion string `json:"minSupportedVersion"`
	Platform            string `json:"platform"`
	URL                 string `json:"url"`
	SHA256              string `json:"sha256"`
	Signature           string `json:"signature"`
	Size                int64  `json:"size"`
	Mandatory           bool   `json:"mandatory"`
}

func ParseManifest(data []byte) (*Manifest, error) {
	var m Manifest
	d := json.NewDecoder(bytesReader(data))
	d.DisallowUnknownFields()
	if err := d.Decode(&m); err != nil {
		return nil, err
	}
	return &m, nil
}
func bytesReader(b []byte) io.Reader { return &byteReader{b: b} }

type byteReader struct {
	b []byte
	i int
}

func (r *byteReader) Read(p []byte) (int, error) {
	if r.i >= len(r.b) {
		return 0, io.EOF
	}
	n := copy(p, r.b[r.i:])
	r.i += n
	return n, nil
}
func (m Manifest) Validate() error {
	if _, err := time.Parse(time.RFC3339, m.PublishedAt); err != nil {
		return fmt.Errorf("publishedAt: %w", err)
	}
	if m.Version == "" || m.Channel == "" || m.Platform == "" || m.URL == "" || m.Size <= 0 {
		return fmt.Errorf("missing required manifest field")
	}
	u, err := url.Parse(m.URL)
	if err != nil || u.Scheme != "https" {
		return fmt.Errorf("artifact URL must use HTTPS")
	}
	h, err := hex.DecodeString(m.SHA256)
	if err != nil || len(h) != 32 {
		return fmt.Errorf("sha256 must be 64 hex characters")
	}
	if m.Signature == "" {
		return fmt.Errorf("signature is required")
	}
	return nil
}
func FetchManifest(ctx context.Context, client *http.Client, endpoint string) (*Manifest, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return nil, err
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("fetch manifest: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("fetch manifest: HTTP %d", resp.StatusCode)
	}
	body, err := io.ReadAll(io.LimitReader(resp.Body, 1<<20))
	if err != nil {
		return nil, err
	}
	m, err := ParseManifest(body)
	if err != nil {
		return nil, fmt.Errorf("parse manifest: %w", err)
	}
	return m, nil
}
