// Command sign-artifact ký artifact bằng private key Minisign và tự verify lại
// bằng public key đang pin trong updater/keys.go.
//
// Ký và verify dùng cùng một implementation (aead.dev/minisign) nên không có
// rủi ro lệch format giữa tool ký và client verify.
//
//	go run ./cmd/sign-artifact -key ../keys/demo-signing.key -password demo-password \
//	  -in build/bin/go-demo-windows-x64.exe
package main

import (
	"encoding/base64"
	"flag"
	"fmt"
	"os"
	"strings"

	"aead.dev/minisign"

	"go-demo/updater"
)

func main() {
	keyPath := flag.String("key", "", "file private key Minisign (da ma hoa)")
	password := flag.String("password", "", "password cua private key")
	in := flag.String("in", "", "file can ky")
	out := flag.String("out", "", "file signature xuat ra (default: <in>.sig)")
	flag.Parse()

	if *keyPath == "" || *in == "" {
		fail("can -key va -in")
	}
	if *out == "" {
		*out = *in + ".sig"
	}

	keyBytes, err := os.ReadFile(*keyPath)
	if err != nil {
		fail("doc private key: %v", err)
	}
	priv, err := minisign.DecryptKey(*password, []byte(strings.TrimSpace(string(keyBytes))))
	if err != nil {
		fail("giai ma private key (sai password?): %v", err)
	}
	data, err := os.ReadFile(*in)
	if err != nil {
		fail("doc file can ky: %v", err)
	}

	encoded := base64.StdEncoding.EncodeToString(minisign.Sign(priv, data))
	if err := updater.VerifyPinnedSignature(data, encoded); err != nil {
		fail("signature khong verify duoc bang public key pin trong updater/keys.go "+
			"-> private key khong khop public key da pin: %v", err)
	}
	if err := os.WriteFile(*out, []byte(encoded), 0o644); err != nil {
		fail("ghi signature: %v", err)
	}
	fmt.Printf("da ky %s -> %s (key ID %X)\n", *in, *out, priv.ID())
}

func fail(format string, args ...any) {
	fmt.Fprintf(os.Stderr, "::error::"+format+"\n", args...)
	os.Exit(1)
}
