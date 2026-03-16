# internal/config

MCP server configuration — loads and exposes all settings needed at startup.

## Requirements

- `FSMCPServerConfig` MUST be loaded from environment variables.
- `FSAPIClientConfig` MUST be loaded from environment variables.
- `FSAPIUserConfig` is optional; MUST be loaded from a JSON file located in the user's home directory by default.

## Files

- **config.go** — `Config` struct, all load helpers, and `MustLoadConfig`

## Config structure

`Config` composes three sub-configs:

| Struct | Env vars | Purpose |
|---|---|---|
| `FSMCPServerConfig` | `FATSECRET_MCP_SERVER_LOG_LEVEL` | MCP server settings |
| `FSAPIClientConfig` | `FATSECRET_MCP_API_CONSUMER_ID`, `FATSECRET_MCP_API_CONSUMER_SECRET`, `FATSECRET_MCP_API_USER_CONFIG_PATH` | FatSecret API client credentials |
| `FSAPIUserConfig` | file only | Per-user OAuth tokens loaded from the user config file |

## Load order (`MustLoadConfig`)

1. `godotenv.Load()` — silently ignored if no `.env` file present
2. `loadMCPConfig()` — parses `FSMCPServerConfig` from env
3. `loadAPIClientConfig()` — parses `FSAPIClientConfig` from env; `ConsumerID` and `ConsumerSecret` are required, will fatal if missing
4. `loadUserConfig()` — reads and JSON-decodes the user config file; absence of the file is not an error

`MustLoadConfig` calls `log.Fatal` on any error — it is only called at process startup.

## User config file

Default path: `~/.fatsecret-mcp-user-creds.json` (overridable via `FATSECRET_MCP_API_USER_CONFIG_PATH`).

`userConfigPath()` resolves `~/` to the real home directory and stats the file. It returns:
- `filePath, nil` — file exists
- `"", nil` — file does not exist (not an error)
- `"", err` — unexpected filesystem error

## Public API

```go
UserConfigExists() bool  // false if file absent or on stat error
UserConfigEmpty() bool   // true if any of UserID / AccessToken / SecretToken is blank
UserConfigValid() bool   // true if file exists and all credential fields are populated
```

`UserConfigValid` is the preferred gate for enabling user-specific MCP endpoints — it combines both checks. `UserConfigEmpty` is only meaningful when `UserConfigExists` is true.

## What this package does NOT do

- OAuth signing or token fetching — that lives in `internal/fatsecret/fsauth` and `internal/fatsecret/fsfetchtoken`.
- Writing the user config file — it is written externally by the token-fetch flow.