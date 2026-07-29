package updater

import "testing"

const validManifest = `{"version":"1.0.1","channel":"stable","publishedAt":"2026-07-29T00:00:00Z","releaseNotes":"Fix","minSupportedVersion":"1.0.0","platform":"windows-x86_64","url":"https://example.com/app.exe","sha256":"aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa","signature":"sig","size":12,"mandatory":false}`

func TestParseAndValidateManifest(t *testing.T) {
	m, err := ParseManifest([]byte(validManifest))
	if err != nil {
		t.Fatal(err)
	}
	if err := m.Validate(); err != nil {
		t.Fatal(err)
	}
}
func TestManifestRejectsHTTP(t *testing.T) {
	m, err := ParseManifest([]byte(validManifest))
	if err != nil {
		t.Fatal(err)
	}
	m.URL = "http://example.com/app.exe"
	if m.Validate() == nil {
		t.Fatal("expected HTTPS rejection")
	}
}
func TestParseManifestInvalid(t *testing.T) {
	if _, err := ParseManifest([]byte("not json")); err == nil {
		t.Fatal("expected error")
	}
}
