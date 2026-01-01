package service

import (
	"context"
	"fmt"

	"github.com/rusgainew/tunduck-app-mk/auth-service/internal/application/dto"
	"github.com/rusgainew/tunduck-app-mk/auth-service/internal/domain/entity"
	"github.com/rusgainew/tunduck-app-mk/auth-service/internal/domain/repository"
)

// GetUserService - Application Service для получения профиля пользователя
type GetUserService struct {
	userRepo repository.UserRepository
}

// NewGetUserService - Factory
func NewGetUserService(userRepo repository.UserRepository) *GetUserService {
	return &GetUserService{
		userRepo: userRepo,
	}
}

// Execute - Получить пользователя по ID
func (s *GetUserService) Execute(ctx context.Context, userID string) (*dto.GetMeResponse, error) {
	user, err := s.userRepo.GetUserByID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	if !user.IsActive() {
		return nil, entity.ErrUserInactive
	}

	return dto.UserToGetMeResponse(user), nil
}

// LogoutUserService - Application Service для выхода пользователя
type LogoutUserService struct {
	tokenBlacklist repository.TokenBlacklist
	eventPublisher repository.EventPublisher
}

// NewLogoutUserService - Factory
func NewLogoutUserService(
	tokenBlacklist repository.TokenBlacklist,
	eventPublisher repository.EventPublisher,
) *LogoutUserService {
	return &LogoutUserService{
		tokenBlacklist: tokenBlacklist,
		eventPublisher: eventPublisher,
	}
}

// Execute - Логаут пользователя (добавить токен в blacklist)
func (s *LogoutUserService) Execute(ctx context.Context, userID, token string) error {
	// Добавить токен в blacklist
	if err := s.tokenBlacklist.AddToBlacklist(ctx, token); err != nil {
		return fmt.Errorf("failed to blacklist token: %w", err)
	}

	// Создать и опубликовать событие
	event := entity.NewUserLoggedOut(userID)
	if err := s.eventPublisher.Publish(ctx, event); err != nil {
		fmt.Printf("failed to publish user_logged_out event: %v\n", err)
	}

	return nil
}

// RefreshTokenService - Application Service для обновления токенов
type RefreshTokenService struct {
	userRepo       repository.UserRepository
	tokenService   *TokenService
	tokenBlacklist repository.TokenBlacklist
}

// NewRefreshTokenService - Factory
func NewRefreshTokenService(
	userRepo repository.UserRepository,
	tokenService *TokenService,
	tokenBlacklist repository.TokenBlacklist,
) *RefreshTokenService {
	return &RefreshTokenService{
		userRepo:       userRepo,
		tokenService:   tokenService,
		tokenBlacklist: tokenBlacklist,
	}
}

// Execute - Обновить токены с помощью refresh token
func (s *RefreshTokenService) Execute(ctx context.Context, refreshToken string) (*dto.LoginResponse, error) {
	// 1. Валидировать refresh token
	userID, err := s.tokenService.ValidateRefreshToken(ctx, refreshToken)
	if err != nil {
		return nil, entity.ErrTokenInvalid
	}

	// 2. Проверить, что токен не в blacklist
	isBlacklisted, err := s.tokenBlacklist.IsBlacklisted(ctx, refreshToken)
	if err != nil {
		return nil, fmt.Errorf("failed to check token blacklist: %w", err)
	}
	if isBlacklisted {
		return nil, entity.ErrTokenRevoked
	}

	// 3. Получить пользователя
	user, err := s.userRepo.GetUserByID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	// 4. Проверить активность пользователя
	if !user.IsActive() {
		return nil, entity.ErrUserInactive
	}

	// 5. Добавить старый refresh token в blacklist
	if err := s.tokenBlacklist.AddToBlacklist(ctx, refreshToken); err != nil {
		return nil, fmt.Errorf("failed to blacklist old token: %w", err)
	}

	// 6. Генерировать новую пару токенов
	token, err := s.tokenService.GenerateTokens(ctx, user.ID, user.Email)
	if err != nil {
		return nil, fmt.Errorf("failed to generate tokens: %w", err)
	}

	return dto.TokenToLoginResponse(token), nil
}

// ValidateTokenService - Application Service для валидации токенов
type ValidateTokenService struct {
	tokenService   *TokenService
	tokenBlacklist repository.TokenBlacklist
	userRepo       repository.UserRepository
}

// NewValidateTokenService - Factory
func NewValidateTokenService(
	tokenService *TokenService,
	tokenBlacklist repository.TokenBlacklist,
	userRepo repository.UserRepository,
) *ValidateTokenService {
	return &ValidateTokenService{
		tokenService:   tokenService,
		tokenBlacklist: tokenBlacklist,
		userRepo:       userRepo,
	}
}

// Execute - Валидировать access token и вернуть userID
func (s *ValidateTokenService) Execute(ctx context.Context, accessToken string) (string, error) {
	// 1. Валидировать токен
	userID, err := s.tokenService.ValidateToken(ctx, accessToken)
	if err != nil {
		return "", entity.ErrTokenInvalid
	}

	// 2. Проверить blacklist
	isBlacklisted, err := s.tokenBlacklist.IsBlacklisted(ctx, accessToken)
	if err != nil {
		return "", fmt.Errorf("failed to check token blacklist: %w", err)
	}
	if isBlacklisted {
		return "", entity.ErrTokenRevoked
	}

	// 3. Проверить существование и активность пользователя
	user, err := s.userRepo.GetUserByID(ctx, userID)
	if err != nil {
		return "", entity.ErrUserNotFound
	}

	if !user.IsActive() {
		return "", entity.ErrUserInactive
	}

	return userID, nil
}
