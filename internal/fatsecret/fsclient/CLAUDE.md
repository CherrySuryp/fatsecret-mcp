# internal/fatsecret/fsclient

Typed HTTP client for the FatSecret REST API. Wraps `fsauth.FSOAuth1Client` and handles URL construction, `format=json` injection, and response unmarshalling into typed Go structs.

## Files

- **client.go** — `FSClient` struct, constructor, and `get`/`post` transport helpers
- **\*.go** — one file per FatSecret resource group (e.g. `foods.go`, `profile.go`)

## FSClient

`NewFSClient(cfg)` wraps a `*fsauth.Config` and exposes typed methods per resource:

```go
client := fsclient.NewFSClient(cfg)
```

Internally, all requests go through two private helpers:

```go
get(path string, params map[string]string) (map[string]interface{}, error)
post(path string, params map[string]string) (map[string]interface{}, error)
```

Both inject `format=json` automatically and delegate signing to `fsauth.FSOAuth1Client`. Resource methods call these helpers with the endpoint-specific path and params, then unmarshal the result into typed structs.

## Adding a new endpoint

1. Create a file named after the resource group (e.g. `foods.go`).
2. Define Go structs matching the FatSecret response schema.
3. Add a method on `*FSClient` that calls `c.get` or `c.post` with the versioned path (e.g. `"foods/search/v1"`), then unmarshals the relevant key from the response.

Example:

```go
func (c *FSClient) SearchFoods(query string) ([]Food, error) {
    result, err := c.get("foods/search/v1", map[string]string{
        "search_expression": query,
    })
    if err != nil {
        return nil, err
    }
    // unmarshal result["foods"] into []Food
}
```

## Key invariants

- Base URL is `https://platform.fatsecret.com/rest/` — paths must not include a leading `/`.
- `format=json` is always injected; do not pass it in `params`.
- OAuth signing and token injection are handled entirely by `fsauth` — this package has no signing logic.

## What this package does NOT do

- OAuth signing or token management — that lives in `internal/fatsecret/fsauth`.
- Interactive OAuth flow — that lives in `internal/fatsecret/fsfetchtoken`.