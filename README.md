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
    // Load .env into the process environment (optional, for local dev)
    _ = envcast.Load()

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

## Loading `.env` files

```go
err := envcast.Load()                    // default: .env
err := envcast.Load(".env", ".env.local") // multiple files, first key wins
envcast.MustLoad()

envcast.Overload(".env") // overwrites existing environment variables
```

`Load` never overwrites variables already set in the environment (same idea as godotenv). Use `Overload` to force values from the file.

## Behavior

| Situation              | Strict (`Int`, `Bool`, …) | With fallback (`IntOr`, …) |
|------------------------|---------------------------|----------------------------|
| Variable not set       | panic                     | returns fallback           |
| Variable set, invalid  | panic                     | panic                      |

Panic messages are explicit, e.g. `envcast: missing required env var PORT` or `envcast: invalid value for PORT: expected int, got "abc"`.

## Supported types

`string`, `int`, `int64`, `float64`, `bool`, `time.Duration`, `[]string` (CSV or custom separator).

