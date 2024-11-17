# Variables
APP_NAME := kubeconfig-tool
VERSION := $(shell git describe --tags --always --dirty)
BUILD_DIR := build
GITHUB_REPO := your-github-username/your-repo-name

# Go build parameters
GO := go
GOOS ?= $(shell go env GOOS)
GOARCH ?= $(shell go env GOARCH)
LD_FLAGS := -s -w -X 'main.version=$(VERSION)'

# Default target
.PHONY: all
all: test build

# Test
.PHONY: test
test:
	$(GO) test ./... -v

# Build for multiple platforms
.PHONY: build
build: clean
	@echo "Building $(APP_NAME) for multiple platforms..."
	$(MAKE) build-linux
	$(MAKE) build-macos
	$(MAKE) build-windows

.PHONY: build-linux
build-linux:
	GOOS=linux GOARCH=amd64 $(GO) build -ldflags "$(LD_FLAGS)" -o $(BUILD_DIR)/$(APP_NAME)-linux-amd64 main.go

.PHONY: build-macos
build-macos:
	GOOS=darwin GOARCH=amd64 $(GO) build -ldflags "$(LD_FLAGS)" -o $(BUILD_DIR)/$(APP_NAME)-macos-amd64 main.go

.PHONY: build-windows
build-windows:
	GOOS=windows GOARCH=amd64 $(GO) build -ldflags "$(LD_FLAGS)" -o $(BUILD_DIR)/$(APP_NAME)-windows-amd64.exe main.go

# Clean build artifacts
.PHONY: clean
clean:
	rm -rf $(BUILD_DIR)

# Create a release on GitHub
.PHONY: release
release: build
	@echo "Creating a release for version $(VERSION)..."
	gh release create $(VERSION) \
		--title "$(APP_NAME) $(VERSION)" \
		--notes "Release of version $(VERSION)" \
		$(BUILD_DIR)/* LICENSE

# Include an MIT license
.PHONY: license
license:
	@echo "Adding MIT License..."
	@cat <<EOL > LICENSE
MIT License

Copyright (c) $(shell date +%Y) Your Name

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
EOL

# Default license generation if missing
LICENSE:
	$(MAKE) license
