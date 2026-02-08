rm -rf dist

build() {
  argOS=$1
  argARCH=$2
  fileName="regotest_${argOS}_${argARCH}"
  destFilePath="dist/$fileName"
  echo "Building $argOS $argARCH -> $destFilePath"
  GOOS=$argOS GOARCH=$argARCH go build -o $destFilePath ./cmd/regotest
  (cd dist && sha256sum $fileName > $fileName.sha256)
}

build linux amd64
build linux arm64
build darwin amd64
build darwin arm64
