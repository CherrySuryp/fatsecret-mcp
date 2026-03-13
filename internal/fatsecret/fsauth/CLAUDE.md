# internal/fatsecret/auth

OAuth 1.0a credentials and request-signing for the FatSecret API.

## Files

- **config.go** — `Config` struct, `LoadConfig`, `SaveConfig`, `ConfigPath`
- **oauth1.go** — `OAuth1Client` and its signing/request helpers

## Config

`Config` holds all credential fields. Loading merges two sources with env vars taking precedence over the file:

| Field | Source |
|---|---|
| `ClientID` / `ClientSecret` | `CLIENT_ID` / `CLIENT_SECRET` env vars, fallback to file |
| `AccessToken` / `AccessTokenSecret` | file only (written by the OAuth flow) |
| `UserID` | file only |

The config file lives at `~/.fatsecret-mcp-config.json` (mode 0600). `ConfigPath()` returns the full path.

## OAuth1Client

`NewOAuth1Client(cfg)` wraps a `*Config` and exposes one method:

```go
MakeRequest(method, rawURL string, extraParams map[string]string, token, tokenSecret string) (map[string]string, error)
```

- `extraParams` are business-level params for the endpoint (e.g. `{"format": "json"}`). OAuth protocol params are added automatically.
- Pass empty `token` / `tokenSecret` for app-only calls (e.g. request-token step).
- **GET**: all params appended as query string — used for FatSecret REST API calls.
- **POST**: all params sent as `application/x-www-form-urlencoded` body — used for OAuth token-exchange endpoints only.
- Response is parsed as JSON, falling back to query-string form, and returned as `map[string]string`.

## Key invariants

- Signing follows RFC 5849 §3.4.1: params sorted lexicographically, percent-encoded per RFC 3986.
- Only unreserved characters (`A-Z a-z 0-9 - _ . ~`) are left unencoded — **do not use** `url.QueryEscape` here, it encodes differently.
- Signing key format: `percentEncode(clientSecret) + "&" + percentEncode(tokenSecret)`.

## What this package does NOT do

- Interactive OAuth flow (user prompts, browser opening) — that lives in `internal/fetchtoken`.
- FatSecret REST API calls — those will live in `internal/fatsecret/client.go`.