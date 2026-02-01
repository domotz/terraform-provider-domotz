.PHONY: build install test testacc docs fmt vet clean

# Build the provider
build:
	go build -o terraform-provider-domotz

# Install the provider locally for testing
install: build
	mkdir -p ~/.terraform.d/plugins/registry.terraform.io/domotz/domotz/0.1.0/darwin_arm64
	cp terraform-provider-domotz ~/.terraform.d/plugins/registry.terraform.io/domotz/domotz/0.1.0/darwin_arm64/

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
