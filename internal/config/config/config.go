package config

import (
	"log"
	"os"
	"strconv"
	"strings"
)

type Port string

type Config struct {
	AppEnv             string
	Port               string
	GrpcPort           string
	GrpcWebPort        string
	JwtSecret          string
	TokenTTL           string // Время жизни токена в секунда
	RefreshTTL         string
	GinMode            string
	PostgresDsn        string
	KafkaEnable        bool
	KafkaBroker        string
	corsAllowedOrigins string
}

func Load() *Config {
	cfg := &Config{
		AppEnv:             getEnv("APP_ENV"),
		Port:               getEnv("PORT"),
		corsAllowedOrigins: getEnv("CORS_ALLOWED_ORIGINS"),
		GrpcPort:           getEnv("GRPC_PORT"),
		GrpcWebPort:        getEnv("GRPC_WEB_PORT"),
		KafkaEnable:        getEnvBool("KAFKA_ENABLE"),
		KafkaBroker:        getEnv("KAFKA_BROKER"),
		JwtSecret:          getEnv("JWT_SECRET"),
		TokenTTL:           getEnv("TOKEN_TTL"),
		RefreshTTL:         getEnv("REFRESH_TTL"),
		GinMode:            getEnv("GIN_MODE"),
		PostgresDsn:        getEnv("POSTGRES_DSN"),
	}

	log.Printf("Loaded config: PORT=%s, GRPC_PORT=%s, GRPC_WEB_PORT=%s, TOKEN_TTL=%s", cfg.Port, cfg.GrpcPort, cfg.GrpcWebPort, cfg.TokenTTL)

	return cfg
}

func (cfg *Config) CorsAllowedOrigins() map[string]bool {
	result := make(map[string]bool)
	port := cfg.Port

	for _, value := range strings.Split(cfg.corsAllowedOrigins, ",") {
		value = strings.TrimSpace(value)
		if value != "" {
			key := value + ":" + port
			result[key] = true
		}
	}
	return result
}

func getEnv(key string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}

	panic("missing env var: " + key)
}

func getEnvBool(key string) bool {
	val := getEnv(key)

	parsed, err := strconv.ParseBool(val)
	if err != nil {
		panic("env " + key + " must be true/false, got: " + val)
	}

	return parsed
}
