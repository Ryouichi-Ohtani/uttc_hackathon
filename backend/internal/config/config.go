package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
)

type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	JWT      JWTConfig
	AI       AIConfig
	Storage  StorageConfig
	CORS     CORSConfig
}

type ServerConfig struct {
	Port    string
	GinMode string
	Env     string
}

type DatabaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
	Params   string
}


type JWTConfig struct {
	Secret          string
	ExpirationHours int
}

type AIConfig struct {
	ServiceURL string
}

type StorageConfig struct {
	GCSBucketName string
	CDNBaseURL    string
}

type CORSConfig struct {
	AllowedOrigins []string
}

func Load() (*Config, error) {
	// Load .env file if it exists (development)
	_ = godotenv.Load()

	config := &Config{
		Server: ServerConfig{
			Port:    getEnv("PORT", "8080"),
			GinMode: getEnv("GIN_MODE", "debug"),
			Env:     getEnv("ENV", "development"),
		},
		Database: DatabaseConfig{
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     getEnv("DB_PORT", "3306"),
			User:     getEnv("DB_USER", "ecomate"),
			Password: getEnv("DB_PASSWORD", ""),
			DBName:   getEnv("DB_NAME", "ecomate_db"),
			Params:   getEnv("DB_PARAMS", "parseTime=true&charset=utf8mb4&loc=UTC"),
		},

		JWT: JWTConfig{
			Secret:          getEnv("JWT_SECRET", ""),
			ExpirationHours: getEnvAsInt("JWT_EXPIRATION_HOURS", 72),
		},
		AI: AIConfig{
			ServiceURL: getEnv("AI_SERVICE_URL", "localhost:50051"),
		},
		Storage: StorageConfig{
			GCSBucketName: getEnv("GCS_BUCKET_NAME", ""),
			CDNBaseURL:    getEnv("CDN_BASE_URL", ""),
		},
		CORS: CORSConfig{
			AllowedOrigins: strings.Split(getEnv("ALLOWED_ORIGINS", "http://localhost:3000"), ","),
		},
	}

	// Validate required fields
	if config.JWT.Secret == "" {
		return nil, fmt.Errorf("JWT_SECRET is required")
	}

	if config.Database.Password == "" {
		return nil, fmt.Errorf("DB_PASSWORD is required")
	}

	return config, nil
}

func (c *DatabaseConfig) DSN() string {
	params := c.Params
	if params == "" {
		params = "parseTime=true&charset=utf8mb4&loc=UTC"
	}

	network := "tcp"
	address := fmt.Sprintf("%s:%s", c.Host, c.Port)
	if strings.HasPrefix(c.Host, "/") {
		network = "unix"
		address = c.Host
	}

	return fmt.Sprintf(
		"%s:%s@%s(%s)/%s?%s",
		c.User, c.Password, network, address, c.DBName, params,
	)
}


func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvAsInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}
