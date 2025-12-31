package redis

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

// TokenBlacklistRedis - Реализация TokenBlacklist для Redis
type TokenBlacklistRedis struct {
	client *redis.Client
}

// NewTokenBlacklistRedis - Factory
func NewTokenBlacklistRedis(client *redis.Client) *TokenBlacklistRedis {
	return &TokenBlacklistRedis{client: client}
}

// AddToBlacklist - добавить токен в черный список
func (r *TokenBlacklistRedis) AddToBlacklist(ctx context.Context, token string, expiresIn int64) error {
	// Ключ: "blacklist:token_hash"
	// Значение: "true"
	// TTL: время истечения токена

	key := fmt.Sprintf("blacklist:%s", token)
	ttl := time.Duration(expiresIn) * time.Second

	err := r.client.Set(ctx, key, "true", ttl).Err()
	if err != nil {
		return fmt.Errorf("failed to add token to blacklist: %w", err)
	}

	return nil
}

// IsBlacklisted - проверить в черном списке
func (r *TokenBlacklistRedis) IsBlacklisted(ctx context.Context, token string) (bool, error) {
	key := fmt.Sprintf("blacklist:%s", token)

	val, err := r.client.Get(ctx, key).Result()
	if err == redis.Nil {
		return false, nil
	}
	if err != nil {
		return false, fmt.Errorf("failed to check token blacklist: %w", err)
	}

	return val == "true", nil
}
