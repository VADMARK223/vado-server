package main

import (
	"log"
	"os"
	"vado_server/internal/config/code"

	"github.com/joho/godotenv"
)

func loadEnv() {
	env := os.Getenv("APP_ENV")
	if env == "" {
		env = code.Local // по умолчанию, если не задано
	}
	switch env {
	case code.Local:
		if err := godotenv.Load(".env.local"); err != nil {
			log.Println("⚠️  .env.local not found — using system env")
		} else {
			log.Println("✅ Loaded .env.local")
		}
	default:
		log.Println("ℹ️  Running in", env, "mode — skipping local env")
	}
}
