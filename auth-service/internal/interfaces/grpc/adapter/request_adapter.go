package adapter

import (
	"fmt"

	"github.com/rusgainew/tunduck-app-mk/auth-service/internal/application/dto"
	authpb "github.com/rusgainew/tunduck-app-mk/proto-lib/auth"
)

// RequestAdapter отвечает за адаптацию proto запросов к DTO приложения
type RequestAdapter struct{}

// NewRequestAdapter создает новый адаптер запросов
func NewRequestAdapter() *RequestAdapter {
	return &RequestAdapter{}
}

// ToRegisterDTO преобразует proto RegisterRequest в DTO
func (a *RequestAdapter) ToRegisterDTO(req *authpb.RegisterRequest) *dto.RegisterRequest {
	fullName := fmt.Sprintf("%s %s", req.FirstName, req.LastName)
	if req.FirstName == "" && req.LastName == "" {
		fullName = req.Email // fallback
	}

	return &dto.RegisterRequest{
		Email:    req.Email,
		Name:     fullName,
		Password: req.Password,
	}
}

// ToLoginDTO преобразует proto LoginRequest в DTO
func (a *RequestAdapter) ToLoginDTO(req *authpb.LoginRequest) *dto.LoginRequest {
	return &dto.LoginRequest{
		Email:    req.Email,
		Password: req.Password,
	}
}
