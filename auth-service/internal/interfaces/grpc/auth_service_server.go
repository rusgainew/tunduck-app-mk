package grpc

import (
	"context"
	"fmt"
	"net"

	"github.com/rusgainew/tunduck-app-mk/auth-service/internal/application/service"
	"github.com/rusgainew/tunduck-app-mk/auth-service/internal/domain/repository"
	"github.com/rusgainew/tunduck-app-mk/auth-service/internal/interfaces/grpc/adapter"
	"github.com/rusgainew/tunduck-app-mk/auth-service/internal/interfaces/grpc/handler"
	"github.com/rusgainew/tunduck-app-mk/auth-service/internal/interfaces/grpc/mapper"
	"github.com/rusgainew/tunduck-app-mk/auth-service/internal/interfaces/grpc/validator"
	authpb "github.com/rusgainew/tunduck-app-mk/proto-lib/auth"
	commonpb "github.com/rusgainew/tunduck-app-mk/proto-lib/common"
	"google.golang.org/grpc"
)

// AuthServiceServer - gRPC server для AuthService
// Следует принципу SRP: только делегирует вызовы обработчикам
type AuthServiceServer struct {
	authpb.UnimplementedAuthServiceServer
	registerHandler *handler.RegisterHandler
	loginHandler    *handler.LoginHandler
	tokenHandler    *handler.TokenHandler
	userHandler     *handler.UserHandler
}

// NewAuthServiceServer создает новый gRPC сервер с применением чистой архитектуры
func NewAuthServiceServer(
	registerService *service.RegisterUserService,
	loginService *service.LoginUserService,
	tokenService *service.TokenService,
	userRepo repository.UserRepository,
	tokenBlacklist repository.TokenBlacklist,
) *AuthServiceServer {
	// Создание компонентов согласно SRP
	validator := validator.NewRequestValidator()
	adapter := adapter.NewRequestAdapter()

	// Mappers
	userMapper := mapper.NewUserMapper()
	tokenMapper := mapper.NewTokenMapper()
	responseMapper := mapper.NewAuthResponseMapper(userMapper, tokenMapper)

	// Handlers
	registerHandler := handler.NewRegisterHandler(
		validator,
		adapter,
		registerService,
		responseMapper,
		tokenMapper,
	)

	loginHandler := handler.NewLoginHandler(
		validator,
		adapter,
		loginService,
		userRepo,
		responseMapper,
		tokenMapper,
	)

	tokenHandler := handler.NewTokenHandler(
		validator,
		tokenService,
		userRepo,
		tokenBlacklist,
		userMapper,
		tokenMapper,
	)

	userHandler := handler.NewUserHandler(
		validator,
		tokenService,
		userRepo,
		tokenBlacklist,
		userMapper,
	)

	return &AuthServiceServer{
		registerHandler: registerHandler,
		loginHandler:    loginHandler,
		tokenHandler:    tokenHandler,
		userHandler:     userHandler,
	}
}

// Register делегирует обработку запроса регистрации
func (s *AuthServiceServer) Register(ctx context.Context, req *authpb.RegisterRequest) (*authpb.AuthResponse, error) {
	return s.registerHandler.Handle(ctx, req)
}

// Login делегирует обработку запроса входа
func (s *AuthServiceServer) Login(ctx context.Context, req *authpb.LoginRequest) (*authpb.AuthResponse, error) {
	return s.loginHandler.Handle(ctx, req)
}

// ValidateToken делегирует валидацию токена
func (s *AuthServiceServer) ValidateToken(ctx context.Context, req *authpb.ValidateTokenRequest) (*authpb.User, error) {
	return s.tokenHandler.ValidateToken(ctx, req)
}

// GetUser делегирует получение пользователя
func (s *AuthServiceServer) GetUser(ctx context.Context, req *authpb.GetUserRequest) (*authpb.User, error) {
	return s.userHandler.GetUser(ctx, req)
}

// Logout делегирует выход пользователя
func (s *AuthServiceServer) Logout(ctx context.Context, req *authpb.LogoutRequest) (*commonpb.Empty, error) {
	return s.userHandler.Logout(ctx, req)
}

// RefreshToken делегирует обновление токена
func (s *AuthServiceServer) RefreshToken(ctx context.Context, req *authpb.RefreshTokenRequest) (*authpb.Token, error) {
	return s.tokenHandler.RefreshToken(ctx, req)
}

// GRPCServer инкапсулирует запуск gRPC сервера
// Следует SRP: отвечает только за жизненный цикл сервера
type GRPCServer struct {
	server *grpc.Server
	addr   string
}

// NewGRPCServer создает новый gRPC сервер
func NewGRPCServer(addr string, authServiceServer *AuthServiceServer) *GRPCServer {
	server := grpc.NewServer()
	authpb.RegisterAuthServiceServer(server, authServiceServer)

	return &GRPCServer{
		server: server,
		addr:   addr,
	}
}

// Start запускает gRPC сервер
func (s *GRPCServer) Start() error {
	listener, err := net.Listen("tcp", s.addr)
	if err != nil {
		return fmt.Errorf("failed to listen on %s: %w", s.addr, err)
	}

	fmt.Printf("Starting gRPC server on %s\n", s.addr)
	return s.server.Serve(listener)
}

// Stop останавливает gRPC сервер
func (s *GRPCServer) Stop() {
	s.server.GracefulStop()
}
