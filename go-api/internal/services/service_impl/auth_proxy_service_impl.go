package service_impl

import (
	"context"
	"fmt"

	grpcClient "github.com/rusgainew/tunduck-app/internal/clients/grpc"
	"github.com/rusgainew/tunduck-app/internal/models"
	"github.com/rusgainew/tunduck-app/internal/services"
	"github.com/rusgainew/tunduck-app/pkg/apperror"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type authProxyService struct {
	authClient *grpcClient.AuthClient
	logger     *logrus.Logger
}

// NewAuthProxyService создает новый Auth Proxy Service
func NewAuthProxyService(authClient *grpcClient.AuthClient, logger *logrus.Logger) services.AuthProxyService {
	return &authProxyService{
		authClient: authClient,
		logger:     logger,
	}
}

// Register проксирует регистрацию на auth-service
func (s *authProxyService) Register(ctx context.Context, req *models.RegisterRequest) (*models.AuthResponse, error) {
	s.logger.WithFields(logrus.Fields{
		"email": req.Email,
	}).Info("Proxying register request to auth-service")

	// Валидация входных данных
	if err := req.Validate(); err != nil {
		return nil, apperror.NewBadRequestError(fmt.Sprintf("Validation error: %v", err))
	}

	// Вызов auth-service через gRPC
	grpcResp, err := s.authClient.Register(ctx, req.Email, req.Password, req.FullName, req.FullName)
	if err != nil {
		return nil, s.handleGRPCError(err, "register")
	}

	// Конвертация gRPC response → HTTP response
	response := &models.AuthResponse{
		Token: grpcResp.Token.AccessToken,
		User: models.UserInfo{
			ID:       grpcResp.User.Id,
			Email:    grpcResp.User.Email,
			FullName: grpcResp.User.FirstName + " " + grpcResp.User.LastName,
			Role:     grpcResp.User.Role, // RBAC: роль пользователя
		},
	}

	s.logger.WithField("user_id", response.User.ID).Info("User registered successfully via auth-service")
	return response, nil
}

// Login проксирует вход на auth-service
func (s *authProxyService) Login(ctx context.Context, req *models.LoginRequest) (*models.AuthResponse, error) {
	s.logger.WithFields(logrus.Fields{
		"email": req.Email,
	}).Info("Proxying login request to auth-service")

	// Валидация входных данных
	if req.Email == "" || req.Password == "" {
		return nil, apperror.NewBadRequestError("Email and password are required")
	}

	// Вызов auth-service через gRPC
	grpcResp, err := s.authClient.Login(ctx, req.Email, req.Password)
	if err != nil {
		return nil, s.handleGRPCError(err, "login")
	}

	// Конвертация gRPC response → HTTP response
	response := &models.AuthResponse{
		Token: grpcResp.Token.AccessToken,
		User: models.UserInfo{
			ID:       grpcResp.User.Id,
			Email:    grpcResp.User.Email,
			FullName: grpcResp.User.FirstName + " " + grpcResp.User.LastName,
			Role:     grpcResp.User.Role, // RBAC: роль пользователя
		},
	}

	s.logger.WithField("user_id", response.User.ID).Info("User logged in successfully via auth-service")
	return response, nil
}

// ValidateToken проверяет токен через auth-service
func (s *authProxyService) ValidateToken(ctx context.Context, token string) (*models.UserInfo, error) {
	if token == "" {
		return nil, apperror.NewUnauthorizedError("Token is required")
	}

	// Вызов auth-service через gRPC
	grpcUser, err := s.authClient.ValidateToken(ctx, token)
	if err != nil {
		return nil, s.handleGRPCError(err, "validate_token")
	}

	// Конвертация gRPC response → UserInfo
	userInfo := &models.UserInfo{
		ID:       grpcUser.Id,
		Email:    grpcUser.Email,
		FullName: grpcUser.FirstName + " " + grpcUser.LastName,
		Role:     grpcUser.Role, // RBAC: роль пользователя
	}

	return userInfo, nil
}

// RefreshToken обновляет токен через auth-service
func (s *authProxyService) RefreshToken(ctx context.Context, refreshToken string) (*models.TokenResponse, error) {
	s.logger.Info("Proxying refresh token request to auth-service")

	if refreshToken == "" {
		return nil, apperror.NewBadRequestError("Refresh token is required")
	}

	// Вызов auth-service через gRPC
	grpcToken, err := s.authClient.RefreshToken(ctx, refreshToken)
	if err != nil {
		return nil, s.handleGRPCError(err, "refresh_token")
	}

	// Конвертация gRPC response → TokenResponse
	response := &models.TokenResponse{
		AccessToken:  grpcToken.AccessToken,
		RefreshToken: grpcToken.RefreshToken,
		ExpiresIn:    int(grpcToken.ExpiresIn),
		TokenType:    grpcToken.TokenType,
	}

	s.logger.Info("Token refreshed successfully via auth-service")
	return response, nil
}

// Logout выполняет выход через auth-service
func (s *authProxyService) Logout(ctx context.Context, token string) error {
	s.logger.Info("Proxying logout request to auth-service")

	if token == "" {
		return apperror.NewBadRequestError("Token is required")
	}

	// Вызов auth-service через gRPC
	err := s.authClient.Logout(ctx, token)
	if err != nil {
		return s.handleGRPCError(err, "logout")
	}

	s.logger.Info("User logged out successfully via auth-service")
	return nil
}

// GetUserByID получает пользователя через auth-service
func (s *authProxyService) GetUserByID(ctx context.Context, userID string) (*models.UserInfo, error) {
	s.logger.WithField("user_id", userID).Info("Proxying get user request to auth-service")

	if userID == "" {
		return nil, apperror.NewBadRequestError("User ID is required")
	}

	// Вызов auth-service через gRPC
	grpcUser, err := s.authClient.GetUser(ctx, userID)
	if err != nil {
		return nil, s.handleGRPCError(err, "get_user")
	}

	// Конвертация gRPC response → UserInfo
	userInfo := &models.UserInfo{
		ID:       grpcUser.Id,
		Email:    grpcUser.Email,
		FullName: grpcUser.FirstName + " " + grpcUser.LastName,
	}

	return userInfo, nil
}

// handleGRPCError конвертирует gRPC ошибки в apperror для унифицированной обработки
func (s *authProxyService) handleGRPCError(err error, operation string) error {
	s.logger.WithError(err).WithField("operation", operation).Error("gRPC call to auth-service failed")

	// Извлечение статуса gRPC
	st, ok := status.FromError(err)
	if !ok {
		return apperror.NewInternalError(fmt.Sprintf("Auth service error: %v", err))
	}

	// Маппинг gRPC статусов на HTTP ошибки
	switch st.Code() {
	case codes.InvalidArgument:
		return apperror.NewBadRequestError(st.Message())
	case codes.Unauthenticated:
		return apperror.NewUnauthorizedError(st.Message())
	case codes.PermissionDenied:
		return apperror.NewForbiddenError(st.Message())
	case codes.NotFound:
		return apperror.NewNotFoundError(st.Message())
	case codes.AlreadyExists:
		return apperror.NewConflictError(st.Message())
	case codes.Unavailable:
		return apperror.NewServiceUnavailableError("Auth service is temporarily unavailable")
	case codes.Internal:
		return apperror.NewInternalError("Internal auth service error")
	case codes.DeadlineExceeded:
		return apperror.NewServiceUnavailableError("Auth service request timeout")
	default:
		return apperror.NewInternalError(fmt.Sprintf("Auth service error: %v", st.Message()))
	}
}
