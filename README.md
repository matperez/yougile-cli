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
- **projects:** `projects list` / `projects get <id>` / `projects create --title "…"` / `projects update <id> [--title "…"]`; **roles:** `projects roles list --project-id <id>` / `projects roles get --project-id <id> <role-id>` / `projects roles create --project-id <id> --name "…"` / `projects roles update --project-id <id> <role-id> [--name "…"]` / `projects roles delete --project-id <id> <role-id>`
- **boards:** `boards list` / `boards get <id>` / `boards create --title "…" --project-id <id>` / `boards update <id> [--title "…"]`
- **columns:** `columns list` / `columns get <id>` / `columns create --title "…" --board-id <id>` / `columns update <id> [--title "…"]`
- **tasks:** `tasks list` / `tasks get <id>` / `tasks create --title "…" [--column-id <id>]` / `tasks update <id>` with optional `--title`, `--column-id`, `--description`, `--color`, `--assigned <id1,id2>`, `--completed true|false`, `--archived true|false`, `--deleted true|false` / `tasks chat-subscribers get <task-id>` / `tasks chat-subscribers update <task-id> --user-ids "id1,id2"`
- **departments:** `departments list` / `departments get <id>` / `departments create --title "…" [--parent-id <id>]` / `departments update <id> [--title "…"]`
- **webhooks:** `webhooks list` / `webhooks create --event "…" --url "…"`
- `yougile files upload <path>`
- **chats:** `chats list` / `chats get <id>` / `chats create --title "…"` / `chats update <id> [--title "…"]`; **messages:** `chats messages list <chat-id>`, `chats messages send <chat-id> --text "…"`, `chats messages update <chat-id> <message-id> [--label "…"]`
- **stickers:** `stickers string list` / `stickers string get <id>` / `stickers string create --name "…"` / `stickers string update <id> [--name "…"]`; **string states:** `stickers string states list <sticker-id>` / `stickers string states get <sticker-id> <state-id>` / `stickers string states create <sticker-id> --name "…"` / `stickers string states update <sticker-id> <state-id> [--name "…"]`; `stickers sprint list` / `stickers sprint get <id>` / `stickers sprint create --name "…"` / `stickers sprint update <id> [--name "…"]`; **sprint states:** `stickers sprint states list <sticker-id>` / `stickers sprint states get <sticker-id> <state-id>` / `stickers sprint states create <sticker-id> --name "…"` / `stickers sprint states update <sticker-id> <state-id> [--name "…"]` (--include-deleted for list)
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

## License

MIT License — see [LICENSE](LICENSE) for details.
