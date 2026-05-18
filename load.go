package envcast

// Load reads one or more .env files and sets variables in the process environment.
// If no filenames are given, ".env" is used.
//
// Existing environment variables are never overwritten. When multiple files
// define the same key, the first file wins.
//
// Returns an error if a file cannot be read or parsed.
func Load(filenames ...string) error {
	return load(false, filenames...)
}

// Overload is like Load but overwrites variables already set in the environment.
func Overload(filenames ...string) error {
	return load(true, filenames...)
}

// MustLoad calls Load and panics on error.
func MustLoad(filenames ...string) {
	if err := Load(filenames...); err != nil {
		panic(err.Error())
	}
}

// MustOverload calls Overload and panics on error.
func MustOverload(filenames ...string) {
	if err := Overload(filenames...); err != nil {
		panic(err.Error())
	}
}

func load(overload bool, filenames ...string) error {
	if len(filenames) == 0 {
		filenames = []string{defaultEnvFile}
	}
	vars, err := mergeDotEnvFiles(filenames)
	if err != nil {
		return err
	}
	applyDotEnv(vars, overload)
	return nil
}
