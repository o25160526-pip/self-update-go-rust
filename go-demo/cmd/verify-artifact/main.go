// Command verify-artifact kiểm tra artifact đúng size + SHA-256 theo manifest và
// signature khớp public key đang pin. Dùng trong CI và để user tự kiểm tra file
// tải từ GitHub Releases.
//
//	go run ./cmd/verify-artifact -manifest manifest-go.json -file go-demo-windows-x64.exe
package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"go-demo/updater"
)

func main() {
	manifestPath := flag.String("manifest", "", "duong dan manifest-go.json")
	filePath := flag.String("file", "", "artifact can kiem tra")
	sigPath := flag.String("signature", "", "file .sig (mac dinh lay signature tu manifest)")
	flag.Parse()

	if *filePath == "" {
		fail("can -file")
	}
	data, err := os.ReadFile(*filePath)
	if err != nil {
		fail("doc artifact: %v", err)
	}

	signature := ""
	if *manifestPath != "" {
		raw, err := os.ReadFile(*manifestPath)
		if err != nil {
			fail("doc manifest: %v", err)
		}
		m, err := updater.ParseManifest(raw)
		if err != nil {
			fail("parse manifest: %v", err)
		}
		if err := m.Validate(); err != nil {
			fail("manifest khong hop le: %v", err)
		}
		if m.Size != int64(len(data)) {
			fail("size khong khop manifest: file %d byte, manifest %d byte", len(data), m.Size)
		}
		if !updater.VerifySHA256(data, m.SHA256) {
			fail("SHA-256 khong khop manifest")
		}
		signature = m.Signature
		fmt.Printf("manifest OK: version=%s platform=%s url=%s\n", m.Version, m.Platform, m.URL)
	}
	if *sigPath != "" {
		raw, err := os.ReadFile(*sigPath)
		if err != nil {
			fail("doc signature: %v", err)
		}
		signature = strings.TrimSpace(string(raw))
	}
	if signature == "" {
		fail("khong co signature de kiem tra (can -manifest hoac -signature)")
	}
	if err := updater.VerifyPinnedSignature(data, signature); err != nil {
		fail("verify signature: %v", err)
	}
	fmt.Printf("signature OK (pinned key %s)\n", updater.PinnedKeyID)
}

func fail(format string, args ...any) {
	fmt.Fprintf(os.Stderr, "::error::"+format+"\n", args...)
	os.Exit(1)
}
