package utility

import (
	"os"
	"strconv"
	"time"
)

func GetEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

func GetEnvInt(key string, fallback int) int {
	if value, ok := os.LookupEnv(key); ok {
		if i, err := strconv.Atoi(value); err == nil {
			return i
		}
	}
	return fallback
}

func GetEnvBool(key string, fallback bool) bool {
	if value, ok := os.LookupEnv(key); ok {
		if b, err := strconv.ParseBool(value); err == nil {
			return b
		}
	}
	return fallback
}

func GetEnvDuration(key string, fallback time.Duration) time.Duration {
	if value, ok := os.LookupEnv(key); ok {
		if dur, err := time.ParseDuration(value); err == nil {
			return dur
		}
	}
	return fallback
}
