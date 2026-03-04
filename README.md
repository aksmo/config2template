# config2template

Convert config files into templates and extract environment-specific values into a separate env file.

Useful for safely committing config templates to version control (no secrets) while keeping real values in an env file that stays out of the repo.

## How it works

- **String** leaf values → replaced with `${KEY_NAME}` placeholder and recorded in the env file
- **Numbers**, **booleans**, **null** → left as-is in the template (not environment-specific)
- Key names are derived from the JSON path (`database.password` → `DATABASE_PASSWORD`)

## Install

```bash
go install config2template@latest
```

Or build from source:

```bash
git clone https://github.com/aksmo/config2template
cd config2template
go build -o config2template .
```

## Usage

```bash
config2template --input config.json [--output config.json.tpl] [--env config.json.env]
```

| Flag | Default | Description |
|------|---------|-------------|
| `--input` | *(required)* | Input JSON config file |
| `--output` | `<input>.tpl` | Output template file |
| `--env` | `<input>.env` | Output env vars file |

## Example

**Input** `config.json`:
```json
{
  "database": {
    "host": "db.example.com",
    "port": 5432,
    "password": "s3cr3t",
    "name": "myapp"
  },
  "api_key": "sk-abc123",
  "debug": false
}
```

**Output** `config.json.tpl`:
```json
{
  "api_key": "${API_KEY}",
  "database": {
    "host": "${DATABASE_HOST}",
    "name": "${DATABASE_NAME}",
    "password": "${DATABASE_PASSWORD}",
    "port": 5432
  },
  "debug": false
}
```

**Output** `config.json.env`:
```
API_KEY=sk-abc123
DATABASE_HOST=db.example.com
DATABASE_NAME=myapp
DATABASE_PASSWORD=s3cr3t
```

Commit `config.json.tpl` to version control. Keep `config.json.env` out of the repo (add it to `.gitignore`).

## Supported formats

- [x] JSON
- [ ] YAML *(planned)*
- [ ] TOML *(planned)*
- [ ] .env *(planned)*
