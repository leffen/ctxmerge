# Variables
APP_NAME := ctxmerge
VERSION := $(shell git describe --tags --always --dirty)
BUILD_DIR := build
GITHUB_REPO := your-github-username/your-repo-name
DOCKER_IMAGE := $(APP_NAME):$(VERSION)

# Go build parameters
GO := go
GOOS ?= $(shell go env GOOS)
GOARCH ?= $(shell go env GOARCH)
LD_FLAGS := -s -w -X 'main.version=$(VERSION)'

# Default target
.PHONY: all
all: test build docker-build

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

# Docker build
.PHONY: docker-build
buildd:
	@echo "Building Docker image $(DOCKER_IMAGE)..."
	docker build -t $(DOCKER_IMAGE) .

# Create a release on GitHub
.PHONY: release
release: build docker-build
	@echo "Creating a release for version $(VERSION)..."
	gh release create $(VERSION) \
		--title "$(APP_NAME) $(VERSION)" \
		--notes "Release of version $(VERSION)" \
		$(BUILD_DIR)/* LICENSE
