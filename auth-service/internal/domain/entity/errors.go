package entity

import (
	"errors"
	"fmt"
)

// Domain Errors
var (
	// User errors
	ErrUserNotFound      = errors.New("user not found")
	ErrUserAlreadyExists = errors.New("user already exists")
	ErrUserBlocked       = errors.New("user is blocked")
	ErrUserInactive      = errors.New("user is inactive")

	// Credential errors
	ErrInvalidEmail       = errors.New("invalid email format")
	ErrInvalidPassword    = errors.New("invalid password")
	ErrPasswordTooShort   = errors.New("password too short (min 8 characters)")
	ErrPasswordTooWeak    = errors.New("password too weak")
	ErrInvalidCredentials = errors.New("invalid email or password")

	// Token errors
	ErrTokenExpired = errors.New("token has expired")
	ErrTokenInvalid = errors.New("token is invalid")
	ErrTokenRevoked = errors.New("token has been revoked")
)

// DomainError - Custom error with context
type DomainError struct {
	Code    string
	Message string
	Err     error
}

func (e *DomainError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("[%s] %s: %v", e.Code, e.Message, e.Err)
	}
	return fmt.Sprintf("[%s] %s", e.Code, e.Message)
}

func (e *DomainError) Unwrap() error {
	return e.Err
}

// NewDomainError - Factory
func NewDomainError(code, message string, err error) *DomainError {
	return &DomainError{
		Code:    code,
		Message: message,
		Err:     err,
	}
}
