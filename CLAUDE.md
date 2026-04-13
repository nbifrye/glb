# glb - GitLab CLI with MCP Server

## Build & Test

```bash
make build    # Build binary to bin/glb
make test     # Run all tests with race detector
make clean    # Remove build artifacts
```

## Architecture

- `cmd/glb/` — Entry point
- `internal/gitlabop/` — Shared service layer (CLI and MCP both call this)
- `internal/commands/` — CLI commands (Cobra)
- `internal/commands/mcp/serve/tools/` — MCP tool registrations
- `internal/auth/` — Token resolution (env vars > config file)
- `internal/config/` — YAML config management (~/.config/glb/config.yml)
- `internal/api/` — GitLab client factory
- `internal/glinstance/` — Hostname/URL utilities

## Adding a New Command

1. Create `internal/gitlabop/<resource>.go` with the business logic
2. Create `internal/commands/<resource>/<action>/<resource>_<action>.go` for CLI
3. Add MCP tool in `internal/commands/mcp/serve/tools/<resource>.go`
4. Register in `tools/registry.go` and parent command's `NewCmd`
5. Add tests

## Error Handling

- CLI commands: return `fmt.Errorf(...)` from `RunE`
- MCP tools: return `errorResult(...)` for API errors (sets `IsError: true`)
- Always check `json.Marshal` errors
- MCP write tools must have `Annotations` with `DestructiveHint`

## Auth Resolution Order

1. `GITLAB_TOKEN` env var
2. `GLB_TOKEN` env var
3. Config file token (per-host)

Host resolution: `GITLAB_HOST` > `GLB_HOST` > single configured host > gitlab.com
