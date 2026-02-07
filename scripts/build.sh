rm -rf dist

echo "Building Linux amd64"
GOOS=linux GOARCH=amd64 go build -o dist/regotest_linux_amd64 ./cmd/regotest

echo "Building Linux arm64"
GOOS=linux GOARCH=arm64 go build -o dist/regotest_linux-arm64 ./cmd/regotest

echo "Building macOS amd64"
GOOS=darwin GOARCH=amd64 go build -o dist/regotest_macos_amd64 ./cmd/regotest

echo "Building macOS arm64"
GOOS=darwin GOARCH=arm64 go build -o dist/regotest_macos_arm64 ./cmd/regotest
