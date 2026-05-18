package envcast_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/cdrouet/envcast"
)

func writeEnvFile(t *testing.T, dir, name, content string) string {
	t.Helper()
	path := filepath.Join(dir, name)
	if err := os.WriteFile(path, []byte(content), 0o600); err != nil {
		t.Fatal(err)
	}
	return path
}

func TestLoad(t *testing.T) {
	dir := t.TempDir()
	writeEnvFile(t, dir, ".env", "ENVCAST_LOAD_PORT=9090\nENVCAST_LOAD_DEBUG=true\n")

	t.Chdir(dir)
	os.Unsetenv("ENVCAST_LOAD_PORT")
	os.Unsetenv("ENVCAST_LOAD_DEBUG")

	if err := envcast.Load(); err != nil {
		t.Fatal(err)
	}
	if got := envcast.Int("ENVCAST_LOAD_PORT"); got != 9090 {
		t.Fatalf("Int() = %d, want 9090", got)
	}
	if !envcast.Bool("ENVCAST_LOAD_DEBUG") {
		t.Fatal("Bool() = false, want true")
	}
}

func TestLoad_doesNotOverrideExisting(t *testing.T) {
	dir := t.TempDir()
	writeEnvFile(t, dir, ".env", "ENVCAST_LOAD_OVERRIDE=from_file\n")

	t.Chdir(dir)
	t.Setenv("ENVCAST_LOAD_OVERRIDE", "from_env")

	if err := envcast.Load(); err != nil {
		t.Fatal(err)
	}
	if got := envcast.String("ENVCAST_LOAD_OVERRIDE"); got != "from_env" {
		t.Fatalf("String() = %q, want from_env", got)
	}
}

func TestOverload_overridesExisting(t *testing.T) {
	dir := t.TempDir()
	writeEnvFile(t, dir, ".env", "ENVCAST_LOAD_OVERLOAD=from_file\n")

	t.Chdir(dir)
	t.Setenv("ENVCAST_LOAD_OVERLOAD", "from_env")

	if err := envcast.Overload(); err != nil {
		t.Fatal(err)
	}
	if got := envcast.String("ENVCAST_LOAD_OVERLOAD"); got != "from_file" {
		t.Fatalf("String() = %q, want from_file", got)
	}
}

func TestLoad_multipleFilesFirstWins(t *testing.T) {
	dir := t.TempDir()
	f1 := writeEnvFile(t, dir, ".env", "ENVCAST_LOAD_MULTI=first\n")
	f2 := writeEnvFile(t, dir, ".env.local", "ENVCAST_LOAD_MULTI=second\n")

	os.Unsetenv("ENVCAST_LOAD_MULTI")
	if err := envcast.Load(f1, f2); err != nil {
		t.Fatal(err)
	}
	if got := envcast.String("ENVCAST_LOAD_MULTI"); got != "first" {
		t.Fatalf("String() = %q, want first", got)
	}
}

func TestLoad_exportAndQuotes(t *testing.T) {
	dir := t.TempDir()
	content := `export ENVCAST_LOAD_QUOTED="hello world"
ENVCAST_LOAD_SINGLE='raw value'
# comment
ENVCAST_LOAD_EMPTY=
`
	writeEnvFile(t, dir, ".env", content)

	t.Chdir(dir)
	os.Unsetenv("ENVCAST_LOAD_QUOTED")
	os.Unsetenv("ENVCAST_LOAD_SINGLE")
	os.Unsetenv("ENVCAST_LOAD_EMPTY")

	if err := envcast.Load(); err != nil {
		t.Fatal(err)
	}
	if got := envcast.String("ENVCAST_LOAD_QUOTED"); got != "hello world" {
		t.Fatalf("String() = %q", got)
	}
	if got := envcast.String("ENVCAST_LOAD_SINGLE"); got != "raw value" {
		t.Fatalf("String() = %q", got)
	}
	if got := envcast.String("ENVCAST_LOAD_EMPTY"); got != "" {
		t.Fatalf("String() = %q, want empty", got)
	}
}

func TestLoad_missingFile(t *testing.T) {
	dir := t.TempDir()
	t.Chdir(dir)
	if err := envcast.Load(); err == nil {
		t.Fatal("expected error for missing .env")
	}
}

func TestLoad_invalidLine(t *testing.T) {
	dir := t.TempDir()
	writeEnvFile(t, dir, ".env", "NOT_A_VALID_LINE\n")
	t.Chdir(dir)
	if err := envcast.Load(); err == nil {
		t.Fatal("expected parse error")
	}
}

func TestMustLoad_panics(t *testing.T) {
	dir := t.TempDir()
	t.Chdir(dir)
	assertPanicContains(t, "envcast: read .env", func() {
		envcast.MustLoad()
	})
}
