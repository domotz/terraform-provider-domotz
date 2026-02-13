.PHONY: build install test testacc docs fmt vet clean

VERSION ?= 1.1.0
PLATFORM := $(shell go env GOOS)_$(shell go env GOARCH)

# Build the provider
build:
	go build -o terraform-provider-domotz

# Install the provider locally for testing
install: build
	mkdir -p ~/.terraform.d/plugins/registry.terraform.io/domotz/domotz/$(VERSION)/$(PLATFORM)
	cp terraform-provider-domotz ~/.terraform.d/plugins/registry.terraform.io/domotz/domotz/$(VERSION)/$(PLATFORM)/

# Run unit tests
test:
	go test -v ./...

# Run acceptance tests (requires DOMOTZ_API_KEY environment variable)
testacc:
	TF_ACC=1 go test -v ./internal/provider -timeout 10m

# Generate documentation
docs:
	go run github.com/hashicorp/terraform-plugin-docs/cmd/tfplugindocs

# Format code
fmt:
	go fmt ./...

# Run go vet
vet:
	go vet ./...

# Clean build artifacts
clean:
	rm -f terraform-provider-domotz
	rm -rf dist/

# Lint code
lint:
	golangci-lint run

# Download dependencies
deps:
	go mod download
	go mod tidy

# Run all checks before commit
check: fmt vet test

# Build binaries for all platforms
build-all:
	./build-all-platforms.sh

# Create release (requires gh CLI)
release: build-all
	@echo "Creating GitHub release v$(VERSION)..."
	gh release create v$(VERSION) \
		--title "Terraform Provider for Domotz v$(VERSION)" \
		--notes-file CHANGELOG.md \
		dist/*.tar.gz dist/*.zip dist/*SHA256SUMS
