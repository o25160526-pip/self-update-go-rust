package updater

import "testing"

func TestParseManifest(t *testing.T) {
	json := `{
		"version": "1.0.1",
		"channel": "stable",
		"publishedAt": "2026-07-29T00:00:00Z",
		"releaseNotes": "Fix update flow",
		"minSupportedVersion": "1.0.0",
		"platform": "windows-x86_64",
		"url": "https://example.com/app.exe",
		"sha256": "abc123",
		"signature": "sig",
		"size": 12345678,
		"mandatory": false
	}`
	m, err := ParseManifest([]byte(json))
	if err != nil {
		t.Fatalf("parse error: %v", err)
	}
	if m.Version != "1.0.1" {
		t.Errorf("expected 1.0.1, got %s", m.Version)
	}
	if m.Channel != "stable" {
		t.Errorf("expected stable, got %s", m.Channel)
	}
	if m.Size != 12345678 {
		t.Errorf("expected 12345678, got %d", m.Size)
	}
}

func TestParseManifestInvalid(t *testing.T) {
	_, err := ParseManifest([]byte("not json"))
	if err == nil {
		t.Error("expected error for invalid JSON")
	}
}
