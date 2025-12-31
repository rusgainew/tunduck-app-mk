package middleware

import (
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"
)

// LoggingMiddleware - HTTP request/response logging
func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		// Wrap response writer to capture status code
		wrapped := &responseWriter{ResponseWriter: w, statusCode: http.StatusOK}

		fmt.Printf("[%s] %s %s %s\n",
			time.Now().Format("2006-01-02 15:04:05"),
			r.Method,
			r.RequestURI,
			r.RemoteAddr,
		)

		next.ServeHTTP(wrapped, r)

		duration := time.Since(start)
		fmt.Printf("[%s] %s %s - Status: %d - Duration: %v\n",
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
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
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

// RecoveryMiddleware - Panic recovery
func RecoveryMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				log.Printf("Panic: %v", err)
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusInternalServerError)
				fmt.Fprintf(w, `{"code":"INTERNAL_ERROR","message":"Internal server error"}`)
			}
		}()

		next.ServeHTTP(w, r)
	})
}

// AuthMiddleware - Bearer token validation
type AuthMiddleware struct {
	excludedPaths map[string]bool
	validateToken func(string) (string, error)
}

// NewAuthMiddleware - Factory
func NewAuthMiddleware(validateToken func(string) (string, error)) *AuthMiddleware {
	return &AuthMiddleware{
		excludedPaths: map[string]bool{
			"/health":        true,
			"/auth/register": true,
			"/auth/login":    true,
			"/auth/refresh":  true,
		},
		validateToken: validateToken,
	}
}

// Middleware - Middleware function
func (am *AuthMiddleware) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Skip auth for excluded paths
		if am.excludedPaths[r.RequestURI] {
			next.ServeHTTP(w, r)
			return
		}

		// Extract token from Authorization header
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnauthorized)
			fmt.Fprintf(w, `{"code":"UNAUTHORIZED","message":"Missing authorization header"}`)
			return
		}

		// Parse Bearer token
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnauthorized)
			fmt.Fprintf(w, `{"code":"UNAUTHORIZED","message":"Invalid authorization header"}`)
			return
		}

		token := parts[1]

		// Validate token
		userID, err := am.validateToken(token)
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnauthorized)
			fmt.Fprintf(w, `{"code":"UNAUTHORIZED","message":"%s"}`, err.Error())
			return
		}

		// Add user ID to request context for downstream handlers
		r.Header.Set("X-User-ID", userID)

		next.ServeHTTP(w, r)
	})
}

// responseWriter - Wrapper to capture status code
type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (w *responseWriter) WriteHeader(statusCode int) {
	w.statusCode = statusCode
	w.ResponseWriter.WriteHeader(statusCode)
}
