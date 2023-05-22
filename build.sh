# Set environment variables
version="2.0.0"
buildTime=$(date '+%Y-%m-%d_%H:%M:%S')

# Build Windows x86 binary
echo 'Building Windows x86 binary...'
GOOS=windows GOARCH=386 go build -o "FileSlice_windows_386.exe" -ldflags "-X main.version=$version -X 'main.buildTime=$buildTime' -X 'main.arch=windows_386'" main.go

# Build Windows amd64 binary
echo 'Building Windows amd64 binary...'
GOOS=windows GOARCH=amd64 go build -o "FileSlice_windows_amd64.exe" -ldflags "-X main.version=$version -X 'main.buildTime=$buildTime' -X 'main.arch=windows_amd64'" main.go

# Build Mac OS X x86 binary
echo 'Building Mac OS X x86 binary...'
GOOS=darwin GOARCH=386 go build -o "FileSlice_darwin_386" -ldflags "-X main.version=$version -X 'main.buildTime=$buildTime' -X 'main.arch=darwin_386'" main.go

# Build Mac OS X amd64 binary
echo 'Building Mac OS X amd64 binary...'
GOOS=darwin GOARCH=amd64 go build -o "FileSlice_darwin_amd64" -ldflags "-X main.version=$version -X 'main.buildTime=$buildTime' -X 'main.arch=darwin_amd64'" main.go

# Build Linux x86 binary
echo 'Building Linux x86 binary...'
GOOS=linux GOARCH=386 go build -o "FileSlice_linux_386" -ldflags "-X main.version=$version -X 'main.buildTime=$buildTime' -X 'main.arch=linux_386'" main.go

# Build Linux amd64 binary
echo 'Building Linux amd64 binary...'
GOOS=linux GOARCH=amd64 go build -o "FileSlice_linux_amd64" -ldflags "-X main.version=$version -X 'main.buildTime=$buildTime' -X 'main.arch=linux_amd64'" main.go

# Build Linux ARM64 binary
echo 'Building Linux ARM64 binary...'
GOOS=linux GOARCH=arm64 go build -o "FileSlice_linux_arm64" -ldflags "-X main.version=$version -X 'main.buildTime=$buildTime' -X 'main.arch=linux_arm64'" main.go

# Build Mac OS X ARM64 binary
echo 'Building Mac OS X ARM64 binary...'
GOOS=darwin GOARCH=arm64 go build -o "FileSlice_darwin_arm64" -ldflags "-X main.version=$version -X 'main.buildTime=$buildTime' -X 'main.arch=darwin_arm64'" main.go

echo 'Build completed!'