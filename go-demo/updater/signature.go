package updater

import (
	"encoding/base64"
	"fmt"
	"strings"

	"aead.dev/minisign"
)

// VerifyPinnedSignature xác minh signature (base64 của file .sig minisign) bằng
// public key pin trong keys.go. Không nhận key từ bên ngoài để tránh bị thay key.
func VerifyPinnedSignature(data []byte, encodedSignature string) error {
	publicText, err := base64.StdEncoding.DecodeString(PublicKey)
	if err != nil {
		return fmt.Errorf("decode pinned public key: %w", err)
	}
	var key minisign.PublicKey
	if err := key.UnmarshalText(publicText); err != nil {
		return fmt.Errorf("parse pinned public key: %w", err)
	}
	signatureText, err := base64.StdEncoding.DecodeString(strings.TrimSpace(encodedSignature))
	if err != nil {
		return fmt.Errorf("decode signature: %w", err)
	}
	if !minisign.Verify(key, data, signatureText) {
		return fmt.Errorf("invalid Minisign signature")
	}
	return nil
}
