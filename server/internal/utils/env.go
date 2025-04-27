package utils

import (
	"log"
	"os"
)

func RequireEnv(key string) string {
	v := os.Getenv(key)
	if v == "" {
		log.Fatalf("Environment variable %s is required and cannot be empty", key)
	}
	return v
}
