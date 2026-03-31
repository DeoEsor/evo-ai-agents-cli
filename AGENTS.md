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
- **Pre-existing lint issues**: The linter reports ~950 issues on `main`. These are not caused by your changes.
- **No external services needed for tests**: All API calls are mocked in tests. Docker is only needed for `--build-image` deploy flows.

### Cloud.ru API credentials

The CLI uses two distinct credential sets:

| Credential | Purpose | Where to get |
|-----------|---------|-------------|
| `IAM_KEY_ID` + `IAM_SECRET` | AI Agents API (agents, MCP servers, systems) | Console → Service Accounts → API Keys (service: AI Agents) |
| `FM_API_KEY` | Foundation Models API (LLM inference) | Console → Service Accounts → API Keys (service: Foundation Models) |
| `PROJECT_ID` | Scopes all API calls to a project | Console → Project Settings |

The agent templates support both auth methods for Foundation Models: `FM_API_KEY` (direct) or `IAM_KEY_ID`+`IAM_SECRET` (IAM token exchange via `/api/v1/auth/token`).

### Scaffolded project structure

`create agent` generates a FastAPI app with:
- `docker-compose.yml` with `pgvector/pgvector:pg16` (always included)
- `src/agent.py` — Foundation Models client (OpenAI-compatible), endpoints: `POST /agent/process`, `GET /agent/models`, `GET /health`
- `.github/workflows/ci.yml` and `.gitlab-ci.yml` — lint/test/build/deploy to `cr.cloud.ru`
- Auth: IAM endpoint at `https://iam.api.cloud.ru/api/v1/auth/token` (not `/iam/v1/`)

### Running scaffolded projects

```bash
cd /tmp && ai-agents-cli create agent my-agent  # interactive form (needs TTY)
cd my-agent
# edit .env with your credentials
docker compose up -d --build
curl http://localhost:8000/health
curl -X POST http://localhost:8000/agent/process -H "Content-Type: application/json" -d '{"message": "Hello"}'
```
