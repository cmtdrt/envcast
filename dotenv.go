package envcast

import (
	"bytes"
	"fmt"
	"os"
	"strings"
)

const defaultEnvFile = ".env"

// Parses the contents of a .env file into key-value pairs.
// Later duplicate keys in the same file are ignored (first wins).
func parseDotEnv(data []byte) (map[string]string, error) {
	vars := make(map[string]string)
	lines := bytes.Split(data, []byte{'\n'})
	for i, line := range lines {
		line = bytes.TrimSuffix(line, []byte{'\r'})
		key, value, ok, err := parseDotEnvLine(string(line))
		if err != nil {
			return nil, fmt.Errorf("envcast: parse .env line %d: %w", i+1, err)
		}
		if !ok {
			continue
		}
		if _, exists := vars[key]; !exists {
			vars[key] = value
		}
	}
	return vars, nil
}

func parseDotEnvLine(line string) (key, value string, ok bool, err error) {
	line = strings.TrimSpace(line)
	// Empty line or comment
	if line == "" || strings.HasPrefix(line, "#") {
		return "", "", false, nil
	}

	if strings.HasPrefix(line, "export ") {
		line = strings.TrimSpace(strings.TrimPrefix(line, "export "))
	}

	eq := strings.Index(line, "=")
	if eq < 0 {
		return "", "", false, fmt.Errorf("expected KEY=VALUE, got %q", line)
	}

	key = strings.TrimSpace(line[:eq])
	if key == "" {
		return "", "", false, fmt.Errorf("empty key in %q", line)
	}

	value, err = parseDotEnvValue(strings.TrimSpace(line[eq+1:]))
	if err != nil {
		return "", "", false, err
	}
	return key, value, true, nil
}

func parseDotEnvValue(raw string) (string, error) {
	if raw == "" {
		return "", nil
	}

	switch raw[0] {
	case '"':
		return parseDoubleQuoted(raw)
	case '\'':
		return parseSingleQuoted(raw)
	default:
		return parseUnquoted(raw), nil
	}
}

func parseUnquoted(raw string) string {
	if i := strings.Index(raw, " #"); i >= 0 {
		raw = raw[:i]
	}
	return strings.TrimSpace(raw)
}

func parseSingleQuoted(raw string) (string, error) {
	if !strings.HasSuffix(raw, "'") || len(raw) < 2 {
		return "", fmt.Errorf("unterminated single-quoted value %q", raw)
	}
	return raw[1 : len(raw)-1], nil
}

func parseDoubleQuoted(raw string) (string, error) {
	if len(raw) < 2 || raw[0] != '"' {
		return "", fmt.Errorf("invalid double-quoted value %q", raw)
	}
	var b strings.Builder
	for i := 1; i < len(raw); i++ {
		if raw[i] == '\\' && i+1 < len(raw) {
			switch raw[i+1] {
			case 'n':
				b.WriteByte('\n')
			case 'r':
				b.WriteByte('\r')
			case 't':
				b.WriteByte('\t')
			case '"', '\\', '\'':
				b.WriteByte(raw[i+1])
			default:
				b.WriteByte(raw[i+1])
			}
			i++
			continue
		}
		if raw[i] == '"' {
			if i != len(raw)-1 {
				return "", fmt.Errorf("invalid double-quoted value %q", raw)
			}
			return b.String(), nil
		}
		b.WriteByte(raw[i])
	}
	return "", fmt.Errorf("unterminated double-quoted value %q", raw)
}

func mergeDotEnvFiles(filenames []string) (map[string]string, error) {
	merged := make(map[string]string)
	for _, name := range filenames {
		data, err := os.ReadFile(name)
		if err != nil {
			return nil, fmt.Errorf("envcast: read %s: %w", name, err)
		}
		vars, err := parseDotEnv(data)
		if err != nil {
			return nil, err
		}
		for k, v := range vars {
			if _, exists := merged[k]; !exists {
				merged[k] = v
			}
		}
	}
	return merged, nil
}

func applyDotEnv(vars map[string]string, overload bool) {
	for key, value := range vars {
		if !overload {
			if _, exists := os.LookupEnv(key); exists {
				continue
			}
		}
		_ = os.Setenv(key, value)
	}
}
