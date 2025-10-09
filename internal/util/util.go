package util

import (
	"log"
	"os"
)

func GetEnv(key string) string {
	value := os.Getenv(key)
	if value == "" {
		log.Fatalf("%s is not set in enviroment!", key)
	}

	return value
}
