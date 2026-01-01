package middleware

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/rusgainew/tunduck-app-mk/auth-service/internal/application/service"
)

// responseWriter - обертка для response writer для захвата status code
type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

// LoggingMiddleware - HTTP request/response logging
func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		wrapped := &responseWriter{ResponseWriter: w, statusCode: http.StatusOK}

		log.Printf("[%s] %s %s from %s\n",
			time.Now().Format("2006-01-02 15:04:05"),
			r.Method,
			r.RequestURI,
			r.RemoteAddr,
		)

		next.ServeHTTP(wrapped, r)

		duration := time.Since(start)
		log.Printf("[%s] %s %s - Status: %d - Duration: %v\n",
			time.Now().Format("2006-01-02 15:04:05"),
			r.Method,
			r.RequestURI,
			wrapped.statusCode,
			duration,
		)
	})
}

// CORSMiddleware - CORS headers
func CORSMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS, PATCH")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		w.Header().Set("Access-Control-Max-Age", "3600")

		// Handle preflight requests
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

// AuthMiddleware - JWT authentication middleware
type AuthMiddleware struct {
	validateService *service.ValidateTokenService
	skipPaths       []string // Пути, которые не требуют аутентификации
}

// NewAuthMiddleware - Factory
func NewAuthMiddleware(validateService *service.ValidateTokenService) *AuthMiddleware {
	return &AuthMiddleware{
		validateService: validateService,
		skipPaths: []string{
			"/api/v1/auth/register",
			"/api/v1/auth/login",
			"/health",
		},
	}
}

// Handler - Middleware handler
func (m *AuthMiddleware) Handler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Пропустить проверку для определенных путей
		if m.shouldSkip(r.URL.Path) {
			next.ServeHTTP(w, r)
			return
		}

		// Извлечь токен
		token := m.extractToken(r)
		if token == "" {
			http.Error(w, `{"code":"UNAUTHORIZED","message":"Missing token"}`, http.StatusUnauthorized)
			return
		}

		// Валидировать токен
		userID, err := m.validateService.Execute(r.Context(), token)
		if err != nil {
			http.Error(w, `{"code":"UNAUTHORIZED","message":"Invalid token"}`, http.StatusUnauthorized)
			return
		}

		// Добавить userID и token в контекст
		ctx := context.WithValue(r.Context(), "userID", userID)
		ctx = context.WithValue(ctx, "token", token)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// shouldSkip - проверить, нужно ли пропустить проверку аутентификации
func (m *AuthMiddleware) shouldSkip(path string) bool {
	for _, skipPath := range m.skipPaths {
		if path == skipPath {
			return true
		}
	}
	return false
}

// extractToken - извлечь JWT токен из Authorization header
func (m *AuthMiddleware) extractToken(r *http.Request) string {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return ""
	}

	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || parts[0] != "Bearer" {
		return ""
	}

	return parts[1]
}

// ErrorHandlerMiddleware - обработка ошибок
func ErrorHandlerMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Установить content-type по умолчанию
		w.Header().Set("Content-Type", "application/json")

		next.ServeHTTP(w, r)
	})
}

// RecoveryMiddleware - обработка panic
func RecoveryMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				log.Printf("PANIC: %v\n", err)
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusInternalServerError)
				fmt.Fprintf(w, `{"code":"INTERNAL_ERROR","message":"Internal server error"}`)
			}
		}()

		next.ServeHTTP(w, r)
	})
}

// AuthMiddlewareFunc - функция-обертка для создания AuthMiddleware handler
func AuthMiddlewareFunc(next http.Handler, validateService *service.ValidateTokenService) http.Handler {
	m := NewAuthMiddleware(validateService)
	return m.Handler(next)
}
