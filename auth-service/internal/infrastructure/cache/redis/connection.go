package redis

import (
	"context"
	"fmt"

	"github.com/redis/go-redis/v9"
)

// Config - конфигурация для Redis
type Config struct {
	Host     string
	Port     int
	Password string
	DB       int
}

// NewConnection - создает новое подключение к Redis
func NewConnection(cfg Config) (*redis.Client, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", cfg.Host, cfg.Port),
		Password: cfg.Password,
		DB:       cfg.DB,
	})

	// Проверка соединения
	ctx := context.Background()
	if err := client.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("failed to connect to Redis: %w", err)
	}

	return client, nil
}

// DefaultConfig - конфигурация по умолчанию
func DefaultConfig() Config {
	return Config{
		Host:     "localhost",
		Port:     6379,
		Password: "",
		DB:       0,
	}
}
