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
		Port:        getEnv("PORT"),
		GrpcPort:    getEnv("GRPC_PORT"),
		GrpcWebPort: getEnv("GRPC_WEB_PORT"),
		KafkaPort:   getEnv(string(kafkaPort)),
		JwtSecret:   getEnv("JWT_SECRET"),
		TokenTTL:    getEnv("TOKEN_TTL"),
		RefreshTTL:  getEnv("REFRESH_TTL"),
		GinMode:     getEnv("GIN_MODE"),
		PostgresDsn: getEnv("POSTGRES_DSN"),
	}

	log.Printf("Loaded config: PORT=%s, GRPC_PORT=%s, GRPC_WEB_PORT=%s, %s=%s, TOKEN_TTL=%s", cfg.Port, cfg.GrpcPort, cfg.GrpcWebPort, kafkaPort, cfg.KafkaPort, cfg.TokenTTL)

	return cfg
}

func getEnv(key string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}

	panic("missing env var: " + key)
}
