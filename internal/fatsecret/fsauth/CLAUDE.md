# internal/fatsecret/fsauth

OAuth 1.0 signing utilities for the FatSecret API. Produces `Authorization` header values — does not make HTTP requests.

## Requirements

- **Isolated**: must not import `internal/config` or `internal/logging`.
- **No HTTP**: no `net/http` usage; signing only.
- **Stdlib only**: no external OAuth libraries.

## Files

| File | Purpose |
|---|---|
| `base_auth.go` | `BaseAuth` struct with shared consumer credentials |
| `signed_auth.go` | `SignedAuth` for non-delegated (public) requests |
| `delegated_auth.go` | `DelegatedAuth` for user-specific (profile) requests |
| `common.go` | Shared signing helpers (nonce, encoding, header building) |

## Auth types

### `SignedAuth`

For endpoints that do not require a user identity (e.g. food search, food lookup).

```go
client := fsauth.NewSignedAuth(consumerKey, consumerSecret)
header, err := client.AuthorizationHeader("GET", baseURL, requestParams)
```

Signing key: `percentEncode(consumerSecret) + "&"`

### `DelegatedAuth`

For endpoints that operate on behalf of a user profile (e.g. diary, favorites, exercise log). Requires per-user tokens from `FSAPIUserConfig`.

```go
client := fsauth.NewDelegatedAuth(consumerKey, consumerSecret, accessToken, secretToken)
header, err := client.AuthorizationHeader("POST", baseURL, requestParams)
```

Signing key: `percentEncode(consumerSecret) + "&" + percentEncode(tokenSecret)`

Mapping from `internal/config.FSAPIUserConfig`:

| `FSAPIUserConfig` field | `NewDelegatedAuth` param |
|---|---|
| `AccessToken` | `oauthToken` |
| `SecretToken` | `oauthTokenSecret` |
| `UserID` | not used for signing |

## `AuthorizationHeader` algorithm

1. Build OAuth params map (nonce, timestamp, consumer key, signature method, version; `oauth_token` added for delegated)
2. Merge OAuth params + `requestParams` into one map for signing
3. Percent-encode and sort all params → normalized parameter string
4. Build signature base string: `METHOD&percentEncode(url)&percentEncode(params)`
5. HMAC-SHA1 sign → base64 encode → `oauth_signature`
6. Build and return `OAuth key="value", ...` header string

## What this package does NOT do

- HTTP requests — that lives in `internal/fatsecret/fsclient`.
- Token fetching or profile creation — that lives in `internal/fatsecret/fsfetchtoken`.
- Config loading — callers extract raw strings from `internal/config` and pass them in.