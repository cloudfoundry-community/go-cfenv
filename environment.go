package cfenv

import (
	"os"
)

func CurrentEnv() map[string]string {
	return Env(os.Environ())
}

func Env(env []string) map[string]string {
	vars := mapEnv(env, splitEnv())
	return vars
}
