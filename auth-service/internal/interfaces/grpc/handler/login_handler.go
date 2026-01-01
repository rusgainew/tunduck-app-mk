package handler

import (
	"context"

	"github.com/rusgainew/tunduck-app-mk/auth-service/internal/application/service"
	"github.com/rusgainew/tunduck-app-mk/auth-service/internal/domain/repository"
	"github.com/rusgainew/tunduck-app-mk/auth-service/internal/interfaces/grpc/adapter"
	"github.com/rusgainew/tunduck-app-mk/auth-service/internal/interfaces/grpc/mapper"
	"github.com/rusgainew/tunduck-app-mk/auth-service/internal/interfaces/grpc/validator"
	authpb "github.com/rusgainew/tunduck-app-mk/proto-lib/auth"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/peer"
	"google.golang.org/grpc/status"
)

// LoginHandler обрабатывает вход пользователя
type LoginHandler struct {
	validator      *validator.RequestValidator
	adapter        *adapter.RequestAdapter
	loginService   *service.LoginUserService
	userRepo       repository.UserRepository
	responseMapper *mapper.AuthResponseMapper
	tokenMapper    *mapper.TokenMapper
}

// NewLoginHandler создает обработчик входа
func NewLoginHandler(
	validator *validator.RequestValidator,
	adapter *adapter.RequestAdapter,
	loginService *service.LoginUserService,
	userRepo repository.UserRepository,
	responseMapper *mapper.AuthResponseMapper,
	tokenMapper *mapper.TokenMapper,
) *LoginHandler {
	return &LoginHandler{
		validator:      validator,
		adapter:        adapter,
		loginService:   loginService,
		userRepo:       userRepo,
		responseMapper: responseMapper,
		tokenMapper:    tokenMapper,
	}
}

// Handle обрабатывает запрос входа
func (h *LoginHandler) Handle(ctx context.Context, req *authpb.LoginRequest) (*authpb.AuthResponse, error) {
	// Валидация
	if err := h.validator.ValidateLoginRequest(req); err != nil {
		return nil, err
	}

	// Адаптация к DTO
	loginDTO := h.adapter.ToLoginDTO(req)

	// Извлечение IP адреса
	ipAddress := h.extractIPAddress(ctx)

	// Выполнение бизнес-логики
	tokenResp, err := h.loginService.Execute(ctx, loginDTO, ipAddress)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, err.Error())
	}

	// Получение пользователя
	user, err := h.userRepo.GetUserByEmail(ctx, req.Email)
	if err != nil {
		return nil, status.Error(codes.NotFound, "user not found")
	}

	// Преобразование токена
	protoToken := h.tokenMapper.ToProtoTokenFromLogin(
		tokenResp.AccessToken,
		tokenResp.RefreshToken,
		tokenResp.ExpiresIn,
	)

	// Создание ответа
	return h.responseMapper.ToAuthResponseWithToken(user, protoToken), nil
}

// extractIPAddress извлекает IP адрес из контекста
func (h *LoginHandler) extractIPAddress(ctx context.Context) string {
	if p, ok := peer.FromContext(ctx); ok && p.Addr != nil {
		return p.Addr.String()
	}
	return ""
}
