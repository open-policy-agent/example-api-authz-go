VERSION := "0.2"
GIT_COMMIT := $(shell git rev-parse --short HEAD)

all: build

LDFLAGS := -X github.com/open-policy-agent/example-api-authz-go/internal/version.Vcs=$(GIT_COMMIT) \
	-X github.com/open-policy-agent/example-api-authz-go/internal/version.Version=$(VERSION)

build:
	go build -ldflags "$(LDFLAGS)" ./cmd/example-api-authz-go/...
