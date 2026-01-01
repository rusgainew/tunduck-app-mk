package handler

import (
	"context"

	"github.com/rusgainew/tunduck-app-mk/auth-service/internal/application/service"
	"github.com/rusgainew/tunduck-app-mk/auth-service/internal/interfaces/grpc/adapter"
	"github.com/rusgainew/tunduck-app-mk/auth-service/internal/interfaces/grpc/mapper"
	"github.com/rusgainew/tunduck-app-mk/auth-service/internal/interfaces/grpc/validator"
	authpb "github.com/rusgainew/tunduck-app-mk/proto-lib/auth"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// RegisterHandler обрабатывает регистрацию пользователя
type RegisterHandler struct {
	validator       *validator.RequestValidator
	adapter         *adapter.RequestAdapter
	registerService *service.RegisterUserService
	responseMapper  *mapper.AuthResponseMapper
	tokenMapper     *mapper.TokenMapper
}

// NewRegisterHandler создает обработчик регистрации
func NewRegisterHandler(
	validator *validator.RequestValidator,
	adapter *adapter.RequestAdapter,
	registerService *service.RegisterUserService,
	responseMapper *mapper.AuthResponseMapper,
	tokenMapper *mapper.TokenMapper,
) *RegisterHandler {
	return &RegisterHandler{
		validator:       validator,
		adapter:         adapter,
		registerService: registerService,
		responseMapper:  responseMapper,
		tokenMapper:     tokenMapper,
	}
}

// Handle обрабатывает запрос регистрации
func (h *RegisterHandler) Handle(ctx context.Context, req *authpb.RegisterRequest) (*authpb.AuthResponse, error) {
	// Валидация
	if err := h.validator.ValidateRegisterRequest(req); err != nil {
		return nil, err
	}

	// Адаптация к DTO
	registerDTO := h.adapter.ToRegisterDTO(req)

	// Выполнение бизнес-логики
	resp, err := h.registerService.Execute(ctx, registerDTO)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	// Создание ответа (токен выдаётся при логине)
	// TODO: переделать на получение entity.User из сервиса
	// Пока используем временное решение
	user := &authpb.User{
		Id:        resp.ID,
		Email:     resp.Email,
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Status:    "active",
	}

	return &authpb.AuthResponse{
		User:      user,
		Token:     h.tokenMapper.ToEmptyToken(),
		Timestamp: 0, // будет установлен в responseMapper
	}, nil
}
