package handler

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/rusgainew/tunduck-app-mk/auth-service/internal/application/dto"
	"github.com/rusgainew/tunduck-app-mk/auth-service/internal/application/service"
	"github.com/rusgainew/tunduck-app-mk/auth-service/internal/domain/repository"
)

// RegisterHandler - HTTP handler для регистрации
type RegisterHandler struct {
	registerService *service.RegisterUserService
}

// NewRegisterHandler - Factory
func NewRegisterHandler(registerService *service.RegisterUserService) *RegisterHandler {
	return &RegisterHandler{
		registerService: registerService,
	}
}

// Handle - Handle POST /auth/register
func (h *RegisterHandler) Handle(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req dto.RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	resp, err := h.registerService.Execute(r.Context(), &req)
	if err != nil {
		errResp := &dto.ErrorResponse{
			Code:    "REGISTRATION_FAILED",
			Message: err.Error(),
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(errResp)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(resp)
}

// LoginHandler - HTTP handler для входа
type LoginHandler struct {
	loginService *service.LoginUserService
}

// NewLoginHandler - Factory
func NewLoginHandler(loginService *service.LoginUserService) *LoginHandler {
	return &LoginHandler{
		loginService: loginService,
	}
}

// Handle - Handle POST /auth/login
func (h *LoginHandler) Handle(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req dto.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	resp, err := h.loginService.Execute(r.Context(), &req)
	if err != nil {
		errResp := &dto.ErrorResponse{
			Code:    "LOGIN_FAILED",
			Message: err.Error(),
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(errResp)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)
}

// GetMeHandler - HTTP handler для получения профиля
type GetMeHandler struct {
	userRepo     repository.UserRepository
	tokenService *service.TokenService
}

// NewGetMeHandler - Factory
func NewGetMeHandler(
	userRepo repository.UserRepository,
	tokenService *service.TokenService,
) *GetMeHandler {
	return &GetMeHandler{
		userRepo:     userRepo,
		tokenService: tokenService,
	}
}

// Handle - Handle GET /auth/me
func (h *GetMeHandler) Handle(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Извлечь токен из Authorization заголовка
	token := extractToken(r)
	if token == "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Валидировать токен
	userID, err := h.tokenService.ValidateToken(r.Context(), token)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Получить пользователя
	user, err := h.userRepo.GetUserByID(r.Context(), userID)
	if err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(dto.UserToGetMeResponse(user))
}

// HealthHandler - HTTP handler для health check
type HealthHandler struct{}

// Handle - Handle GET /health
func (h *HealthHandler) Handle(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, `{"status":"ok"}`)
}

// Вспомогательные функции
func extractToken(r *http.Request) string {
	auth := r.Header.Get("Authorization")
	if len(auth) > 7 && auth[:7] == "Bearer " {
		return auth[7:]
	}
	return ""
}
