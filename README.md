# depwatch

> Monitors dependency files across repos and alerts on outdated or vulnerable packages

## Installation

```bash
go install github.com/yourorg/depwatch@latest
```

Or build from source:

```bash
git clone https://github.com/yourorg/depwatch.git && cd depwatch && go build ./...
```

## Usage

Point `depwatch` at one or more repositories and it will scan dependency files for outdated or vulnerable packages.

```bash
# Scan a single local repo
depwatch scan ./my-project

# Watch multiple repos and alert on issues
depwatch watch --repos ./service-a,./service-b --interval 24h

# Output results as JSON
depwatch scan ./my-project --format json
```

### Example Output

```
[WARN]  lodash 4.17.15 → 4.17.21  (CVE-2021-23337)
[INFO]  express 4.17.1 → 4.18.2   (outdated)
[OK]    golang.org/x/net           (up to date)
```

### Supported Dependency Files

- `go.mod` / `go.sum`
- `package.json` / `package-lock.json`
- `requirements.txt`
- `Gemfile`

## Configuration

Create a `.depwatch.yaml` in your project root to customize behavior:

```yaml
interval: 12h
severity: warn
ignore:
  - lodash
```

## Contributing

Pull requests are welcome. Please open an issue first to discuss any significant changes.

## License

[MIT](LICENSE)