package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/rusgainew/tunduck-app-mk/auth-service/internal/application/dto"
	"github.com/rusgainew/tunduck-app-mk/auth-service/internal/application/service"
	"github.com/rusgainew/tunduck-app-mk/auth-service/internal/domain/entity"
)

// AuthHandler - контейнер для всех auth обработчиков
type AuthHandler struct {
	registerService     *service.RegisterUserService
	loginService        *service.LoginUserService
	validateService     *service.ValidateTokenService
	getUserService      *service.GetUserService
	logoutService       *service.LogoutUserService
	refreshTokenService *service.RefreshTokenService
}

// NewAuthHandler - Factory
func NewAuthHandler(
	registerService *service.RegisterUserService,
	loginService *service.LoginUserService,
	validateService *service.ValidateTokenService,
	getUserService *service.GetUserService,
	logoutService *service.LogoutUserService,
	refreshTokenService *service.RefreshTokenService,
) *AuthHandler {
	return &AuthHandler{
		registerService:     registerService,
		loginService:        loginService,
		validateService:     validateService,
		getUserService:      getUserService,
		logoutService:       logoutService,
		refreshTokenService: refreshTokenService,
	}
}

// Register - POST /api/v1/auth/register
func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeErrorResponse(w, http.StatusMethodNotAllowed, "METHOD_NOT_ALLOWED", "Only POST requests are allowed")
		return
	}

	var req dto.RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeErrorResponse(w, http.StatusBadRequest, "INVALID_REQUEST", "Invalid request format")
		return
	}

	resp, err := h.registerService.Execute(r.Context(), &req)
	if err != nil {
		handleServiceError(w, err, "registration")
		return
	}

	writeJSONResponse(w, http.StatusCreated, resp)
}

// Login - POST /api/v1/auth/login
func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeErrorResponse(w, http.StatusMethodNotAllowed, "METHOD_NOT_ALLOWED", "Only POST requests are allowed")
		return
	}

	var req dto.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeErrorResponse(w, http.StatusBadRequest, "INVALID_REQUEST", "Invalid request format")
		return
	}

	// Получить IP адрес клиента
	ipAddress := getClientIP(r)

	resp, err := h.loginService.Execute(r.Context(), &req, ipAddress)
	if err != nil {
		handleServiceError(w, err, "login")
		return
	}

	writeJSONResponse(w, http.StatusOK, resp)
}

// GetMe - GET /api/v1/auth/me
func (h *AuthHandler) GetMe(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeErrorResponse(w, http.StatusMethodNotAllowed, "METHOD_NOT_ALLOWED", "Only GET requests are allowed")
		return
	}

	// Извлечь userID из контекста (установлено middleware)
	userID, ok := r.Context().Value("userID").(string)
	if !ok {
		writeErrorResponse(w, http.StatusUnauthorized, "UNAUTHORIZED", "Missing user ID")
		return
	}

	resp, err := h.getUserService.Execute(r.Context(), userID)
	if err != nil {
		handleServiceError(w, err, "get user")
		return
	}

	writeJSONResponse(w, http.StatusOK, resp)
}

// Logout - POST /api/v1/auth/logout
func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeErrorResponse(w, http.StatusMethodNotAllowed, "METHOD_NOT_ALLOWED", "Only POST requests are allowed")
		return
	}

	userID, ok := r.Context().Value("userID").(string)
	if !ok {
		writeErrorResponse(w, http.StatusUnauthorized, "UNAUTHORIZED", "Missing user ID")
		return
	}

	token, ok := r.Context().Value("token").(string)
	if !ok {
		writeErrorResponse(w, http.StatusUnauthorized, "UNAUTHORIZED", "Missing token")
		return
	}

	err := h.logoutService.Execute(r.Context(), userID, token)
	if err != nil {
		handleServiceError(w, err, "logout")
		return
	}

	writeJSONResponse(w, http.StatusOK, map[string]string{
		"message": "Successfully logged out",
	})
}

// Refresh - POST /api/v1/auth/refresh
func (h *AuthHandler) Refresh(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeErrorResponse(w, http.StatusMethodNotAllowed, "METHOD_NOT_ALLOWED", "Only POST requests are allowed")
		return
	}

	var req dto.RefreshTokenRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeErrorResponse(w, http.StatusBadRequest, "INVALID_REQUEST", "Invalid request format")
		return
	}

	resp, err := h.refreshTokenService.Execute(r.Context(), req.RefreshToken)
	if err != nil {
		handleServiceError(w, err, "token refresh")
		return
	}

	writeJSONResponse(w, http.StatusOK, resp)
}

// Health - GET /health
func (h *AuthHandler) Health(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeErrorResponse(w, http.StatusMethodNotAllowed, "METHOD_NOT_ALLOWED", "Only GET requests are allowed")
		return
	}

	writeJSONResponse(w, http.StatusOK, map[string]string{
		"status":  "ok",
		"service": "auth-service",
	})
}

// ============= Вспомогательные функции =============

// writeJSONResponse - записывает JSON ответ
func writeJSONResponse(w http.ResponseWriter, statusCode int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(data)
}

// writeErrorResponse - записывает ошибку в JSON формате
func writeErrorResponse(w http.ResponseWriter, statusCode int, code, message string) {
	errResp := &dto.ErrorResponse{
		Code:    code,
		Message: message,
	}
	writeJSONResponse(w, statusCode, errResp)
}

// handleServiceError - обрабатывает ошибки от сервисов
func handleServiceError(w http.ResponseWriter, err error, operation string) {
	switch err {
	case entity.ErrUserNotFound:
		writeErrorResponse(w, http.StatusNotFound, "USER_NOT_FOUND", "User not found")
	case entity.ErrUserAlreadyExists:
		writeErrorResponse(w, http.StatusConflict, "USER_EXISTS", "User with this email already exists")
	case entity.ErrUserBlocked:
		writeErrorResponse(w, http.StatusForbidden, "USER_BLOCKED", "User account is blocked")
	case entity.ErrUserInactive:
		writeErrorResponse(w, http.StatusForbidden, "USER_INACTIVE", "User account is not active")
	case entity.ErrInvalidCredentials:
		writeErrorResponse(w, http.StatusUnauthorized, "INVALID_CREDENTIALS", "Invalid email or password")
	case entity.ErrInvalidEmail:
		writeErrorResponse(w, http.StatusBadRequest, "INVALID_EMAIL", "Invalid email format")
	case entity.ErrInvalidPassword:
		writeErrorResponse(w, http.StatusBadRequest, "INVALID_PASSWORD", "Invalid password")
	case entity.ErrPasswordTooShort:
		writeErrorResponse(w, http.StatusBadRequest, "PASSWORD_TOO_SHORT", "Password must be at least 8 characters")
	case entity.ErrPasswordTooWeak:
		writeErrorResponse(w, http.StatusBadRequest, "PASSWORD_WEAK", "Password must contain uppercase, lowercase, and numbers")
	case entity.ErrTokenExpired:
		writeErrorResponse(w, http.StatusUnauthorized, "TOKEN_EXPIRED", "Token has expired")
	case entity.ErrTokenInvalid:
		writeErrorResponse(w, http.StatusUnauthorized, "TOKEN_INVALID", "Token is invalid")
	case entity.ErrTokenRevoked:
		writeErrorResponse(w, http.StatusUnauthorized, "TOKEN_REVOKED", "Token has been revoked")
	default:
		writeErrorResponse(w, http.StatusInternalServerError, "INTERNAL_ERROR", fmt.Sprintf("Failed to complete %s: %v", operation, err))
	}
}

// getClientIP - получает IP адрес клиента
func getClientIP(r *http.Request) string {
	// Проверить X-Forwarded-For header (для proxy)
	if xForwardedFor := r.Header.Get("X-Forwarded-For"); xForwardedFor != "" {
		ips := strings.Split(xForwardedFor, ",")
		if len(ips) > 0 {
			return strings.TrimSpace(ips[0])
		}
	}

	// Проверить X-Real-IP header
	if xRealIP := r.Header.Get("X-Real-IP"); xRealIP != "" {
		return xRealIP
	}

	// Использовать RemoteAddr
	if colonIdx := strings.LastIndex(r.RemoteAddr, ":"); colonIdx != -1 {
		return r.RemoteAddr[:colonIdx]
	}

	return r.RemoteAddr
}
