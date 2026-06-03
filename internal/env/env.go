package env

import (
	"os"
	"strconv"
)


func GetEnvString(key string, defaultValue string) string {
	value, exists := os.LookupEnv(key)
	if exists {
		return value
	}
	return  defaultValue
}

func GetEnvInt(key string, defaultValue int) int {
	value, exists := os.LookupEnv(key)
	if exists {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}