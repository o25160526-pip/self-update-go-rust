package version

import "testing"

func TestVersionDefault(t *testing.T) {
	// Default version should be 1.0.0
	if Version != "1.0.0" {
		t.Errorf("expected 1.0.0, got %s", Version)
	}
}
