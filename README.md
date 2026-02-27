# yougile-cli

CLI for [YouGile](https://ru.yougile.com/) — project management and CRM.

## Install

```bash
go build -o yougile ./cmd/yougile
# or
make build
```

## Config

Config file path:

- Default: `~/.config/yougile-cli/config.yaml`
- Override: `yougile -c /path/to/config.yaml`

Example config:

```yaml
base_url: "https://ru.yougile.com"
api_key: "your-api-key"
```

## Auth

Get an API key and save it to config:

```bash
yougile auth login --email your@email.com --password yourpassword
```

List companies (no saved key needed):

```bash
yougile auth companies --email your@email.com --password yourpassword
```

List API keys:

```bash
yougile auth keys list --email your@email.com --password yourpassword
```

Create/delete API keys (no saved key; use email/password + company-id):

```bash
yougile auth keys create --email your@email.com --password yourpassword --company-id <id>
yougile auth keys delete <key>
```

## Commands

- `yougile config path` — print config file path
- `yougile config show` — show config (api_key masked in human output)
- `yougile company get` — current company details
- **users:** `users list` / `users get <id>` / `users create --email … [--admin]` / `users update <id> [--admin]` / `users delete <id>`
- **projects:** `projects list` / `projects get <id>` / `projects create --title "…"` / `projects update <id> [--title "…"]`
- **boards:** `boards list` / `boards get <id>` / `boards create --title "…" --project-id <id>` / `boards update <id> [--title "…"]`
- **columns:** `columns list` / `columns get <id>` / `columns create --title "…" --board-id <id>` / `columns update <id> [--title "…"]`
- **tasks:** `tasks list` / `tasks get <id>` / `tasks create --title "…" [--column-id <id>]` / `tasks update <id> [--title "…"]` / `tasks chat-subscribers get <task-id>` / `tasks chat-subscribers update <task-id> --user-ids "id1,id2"`
- **departments:** `departments list` / `departments get <id>` / `departments create --title "…" [--parent-id <id>]` / `departments update <id> [--title "…"]`
- **webhooks:** `webhooks list` / `webhooks create --event "…" --url "…"`
- `yougile files upload <path>`
- **chats:** `chats list` (--limit, --offset, --title), `chats get <id>`, `chats messages list <chat-id>`, `chats messages send <chat-id> --text "…"`
- **stickers:** `stickers string list` / `stickers string get <id>`, `stickers sprint list` / `stickers sprint get <id>` (--include-deleted)
- **crm:** `crm contact-persons create --title "…" --project-id <id>` (optional: --email, --phone, --address, --position, --additional-phone), `crm contacts by-external-id --provider <name> --chat-id <id>`

Global flags:

- `-c, --config` — config file path
- `--json` — output as JSON

## Regenerate API client

After changing `docs/api.json`:

```bash
make generate
go build ./...
```

## Lint and test

```bash
make lint
make test
```
