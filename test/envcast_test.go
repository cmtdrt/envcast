package envcast_test

import (
	"strings"
	"testing"
	"time"

	"github.com/cdrouet/envcast"
)

func assertPanicContains(t *testing.T, substr string, fn func()) {
	t.Helper()
	defer func() {
		r := recover()
		if r == nil {
			t.Fatalf("expected panic containing %q", substr)
		}
		msg, ok := r.(string)
		if !ok {
			t.Fatalf("panic value type %T, want string", r)
		}
		if !strings.Contains(msg, substr) {
			t.Fatalf("panic %q does not contain %q", msg, substr)
		}
	}()
	fn()
}

func TestString(t *testing.T) {
	t.Setenv("ENVCAST_STRING", "hello")
	if got := envcast.String("ENVCAST_STRING"); got != "hello" {
		t.Fatalf("String() = %q, want hello", got)
	}
}

func TestString_missing(t *testing.T) {
	assertPanicContains(t, "missing required env var ENVCAST_STRING_MISSING", func() {
		_ = envcast.String("ENVCAST_STRING_MISSING")
	})
}

func TestStringOr(t *testing.T) {
	t.Setenv("ENVCAST_STRING_OR", "custom")
	if got := envcast.StringOr("ENVCAST_STRING_OR", "default"); got != "custom" {
		t.Fatalf("StringOr() = %q, want custom", got)
	}
	if got := envcast.StringOr("ENVCAST_STRING_OR_ABSENT", "default"); got != "default" {
		t.Fatalf("StringOr() = %q, want default", got)
	}
}

func TestInt(t *testing.T) {
	t.Setenv("ENVCAST_INT", "42")
	if got := envcast.Int("ENVCAST_INT"); got != 42 {
		t.Fatalf("Int() = %d, want 42", got)
	}
}

func TestInt_missing(t *testing.T) {
	assertPanicContains(t, "missing required env var ENVCAST_INT_MISSING", func() {
		_ = envcast.Int("ENVCAST_INT_MISSING")
	})
}

func TestInt_invalid(t *testing.T) {
	t.Setenv("ENVCAST_INT_BAD", "abc")
	assertPanicContains(t, "invalid value for ENVCAST_INT_BAD: expected int, got \"abc\"", func() {
		_ = envcast.Int("ENVCAST_INT_BAD")
	})
}

func TestIntOr(t *testing.T) {
	if got := envcast.IntOr("ENVCAST_INT_OR_ABSENT", 8080); got != 8080 {
		t.Fatalf("IntOr() = %d, want 8080", got)
	}
	t.Setenv("ENVCAST_INT_OR", "3000")
	if got := envcast.IntOr("ENVCAST_INT_OR", 8080); got != 3000 {
		t.Fatalf("IntOr() = %d, want 3000", got)
	}
	t.Setenv("ENVCAST_INT_OR_BAD", "nope")
	assertPanicContains(t, "invalid value for ENVCAST_INT_OR_BAD", func() {
		_ = envcast.IntOr("ENVCAST_INT_OR_BAD", 8080)
	})
}

func TestInt64(t *testing.T) {
	t.Setenv("ENVCAST_INT64", "9223372036854775807")
	if got := envcast.Int64("ENVCAST_INT64"); got != 9223372036854775807 {
		t.Fatalf("Int64() = %d", got)
	}
}

func TestInt64_missing(t *testing.T) {
	assertPanicContains(t, "missing required env var ENVCAST_INT64_MISSING", func() {
		_ = envcast.Int64("ENVCAST_INT64_MISSING")
	})
}

func TestInt64_invalid(t *testing.T) {
	t.Setenv("ENVCAST_INT64_BAD", "x")
	assertPanicContains(t, "expected int64", func() {
		_ = envcast.Int64("ENVCAST_INT64_BAD")
	})
}

func TestInt64Or(t *testing.T) {
	if got := envcast.Int64Or("ENVCAST_INT64_OR_ABSENT", 1); got != 1 {
		t.Fatalf("Int64Or() = %d", got)
	}
}

func TestFloat64(t *testing.T) {
	t.Setenv("ENVCAST_FLOAT", "3.14")
	if got := envcast.Float64("ENVCAST_FLOAT"); got != 3.14 {
		t.Fatalf("Float64() = %v", got)
	}
}

func TestFloat64_missing(t *testing.T) {
	assertPanicContains(t, "missing required env var ENVCAST_FLOAT_MISSING", func() {
		_ = envcast.Float64("ENVCAST_FLOAT_MISSING")
	})
}

func TestFloat64_invalid(t *testing.T) {
	t.Setenv("ENVCAST_FLOAT_BAD", "pi")
	assertPanicContains(t, "expected float64", func() {
		_ = envcast.Float64("ENVCAST_FLOAT_BAD")
	})
}

func TestFloat64Or(t *testing.T) {
	if got := envcast.Float64Or("ENVCAST_FLOAT_OR_ABSENT", 1.5); got != 1.5 {
		t.Fatalf("Float64Or() = %v", got)
	}
}

func TestBool(t *testing.T) {
	t.Setenv("ENVCAST_BOOL", "true")
	if got := envcast.Bool("ENVCAST_BOOL"); !got {
		t.Fatal("Bool() = false, want true")
	}
}

func TestBool_missing(t *testing.T) {
	assertPanicContains(t, "missing required env var ENVCAST_BOOL_MISSING", func() {
		_ = envcast.Bool("ENVCAST_BOOL_MISSING")
	})
}

func TestBool_invalid(t *testing.T) {
	t.Setenv("ENVCAST_BOOL_BAD", "maybe")
	assertPanicContains(t, "expected bool", func() {
		_ = envcast.Bool("ENVCAST_BOOL_BAD")
	})
}

func TestBoolOr(t *testing.T) {
	if got := envcast.BoolOr("ENVCAST_BOOL_OR_ABSENT", false); got {
		t.Fatal("BoolOr() = true, want false")
	}
	t.Setenv("ENVCAST_BOOL_OR", "1")
	if got := envcast.BoolOr("ENVCAST_BOOL_OR", false); !got {
		t.Fatal("BoolOr() = false, want true")
	}
}

func TestDuration(t *testing.T) {
	t.Setenv("ENVCAST_DURATION", "5s")
	if got := envcast.Duration("ENVCAST_DURATION"); got != 5*time.Second {
		t.Fatalf("Duration() = %v", got)
	}
}

func TestDuration_missing(t *testing.T) {
	assertPanicContains(t, "missing required env var ENVCAST_DURATION_MISSING", func() {
		_ = envcast.Duration("ENVCAST_DURATION_MISSING")
	})
}

func TestDuration_invalid(t *testing.T) {
	t.Setenv("ENVCAST_DURATION_BAD", "forever")
	assertPanicContains(t, "expected duration", func() {
		_ = envcast.Duration("ENVCAST_DURATION_BAD")
	})
}

func TestDurationOr(t *testing.T) {
	want := 2 * time.Minute
	if got := envcast.DurationOr("ENVCAST_DURATION_OR_ABSENT", want); got != want {
		t.Fatalf("DurationOr() = %v", got)
	}
}

func TestStringSlice(t *testing.T) {
	t.Setenv("ENVCAST_SLICE", "a, b ,c")
	got := envcast.StringSlice("ENVCAST_SLICE")
	want := []string{"a", "b", "c"}
	if len(got) != len(want) {
		t.Fatalf("StringSlice() = %v", got)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("StringSlice()[%d] = %q, want %q", i, got[i], want[i])
		}
	}
}

func TestStringSlice_empty(t *testing.T) {
	t.Setenv("ENVCAST_SLICE_EMPTY", "")
	got := envcast.StringSlice("ENVCAST_SLICE_EMPTY")
	if len(got) != 0 {
		t.Fatalf("StringSlice() = %v, want empty slice", got)
	}
}

func TestStringSlice_missing(t *testing.T) {
	assertPanicContains(t, "missing required env var ENVCAST_SLICE_MISSING", func() {
		_ = envcast.StringSlice("ENVCAST_SLICE_MISSING")
	})
}

func TestStringSliceOr(t *testing.T) {
	fallback := []string{"localhost"}
	if got := envcast.StringSliceOr("ENVCAST_SLICE_OR_ABSENT", fallback); len(got) != 1 || got[0] != "localhost" {
		t.Fatalf("StringSliceOr() = %v", got)
	}
}

func TestStringSliceSep(t *testing.T) {
	t.Setenv("ENVCAST_SLICE_SEP", "a|b|c")
	got := envcast.StringSliceSep("ENVCAST_SLICE_SEP", "|")
	want := []string{"a", "b", "c"}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("StringSliceSep()[%d] = %q", i, got[i])
		}
	}
}

func TestStringSliceSepOr(t *testing.T) {
	fallback := []string{"x"}
	if got := envcast.StringSliceSepOr("ENVCAST_SLICE_SEP_OR_ABSENT", ";", fallback); len(got) != 1 {
		t.Fatalf("StringSliceSepOr() = %v", got)
	}
}

func TestGet(t *testing.T) {
	t.Setenv("ENVCAST_GET_INT", "7")
	if got := envcast.Get[int]("ENVCAST_GET_INT"); got != 7 {
		t.Fatalf("Get[int]() = %d", got)
	}
	t.Setenv("ENVCAST_GET_STRING", "ok")
	if got := envcast.Get[string]("ENVCAST_GET_STRING"); got != "ok" {
		t.Fatalf("Get[string]() = %q", got)
	}
}

func TestGet_unsupported(t *testing.T) {
	assertPanicContains(t, "unsupported type", func() {
		_ = envcast.Get[struct{}]("ENVCAST_GET_STRUCT")
	})
}

func TestGet_missing(t *testing.T) {
	assertPanicContains(t, "missing required env var ENVCAST_GET_MISSING", func() {
		_ = envcast.Get[int]("ENVCAST_GET_MISSING")
	})
}
