package redis

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

// TokenBlacklistRedis - реализация интерфейса TokenBlacklist для Redis
// Используется для хранения заблокированных токенов при логауте пользователя
type TokenBlacklistRedis struct {
	client *redis.Client // подключение к Redis
}

// NewTokenBlacklistRedis - конструктор для создания хранилища черного списка токенов
// Параметры:
//   - client: клиент Redis
//
// Возвращает: указатель на структуру TokenBlacklistRedis
func NewTokenBlacklistRedis(client *redis.Client) *TokenBlacklistRedis {
	return &TokenBlacklistRedis{client: client}
}

// AddToBlacklist - добавляет токен в черный список (список заблокированных токенов)
// При выходе пользователя его токен добавляется, чтобы его нельзя было использовать
// Параметры:
//   - ctx: контекст запроса
//   - token: JWT токен для добавления в черный список
//
// Возвращает: ошибка, если не удалось добавить токен
func (r *TokenBlacklistRedis) AddToBlacklist(ctx context.Context, token string) error {
	// Ключ: "blacklist:token_hash" - для быстрого поиска
	// TTL: 24 часа (стандартное время жизни JWT токена доступа)
	key := fmt.Sprintf("blacklist:%s", token)
	ttl := 24 * time.Hour

	err := r.client.Set(ctx, key, "true", ttl).Err()
	if err != nil {
		return fmt.Errorf("failed to add token to blacklist: %w", err)
	}

	return nil
}

// IsBlacklisted - проверяет, находится ли токен в черном списке
// Используется при валидации токена, чтобы убедиться, что пользователь вышел
// Параметры:
//   - ctx: контекст запроса
//   - token: JWT токен для проверки
//
// Возвращает: true если токен в черном списке, false если нет, ошибка при ошибке Redis
func (r *TokenBlacklistRedis) IsBlacklisted(ctx context.Context, token string) (bool, error) {
	key := fmt.Sprintf("blacklist:%s", token)

	val, err := r.client.Get(ctx, key).Result()
	if err == redis.Nil {
		// Ключ не найден - токен не в черном списке
		return false, nil
	}
	if err != nil {
		return false, fmt.Errorf("failed to check token blacklist: %w", err)
	}

	return val == "true", nil
}
