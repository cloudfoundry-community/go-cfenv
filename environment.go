package cfenv

import (
	"os"
)

// Translate the current environment to a map[string]string.
func CurrentEnv() map[string]string {
	return Env(os.Environ())
}

// Translate the provided environment to a map[string]string.
func Env(env []string) map[string]string {
	vars := mapEnv(env, splitEnv())
	return vars
}
