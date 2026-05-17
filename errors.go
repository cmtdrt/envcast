package envcast

import "fmt"

func panicMissing(key string) {
	panic(fmt.Sprintf("envcast: missing required env var %s", key))
}

func panicInvalid(key, expected, got string) {
	panic(fmt.Sprintf("envcast: invalid value for %s: expected %s, got %q", key, expected, got))
}

func panicUnsupportedType(typ string) {
	panic(fmt.Sprintf("envcast: unsupported type %s for Get", typ))
}
