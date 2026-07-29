package updater

import (
	"encoding/base64"
	"fmt"
	"strings"

	"aead.dev/minisign"
)

// PublicKey is the Tauri-compatible Minisign public key pinned in both clients.
const PublicKey = "dW50cnVzdGVkIGNvbW1lbnQ6IG1pbmlzaWduIHB1YmxpYyBrZXk6IDQwOUZEREFCMkJBMDMwNzgKUldSNE1LQXJxOTJmUVBTSElkOHRwY2gxOXp1dmgrRVpGdXNpUGdCMkhRR25vSmY1dGo5US91K3YK"

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
