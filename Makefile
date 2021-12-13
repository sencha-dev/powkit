SHELL=/bin/bash -o pipefail

BUILDARGS=CGO_ENABLED=0

.PHONY: test
test:
	$(BUILD_ARGS) go test ./...
