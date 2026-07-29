package updater

import (
	"strings"
	"testing"
)

func TestPlatformString(t *testing.T) {
	if PlatformString() == "" {
		t.Fatal("empty platform")
	}
}
func TestIsValidState(t *testing.T) {
	for _, s := range []string{"checking", "up-to-date", "update-available", "downloading", "verifying", "installing", "restarting", "failed", "rolled-back"} {
		if !IsValidState(s) {
			t.Fatalf("invalid %s", s)
		}
	}
	if IsValidState("pending") {
		t.Fatal("pending must not be UI state")
	}
}
func TestMockCheckForUpdate(t *testing.T) {
	if !MockCheckForUpdate("1.0.0", true).HasUpdate {
		t.Fatal("expected update")
	}
	if MockCheckForUpdate("1.0.1", false).HasUpdate {
		t.Fatal("unexpected update")
	}
}
func TestParseReleaseNotes(t *testing.T) {
	got := ParseReleaseNotes(strings.Repeat("a", 250))
	if len(got) != 200 || !strings.HasSuffix(got, "...") {
		t.Fatalf("unexpected notes %q", got)
	}
}
