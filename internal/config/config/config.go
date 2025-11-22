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
	RefreshTTL         string
	GinMode            string
	PostgresDsn        string
	KafkaEnable        bool
	KafkaBroker        string
	corsAllowedOrigins string
	tokenTTL           string // Время жизни токена в секунда
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
		tokenTTL:           getEnv("TOKEN_TTL"),
		RefreshTTL:         getEnv("REFRESH_TTL"),
		GinMode:            getEnv("GIN_MODE"),
		PostgresDsn:        getEnv("POSTGRES_DSN"),
	}

	log.Printf("Loaded config: PORT=%s, GRPC_PORT=%s, GRPC_WEB_PORT=%s, TOKEN_TTL=%s", cfg.Port, cfg.GrpcPort, cfg.GrpcWebPort, cfg.tokenTTL)

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

func (cfg *Config) TokenTTL() string {
	return cfg.tokenTTL
}

/*func (cfg *Config) Test() time.Duration {
	ttl, err := strconv.Atoi(cfg.tokenTTL)
	if err != nil {
		panic(pp.Sprintf("Failed to parse token TTL: %v", err))
	}
	return time.Duration(ttl) * time.Second
}*/

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
