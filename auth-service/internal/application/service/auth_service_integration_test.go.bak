package service

import (
	"context"
	"testing"

	"github.com/rusgainew/tunduck-app-mk/auth-service/internal/application/dto"
	"github.com/rusgainew/tunduck-app-mk/auth-service/internal/domain/entity"
	"golang.org/x/crypto/bcrypt"
)

// TestRegisterAndLoginFlow - Integration test: register then login
func TestRegisterAndLoginFlow(t *testing.T) {
	mockRepo := NewMockUserRepository()
	mockPublisher := NewMockEventPublisher()

	registerService := NewRegisterUserService(mockRepo, mockPublisher)
	loginService := NewLoginUserService(mockRepo, nil, mockPublisher) // Will set token service later

	// Step 1: Register user
	registerReq := &dto.RegisterRequest{
		Email:    "test@example.com",
		Name:     "Test User",
		Password: "SecurePassword123",
	}

	registerResp, err := registerService.Execute(context.Background(), registerReq)
	if err != nil {
		t.Fatalf("Registration failed: %v", err)
	}

	if registerResp.Email != "test@example.com" {
		t.Errorf("Expected email test@example.com, got %s", registerResp.Email)
	}

	// Verify user was created
	exists, _ := mockRepo.UserExists(context.Background(), "test@example.com")
	if !exists {
		t.Error("User should exist after registration")
	}

	// Verify event was published
	if len(mockPublisher.events) != 1 {
		t.Errorf("Expected 1 event, got %d", len(mockPublisher.events))
	}
}

// TestLoginWithInvalidCredentials - Test login with wrong password
func TestLoginWithInvalidCredentials(t *testing.T) {
	ctx := context.Background()
	mockRepo := NewMockUserRepository()
	mockPublisher := NewMockEventPublisher()

	// Create user manually
	passwordHash, _ := bcrypt.GenerateFromPassword([]byte("CorrectPassword123"), bcrypt.DefaultCost)
	user, _ := entity.NewUser("test_id", "test@example.com", "Test User", string(passwordHash))
	mockRepo.CreateUser(ctx, user)

	loginService := NewLoginUserService(mockRepo, nil, mockPublisher)

	// Try to login with wrong password
	loginReq := &dto.LoginRequest{
		Email:    "test@example.com",
		Password: "WrongPassword123",
	}

	_, err := loginService.Execute(ctx, loginReq)
	if err == nil {
		t.Error("Expected error for wrong password, got nil")
	}
}

// TestDuplicateEmailRegistration - Test cannot register twice with same email
func TestDuplicateEmailRegistration(t *testing.T) {
	ctx := context.Background()
	mockRepo := NewMockUserRepository()
	mockPublisher := NewMockEventPublisher()

	registerService := NewRegisterUserService(mockRepo, mockPublisher)

	registerReq := &dto.RegisterRequest{
		Email:    "test@example.com",
		Name:     "Test User",
		Password: "Password123",
	}

	// Register first time
	_, err1 := registerService.Execute(ctx, registerReq)
	if err1 != nil {
		t.Fatalf("First registration failed: %v", err1)
	}

	// Try to register again with same email
	_, err2 := registerService.Execute(ctx, registerReq)
	if err2 == nil {
		t.Error("Expected error for duplicate email, got nil")
	}
}

// TestUserStatusAfterRegistration - Verify user is active after registration
func TestUserStatusAfterRegistration(t *testing.T) {
	ctx := context.Background()
	mockRepo := NewMockUserRepository()
	mockPublisher := NewMockEventPublisher()

	registerService := NewRegisterUserService(mockRepo, mockPublisher)

	registerReq := &dto.RegisterRequest{
		Email:    "test@example.com",
		Name:     "Test User",
		Password: "Password123",
	}

	registerService.Execute(ctx, registerReq)

	// Verify user status
	user, _ := mockRepo.GetUserByEmail(ctx, "test@example.com")
	if !user.IsActive() {
		t.Error("User should be active after registration")
	}
}

// TestMultipleLoginsUpdateLastLogin - Each login should update last_login timestamp
func TestMultipleLoginsUpdateLastLogin(t *testing.T) {
	ctx := context.Background()
	mockRepo := NewMockUserRepository()
	mockPublisher := NewMockEventPublisher()

	// Create user
	passwordHash, _ := bcrypt.GenerateFromPassword([]byte("Password123"), bcrypt.DefaultCost)
	user, _ := entity.NewUser("test_id", "test@example.com", "Test User", string(passwordHash))
	mockRepo.CreateUser(ctx, user)

	loginService := NewLoginUserService(mockRepo, nil, mockPublisher)

	loginReq := &dto.LoginRequest{
		Email:    "test@example.com",
		Password: "Password123",
	}

	// First login
	loginService.Execute(ctx, loginReq)
	user1, _ := mockRepo.GetUserByEmail(ctx, "test@example.com")
	lastLogin1 := user1.LastLogin

	// Second login
	loginService.Execute(ctx, loginReq)
	user2, _ := mockRepo.GetUserByEmail(ctx, "test@example.com")
	lastLogin2 := user2.LastLogin

	if lastLogin1 == lastLogin2 {
		t.Error("LastLogin should be updated on each login")
	}

	if lastLogin2 == nil {
		t.Error("LastLogin should not be nil")
	}
}
