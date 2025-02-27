package utils

import (
	"os"
)

func GetEnvDefault(key, deftVal string) string {
	val, ok := os.LookupEnv(key)
	if !ok {
		return deftVal
	}
	return val
}
