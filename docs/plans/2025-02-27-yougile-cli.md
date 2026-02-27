# YouGile CLI Implementation Plan

> **For Claude:** REQUIRED SUB-SKILL: Use superpowers:executing-plans to implement this plan task-by-task.

**Goal:** Deliver a Go CLI for YouGile with full API coverage, generated client from OpenAPI, Cobra commands, YAML config, table/JSON output, auth login, and CI with lint and tests.

**Architecture:** Client and types generated from `docs/api.json` into `pkg/client/`. Cobra root in `cmd/yougile`; subcommands in `internal/cmd`. Config and output in `internal/config`, `internal/output`. Auth login in `internal/auth`. All commands use generated client; output layer formats for table or `--json`.

**Tech Stack:** Go, Cobra, oapi-codegen (or go-swagger), YAML config, golangci-lint.

---

## Phase 0: Bootstrap and tooling

### Task 0.1: Go module and directory layout

**Files:**
- Create: `yougile-cli/go.mod`
- Create: `yougile-cli/cmd/yougile/main.go` (empty or minimal)
- Create: `yougile-cli/Makefile` (targets: build, test, lint, generate)

**Steps:**
1. Run `go mod init github.com/<your-org>/yougile-cli` (or similar) in project root.
2. Create directories: `cmd/yougile`, `internal/config`, `internal/output`, `internal/auth`, `internal/cmd`, `pkg/client`.
3. Add minimal `main.go` that prints "yougile" and exits.
4. Add Makefile with `build` (go build -o bin/yougile ./cmd/yougile), `test` (go test ./...), `lint` (golangci-lint run), `generate` (placeholder: echo "run codegen").
5. Run `go build ./cmd/yougile` and `go test ./...` (no tests yet). Commit: "chore: bootstrap Go module and layout".

---

### Task 0.2: golangci-lint config and CI placeholder

**Files:**
- Create: `yougile-cli/.golangci.yml`

**Steps:**
1. Add .golangci.yml with linters: govet, errcheck, staticcheck, ineffassign, gofmt.
2. Run `golangci-lint run ./...` (install if needed); fix any issues in existing code.
3. Ensure Makefile `lint` target runs `golangci-lint run ./...`.
4. Commit: "chore: add golangci-lint config".

---

### Task 0.3: OpenAPI code generation setup

**Files:**
- Modify: `yougile-cli/go.mod` (add oapi-codegen or go-swagger)
- Create: `yougile-cli/scripts/gen.sh` or use `//go:generate` in pkg/client
- Modify: `yougile-cli/Makefile` (generate target invokes codegen)

**Steps:**
1. Add codegen tool: `go install github.com/deepmap/oapi-codegen/cmd/oapi-codegen@latest` (or go-swagger). Add to go.mod if needed.
2. Create config for codegen (e.g. `pkg/client/gen.yaml` or inline flags) to generate from `docs/api.json`: types + client. Target package `client`, output `pkg/client/`.
3. Run codegen once; commit generated files to `pkg/client/` (e.g. `client.gen.go`, `types.gen.go`). Ensure base URL and auth are injectable (e.g. client constructor accepts base URL and http.Client with auth roundtripper).
4. Update Makefile `generate` to run the codegen command.
5. Run `make generate` and `go build ./...`. Commit: "feat: add OpenAPI codegen and generated client".

---

## Phase 1: Config and global flags

### Task 1.1: Config struct and file loading

**Files:**
- Create: `internal/config/config.go` (struct Config with BaseURL, APIKey; Load(path string) (*Config, error))
- Create: `internal/config/config_test.go`

**Steps:**
1. Write failing test: Load with missing file returns error; Load with valid YAML returns struct with BaseURL and APIKey.
2. Run test: FAIL.
3. Implement: read file, parse YAML into Config; default BaseURL to "https://ru.yougile.com" if empty; return error if file missing or api_key empty (or allow empty for auth login only).
4. Run test: PASS. Run lint. Commit: "feat: config load and parse".

---

### Task 1.2: Default config path and -c flag

**Files:**
- Create: `cmd/yougile/root.go` (Cobra root command, PersistentFlags for -c and --json)
- Modify: `cmd/yougile/main.go` (execute root)

**Steps:**
1. Root command: PersistentFlag `-c, --config` (string); PersistentFlag `--json` (bool). Resolve config path: if -c set use it, else use ~/.config/yougile-cli/config.yaml. Store --json in a root-level variable or context for output layer.
2. Wire root to main. Run `go run ./cmd/yougile --help`; ensure -c and --json appear. Commit: "feat: root command and config path resolution".

---

### Task 1.3: Config show and config path commands

**Files:**
- Create: `internal/cmd/config.go` (subcommands: show, path)

**Steps:**
1. `yougile config path`: print resolved config path (from same logic as root). No API call.
2. `yougile config show`: load config and print (mask api_key in human output; full in --json if needed). Exit with clear error if config missing or invalid.
3. Register subcommands under root. Tests: unit test with temp dir and -c pointing to temp config. Commit: "feat: config show and path commands".

---

## Phase 2: Auth and API client wiring

### Task 2.1: Auth login flow

**Files:**
- Create: `internal/auth/login.go` (Login(ctx, baseURL, email, password) (apiKey string, err error))
- Create: `internal/auth/login_test.go` (mock HTTP or use httptest)

**Steps:**
1. Implement Login: call API getCompanies (or create key) with credentials; parse response to get key (and optionally company id). Use generated client or http.Post with JSON body. Accept baseURL.
2. Write test: mock server returns 200 and key; Login returns key. Test failure on 401.
3. Run tests and lint. Commit: "feat: auth login API flow".

---

### Task 2.2: Auth login command and config write

**Files:**
- Create: `internal/cmd/auth.go` (auth login; auth companies; auth keys list/create/delete)
- Modify: `internal/config/config.go` (Save(path string, cfg *Config) error — write YAML, create dir if needed)

**Steps:**
1. `yougile auth login`: prompt or flags for email/password; call auth.Login; build Config with APIKey and BaseURL; config.Save to resolved path. Create ~/.config/yougile-cli if needed.
2. Add config Save and directory creation. Test with temp dir.
3. Register `auth login` under root. Commit: "feat: auth login command and config save".

---

### Task 2.3: Inject config into generated client and run first API call

**Files:**
- Modify: `cmd/yougile/root.go` or `internal/cmd` (create client from config after load)
- Create: `internal/cmd/company.go` (company get) as first real command

**Steps:**
1. After loading config in root (or in RunE of commands that need API), instantiate generated client with config.BaseURL and config.APIKey (e.g. custom Client with Bearer transport). Pass client into command constructors or use a factory.
2. Implement `yougile company get`: call generated company get endpoint; output table or JSON per --json.
3. Manual test with valid config and key. Unit test with mocked client. Commit: "feat: wire generated client and company get".

---

## Phase 3: Output layer

### Task 3.1: Table and JSON printer

**Files:**
- Create: `internal/output/printer.go` (PrintTable(headers []string, rows [][]string), PrintJSON(v interface{}))
- Create: `internal/output/printer_test.go`

**Steps:**
1. PrintTable: format rows as table (e.g. tabwriter or tablewriter). PrintJSON: json.Encoder with indent when not --json for single object, or no indent for list (or follow design: --json = raw).
2. Tests: capture stdout, assert table has expected columns; assert JSON is valid and contains key fields.
3. Wire global --json flag to printer (e.g. if --json then PrintJSON else PrintTable for lists). Commit: "feat: output table and JSON printer".

---

## Phase 4: Commands (by resource)

Implement each resource group with the same pattern: add subcommand file in `internal/cmd`, use generated client, call output printer. Add unit tests with mocked client. Run lint and tests after each group.

### Task 4.1: Users commands

**Files:**
- Create or extend: `internal/cmd/users.go` (list, get, create, update, delete)

**Steps:**
1. Implement users list (with --limit, --offset), get &lt;id&gt;, create (flags or stdin for body), update &lt;id&gt;, delete &lt;id&gt;. Use generated types and client.
2. Table columns: id, title/name, email as in API. JSON: raw response.
3. Tests: mock client, assert correct methods and args; assert output format. Commit: "feat: users commands".

---

### Task 4.2: Projects and project roles

**Files:**
- Create: `internal/cmd/projects.go` (projects list, get, create, update; roles list, get, create, update, delete under project)

**Steps:**
1. Projects: list, get &lt;id&gt;, create, update &lt;id&gt;. Roles: `yougile projects roles list --project-id=&lt;id&gt;` etc.
2. Tests and lint. Commit: "feat: projects and roles commands".

---

### Task 4.3: Departments, boards, columns

**Files:**
- Create: `internal/cmd/departments.go`, `internal/cmd/boards.go`, `internal/cmd/columns.go`

**Steps:**
1. list, get, create, update (and delete if in API) for each. Boards/columns: filter by project/board where applicable (--project-id, --board-id).
2. Tests and lint. Commit: "feat: departments, boards, columns commands".

---

### Task 4.4: Tasks and task-list

**Files:**
- Create: `internal/cmd/tasks.go` (list, get, create, update; task-list if separate endpoint; chat-subscribers get/update)

**Steps:**
1. Tasks list with --project-id, --board-id, --column-id, --limit, --offset. get, create, update. task-list endpoint if different. Chat-subscribers: get and update for task id.
2. Tests and lint. Commit: "feat: tasks and task-list commands".

---

### Task 4.5: String and sprint stickers (and states)

**Files:**
- Create: `internal/cmd/stickers.go` or split string_stickers.go / sprint_stickers.go

**Steps:**
1. string-stickers: list, get, create, update; states list, get, create, update. Same for sprint-stickers.
2. Tests and lint. Commit: "feat: stickers commands".

---

### Task 4.6: Chats and messages

**Files:**
- Create: `internal/cmd/chats.go` (group-chats list, get, create, update; messages list, send, update)

**Steps:**
1. group-chats: list, get &lt;id&gt;, create, update. messages: list for chat &lt;chatId&gt;, send, update &lt;id&gt;.
2. Tests and lint. Commit: "feat: chats and messages commands".

---

### Task 4.7: Webhooks and CRM

**Files:**
- Create: `internal/cmd/webhooks.go`, `internal/cmd/crm.go`

**Steps:**
1. Webhooks: list, get; create/delete if in API. CRM: contact-persons list/get/...; contacts by-external-id per API.
2. Tests and lint. Commit: "feat: webhooks and crm commands".

---

### Task 4.8: File upload

**Files:**
- Create: `internal/cmd/files.go` (upload &lt;path&gt;)

**Steps:**
1. Call upload-file endpoint with multipart form; output URL or success message. Table/JSON as per --json.
2. Tests and lint. Commit: "feat: file upload command".

---

### Task 4.9: Auth companies and keys commands

**Files:**
- Modify: `internal/cmd/auth.go` (add companies, keys list, keys create, keys delete)

**Steps:**
1. auth companies: list companies (may require credentials in request body per API). auth keys list: list keys; auth keys create: create key; auth keys delete: delete by key id. Use generated client.
2. Tests and lint. Commit: "feat: auth companies and keys commands".

---

## Phase 5: Error handling and polish

### Task 5.1: Centralized error handling and exit codes

**Files:**
- Create: `internal/errors/errors.go` or use in each RunE (e.g. print message, os.Exit(1))

**Steps:**
1. Define exit codes: 0 success, 1 usage/config/API error. In RunE return error or set exit code; Cobra can run SilenceErrors and handle in root.
2. Ensure API 4xx/5xx and network errors print clear message and exit non-zero. Commit: "fix: error handling and exit codes".

---

### Task 5.2: README and usage docs

**Files:**
- Modify: `README.md` (install, config, auth login, examples for main commands)

**Steps:**
1. Document: build/install, config file location and -c, auth login, list of command groups, --json and -c. Commit: "docs: README and usage".

---

## Phase 6: CI and final checks

### Task 6.1: Ensure generate is idempotent and CI

**Files:**
- Modify: `Makefile` or add script to verify generated files (e.g. make generate && git diff --exit-code pkg/client)

**Steps:**
1. Add target or CI step: run generate, then fail if pkg/client has uncommitted changes (so PRs keep generated code in sync with api.json).
2. Document in README: run make generate after changing docs/api.json. Commit: "ci: verify generated code".

---

## Execution summary

- **Phase 0:** Bootstrap, lint, codegen (Tasks 0.1–0.3).
- **Phase 1:** Config load, root flags, config show/path (Tasks 1.1–1.3).
- **Phase 2:** Auth login, client wiring, company get (Tasks 2.1–2.3).
- **Phase 3:** Output printer (Task 3.1).
- **Phase 4:** All resource commands (Tasks 4.1–4.9).
- **Phase 5:** Errors and README (Tasks 5.1–5.2).
- **Phase 6:** CI and generate check (Task 6.1).

Run tests and lint after each task. Commit after each task. Use TDD where applicable (config, auth, output, command mocks).
