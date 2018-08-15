all: build

build:
	@go build ./cmd/example-api-authz-go/...

update-opa:
	@./revendor_opa.sh $(TAG)