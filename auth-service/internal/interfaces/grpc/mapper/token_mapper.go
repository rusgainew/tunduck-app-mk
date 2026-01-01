package mapper

import (
	authpb "github.com/rusgainew/tunduck-app-mk/proto-lib/auth"
)

// TokenMapper отвечает за преобразование токенов между слоями
type TokenMapper struct{}

// NewTokenMapper создает новый маппер токенов
func NewTokenMapper() *TokenMapper {
	return &TokenMapper{}
}

// ToProtoTokenFromLogin преобразует LoginResponse в proto Token
func (m *TokenMapper) ToProtoTokenFromLogin(accessToken, refreshToken string, expiresIn int64) *authpb.Token {
	return &authpb.Token{
		AccessToken:      accessToken,
		RefreshToken:     refreshToken,
		ExpiresIn:        expiresIn,
		TokenType:        "Bearer",
		RefreshExpiresIn: 0,
	}
}

// ToEmptyToken возвращает пустой токен
func (m *TokenMapper) ToEmptyToken() *authpb.Token {
	return &authpb.Token{}
}
