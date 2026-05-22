package envcast

import (
	"os"
	"time"
)

func read[T any](key string, fallback *T, parse func(string) (T, error), expected string, tolerant bool) T {
	raw, ok := os.LookupEnv(key)
	if !ok {
		if fallback == nil {
			panicMissing(key)
		}
		return *fallback
	}

	v, err := parse(raw)
	if err != nil {
		if fallback != nil && tolerant {
			return *fallback
		}
		panicInvalid(key, expected, raw)
	}
	return v
}

func readSlice(key, sep string, fallback *[]string, tolerant bool) []string {
	raw, ok := os.LookupEnv(key)
	if !ok {
		if fallback == nil {
			panicMissing(key)
		}
		return *fallback
	}

	v, err := parseStringSlice(raw, sep)
	if err != nil {
		if fallback != nil && tolerant {
			return *fallback
		}
		panicInvalid(key, "[]string", raw)
	}
	return v
}

// String returns the value of the environment variable key.
// It panics if the variable is not set.
func String(key string) string {
	return read(key, nil, parseString, "string", false)
}

// StringOr returns the value of the environment variable key, or fallback if unset.
func StringOr(key, fallback string) string {
	return read(key, &fallback, parseString, "string", false)
}

// StringDefault returns the value of the environment variable key, or fallback if unset or invalid.
func StringDefault(key, fallback string) string {
	return read(key, &fallback, parseString, "string", true)
}

// Int returns key parsed as an int. It panics if the variable is unset or invalid.
func Int(key string) int {
	return read(key, nil, parseInt, "int", false)
}

// IntOr returns key parsed as an int, or fallback if unset.
// It panics if the value is present but invalid.
func IntOr(key string, fallback int) int {
	return read(key, &fallback, parseInt, "int", false)
}

// IntDefault returns key parsed as an int, or fallback if unset or invalid.
func IntDefault(key string, fallback int) int {
	return read(key, &fallback, parseInt, "int", true)
}

// Int64 returns key parsed as an int64. It panics if the variable is unset or invalid.
func Int64(key string) int64 {
	return read(key, nil, parseInt64, "int64", false)
}

// Int64Or returns key parsed as an int64, or fallback if unset.
// It panics if the value is present but invalid.
func Int64Or(key string, fallback int64) int64 {
	return read(key, &fallback, parseInt64, "int64", false)
}

// Int64Default returns key parsed as an int64, or fallback if unset or invalid.
func Int64Default(key string, fallback int64) int64 {
	return read(key, &fallback, parseInt64, "int64", true)
}

// Float64 returns key parsed as a float64. It panics if the variable is unset or invalid.
func Float64(key string) float64 {
	return read(key, nil, parseFloat64, "float64", false)
}

// Float64Or returns key parsed as a float64, or fallback if unset.
// It panics if the value is present but invalid.
func Float64Or(key string, fallback float64) float64 {
	return read(key, &fallback, parseFloat64, "float64", false)
}

// Float64Default returns key parsed as a float64, or fallback if unset or invalid.
func Float64Default(key string, fallback float64) float64 {
	return read(key, &fallback, parseFloat64, "float64", true)
}

// Bool returns key parsed as a bool. It panics if the variable is unset or invalid.
func Bool(key string) bool {
	return read(key, nil, parseBool, "bool", false)
}

// BoolOr returns key parsed as a bool, or fallback if unset.
// It panics if the value is present but invalid.
func BoolOr(key string, fallback bool) bool {
	return read(key, &fallback, parseBool, "bool", false)
}

// BoolDefault returns key parsed as a bool, or fallback if unset or invalid.
func BoolDefault(key string, fallback bool) bool {
	return read(key, &fallback, parseBool, "bool", true)
}

// Duration returns key parsed with time.ParseDuration.
// It panics if the variable is unset or invalid.
func Duration(key string) time.Duration {
	return read(key, nil, parseDuration, "duration", false)
}

// DurationOr returns key parsed with time.ParseDuration, or fallback if unset.
// It panics if the value is present but invalid.
func DurationOr(key string, fallback time.Duration) time.Duration {
	return read(key, &fallback, parseDuration, "duration", false)
}

// DurationDefault returns key parsed with time.ParseDuration, or fallback if unset or invalid.
func DurationDefault(key string, fallback time.Duration) time.Duration {
	return read(key, &fallback, parseDuration, "duration", true)
}

// StringSlice returns key split as a comma-separated list of strings.
// It panics if the variable is unset.
func StringSlice(key string) []string {
	return readSlice(key, ",", nil, false)
}

// StringSliceOr returns key split as a comma-separated list, or fallback if unset.
// It panics if the value is present but invalid.
func StringSliceOr(key string, fallback []string) []string {
	return readSlice(key, ",", &fallback, false)
}

// StringSliceDefault returns key split as a comma-separated list, or fallback if unset or invalid.
func StringSliceDefault(key string, fallback []string) []string {
	return readSlice(key, ",", &fallback, true)
}

// StringSliceSep returns key split using sep as the delimiter.
// It panics if the variable is unset.
func StringSliceSep(key, sep string) []string {
	return readSlice(key, sep, nil, false)
}

// StringSliceSepOr returns key split using sep, or fallback if unset.
// It panics if the value is present but invalid.
func StringSliceSepOr(key, sep string, fallback []string) []string {
	return readSlice(key, sep, &fallback, false)
}

// StringSliceSepDefault returns key split using sep, or fallback if unset or invalid.
func StringSliceSepDefault(key, sep string, fallback []string) []string {
	return readSlice(key, sep, &fallback, true)
}

// Get returns the environment variable key converted to type T.
// Supported types: string, int, int64, float64, bool, time.Duration, []string.
// It panics if the variable is unset or cannot be converted.
func Get[T any](key string) T {
	var zero T
	switch any(zero).(type) {
	case string:
		return any(String(key)).(T)
	case int:
		return any(Int(key)).(T)
	case int64:
		return any(Int64(key)).(T)
	case float64:
		return any(Float64(key)).(T)
	case bool:
		return any(Bool(key)).(T)
	case time.Duration:
		return any(Duration(key)).(T)
	case []string:
		return any(StringSlice(key)).(T)
	default:
		panicUnsupportedType(typeName[T]())
		return zero
	}
}

// GetDefault returns the environment variable key converted to type T, or fallback if unset or invalid.
// Supported types: string, int, int64, float64, bool, time.Duration, []string.
func GetDefault[T any](key string, fallback T) T {
	var zero T
	switch any(zero).(type) {
	case string:
		return any(StringDefault(key, any(fallback).(string))).(T)
	case int:
		return any(IntDefault(key, any(fallback).(int))).(T)
	case int64:
		return any(Int64Default(key, any(fallback).(int64))).(T)
	case float64:
		return any(Float64Default(key, any(fallback).(float64))).(T)
	case bool:
		return any(BoolDefault(key, any(fallback).(bool))).(T)
	case time.Duration:
		return any(DurationDefault(key, any(fallback).(time.Duration))).(T)
	case []string:
		return any(StringSliceDefault(key, any(fallback).([]string))).(T)
	default:
		panicUnsupportedType(typeName[T]())
		return zero
	}
}

func typeName[T any]() string {
	var zero T
	return fmtType(any(zero))
}

func fmtType(v any) string {
	switch v.(type) {
	case string:
		return "string"
	case int:
		return "int"
	case int64:
		return "int64"
	case float64:
		return "float64"
	case bool:
		return "bool"
	case time.Duration:
		return "time.Duration"
	case []string:
		return "[]string"
	default:
		return "unknown"
	}
}
