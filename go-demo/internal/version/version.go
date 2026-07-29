// Package version chứa version duy nhất của app.
// Được set bằng -ldflags "-X go-demo/internal/version.Version=X.Y.Z" lúc build.
package version

// Version là version hiện tại của app. Mặc định "1.0.0" cho dev build.
// Override lúc build: go build -ldflags "-X go-demo/internal/version.Version=1.0.1"
var Version = "1.0.0"
