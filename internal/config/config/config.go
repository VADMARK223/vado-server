package config

import (
	"log"
	"os"
)

type Port string

type Config struct {
	AppEnv      string
	Port        string
	GrpcPort    string
	GrpcWebPort string
	JwtSecret   string
	TokenTTL    string
	RefreshTTL  string
	GinMode     string
	PostgresDsn string
	KafkaBroker string
}

func Load() *Config {
	cfg := &Config{
		AppEnv:      getEnv("APP_ENV"),
		Port:        getEnv("PORT"),
		GrpcPort:    getEnv("GRPC_PORT"),
		GrpcWebPort: getEnv("GRPC_WEB_PORT"),
		KafkaBroker: getEnv("KAFKA_BROKER"),
		JwtSecret:   getEnv("JWT_SECRET"),
		TokenTTL:    getEnv("TOKEN_TTL"),
		RefreshTTL:  getEnv("REFRESH_TTL"),
		GinMode:     getEnv("GIN_MODE"),
		PostgresDsn: getEnv("POSTGRES_DSN"),
	}

	log.Printf("Loaded config: PORT=%s, GRPC_PORT=%s, GRPC_WEB_PORT=%s, TOKEN_TTL=%s", cfg.Port, cfg.GrpcPort, cfg.GrpcWebPort, cfg.TokenTTL)

	return cfg
}

func getEnv(key string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}

	panic("missing env var: " + key)
}
