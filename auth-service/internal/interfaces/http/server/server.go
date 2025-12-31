package server

import (
	"fmt"
	"net/http"

	"github.com/rusgainew/tunduck-app-mk/auth-service/internal/application/service"
	"github.com/rusgainew/tunduck-app-mk/auth-service/internal/infrastructure/config"
	"github.com/rusgainew/tunduck-app-mk/auth-service/internal/interfaces/http/handler"
	"github.com/rusgainew/tunduck-app-mk/auth-service/internal/interfaces/http/middleware"
)

// HTTPServer - HTTP server для Auth Service
type HTTPServer struct {
	cfg             *config.Config
	registerHandler *handler.RegisterHandler
	loginHandler    *handler.LoginHandler
	getMeHandler    *handler.GetMeHandler
}

// NewHTTPServer - Factory
func NewHTTPServer(
	cfg *config.Config,
	registerService *service.RegisterUserService,
	loginService *service.LoginUserService,
	tokenService *service.TokenService,
	userRepo interface{}, // будет заменено на interface{}
) *HTTPServer {
	return &HTTPServer{
		cfg:             cfg,
		registerHandler: handler.NewRegisterHandler(registerService),
		loginHandler:    handler.NewLoginHandler(loginService),
		getMeHandler:    handler.NewGetMeHandler(userRepo, tokenService),
	}
}

// Start - запустить HTTP сервер
func (s *HTTPServer) Start() error {
	mux := http.NewServeMux()

	// Регистрировать endpoints
	mux.HandleFunc("/health", (&handler.HealthHandler{}).Handle)
	mux.HandleFunc("/auth/register", s.registerHandler.Handle)
	mux.HandleFunc("/auth/login", s.loginHandler.Handle)
	mux.HandleFunc("/auth/me", s.getMeHandler.Handle)

	// Wrap mux with middleware
	var finalHandler http.Handler = mux
	finalHandler = middleware.RecoveryMiddleware(finalHandler)
	finalHandler = middleware.LoggingMiddleware(finalHandler)
	finalHandler = middleware.CORSMiddleware(finalHandler)

	addr := fmt.Sprintf(":%d", s.cfg.HttpPort)
	fmt.Printf("Starting HTTP server on %s\n", addr)

	return http.ListenAndServe(addr, finalHandler)
}
