package cli

import (
	"os"
	"strconv"
)

func getEnvString(env string, def string) string {
	value := os.Getenv(env)

	if value != "" {
		return value
	} else {
		return def
	}
}

func getEnvInt(env string, def string) int {
	value := getEnvString(env, def)

	// TODO: handle error
	converted, _ := strconv.Atoi(value)

	return converted
}

func getEnvBool(env string, def string) bool {
	value := getEnvString(env, def)

	// TODO: handle error
	converted, _ := strconv.ParseBool(value)

	return converted
}
