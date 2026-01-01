package http

import (
	"net/http"

	"github.com/rusgainew/tunduck-app-mk/auth-service/internal/application/service"
	"github.com/rusgainew/tunduck-app-mk/auth-service/internal/interfaces/http/handler"
	"github.com/rusgainew/tunduck-app-mk/auth-service/internal/interfaces/http/middleware"
)

// Router - структура роутера
type Router struct {
	handler http.Handler
}

// NewRouter - создает новый роутер с цепочкой middleware
func NewRouter(
	registerService *service.RegisterUserService,
	loginService *service.LoginUserService,
	validateService *service.ValidateTokenService,
	getUserService *service.GetUserService,
	logoutService *service.LogoutUserService,
	refreshTokenService *service.RefreshTokenService,
) *Router {
	mux := http.NewServeMux()

	// Инициализировать handlers
	authHandler := handler.NewAuthHandler(
		registerService,
		loginService,
		validateService,
		getUserService,
		logoutService,
		refreshTokenService,
	)

	// Регистрировать endpoints
	mux.HandleFunc("/health", authHandler.Health)
	mux.HandleFunc("/api/v1/auth/register", authHandler.Register)
	mux.HandleFunc("/api/v1/auth/login", authHandler.Login)
	mux.HandleFunc("/api/v1/auth/me", authHandler.GetMe)
	mux.HandleFunc("/api/v1/auth/logout", authHandler.Logout)
	mux.HandleFunc("/api/v1/auth/refresh", authHandler.Refresh)

	// Применить middleware ко всему роутеру в порядке (от последнего применяемого к первому)
	// Порядок выполнения: Recovery -> Error -> Logging -> CORS -> Auth
	var h http.Handler = mux
	h = middleware.AuthMiddlewareFunc(h, validateService)
	h = middleware.CORSMiddleware(h)
	h = middleware.LoggingMiddleware(h)
	h = middleware.ErrorHandlerMiddleware(h)
	h = middleware.RecoveryMiddleware(h)

	return &Router{
		handler: h,
	}
}

// ServeHTTP - имплементирует http.Handler interface
func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	r.handler.ServeHTTP(w, req)
}
