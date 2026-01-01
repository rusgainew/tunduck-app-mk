package service

import (
	"context"
	"fmt"
	"time"

	"github.com/rusgainew/tunduck-app-mk/auth-service/internal/application/dto"
	"github.com/rusgainew/tunduck-app-mk/auth-service/internal/domain/entity"
	"github.com/rusgainew/tunduck-app-mk/auth-service/internal/domain/repository"
	"github.com/rusgainew/tunduck-app-mk/auth-service/internal/infrastructure/jwt"
	"golang.org/x/crypto/bcrypt"
)

// RegisterUserService - Application Service для регистрации
type RegisterUserService struct {
	userRepo       repository.UserRepository
	eventPublisher repository.EventPublisher
}

// NewRegisterUserService - Factory
func NewRegisterUserService(
	userRepo repository.UserRepository,
	eventPublisher repository.EventPublisher,
) *RegisterUserService {
	return &RegisterUserService{
		userRepo:       userRepo,
		eventPublisher: eventPublisher,
	}
}

// Execute - Бизнес-логика регистрации
func (s *RegisterUserService) Execute(ctx context.Context, req *dto.RegisterRequest) (*dto.RegisterResponse, error) {
	// 1. Валидация credential через domain layer
	credential, err := entity.NewCredential(req.Email, req.Password)
	if err != nil {
		return nil, err
	}

	// 2. Проверить, есть ли уже пользователь с таким email
	exists, err := s.userRepo.UserExists(ctx, credential.Email)
	if err != nil {
		return nil, fmt.Errorf("failed to check user existence: %w", err)
	}
	if exists {
		return nil, entity.ErrUserAlreadyExists
	}

	// 3. Хеширование пароля
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(credential.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	// 4. Создать User aggregate
	userID := fmt.Sprintf("user_%d", time.Now().UnixNano())
	user, err := entity.NewUser(
		userID,
		credential.Email,
		req.Name,
		string(passwordHash),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	// 5. Сохранить в БД
	if err := s.userRepo.CreateUser(ctx, user); err != nil {
		return nil, fmt.Errorf("failed to save user: %w", err)
	}

	// 6. Опубликовать domain events
	for _, event := range user.DomainEvents() {
		if err := s.eventPublisher.Publish(ctx, event); err != nil {
			// Логировать ошибку, но не падать (eventual consistency)
			fmt.Printf("failed to publish event %s: %v\n", event.EventName(), err)
		}
	}
	user.ClearDomainEvents()

	// 7. Вернуть response
	return dto.UserToRegisterResponse(user), nil
}

// LoginUserService - Application Service для входа
type LoginUserService struct {
	userRepo       repository.UserRepository
	tokenService   *TokenService
	eventPublisher repository.EventPublisher
}

// NewLoginUserService - Factory
func NewLoginUserService(
	userRepo repository.UserRepository,
	tokenService *TokenService,
	eventPublisher repository.EventPublisher,
) *LoginUserService {
	return &LoginUserService{
		userRepo:       userRepo,
		tokenService:   tokenService,
		eventPublisher: eventPublisher,
	}
}

// Execute - Бизнес-логика входа
func (s *LoginUserService) Execute(ctx context.Context, req *dto.LoginRequest, ipAddress string) (*dto.LoginResponse, error) {
	// 1. Найти пользователя
	user, err := s.userRepo.GetUserByEmail(ctx, req.Email)
	if err != nil {
		return nil, entity.ErrInvalidCredentials
	}

	// 2. Проверить активность
	if !user.IsActive() {
		return nil, entity.ErrUserInactive
	}

	// 3. Проверить блокировку
	if user.IsBlocked() {
		return nil, entity.ErrUserBlocked
	}

	// 4. Проверить пароль
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return nil, entity.ErrInvalidCredentials
	}

	// 5. Обновить последний вход с IP адресом
	user.UpdateLastLogin(ipAddress)
	if err := s.userRepo.UpdateUser(ctx, user); err != nil {
		return nil, fmt.Errorf("failed to update user: %w", err)
	}

	// 6. Создать токены
	token, err := s.tokenService.GenerateTokens(ctx, user.ID, user.Email)
	if err != nil {
		return nil, fmt.Errorf("failed to generate tokens: %w", err)
	}

	// 7. Опубликовать domain events
	for _, event := range user.DomainEvents() {
		if err := s.eventPublisher.Publish(ctx, event); err != nil {
			fmt.Printf("failed to publish event %s: %v\n", event.EventName(), err)
		}
	}
	user.ClearDomainEvents()

	// 8. Вернуть response
	return dto.TokenToLoginResponse(token), nil
}

// TokenService - Service для управления токенами
type TokenService struct {
	jwtManager *jwt.JWTManager
}

// NewTokenService - Factory
func NewTokenService(jwtManager *jwt.JWTManager) *TokenService {
	return &TokenService{
		jwtManager: jwtManager,
	}
}

// GenerateTokens - Генерирует access и refresh токены
func (s *TokenService) GenerateTokens(ctx context.Context, userID string, email string) (*entity.Token, error) {
	accessToken, err := s.jwtManager.GenerateAccessToken(userID, email)
	if err != nil {
		return nil, fmt.Errorf("failed to generate access token: %w", err)
	}

	refreshToken, err := s.jwtManager.GenerateRefreshToken(userID)
	if err != nil {
		return nil, fmt.Errorf("failed to generate refresh token: %w", err)
	}

	return entity.NewToken(accessToken, refreshToken, s.jwtManager.GetExpiresIn()), nil
}

// ValidateToken - Проверяет валидность токена и возвращает userID
func (s *TokenService) ValidateToken(ctx context.Context, token string) (string, error) {
	claims, err := s.jwtManager.ValidateAccessToken(token)
	if err != nil {
		return "", err
	}
	return claims.UserID, nil
}

// ValidateRefreshToken - Проверяет refresh token
func (s *TokenService) ValidateRefreshToken(ctx context.Context, token string) (string, error) {
	claims, err := s.jwtManager.ValidateRefreshToken(token)
	if err != nil {
		return "", err
	}
	return claims.UserID, nil
}

// GetJWTManager - Expose JWT manager for external use
func (s *TokenService) GetJWTManager() *jwt.JWTManager {
	return s.jwtManager
}
