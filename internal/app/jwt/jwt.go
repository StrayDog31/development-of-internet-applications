package config

import (
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type Config struct {
	JWT   JWTConfig
	Redis RedisConfig
}

type JWTConfig struct {
	Token         string
	ExpiresIn     time.Duration
	SigningMethod jwt.SigningMethod
}

type RedisConfig struct {
	Host        string
	Port        string
	Password    string
	User        string
	DialTimeout time.Duration
	ReadTimeout time.Duration
}

func LoadConfig() *Config {
	jwtToken := os.Getenv("JWT_SECRET")
	if jwtToken == "" {
		jwtToken = "default-secret-key"
	}

	jwtExpiresIn := 24 * time.Hour

	return &Config{
		JWT: JWTConfig{
			Token:         jwtToken,
			ExpiresIn:     jwtExpiresIn,
			SigningMethod: jwt.SigningMethodHS256,
		},
		Redis: RedisConfig{
			Host:        getEnv("REDIS_HOST", "localhost"),
			Port:        getEnv("REDIS_PORT", "6379"),
			User:        getEnv("REDIS_USER", ""),
			Password:    getEnv("REDIS_PASSWORD", ""),
			DialTimeout: 5 * time.Second,
			ReadTimeout: 3 * time.Second,
		},
	}
}

func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}