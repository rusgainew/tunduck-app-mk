package entity

import (
	"testing"
	"time"
)

// TestUserCreation tests the User aggregate factory
func TestUserCreation(t *testing.T) {
	user, err := NewUser("1", "test@example.com", "Test User", "hashedpassword")

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if user.Email != "test@example.com" {
		t.Errorf("Expected email test@example.com, got %s", user.Email)
	}

	if user.Status != UserStatusActive {
		t.Errorf("Expected status %s, got %s", UserStatusActive, user.Status)
	}
}

// TestUserCreation_InvalidEmail tests error handling
func TestUserCreation_InvalidEmail(t *testing.T) {
	_, err := NewUser("1", "", "Test User", "hashedpassword")

	if err == nil {
		t.Errorf("Expected error for empty email, got nil")
	}
}

// TestUserIsActive tests the IsActive business logic
func TestUserIsActive(t *testing.T) {
	user, _ := NewUser("1", "test@example.com", "Test User", "hashedpassword")

	if !user.IsActive() {
		t.Errorf("Expected user to be active")
	}

	user.Block()
	if user.IsActive() {
		t.Errorf("Expected user to be inactive after blocking")
	}
}

// TestUserUpdateLastLogin tests UpdateLastLogin method
func TestUserUpdateLastLogin(t *testing.T) {
	user, _ := NewUser("1", "test@example.com", "Test User", "hashedpassword")

	if user.LastLogin != nil {
		t.Errorf("Expected no last login initially")
	}

	before := time.Now()
	user.UpdateLastLogin()
	after := time.Now()

	if user.LastLogin == nil {
		t.Errorf("Expected LastLogin to be set")
	}

	if user.LastLogin.Before(before) || user.LastLogin.After(after.Add(1*time.Second)) {
		t.Errorf("LastLogin time is not within expected range")
	}
}

// TestTokenExpiration tests the Token value object
func TestTokenExpiration(t *testing.T) {
	token := NewToken("access", "refresh", 3600)

	if token.IsExpired() {
		t.Errorf("Expected token to not be expired immediately")
	}

	expiredToken := NewToken("access", "refresh", -1)
	if !expiredToken.IsExpired() {
		t.Errorf("Expected token to be expired")
	}
}

// TestCredentialCreation tests Credential value object
func TestCredentialCreation(t *testing.T) {
	cred, err := NewCredential("test@example.com", "password")

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if cred.Email != "test@example.com" {
		t.Errorf("Expected email test@example.com, got %s", cred.Email)
	}
}

// TestCredentialCreation_InvalidPassword tests error handling
func TestCredentialCreation_InvalidPassword(t *testing.T) {
	_, err := NewCredential("test@example.com", "")

	if err == nil {
		t.Errorf("Expected error for empty password, got nil")
	}
}
