package handler

import (
	"context"

	"github.com/rusgainew/tunduck-app-mk/auth-service/internal/application/service"
	"github.com/rusgainew/tunduck-app-mk/auth-service/internal/domain/repository"
	"github.com/rusgainew/tunduck-app-mk/auth-service/internal/interfaces/grpc/mapper"
	"github.com/rusgainew/tunduck-app-mk/auth-service/internal/interfaces/grpc/validator"
	authpb "github.com/rusgainew/tunduck-app-mk/proto-lib/auth"
	commonpb "github.com/rusgainew/tunduck-app-mk/proto-lib/common"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// UserHandler обрабатывает операции с пользователями
type UserHandler struct {
	validator      *validator.RequestValidator
	tokenService   *service.TokenService
	userRepo       repository.UserRepository
	tokenBlacklist repository.TokenBlacklist
	userMapper     *mapper.UserMapper
}

// NewUserHandler создает обработчик пользователей
func NewUserHandler(
	validator *validator.RequestValidator,
	tokenService *service.TokenService,
	userRepo repository.UserRepository,
	tokenBlacklist repository.TokenBlacklist,
	userMapper *mapper.UserMapper,
) *UserHandler {
	return &UserHandler{
		validator:      validator,
		tokenService:   tokenService,
		userRepo:       userRepo,
		tokenBlacklist: tokenBlacklist,
		userMapper:     userMapper,
	}
}

// GetUser получает пользователя по ID
func (h *UserHandler) GetUser(ctx context.Context, req *authpb.GetUserRequest) (*authpb.User, error) {
	// Валидация
	if err := h.validator.ValidateGetUserRequest(req); err != nil {
		return nil, err
	}

	// Получение пользователя
	user, err := h.userRepo.GetUserByID(ctx, req.UserId)
	if err != nil {
		return nil, status.Error(codes.NotFound, "user not found")
	}

	return h.userMapper.ToProtoUser(user), nil
}

// Logout выполняет выход пользователя
func (h *UserHandler) Logout(ctx context.Context, req *authpb.LogoutRequest) (*commonpb.Empty, error) {
	// Валидация
	if err := h.validator.ValidateLogoutRequest(req); err != nil {
		return nil, err
	}

	// Валидация токена
	_, err := h.tokenService.ValidateToken(ctx, req.AccessToken)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, "invalid token")
	}

	// Добавление в черный список
	if err := h.tokenBlacklist.AddToBlacklist(ctx, req.AccessToken); err != nil {
		return nil, status.Error(codes.Internal, "failed to logout")
	}

	return &commonpb.Empty{}, nil
}
