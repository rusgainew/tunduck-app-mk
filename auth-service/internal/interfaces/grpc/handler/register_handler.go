package handler

import (
	"context"
	"time"

	"github.com/rusgainew/tunduck-app-mk/auth-service/internal/application/service"
	"github.com/rusgainew/tunduck-app-mk/auth-service/internal/domain/entity"
	"github.com/rusgainew/tunduck-app-mk/auth-service/internal/interfaces/grpc/adapter"
	"github.com/rusgainew/tunduck-app-mk/auth-service/internal/interfaces/grpc/mapper"
	"github.com/rusgainew/tunduck-app-mk/auth-service/internal/interfaces/grpc/validator"
	authpb "github.com/rusgainew/tunduck-app-mk/proto-lib/auth"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// RegisterHandler - обработчик запроса регистрации пользователя через gRPC
type RegisterHandler struct {
	validator       *validator.RequestValidator  // валидатор входящих запросов
	adapter         *adapter.RequestAdapter      // адаптер для преобразования в DTO
	registerService *service.RegisterUserService // сервис бизнес-логики регистрации
	responseMapper  *mapper.AuthResponseMapper   // маппер для преобразования ответа
	tokenMapper     *mapper.TokenMapper          // маппер для преобразования токенов
}

// NewRegisterHandler - конструктор для создания обработчика регистрации
// Параметры:
//   - validator: валидатор для проверки входящих запросов
//   - adapter: адаптер для преобразования gRPC запроса в DTO
//   - registerService: сервис для обработки бизнес-логики регистрации
//   - responseMapper: маппер для преобразования ответа
//   - tokenMapper: маппер для преобразования токенов
//
// Возвращает: указатель на структуру RegisterHandler
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

// Handle - обрабатывает запрос регистрации нового пользователя
// Параметры:
//   - ctx: контекст запроса
//   - req: gRPC запрос с данными для регистрации (email, пароль, имя, фамилия)
//
// Возвращает: ответ с информацией о пользователе, или ошибка
func (h *RegisterHandler) Handle(ctx context.Context, req *authpb.RegisterRequest) (*authpb.AuthResponse, error) {
	// Валидация входящего запроса (проверка email, пароля и т.д.)
	if err := h.validator.ValidateRegisterRequest(req); err != nil {
		return nil, err
	}

	// Адаптация gRPC запроса к DTO для использования в сервисе
	registerDTO := h.adapter.ToRegisterDTO(req)

	// Выполнение бизнес-логики регистрации (сохранение в БД, отправка событий)
	resp, err := h.registerService.Execute(ctx, registerDTO)
	if err != nil {
		switch err {
		case entity.ErrUserAlreadyExists:
			return nil, status.Error(codes.AlreadyExists, err.Error())
		default:
			return nil, status.Error(codes.Internal, err.Error())
		}
	}

	createdAt := time.Now()
	if parsedTime, parseErr := time.Parse(time.RFC3339, resp.CreatedAt); parseErr == nil {
		createdAt = parsedTime
	}

	user := &entity.User{
		ID:        resp.ID,
		Email:     resp.Email,
		Name:      resp.Name,
		Status:    entity.UserStatus(resp.Status),
		CreatedAt: createdAt,
		UpdatedAt: createdAt,
	}

	return h.responseMapper.ToAuthResponseWithNames(
		user,
		req.FirstName,
		req.LastName,
		h.tokenMapper.ToEmptyToken(),
	), nil
}
