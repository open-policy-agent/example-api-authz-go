# OPA-Go API Authorization Example

This repository shows how to integrate a service written in Go with the OPA SDK to perform API authorization.

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

The example implementation is hardcoded to assume a policy decision will be generated at path
`system.main`. You **must** define a policy decision at that
path. If your policies use another package, you can include an
entrypoint policy.

**Entrypoint**:

```rego
package system

main = data.example # api queries data.system.main.allow
```

**Your policy**:

```rego
package example

import future.keywords.if

default allow := false

allow if {
    input.method == "GET"
    input.user == "bob"
}
```
