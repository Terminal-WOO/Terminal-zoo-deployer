# Makefile

# Binary name
BINARY := ipo-deployer

# Version info (defaults, can be overridden)
VERSION   := $(shell git describe --tags --always --dirty)
COMMIT    := $(shell git rev-parse --short HEAD)
BUILDDATE := $(shell date -u +%Y-%m-%dT%H:%M:%SZ)

IMAGE := clappform/$(BINARY):$(VERSION)

# LDFLAGS to inject into Go binary
LDFLAGS := -X 'main.Version=$(VERSION)' \
           -X 'main.Commit=$(COMMIT)' \
           -X 'main.BuildDate=$(BUILDDATE)' \
           -s -w

# Build target
.PHONY: build
build:
	@echo "Building $(BINARY)..."
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
	go build -v -ldflags="$(LDFLAGS)" -o $(BINARY) ./cmd/server

# Run target
.PHONY: run
run: build
	./$(BINARY)

# Clean target
.PHONY: clean
clean:
	rm -f $(BINARY)

# Docker build target
.PHONY: docker-build
docker-build: build
	@echo "Building Docker image $(IMAGE)..."
	docker build --platform linux/amd64 -t $(IMAGE) .
