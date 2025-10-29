#!/bin/bash

# This script builds the Go application for macOS.
# It ensures code quality, runs tests, and creates an optimized, versioned executable.

# Exit immediately if a command exits with a non-zero status.
set -e

# --- Configuration ---
BINARY_NAME="villain_couch"

# --- Build Steps ---

echo "--- Starting macOS Build for '${BINARY_NAME}' ---"

# 1. Clean up previous builds to ensure a fresh executable.
echo "=> Cleaning up previous build artifacts..."
rm -f "bin/${BINARY_NAME}"

# 2. Ensure Go modules are tidy and dependencies are correct.
echo "=> Tidying Go modules..."
go mod tidy

# 3. Format all Go files in the project.
echo "=> Formatting Go code..."
go fmt ./...

# 4. Run go vet to catch suspicious constructs and potential bugs.
echo "=> Vetting Go code..."
go vet ./...

# 5. Run all tests to ensure the application is stable.
#echo "=> Running tests..."
#go test ./...

# 6. Build the application for macOS.
# - We will embed the current Git tag/commit info into the binary.
# - We use ldflags to create a smaller, optimized binary.
echo "=> Building the application..."

# Get the version from the latest Git tag or commit hash.
#VERSION=$(git describe --tags --always --dirty)
#echo "   - Version: ${VERSION}"

# The '-X' flag injects the version string into a variable named 'version' in the 'main' package.
# The '-s -w' flags strip debug information, making the binary smaller.
#go build -o "${BINARY_NAME}" -ldflags="-s -w -X 'main.version=${VERSION}'" .

echo "--- Starting macOS Build for '${BINARY_NAME}' ---"
GOOS=darwin GOARCH=arm64 go build -o ./bin/villain_couch_dawin_arm64 ./agent/src/
GOOS=darwin GOARCH=amd64 go build -o ./bin/villain_couch_dawin_amd64 ./agent/src/
chmod +x ./bin/villain_couch_dawin_arm64
chmod +x ./bin/villain_couch_dawin_amd64

echo "--- Starting Linux Build for '${BINARY_NAME}' ---"
GOOS=linux GOARCH=amd64 go build -o ./bin/villain_couch_linux_amd64 ./agent/src/
GOOS=linux GOARCH=arm64 go build -o ./bin/villain_couch_linux_arm64 ./agent/src/
chmod +x ./bin/villain_couch_linux_amd64
chmod +x ./bin/villain_couch_linux_arm64

# 7. Make the resulting binary executable.
#chmod +x "${BINARY_NAME}"
#
#echo ""
#echo "--- Build Complete ---"
#echo "Executable created: ./${BINARY_NAME}"
#echo "Run './${BINARY_NAME} --version' to see the embedded version info."
