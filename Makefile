VERSION ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo dev)
COMMIT  ?= $(shell git rev-parse --short HEAD 2>/dev/null || echo unknown)
DATE    ?= $(shell date -u +%Y-%m-%dT%H:%M:%SZ)
LDFLAGS := -X github.com/openziti/ziti-mcp-server-go/internal/version.Version=$(VERSION) \
           -X github.com/openziti/ziti-mcp-server-go/internal/version.Commit=$(COMMIT) \
           -X github.com/openziti/ziti-mcp-server-go/internal/version.Date=$(DATE)

.PHONY: build test lint clean generate

build:
	go build -ldflags "$(LDFLAGS)" -o ziti-mcp-server ./cmd/ziti-mcp-server

test:
	go test ./internal/...

lint:
	golangci-lint run ./...

clean:
	rm -f ziti-mcp-server

generate:
	swagger generate client -f edge-management.yml -t internal/gen/edge -A ziti-edge-management --skip-validation
	swagger generate client -f fabric-management.yml -t internal/gen/fabric -A ziti-fabric-management --skip-validation
