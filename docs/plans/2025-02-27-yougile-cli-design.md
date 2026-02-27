# YouGile CLI — Design Document

**Date:** 2025-02-27  
**Status:** Approved

---

## 1. Goals and scope

CLI for [YouGile](https://ru.yougile.com/) (project management and CRM) with **full API coverage**: auth, companies, users, projects, boards, columns, tasks, chats, webhooks, CRM, stickers, file upload. Users work from the terminal: list, create, update, delete entities; auth via config or `auth login`.

---

## 2. Tech stack

- **Language:** Go only.
- **API client:** Generated from OpenAPI spec `docs/api.json` (oapi-codegen or go-swagger).
- **CLI framework:** Cobra (subcommands, flags, help).
- **Config:** YAML; default path `~/.config/yougile-cli/config.yaml`; override via `-c/--config` (e.g. `-c local-config.yaml` for testing).

**Config shape (minimal):**
```yaml
base_url: "https://ru.yougile.com"   # optional, default production
api_key: "your-jwt-or-api-key"
# company_id if required by API after auth login
```

**Global flags:** `-c, --config` (config path), `--json` (machine-readable output).

---

## 3. Project structure

```
yougile-cli/
├── cmd/yougile/main.go           # entrypoint, Cobra root
├── internal/
│   ├── config/                   # load and parse config.yaml
│   ├── output/                   # format: table vs JSON (--json)
│   ├── auth/                     # auth login logic (obtain/save key)
│   └── cmd/                      # Cobra subcommands
├── pkg/client/                   # generated API client + types
│   ├── client.gen.go
│   ├── types.gen.go
├── docs/
│   ├── api.json                  # OpenAPI spec (source for generation)
│   └── plans/                    # this design and implementation plan
├── go.mod, go.sum
├── Makefile                      # generate, build, lint, test
├── .golangci.yml                 # golangci-lint config
└── README.md
```

Generated code lives in `pkg/client/`; regeneration via `make generate` (or `go generate`).

---

## 4. Command tree (full coverage)

| Group | Commands |
|-------|----------|
| **auth** | `login`, `companies`, `keys list \| create \| delete` |
| **config** | `show`, `path` |
| **users** | `list`, `get`, `create`, `update`, `delete` |
| **company** | `get`, `update` |
| **projects** | `list`, `get`, `create`, `update` |
| **projects roles** | `list`, `get`, `create`, `update`, `delete` (under `projects/{id}/roles`) |
| **departments** | `list`, `get`, `create`, `update` |
| **boards** | `list`, `get`, `create`, `update` |
| **columns** | `list`, `get`, `create`, `update` |
| **tasks** | `list`, `get`, `create`, `update`, `chat-subscribers get \| update` |
| **task-list** | dedicated list endpoint per API spec |
| **string-stickers** | `list`, `get`, `create`, `update`; states: `list`, `get`, `create`, `update` |
| **sprint-stickers** | same pattern as string-stickers |
| **chats** | group-chats: `list`, `get`, `create`, `update`; messages: `list`, `send`, `update` |
| **webhooks** | `list`, `get` (and create/delete if in API) |
| **crm** | `contact-persons`, `contacts by-external-id`, etc. per API |
| **files** | `upload <path>` |

Subcommand hierarchy and flags (e.g. `--project-id`, `--board-id`) follow `docs/api.json`. Pagination: `--limit`, `--offset` where the API supports it.

---

## 5. Output

- **Default (no `--json`):** Human-readable tables for lists (id, title, key fields); key-value style for single resource.
- **With `--json`:** Raw API response JSON (or stable struct from generated types), no extra formatting.
- **Pagination:** Only current page shown in tables; total/count from API if available.

---

## 6. Auth flow

- **Option A:** User puts `api_key` (and optional `base_url`) in config; CLI uses it for all requests.
- **Option B:** User runs `yougile auth login` (email + password); CLI calls API to get companies / create key, writes key (and optional company context) to config. Config file and directory created if missing.

Both supported; no mandatory login if key is already in config.

---

## 7. Linting and testing

- **Linter:** golangci-lint with repo config (`.golangci.yml`). Enable: govet, errcheck, staticcheck, ineffassign, gofmt; optional revive. Run locally and in CI (`make lint`).
- **Unit tests:** Config load/parse; output formatters; command handlers via mocked API client interface — verify args and output (including `--json`). Do not unit-test generated client; test only our wrappers and commands.
- **Integration tests (optional):** Tagged e.g. `integration`, run with `go test -tags=integration ./...`; require real config or test instance; do not block default `go test ./...`.
- **CI:** `go test ./...` and `make lint` (or equivalent) on every commit/PR. After `make generate`, ensure generated files are either committed or CI checks they are up to date.

---

## 8. Error handling

- Missing or invalid config: clear message and non-zero exit.
- API errors (4xx/5xx): print status and body (or summary); non-zero exit.
- Network errors: clear message; retries optional (later).
- All user-facing messages in English.

---

## 9. Implementation approach

- **Approach chosen:** Full code generation from OpenAPI (approach 2). Generate client and DTOs from `docs/api.json`; CLI calls generated methods and formats output (table or `--json`).
- Generator: oapi-codegen or go-swagger; base URL and Bearer token injected from config.

---

## Next step

Implementation plan: see `docs/plans/2025-02-27-yougile-cli.md` (bite-sized tasks, TDD, lint, tests, frequent commits).
