package config

import (
	"log"
	"os"
)

type Port string

const (
	kafkaPort Port = "KAFKA_PORT"
)

type Config struct {
	Port        string
	GrpcPort    string
	GrpcWebPort string
	KafkaPort   string
	JwtSecret   string
	TokenTTL    string
	RefreshTTL  string
	GinMode     string
	PostgresDsn string
}

func Load() *Config {
	cfg := &Config{
		Port:        getEnv("PORT", "5555"),
		GrpcPort:    getEnv("GRPC_PORT", "50051"),
		GrpcWebPort: getEnv("GRPC_WEB_PORT", "8090"),
		KafkaPort:   getEnv(string(kafkaPort), "9094"),
		JwtSecret:   getEnv("JWT_SECRET", "asdfkjh"),
		TokenTTL:    getEnv("TOKEN_TTL", "900"),
		RefreshTTL:  getEnv("REFRESH_TTL", "604800"),
		GinMode:     getEnv("GIN_MODE", "debug"),
		PostgresDsn: getEnv("POSTGRES_DSN", "DNS"),
	}

	log.Printf("Loaded config: PORT=%s, GRPC_PORT=%s, GRPC_WEB_PORT=%s, %s=%s, TOKEN_TTL=%s", cfg.Port, cfg.GrpcPort, cfg.GrpcWebPort, kafkaPort, cfg.KafkaPort, cfg.TokenTTL)

	return cfg
}

func getEnv(key, def string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}

	return def
}
