# AGENTS.md

## Cursor Cloud specific instructions

This is a Go CLI tool (`ai-agents-cli`) for managing AI agents, MCP servers, and agent systems on the Cloud.ru platform. It is a single standalone binary with no local service dependencies.

### Key commands

| Task | Command |
|------|---------|
| Install deps | `go mod download` |
| Build | `make build` (outputs to `bin/ai-agents-cli`) |
| Run (dev) | `go run . --help` |
| Test | `go test ./...` |
| Lint | `golangci-lint run` (requires v2; config in `.golangci.yml`) |
| Format | `go fmt ./...` |

### Caveats

- **Go version**: The project requires Go 1.25.2+ (`go.mod`). The VM has this pre-installed.
- **golangci-lint v2**: The `.golangci.yml` uses v2 config format. Install with `go install github.com/golangci/golangci-lint/v2/cmd/golangci-lint@latest`. The binary lands in `~/go/bin/`; ensure `$HOME/go/bin` is on `PATH`.
- **Interactive commands**: The `create` subcommand uses `charmbracelet/huh` TUI forms that require a TTY. These cannot be tested from a non-interactive shell; use the Desktop terminal or the `computerUse` subagent.
- **Pre-existing test failures**: Some tests in `cmd/`, `internal/scaffolder/`, and `internal/validator/` fail on `main` due to YAML parsing and template path issues. Packages `internal/api`, `internal/auth`, and `internal/config` pass cleanly.
- **Pre-existing lint issues**: The linter reports ~950 issues on `main`. These are not caused by your changes.
- **No external services needed**: All API calls are mocked in tests. Docker is only needed for `--build-image` deploy flows and is not required for build/test.
