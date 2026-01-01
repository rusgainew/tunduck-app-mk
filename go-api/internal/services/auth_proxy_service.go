package services

import (
	"context"

	"github.com/rusgainew/tunduck-app/internal/models"
)

// AuthProxyService проксирует запросы аутентификации на auth-service через gRPC
// Это позволяет go-api быть API Gateway, делегируя всю auth логику микросервису
type AuthProxyService interface {
	// Register регистрирует нового пользователя через auth-service
	Register(ctx context.Context, req *models.RegisterRequest) (*models.AuthResponse, error)

	// Login выполняет вход пользователя через auth-service
	Login(ctx context.Context, req *models.LoginRequest) (*models.AuthResponse, error)

	// ValidateToken проверяет валидность токена через auth-service
	// Используется middleware для защиты endpoints
	ValidateToken(ctx context.Context, token string) (*models.UserInfo, error)

	// RefreshToken обновляет access token используя refresh token
	RefreshToken(ctx context.Context, refreshToken string) (*models.TokenResponse, error)

	// Logout выходит из системы (добавляет токен в blacklist)
	Logout(ctx context.Context, token string) error

	// GetUserByID получает информацию о пользователе по ID
	GetUserByID(ctx context.Context, userID string) (*models.UserInfo, error)
}
