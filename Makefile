SHELL=/bin/bash -o pipefail

BUILDARGS=CGO_ENABLED=0

.PHONY: generate
generate:
	$(BUILD_ARGS) go build -o .bin/gen-lookup ./cmd/gen-lookup
	go generate ./...

.PHONY: test
test:
	$(BUILD_ARGS) go test ./...
