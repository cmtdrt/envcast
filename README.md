# envcast

Typed, strict environment variable parsing for Go — a small alternative to `os.Getenv` + `strconv`.

Requires **Go 1.26** or later.

## Install

```bash
go get github.com/cdrouet/envcast
```

## Usage

```go
import (
    "time"

    "github.com/cdrouet/envcast"
)

func main() {
    // Strict: panic if missing or invalid
    port := envcast.Int("PORT")
    debug := envcast.Bool("DEBUG")
    timeout := envcast.Duration("TIMEOUT")

    // With fallback: panic only if invalid
    port = envcast.IntOr("PORT", 8080)
    debug = envcast.BoolOr("DEBUG", false)
    timeout = envcast.DurationOr("TIMEOUT", 5*time.Second)

    hosts := envcast.StringSliceOr("HOSTS", []string{"localhost"})
}
```

## Behavior

| Situation              | Strict (`Int`, `Bool`, …) | With fallback (`IntOr`, …) |
|------------------------|---------------------------|----------------------------|
| Variable not set       | panic                     | returns fallback           |
| Variable set, invalid  | panic                     | panic                      |

Panic messages are explicit, e.g. `envcast: missing required env var PORT` or `envcast: invalid value for PORT: expected int, got "abc"`.

## Supported types

`string`, `int`, `int64`, `float64`, `bool`, `time.Duration`, `[]string` (CSV or custom separator).

Generic helper: `envcast.Get[int]("WORKERS")`.

## Development

```bash
go test ./...
go test -cover ./...
```

## License

MIT — see [LICENSE](LICENSE).
