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

## Commands

- `yougile config path` — print config file path
- `yougile config show` — show config (api_key masked in human output)
- `yougile company get` — current company details
- `yougile users list` / `users get <id>`
- `yougile projects list` / `projects get <id>`
- `yougile boards list` / `boards get <id>`
- `yougile columns list` / `columns get <id>`
- `yougile tasks list` / `tasks get <id>`
- `yougile departments list` / `departments get <id>`
- `yougile webhooks list`
- `yougile files upload <path>`

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
