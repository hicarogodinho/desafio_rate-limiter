package config

import (
	"os"
	"strconv"
	"time"
)

type Config struct {
	RateLimitIP           int
	RateLimitTokenDefault int
	BlockDuration         time.Duration
	RedisAddr             string
	RedisPassword         string
	RedisDB               int
}

func Load() Config {
	return Config{
		RateLimitIP:           getEnvAsInt("RATE_LIMIT_IP", 10),
		RateLimitTokenDefault: getEnvAsInt("RATE_LIMIT_TOKEN_DEFAULT", 100),
		BlockDuration:         time.Duration(getEnvAsInt("BLOCK_DURATION_SECONDS", 300)) * time.Second,
		RedisAddr:             getEnv("REDIS_ADDR", "localhost:6379"),
		RedisPassword:         getEnv("REDIS_PASSWORD", ""),
		RedisDB:               getEnvAsInt("REDIS_DB", 0),
	}
}

func getEnv(key string, defaultVal string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultVal
}

func getEnvAsInt(name string, defaultVal int) int {
	valueStr := getEnv(name, "")
	if value, err := strconv.Atoi(valueStr); err == nil {
		return value
	}
	return defaultVal
}
