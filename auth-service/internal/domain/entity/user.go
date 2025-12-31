package entity

import (
	"errors"
	"time"
)

// User - Aggregate Root для домена аутентификации
type User struct {
	ID        string
	Email     string
	Name      string
	Password  string // хеш
	Status    UserStatus
	CreatedAt time.Time
	UpdatedAt time.Time
	LastLogin *time.Time
}

type UserStatus string

const (
	UserStatusActive   UserStatus = "active"
	UserStatusInactive UserStatus = "inactive"
	UserStatusBlocked  UserStatus = "blocked"
)

// NewUser - Factory для создания нового пользователя
func NewUser(id, email, name, passwordHash string) (*User, error) {
	if email == "" || passwordHash == "" {
		return nil, errors.New("email and password are required")
	}

	return &User{
		ID:        id,
		Email:     email,
		Name:      name,
		Password:  passwordHash,
		Status:    UserStatusActive,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}, nil
}

// IsActive - Business logic
func (u *User) IsActive() bool {
	return u.Status == UserStatusActive
}

// UpdateLastLogin - Update domain state
func (u *User) UpdateLastLogin() {
	now := time.Now()
	u.LastLogin = &now
	u.UpdatedAt = now
}

// Block - Change status
func (u *User) Block() {
	u.Status = UserStatusBlocked
	u.UpdatedAt = time.Now()
}

// Credential - Value Object для хранения учетных данных
type Credential struct {
	Email    string
	Password string // plain password for verification
}

// NewCredential - Factory
func NewCredential(email, password string) (*Credential, error) {
	if email == "" || password == "" {
		return nil, errors.New("email and password required")
	}
	return &Credential{Email: email, Password: password}, nil
}

// Token - Value Object для JWT
type Token struct {
	AccessToken  string
	RefreshToken string
	ExpiresIn    int64 // seconds
	IssuedAt     time.Time
}

// IsExpired - Check if token expired
func (t *Token) IsExpired() bool {
	return time.Now().Unix() > (t.IssuedAt.Unix() + t.ExpiresIn)
}

// NewToken - Factory
func NewToken(accessToken, refreshToken string, expiresIn int64) *Token {
	return &Token{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    expiresIn,
		IssuedAt:     time.Now(),
	}
}

// Role - Value Object для ролей
type Role struct {
	ID   string
	Name string
}

// Permission - Value Object для прав доступа
type Permission struct {
	ID     string
	Name   string
	Action string
}
