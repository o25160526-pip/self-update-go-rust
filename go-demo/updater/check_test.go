package updater

import "testing"

func TestPlatformString(t *testing.T) {
	p := PlatformString()
	// Should contain os and arch
	if p == "" {
		t.Error("expected non-empty platform string")
	}
}

func TestIsValidState(t *testing.T) {
	validStates := []string{
		"checking", "up-to-date", "update-available",
		"downloading", "verifying", "installing",
		"restarting", "failed", "rolled-back",
	}
	for _, s := range validStates {
		if !IsValidState(s) {
			t.Errorf("expected %s to be valid", s)
		}
	}
	invalidStates := []string{"", "unknown", "idle", "pending"}
	for _, s := range invalidStates {
		if IsValidState(s) {
			t.Errorf("expected %s to be invalid", s)
		}
	}
}

func TestMockCheckForUpdate_HasUpdate(t *testing.T) {
	result := MockCheckForUpdate("1.0.0", true)
	if !result.HasUpdate {
		t.Error("expected hasUpdate=true")
	}
	if result.LatestVersion != "1.0.1" {
		t.Errorf("expected 1.0.1, got %s", result.LatestVersion)
	}
}

func TestMockCheckForUpdate_NoUpdate(t *testing.T) {
	result := MockCheckForUpdate("1.0.1", false)
	if result.HasUpdate {
		t.Error("expected hasUpdate=false")
	}
}

func TestParseReleaseNotes(t *testing.T) {
	short := "Fix bug"
	if ParseReleaseNotes(short) != short {
		t.Error("short notes should be unchanged")
	}

	long := strings.Repeat("a", 250)
	result := ParseReleaseNotes(long)
	if len(result) != 200 {
		t.Errorf("expected 200 chars, got %d", len(result))
	}
	if result[197:] != "..." {
		t.Error("expected ... at end")
	}
}
