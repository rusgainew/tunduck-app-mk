package entity

import (
	"testing"
	"time"
)

func TestNewUser(t *testing.T) {
	tests := []struct {
		name          string
		id            string
		email         string
		userName      string
		passwordHash  string
		expectedError error
	}{
		{
			name:          "valid user",
			id:            "user-123",
			email:         "test@example.com",
			userName:      "Test User",
			passwordHash:  "hashed_password",
			expectedError: nil,
		},
		{
			name:          "empty email",
			id:            "user-123",
			email:         "",
			userName:      "Test User",
			passwordHash:  "hashed_password",
			expectedError: ErrInvalidEmail,
		},
		{
			name:          "invalid email format",
			id:            "user-123",
			email:         "invalid-email",
			userName:      "Test User",
			passwordHash:  "hashed_password",
			expectedError: ErrInvalidEmail,
		},
		{
			name:          "empty password",
			id:            "user-123",
			email:         "test@example.com",
			userName:      "Test User",
			passwordHash:  "",
			expectedError: ErrInvalidPassword,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			user, err := NewUser(tt.id, tt.email, tt.userName, tt.passwordHash)

			if tt.expectedError != nil {
				if err != tt.expectedError {
					t.Errorf("expected error %v, got %v", tt.expectedError, err)
				}
				return
			}

			if err != nil {
				t.Errorf("unexpected error: %v", err)
				return
			}

			if user.ID != tt.id {
				t.Errorf("expected ID %s, got %s", tt.id, user.ID)
			}

			if user.Email != tt.email {
				t.Errorf("expected email %s, got %s", tt.email, user.Email)
			}

			if user.Status != UserStatusActive {
				t.Errorf("expected status %s, got %s", UserStatusActive, user.Status)
			}

			// Check domain event was added
			events := user.DomainEvents()
			if len(events) != 1 {
				t.Errorf("expected 1 domain event, got %d", len(events))
			}

			if events[0].EventName() != "user.registered" {
				t.Errorf("expected event 'user.registered', got %s", events[0].EventName())
			}
		})
	}
}

func TestUserIsActive(t *testing.T) {
	user, _ := NewUser("user-123", "test@example.com", "Test", "hash")

	if !user.IsActive() {
		t.Error("new user should be active")
	}

	user.Status = UserStatusInactive
	if user.IsActive() {
		t.Error("inactive user should not be active")
	}
}

func TestUserIsBlocked(t *testing.T) {
	user, _ := NewUser("user-123", "test@example.com", "Test", "hash")

	if user.IsBlocked() {
		t.Error("new user should not be blocked")
	}

	user.Status = UserStatusBlocked
	if !user.IsBlocked() {
		t.Error("blocked user should be blocked")
	}
}

func TestUserActivate(t *testing.T) {
	user, _ := NewUser("user-123", "test@example.com", "Test", "hash")
	user.Status = UserStatusInactive

	err := user.Activate()
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if user.Status != UserStatusActive {
		t.Errorf("expected status %s, got %s", UserStatusActive, user.Status)
	}

	// Blocked user cannot be activated
	user.Status = UserStatusBlocked
	err = user.Activate()
	if err != ErrUserBlocked {
		t.Errorf("expected error %v, got %v", ErrUserBlocked, err)
	}
}

func TestUserDeactivate(t *testing.T) {
	user, _ := NewUser("user-123", "test@example.com", "Test", "hash")

	user.Deactivate()

	if user.Status != UserStatusInactive {
		t.Errorf("expected status %s, got %s", UserStatusInactive, user.Status)
	}
}

func TestUserBlock(t *testing.T) {
	user, _ := NewUser("user-123", "test@example.com", "Test", "hash")
	user.ClearDomainEvents() // Clear registration event

	err := user.Block("suspicious activity")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if user.Status != UserStatusBlocked {
		t.Errorf("expected status %s, got %s", UserStatusBlocked, user.Status)
	}

	// Check domain event
	events := user.DomainEvents()
	if len(events) != 1 {
		t.Errorf("expected 1 domain event, got %d", len(events))
	}

	if events[0].EventName() != "user.blocked" {
		t.Errorf("expected event 'user.blocked', got %s", events[0].EventName())
	}

	// Cannot block already blocked user
	err = user.Block("another reason")
	if err != ErrUserBlocked {
		t.Errorf("expected error %v, got %v", ErrUserBlocked, err)
	}
}

func TestUserUpdateLastLogin(t *testing.T) {
	user, _ := NewUser("user-123", "test@example.com", "Test", "hash")
	user.ClearDomainEvents()

	ipAddress := "192.168.1.1"
	user.UpdateLastLogin(ipAddress)

	if user.LastLogin == nil {
		t.Error("LastLogin should be set")
	}

	if time.Since(*user.LastLogin) > time.Second {
		t.Error("LastLogin should be recent")
	}

	// Check domain event
	events := user.DomainEvents()
	if len(events) != 1 {
		t.Errorf("expected 1 domain event, got %d", len(events))
	}

	if events[0].EventName() != "user.logged_in" {
		t.Errorf("expected event 'user.logged_in', got %s", events[0].EventName())
	}
}

func TestUserChangePassword(t *testing.T) {
	user, _ := NewUser("user-123", "test@example.com", "Test", "hash")
	user.ClearDomainEvents()

	newHash := "new_hashed_password"
	err := user.ChangePassword(newHash)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if user.Password != newHash {
		t.Errorf("expected password %s, got %s", newHash, user.Password)
	}

	// Check domain event
	events := user.DomainEvents()
	if len(events) != 1 {
		t.Errorf("expected 1 domain event, got %d", len(events))
	}

	if events[0].EventName() != "user.password_changed" {
		t.Errorf("expected event 'user.password_changed', got %s", events[0].EventName())
	}

	// Empty password should fail
	err = user.ChangePassword("")
	if err != ErrInvalidPassword {
		t.Errorf("expected error %v, got %v", ErrInvalidPassword, err)
	}
}

func TestUserDomainEvents(t *testing.T) {
	user, _ := NewUser("user-123", "test@example.com", "Test", "hash")

	events := user.DomainEvents()
	if len(events) != 1 {
		t.Errorf("expected 1 event, got %d", len(events))
	}

	user.ClearDomainEvents()
	events = user.DomainEvents()
	if len(events) != 0 {
		t.Errorf("expected 0 events after clear, got %d", len(events))
	}
}

func TestNewCredential(t *testing.T) {
	tests := []struct {
		name          string
		email         string
		password      string
		expectedError error
	}{
		{
			name:          "valid credential",
			email:         "test@example.com",
			password:      "ValidPass123!",
			expectedError: nil,
		},
		{
			name:          "empty email",
			email:         "",
			password:      "ValidPass123!",
			expectedError: ErrInvalidEmail,
		},
		{
			name:          "invalid email",
			email:         "invalid-email",
			password:      "ValidPass123!",
			expectedError: ErrInvalidEmail,
		},
		{
			name:          "empty password",
			email:         "test@example.com",
			password:      "",
			expectedError: ErrInvalidPassword,
		},
		{
			name:          "password too short",
			email:         "test@example.com",
			password:      "Pass1!",
			expectedError: ErrPasswordTooShort,
		},
		{
			name:          "password too weak (no uppercase)",
			email:         "test@example.com",
			password:      "password123!",
			expectedError: ErrPasswordTooWeak,
		},
		{
			name:          "password too weak (no lowercase)",
			email:         "test@example.com",
			password:      "PASSWORD123!",
			expectedError: ErrPasswordTooWeak,
		},
		{
			name:          "password too weak (no number)",
			email:         "test@example.com",
			password:      "Password!",
			expectedError: ErrPasswordTooWeak,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cred, err := NewCredential(tt.email, tt.password)

			if tt.expectedError != nil {
				if err != tt.expectedError {
					t.Errorf("expected error %v, got %v", tt.expectedError, err)
				}
				return
			}

			if err != nil {
				t.Errorf("unexpected error: %v", err)
				return
			}

			if cred.Email != tt.email {
				t.Errorf("expected email %s, got %s", tt.email, cred.Email)
			}

			if cred.Password != tt.password {
				t.Errorf("expected password %s, got %s", tt.password, cred.Password)
			}
		})
	}
}

func TestValidatePassword(t *testing.T) {
	tests := []struct {
		name          string
		password      string
		expectedError error
	}{
		{
			name:          "valid password",
			password:      "ValidPass123!",
			expectedError: nil,
		},
		{
			name:          "valid password without special char",
			password:      "ValidPass123",
			expectedError: nil,
		},
		{
			name:          "too short",
			password:      "Pass1!",
			expectedError: ErrPasswordTooShort,
		},
		{
			name:          "no uppercase",
			password:      "password123",
			expectedError: ErrPasswordTooWeak,
		},
		{
			name:          "no lowercase",
			password:      "PASSWORD123",
			expectedError: ErrPasswordTooWeak,
		},
		{
			name:          "no number",
			password:      "Password",
			expectedError: ErrPasswordTooWeak,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidatePassword(tt.password)
			if err != tt.expectedError {
				t.Errorf("expected error %v, got %v", tt.expectedError, err)
			}
		})
	}
}

func TestTokenIsExpired(t *testing.T) {
	// Non-expired token
	token := NewToken("access", "refresh", 3600)
	if token.IsExpired() {
		t.Error("token should not be expired")
	}

	// Expired token
	expiredToken := &Token{
		AccessToken:  "access",
		RefreshToken: "refresh",
		ExpiresIn:    1,
		IssuedAt:     time.Now().Add(-2 * time.Second),
	}
	if !expiredToken.IsExpired() {
		t.Error("token should be expired")
	}
}
