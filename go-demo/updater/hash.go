package updater

import (
	"crypto/sha256"
	"encoding/hex"
)

// SHA256 tính SHA-256 hash của data.
func SHA256(data []byte) string {
	h := sha256.Sum256(data)
	return hex.EncodeToString(h[:])
}

// VerifySHA256 so sánh hash tính được với expected (constant-time).
func VerifySHA256(data []byte, expected string) bool {
	actual := SHA256(data)
	if len(actual) != len(expected) {
		return false
	}
	var diff byte
	for i := 0; i < len(actual); i++ {
		diff |= actual[i] ^ expected[i]
	}
	return diff == 0
}
