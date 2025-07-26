package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	Port      string
	Database  DatabaseConfig
	JWT       JWTConfig
	CORS      CORSConfig
	RateLimit RateLimitConfig
	LogLevel  string
}

type DatabaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	Name     string
}

type JWTConfig struct {
	Secret      string
	ExpireHours int
}

type CORSConfig struct {
	Origin string
}

type RateLimitConfig struct {
	Max        int
	Expiration int
}

func Load() *Config {
	// Load .env file if it exists
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
	}

	expireHours, _ := strconv.Atoi(getEnv("JWT_EXPIRE_HOURS", "24"))
	rateMax, _ := strconv.Atoi(getEnv("RATE_LIMIT_MAX", "100"))
	rateExp, _ := strconv.Atoi(getEnv("RATE_LIMIT_EXPIRATION", "1"))

	return &Config{
		Port: getEnv("PORT", "8000"),
		Database: DatabaseConfig{
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     getEnv("DB_PORT", "3306"),
			User:     getEnv("DB_USER", "root"),
			Password: getEnv("DB_PASSWORD", ""),
			Name:     getEnv("DB_NAME", "code_valley"),
		},
		JWT: JWTConfig{
			Secret:      getEnv("JWT_SECRET", "your-secret-key"),
			ExpireHours: expireHours,
		},
		CORS: CORSConfig{
			Origin: getEnv("CORS_ORIGIN", "*"),
		},
		RateLimit: RateLimitConfig{
			Max:        rateMax,
			Expiration: rateExp,
		},
		LogLevel: getEnv("LOG_LEVEL", "info"),
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
