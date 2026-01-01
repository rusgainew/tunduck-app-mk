package mapper

import (
	"time"

	"github.com/rusgainew/tunduck-app-mk/auth-service/internal/domain/entity"
	authpb "github.com/rusgainew/tunduck-app-mk/proto-lib/auth"
)

// AuthResponseMapper отвечает за создание AuthResponse
type AuthResponseMapper struct {
	userMapper  *UserMapper
	tokenMapper *TokenMapper
}

// NewAuthResponseMapper создает новый маппер для AuthResponse
func NewAuthResponseMapper(userMapper *UserMapper, tokenMapper *TokenMapper) *AuthResponseMapper {
	return &AuthResponseMapper{
		userMapper:  userMapper,
		tokenMapper: tokenMapper,
	}
}

// ToAuthResponseWithToken создает AuthResponse с пользователем и токеном
func (m *AuthResponseMapper) ToAuthResponseWithToken(
	user *entity.User,
	token *authpb.Token,
) *authpb.AuthResponse {
	return &authpb.AuthResponse{
		User:      m.userMapper.ToProtoUser(user),
		Token:     token,
		Timestamp: time.Now().Unix(),
	}
}

// ToAuthResponseWithNames создает AuthResponse с разделенными именами
func (m *AuthResponseMapper) ToAuthResponseWithNames(
	user *entity.User,
	firstName, lastName string,
	token *authpb.Token,
) *authpb.AuthResponse {
	return &authpb.AuthResponse{
		User:      m.userMapper.ToProtoUserWithNames(user, firstName, lastName),
		Token:     token,
		Timestamp: time.Now().Unix(),
	}
}
