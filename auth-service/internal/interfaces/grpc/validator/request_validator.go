package validator

import (
	authpb "github.com/rusgainew/tunduck-app-mk/proto-lib/auth"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// RequestValidator отвечает за валидацию входящих запросов
type RequestValidator struct{}

// NewRequestValidator создает новый валидатор
func NewRequestValidator() *RequestValidator {
	return &RequestValidator{}
}

// ValidateRegisterRequest валидирует запрос регистрации
func (v *RequestValidator) ValidateRegisterRequest(req *authpb.RegisterRequest) error {
	if req == nil {
		return status.Error(codes.InvalidArgument, "request cannot be nil")
	}

	if req.Email == "" {
		return status.Error(codes.InvalidArgument, "email is required")
	}

	if req.Password == "" {
		return status.Error(codes.InvalidArgument, "password is required")
	}

	return nil
}

// ValidateLoginRequest валидирует запрос входа
func (v *RequestValidator) ValidateLoginRequest(req *authpb.LoginRequest) error {
	if req == nil {
		return status.Error(codes.InvalidArgument, "request cannot be nil")
	}

	if req.Email == "" {
		return status.Error(codes.InvalidArgument, "email is required")
	}

	if req.Password == "" {
		return status.Error(codes.InvalidArgument, "password is required")
	}

	return nil
}

// ValidateTokenRequest валидирует запрос с токеном
func (v *RequestValidator) ValidateTokenRequest(req *authpb.ValidateTokenRequest) error {
	if req == nil || req.AccessToken == "" {
		return status.Error(codes.InvalidArgument, "access_token is required")
	}
	return nil
}

// ValidateGetUserRequest валидирует запрос получения пользователя
func (v *RequestValidator) ValidateGetUserRequest(req *authpb.GetUserRequest) error {
	if req == nil || req.UserId == "" {
		return status.Error(codes.InvalidArgument, "user_id is required")
	}
	return nil
}

// ValidateLogoutRequest валидирует запрос выхода
func (v *RequestValidator) ValidateLogoutRequest(req *authpb.LogoutRequest) error {
	if req == nil || req.AccessToken == "" {
		return status.Error(codes.InvalidArgument, "access_token is required")
	}
	return nil
}

// ValidateRefreshTokenRequest валидирует запрос обновления токена
func (v *RequestValidator) ValidateRefreshTokenRequest(req *authpb.RefreshTokenRequest) error {
	if req == nil || req.RefreshToken == "" {
		return status.Error(codes.InvalidArgument, "refresh_token is required")
	}
	return nil
}
