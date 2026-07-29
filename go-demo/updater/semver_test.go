package updater

import "testing"

func TestParseVersion(t *testing.T) {
	v, err := ParseVersion("v1.0.0")
	if err != nil {
		t.Fatalf("error: %v", err)
	}
	if v.Original() != "1.0.0" {
		t.Errorf("expected 1.0.0, got %s", v.Original())
	}
}

func TestIsNewer(t *testing.T) {
	newer, err := IsNewer("1.0.0", "1.0.1")
	if err != nil {
		t.Fatalf("error: %v", err)
	}
	if !newer {
		t.Error("expected 1.0.1 > 1.0.0")
	}

	newer, _ = IsNewer("1.0.1", "1.0.0")
	if newer {
		t.Error("expected 1.0.0 not > 1.0.1")
	}
}

func TestMeetsMinimum(t *testing.T) {
	ok, _ := MeetsMinimum("1.0.1", "1.0.0")
	if !ok {
		t.Error("expected 1.0.1 >= 1.0.0")
	}
	ok, _ = MeetsMinimum("0.9.0", "1.0.0")
	if ok {
		t.Error("expected 0.9.0 < 1.0.0")
	}
}
