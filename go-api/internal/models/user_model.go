package models

import "fmt"

// RegisterRequest запрос на регистрацию обычного пользователя
type RegisterRequest struct {
	Email           string `json:"email" validate:"required,email"`
	FullName        string `json:"fullName" validate:"required,min=2,max=100"`
	Password        string `json:"password" validate:"required,min=6,max=100"`
	ConfirmPassword string `json:"confirmPassword" validate:"required,eqfield=Password"`
}

// Validate проверяет валидность данных регистрации
func (r *RegisterRequest) Validate() error {
	if r.Email == "" {
		return fmt.Errorf("email is required")
	}
	if r.FullName == "" {
		return fmt.Errorf("full name is required")
	}
	if r.Password == "" {
		return fmt.Errorf("password is required")
	}
	if len(r.Password) < 6 {
		return fmt.Errorf("password must be at least 6 characters")
	}
	if r.Password != r.ConfirmPassword {
		return fmt.Errorf("passwords do not match")
	}
	return nil
}

// AdminRegisterRequest запрос на регистрацию администратора
type AdminRegisterRequest struct {
	Username        string `json:"username" validate:"required,min=3,max=50"`
	Email           string `json:"email" validate:"required,email"`
	FullName        string `json:"fullName" validate:"required,min=2,max=100"`
	Phone           string `json:"phone" validate:"required,min=10,max=20"`
	Password        string `json:"password" validate:"required,min=6,max=100"`
	ConfirmPassword string `json:"confirmPassword" validate:"required,eqfield=Password"`
	AdminSecret     string `json:"adminSecret" validate:"required"` // Обязательное поле для регистрации админа
}

// LoginRequest запрос на вход
type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"` // Changed from username to email for auth-service
	Password string `json:"password" validate:"required"`
}

// AuthResponse ответ при успешной аутентификации
type AuthResponse struct {
	Token string    `json:"token"`
	User  *UserInfo `json:"user"`
}

// UserInfo информация о пользователе
type UserInfo struct {
	ID       string `json:"id"` // Changed from uuid.UUID to string for auth-service compatibility
	Username string `json:"username,omitempty"`
	Email    string `json:"email"`
	FullName string `json:"fullName"`
	Phone    string `json:"phone,omitempty"`
	Role     string `json:"role,omitempty"`
	IsActive bool   `json:"isActive,omitempty"`
}

// TokenResponse ответ с токенами
type TokenResponse struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken,omitempty"`
	ExpiresIn    int    `json:"expiresIn"`
	TokenType    string `json:"tokenType"`
}

// ErrorResponse ответ с ошибкой
type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message"`
}
