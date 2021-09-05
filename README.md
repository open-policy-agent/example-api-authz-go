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

The example implementation is hardcoded to assume a policy decision will be generated at path
`system.main`. You **must** define a policy decision at that
path. If your policies use another package, you can include an
entrypoint policy.

**Entrypoint**:

```rego
package system

main = data.example.allow
```

**Your policy**:

```rego
package example

default allow = false

allow {
    input.method == "GET"
    input.user == "bob"
}
```

## Trying the example
- As a manager, create a car (this should be allowed):
```
curl -H 'Authorization: alice' -H 'Content-Type: application/json' \
    -X PUT localhost:8080/cars/test-car \
    -d '{"model": "Toyota", "vehicle_id": "357192", "owner_id": "4821", "id": "test-car"}'
```

- As a car admin, try to delete a car (this should be denied):
```
curl -H 'Authorization: kelly' \
    -X DELETE localhost:8080/cars/test-car
```

### Note
- Example rego in [playground](https://play.openpolicyagent.org/p/QYwdV70Mac)

- If it's not available in playground anymore, push the files in `/rego` to [playground](https://play.openpolicyagent.org/) and update the `config.yaml` to run the example