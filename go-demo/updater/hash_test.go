package updater

import "testing"

func TestSHA256Empty(t *testing.T) {
	expected := "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855"
	if SHA256(nil) != expected {
		t.Errorf("expected %s, got %s", expected, SHA256(nil))
	}
}

func TestVerifySHA256Match(t *testing.T) {
	data := []byte("hello world")
	hash := SHA256(data)
	if !VerifySHA256(data, hash) {
		t.Error("expected match")
	}
}

func TestVerifySHA256Mismatch(t *testing.T) {
	data := []byte("hello world")
	if VerifySHA256(data, "abc") {
		t.Error("expected mismatch")
	}
}

func TestVerifySHA256Tampered(t *testing.T) {
	original := []byte("hello world")
	tampered := []byte("hello World")
	hash := SHA256(original)
	if VerifySHA256(tampered, hash) {
		t.Error("expected mismatch for tampered data")
	}
}
