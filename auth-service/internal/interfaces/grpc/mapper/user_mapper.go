package mapper

import (
	"github.com/rusgainew/tunduck-app-mk/auth-service/internal/domain/entity"
	authpb "github.com/rusgainew/tunduck-app-mk/proto-lib/auth"
)

// UserMapper отвечает за преобразование между доменными моделями и proto сообщениями
type UserMapper struct{}

// NewUserMapper создает новый маппер пользователей
func NewUserMapper() *UserMapper {
	return &UserMapper{}
}

// ToProtoUser преобразует доменную модель User в proto User
func (m *UserMapper) ToProtoUser(user *entity.User) *authpb.User {
	if user == nil {
		return nil
	}

	return &authpb.User{
		Id:        user.ID,
		Email:     user.Email,
		FirstName: user.Name,
		Status:    string(user.Status),
		CreatedAt: user.CreatedAt.Unix(),
		UpdatedAt: user.UpdatedAt.Unix(),
	}
}

// ToProtoUserWithNames преобразует доменную модель с разделением имени
func (m *UserMapper) ToProtoUserWithNames(user *entity.User, firstName, lastName string) *authpb.User {
	if user == nil {
		return nil
	}

	return &authpb.User{
		Id:        user.ID,
		Email:     user.Email,
		FirstName: firstName,
		LastName:  lastName,
		Status:    string(user.Status),
		CreatedAt: user.CreatedAt.Unix(),
		UpdatedAt: user.UpdatedAt.Unix(),
	}
}
