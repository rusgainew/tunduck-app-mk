package handler

import (
	"context"

	"github.com/rusgainew/tunduck-app-mk/auth-service/internal/application/service"
	"github.com/rusgainew/tunduck-app-mk/auth-service/internal/domain/repository"
	"github.com/rusgainew/tunduck-app-mk/auth-service/internal/interfaces/grpc/mapper"
	"github.com/rusgainew/tunduck-app-mk/auth-service/internal/interfaces/grpc/validator"
	authpb "github.com/rusgainew/tunduck-app-mk/proto-lib/auth"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// TokenHandler обрабатывает операции с токенами
type TokenHandler struct {
	validator      *validator.RequestValidator
	tokenService   *service.TokenService
	userRepo       repository.UserRepository
	tokenBlacklist repository.TokenBlacklist
	userMapper     *mapper.UserMapper
	tokenMapper    *mapper.TokenMapper
}

// NewTokenHandler создает обработчик токенов
func NewTokenHandler(
	validator *validator.RequestValidator,
	tokenService *service.TokenService,
	userRepo repository.UserRepository,
	tokenBlacklist repository.TokenBlacklist,
	userMapper *mapper.UserMapper,
	tokenMapper *mapper.TokenMapper,
) *TokenHandler {
	return &TokenHandler{
		validator:      validator,
		tokenService:   tokenService,
		userRepo:       userRepo,
		tokenBlacklist: tokenBlacklist,
		userMapper:     userMapper,
		tokenMapper:    tokenMapper,
	}
}

// ValidateToken валидирует токен и возвращает пользователя
func (h *TokenHandler) ValidateToken(ctx context.Context, req *authpb.ValidateTokenRequest) (*authpb.User, error) {
	// Валидация запроса
	if err := h.validator.ValidateTokenRequest(req); err != nil {
		return nil, err
	}

	// Валидация токена
	userID, err := h.tokenService.ValidateToken(ctx, req.AccessToken)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, "invalid or expired token")
	}

	// Получение пользователя
	user, err := h.userRepo.GetUserByID(ctx, userID)
	if err != nil {
		return nil, status.Error(codes.NotFound, "user not found")
	}

	return h.userMapper.ToProtoUser(user), nil
}

// RefreshToken обновляет токен доступа
func (h *TokenHandler) RefreshToken(ctx context.Context, req *authpb.RefreshTokenRequest) (*authpb.Token, error) {
	// Валидация запроса
	if err := h.validator.ValidateRefreshTokenRequest(req); err != nil {
		return nil, err
	}

	// Валидация refresh токена
	userID, err := h.tokenService.ValidateRefreshToken(ctx, req.RefreshToken)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, "invalid or expired refresh token")
	}

	// Получение пользователя
	user, err := h.userRepo.GetUserByID(ctx, userID)
	if err != nil {
		return nil, status.Error(codes.NotFound, "user not found")
	}

	// Генерация новых токенов
	tokenResp, err := h.tokenService.GenerateTokens(ctx, user.ID, user.Email)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to generate new token")
	}

	return h.tokenMapper.ToProtoTokenFromLogin(
		tokenResp.AccessToken,
		tokenResp.RefreshToken,
		tokenResp.ExpiresIn,
	), nil
}
