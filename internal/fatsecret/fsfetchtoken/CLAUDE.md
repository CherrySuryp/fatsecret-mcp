# internal/fatsecret/fsfetchtoken

Interactive CLI utility for completing the FatSecret OAuth 1.0a three-legged flow and persisting the resulting access token to disk.

## Files

- **fetchtoken.go** — `OAuthClient` struct, OAuth flow steps, credential setup, and `FetchToken` entrypoint

## Usage

Called from `cmd/fetchtoken/main.go`:

```go
fsfetchtoken.FetchToken()
```

Runs an interactive terminal session that guides the user through:

1. **Credential setup** — reads `CLIENT_ID` / `CLIENT_SECRET` from env or prompts the user, saves to config file.
2. **Request token** — POST to `https://authentication.fatsecret.com/oauth/request_token`.
3. **User authorization** — opens the FatSecret authorization URL in the browser and prompts for the verifier code.
4. **Access token** — GET to `https://authentication.fatsecret.com/oauth/access_token` using the verifier.
5. **Verification** — GET to `https://platform.fatsecret.com/rest/profile/v1` to confirm the token works.

On success, `AccessToken`, `AccessTokenSecret`, and `UserID` are written to the config file via `fsauth.SaveConfig`.

## Key invariants

- This package is a one-shot CLI tool — it is not imported by the MCP server or `fsclient`.
- All HTTP signing is delegated to `fsauth.FSOAuth1Client`; this package has no signing logic.
- Token exchange endpoints use plain query-string / form responses (not JSON) — `MakeRequest` handles the fallback parsing transparently.

## What this package does NOT do

- OAuth signing — that lives in `internal/fatsecret/fsauth`.
- FatSecret REST API calls — those live in `internal/fatsecret/fsclient`.