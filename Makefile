SHELL=/bin/bash -o pipefail

generate:
	go build -o .bin/gen-lookup ./cmd/gen-lookup
	go generate ./...

fmt:
	go fmt ./...

test:
	go test ./...

.PHONY: generate fmt test
