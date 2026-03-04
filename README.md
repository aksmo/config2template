# config2template

Convert config files into templates and extract environment-specific values into a separate env file.

Useful for safely committing config templates to version control (no secrets) while keeping real values in an env file that stays out of the repo.

## How it works

- **String** leaf values → replaced with `${KEY_NAME}` placeholder and recorded in the env file
- **Numbers**, **booleans**, **null** → left as-is in the template (not environment-specific)
- Key names are derived from the key path (`database.password` → `DATABASE_PASSWORD`)

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
config2template --input config.json [--output config.json.tpl] [--env config.json.env] [--format json|yaml|toml]
```

| Flag | Default | Description |
|------|---------|-------------|
| `--input` | *(required)* | Input config file (.json, .yaml/.yml, .toml, .env) |
| `--output` | `<input>.tpl` | Output template file |
| `--env` | `<input>.env` | Output env vars file |
| `--format` | *(auto from ext)* | Force format: `json`, `yaml`, `toml`, `env` |

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

## YAML example

**Input** `config.yaml`:
```yaml
database:
  host: db.example.com
  port: 5432
  password: s3cr3t
api_key: sk-abc123
debug: false
```

**Output** `config.yaml.tpl`:
```yaml
api_key: ${API_KEY}
database:
    host: ${DATABASE_HOST}
    password: ${DATABASE_PASSWORD}
    port: 5432
debug: false
```

## TOML example

**Input** `config.toml`:
```toml
api_key = "sk-abc123"
debug = false

[database]
host = "db.example.com"
port = 5432
password = "s3cr3t"
```

**Output** `config.toml.tpl`:
```toml
api_key = "${API_KEY}"
debug = false

[database]
  host = "${DATABASE_HOST}"
  password = "${DATABASE_PASSWORD}"
  port = 5432
```

## .env example

**Input** `.env`:
```
# App secrets
DB_HOST=db.example.com
DB_PASSWORD=s3cr3t
API_KEY="sk-abc123"
export STRIPE_KEY=sk_live_xyz
DEBUG=false
```

**Output** `.env.tpl`:
```
API_KEY=${API_KEY}
DB_HOST=${DB_HOST}
DB_PASSWORD=${DB_PASSWORD}
DEBUG=${DEBUG}
STRIPE_KEY=${STRIPE_KEY}
```

Parser supports: bare values, `"double"` and `'single'` quotes, `export KEY=value`, inline `# comments`, and blank lines.

## Supported formats

- [x] JSON
- [x] YAML
- [x] TOML
- [x] .env
