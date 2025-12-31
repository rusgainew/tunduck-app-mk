package repository

import (
	"context"

	"github.com/rusgainew/tunduck-app-mk/auth-service/internal/domain/entity"
)

// UserRepository - Repository Interface (dependency injection point)
type UserRepository interface {
	// CreateUser - сохранить нового пользователя
	CreateUser(ctx context.Context, user *entity.User) error

	// GetUserByEmail - найти пользователя по email
	GetUserByEmail(ctx context.Context, email string) (*entity.User, error)

	// GetUserByID - найти пользователя по ID
	GetUserByID(ctx context.Context, id string) (*entity.User, error)

	// UpdateUser - обновить пользователя
	UpdateUser(ctx context.Context, user *entity.User) error

	// DeleteUser - удалить пользователя
	DeleteUser(ctx context.Context, id string) error

	// UserExists - проверить существование
	UserExists(ctx context.Context, email string) (bool, error)
}

// EventPublisher - Interface для публикации событий
type EventPublisher interface {
	// PublishUserRegistered - публикует событие регистрации
	PublishUserRegistered(ctx context.Context, user *entity.User) error

	// PublishUserLoggedIn - публикует событие входа
	PublishUserLoggedIn(ctx context.Context, userID string) error

	// PublishUserLoggedOut - публикует событие выхода
	PublishUserLoggedOut(ctx context.Context, userID string) error
}

// TokenBlacklist - Interface для управления черным списком токенов
type TokenBlacklist interface {
	// AddToBlacklist - добавить токен в черный список
	AddToBlacklist(ctx context.Context, token string, expiresIn int64) error

	// IsBlacklisted - проверить в черном списке
	IsBlacklisted(ctx context.Context, token string) (bool, error)
}
