package config

import (
	"fmt"
	"os"
	"strconv"
)

// Config - Application configuration
type Config struct {
	// Server
	HttpPort int
	GrpcPort int

	// Database
	DatabaseURL string

	// Redis
	RedisAddr string

	// RabbitMQ
	RabbitMQURL string

	// JWT
	JwtSecret  string
	JwtExpires int64 // seconds

	// Environment
	Environment string
}

// LoadConfig - Load configuration from environment
func LoadConfig() *Config {
	return &Config{
		HttpPort:    parseInt(os.Getenv("HTTP_PORT"), 8001),
		GrpcPort:    parseInt(os.Getenv("GRPC_PORT"), 9001),
		DatabaseURL: getEnv("DATABASE_URL", "postgres://tunduck_user:tunduck_password@localhost:5432/tunduck_auth"),
		RedisAddr:   getEnv("REDIS_ADDR", "localhost:6379"),
		RabbitMQURL: getEnv("RABBITMQ_URL", "amqp://guest:guest@localhost:5672/"),
		JwtSecret:   getEnv("JWT_SECRET", "your-secret-key-change-in-production"),
		JwtExpires:  int64(parseInt(os.Getenv("JWT_EXPIRES"), 3600)), // 1 hour
		Environment: getEnv("ENVIRONMENT", "development"),
	}
}

// Validate - Validate configuration
func (c *Config) Validate() error {
	if c.JwtSecret == "your-secret-key-change-in-production" && c.Environment == "production" {
		return fmt.Errorf("JWT_SECRET must be changed in production")
	}
	if c.HttpPort <= 0 || c.HttpPort > 65535 {
		return fmt.Errorf("invalid HTTP_PORT: %d", c.HttpPort)
	}
	if c.GrpcPort <= 0 || c.GrpcPort > 65535 {
		return fmt.Errorf("invalid GRPC_PORT: %d", c.GrpcPort)
	}
	return nil
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func parseInt(value string, defaultValue int) int {
	if value == "" {
		return defaultValue
	}
	intVal, err := strconv.Atoi(value)
	if err != nil {
		return defaultValue
	}
	return intVal
}
