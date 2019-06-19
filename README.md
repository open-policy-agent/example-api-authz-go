# OPA-Go API Authorization Example

This repository shows how to integrate a service written in Go with OPA to perform API authorization.

## Building

Build the example by running `go build ./cmd/example-api-authz-go/...`

## Requirements

This example requires an external HTTP server that serves [OPA
Bundles](https://www.openpolicyagent.org/docs/latest/bundles/). If you
don't provide an OPA configuration that enables bundle downloading,
the server will fail-closed.

## Running the example

Run the example with an [OPA Configuration File](https://www.openpolicyagent.org/docs/configuration.html):

```bash
./example-api-authz-go -config config.yaml
```
