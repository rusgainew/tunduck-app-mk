package entity

import (
	"regexp"
	"time"
	"unicode"
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

	// Domain events to be published
	domainEvents []DomainEvent
}

type UserStatus string

const (
	UserStatusActive   UserStatus = "active"
	UserStatusInactive UserStatus = "inactive"
	UserStatusBlocked  UserStatus = "blocked"
)

var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)

// NewUser - Factory для создания нового пользователя
func NewUser(id, email, name, passwordHash string) (*User, error) {
	if email == "" {
		return nil, ErrInvalidEmail
	}

	if !emailRegex.MatchString(email) {
		return nil, ErrInvalidEmail
	}

	if passwordHash == "" {
		return nil, ErrInvalidPassword
	}

	user := &User{
		ID:           id,
		Email:        email,
		Name:         name,
		Password:     passwordHash,
		Status:       UserStatusActive,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
		domainEvents: make([]DomainEvent, 0),
	}

	// Add domain event
	user.AddDomainEvent(NewUserRegistered(id, email, name))

	return user, nil
}

// IsActive - Business logic
func (u *User) IsActive() bool {
	return u.Status == UserStatusActive
}

// IsBlocked - Check if user is blocked
func (u *User) IsBlocked() bool {
	return u.Status == UserStatusBlocked
}

// Activate - Activate user
func (u *User) Activate() error {
	if u.Status == UserStatusBlocked {
		return ErrUserBlocked
	}
	u.Status = UserStatusActive
	u.UpdatedAt = time.Now()
	return nil
}

// Deactivate - Deactivate user
func (u *User) Deactivate() {
	u.Status = UserStatusInactive
	u.UpdatedAt = time.Now()
}

// UpdateLastLogin - Update domain state
func (u *User) UpdateLastLogin(ipAddress string) {
	now := time.Now()
	u.LastLogin = &now
	u.UpdatedAt = now
	u.AddDomainEvent(NewUserLoggedIn(u.ID, u.Email, ipAddress))
}

// Block - Change status
func (u *User) Block(reason string) error {
	if u.Status == UserStatusBlocked {
		return ErrUserBlocked
	}

	u.Status = UserStatusBlocked
	u.UpdatedAt = time.Now()
	u.AddDomainEvent(NewUserBlocked(u.ID, reason))

	return nil
}

// ChangePassword - Change user password
func (u *User) ChangePassword(newPasswordHash string) error {
	if newPasswordHash == "" {
		return ErrInvalidPassword
	}

	u.Password = newPasswordHash
	u.UpdatedAt = time.Now()
	u.AddDomainEvent(NewPasswordChanged(u.ID))

	return nil
}

// Domain Events management
func (u *User) AddDomainEvent(event DomainEvent) {
	u.domainEvents = append(u.domainEvents, event)
}

func (u *User) DomainEvents() []DomainEvent {
	return u.domainEvents
}

func (u *User) ClearDomainEvents() {
	u.domainEvents = make([]DomainEvent, 0)
}

// Credential - Value Object для хранения учетных данных
type Credential struct {
	Email    string
	Password string // plain password for verification
}

// NewCredential - Factory with validation
func NewCredential(email, password string) (*Credential, error) {
	if email == "" {
		return nil, ErrInvalidEmail
	}

	if !emailRegex.MatchString(email) {
		return nil, ErrInvalidEmail
	}

	if password == "" {
		return nil, ErrInvalidPassword
	}

	if err := ValidatePassword(password); err != nil {
		return nil, err
	}

	return &Credential{Email: email, Password: password}, nil
}

// ValidatePassword - Check password strength
func ValidatePassword(password string) error {
	if len(password) < 8 {
		return ErrPasswordTooShort
	}

	var (
		hasUpper  bool
		hasLower  bool
		hasNumber bool
	)

	for _, char := range password {
		switch {
		case unicode.IsUpper(char):
			hasUpper = true
		case unicode.IsLower(char):
			hasLower = true
		case unicode.IsNumber(char):
			hasNumber = true
		}
	}

	if !hasUpper || !hasLower || !hasNumber {
		return ErrPasswordTooWeak
	}

	return nil
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
