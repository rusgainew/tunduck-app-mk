package dto

import (
	"github.com/rusgainew/tunduck-app-mk/auth-service/internal/domain/entity"
)

// RegisterRequest - Input DTO для регистрации
type RegisterRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Name     string `json:"name" validate:"required,min=2"`
	Password string `json:"password" validate:"required,min=8"`
}

// RegisterResponse - Output DTO для регистрации
type RegisterResponse struct {
	ID        string `json:"id"`
	Email     string `json:"email"`
	Name      string `json:"name"`
	Status    string `json:"status"`
	CreatedAt string `json:"created_at"`
}

func UserToRegisterResponse(user *entity.User) *RegisterResponse {
	return &RegisterResponse{
		ID:        user.ID,
		Email:     user.Email,
		Name:      user.Name,
		Status:    string(user.Status),
		CreatedAt: user.CreatedAt.Format("2006-01-02T15:04:05Z"),
	}
}

// LoginRequest - Input DTO для входа
type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

// LoginResponse - Output DTO для входа
type LoginResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int64  `json:"expires_in"`
	TokenType    string `json:"token_type"`
}

func TokenToLoginResponse(token *entity.Token) *LoginResponse {
	return &LoginResponse{
		AccessToken:  token.AccessToken,
		RefreshToken: token.RefreshToken,
		ExpiresIn:    token.ExpiresIn,
		TokenType:    "Bearer",
	}
}

// GetMeResponse - Output DTO для получения профиля
type GetMeResponse struct {
	ID        string  `json:"id"`
	Email     string  `json:"email"`
	Name      string  `json:"name"`
	Status    string  `json:"status"`
	CreatedAt string  `json:"created_at"`
	UpdatedAt string  `json:"updated_at"`
	LastLogin *string `json:"last_login,omitempty"`
}

func UserToGetMeResponse(user *entity.User) *GetMeResponse {
	resp := &GetMeResponse{
		ID:        user.ID,
		Email:     user.Email,
		Name:      user.Name,
		Status:    string(user.Status),
		CreatedAt: user.CreatedAt.Format("2006-01-02T15:04:05Z"),
		UpdatedAt: user.UpdatedAt.Format("2006-01-02T15:04:05Z"),
	}
	if user.LastLogin != nil {
		lastLogin := user.LastLogin.Format("2006-01-02T15:04:05Z")
		resp.LastLogin = &lastLogin
	}
	return resp
}

// RefreshTokenRequest - Input DTO для refresh token
type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" validate:"required"`
}

// ErrorResponse - Output DTO для ошибок
type ErrorResponse struct {
	Code    string                 `json:"code"`
	Message string                 `json:"message"`
	Details map[string]interface{} `json:"details,omitempty"`
}
