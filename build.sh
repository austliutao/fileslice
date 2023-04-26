#!/bin/bash

# Build Windows x86 binary
GOOS=windows GOARCH=386 go build -o FileSlice_windows_386.exe main.go

# Build Windows amd64 binary
GOOS=windows GOARCH=amd64 go build -o FileSlice_windows_amd64.exe main.go

# Build Mac OS X x86 binary
GOOS=darwin GOARCH=386 go build -o FileSlice_darwin_386 main.go

# Build Mac OS X amd64 binary
GOOS=darwin GOARCH=amd64 go build -o FileSlice_darwin_amd64 main.go

# Build Linux x86 binary
GOOS=linux GOARCH=386 go build -o FileSlice_linux_386 main.go

# Build Linux amd64 binary
GOOS=linux GOARCH=amd64 go build -o FileSlice_linux_amd64 main.go

# Build Linux ARM64 binary
GOOS=linux GOARCH=arm64 go build -o FileSlice_linux_arm64 main.go

# Build Mac OS X ARM64 binary
GOOS=darwin GOARCH=arm64 go build -o FileSlice_darwin_arm64 main.go