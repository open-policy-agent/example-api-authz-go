# OPA-Go API Authorization Example

This repository shows how to integrate a service written in Go with OPA to perform API authorization.

## Building

Build the example by running `go build ./cmd/exmaple-api-authz-go/...`

## Running the example

Run the example with an [OPA Configuration File](https://www.openpolicyagent.org/docs/configuration.html):

```bash
./example-api-authz-go -config config.yaml
```
