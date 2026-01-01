package service

import (
	"context"
	"errors"
	"testing"

	"github.com/rusgainew/tunduck-app-mk/auth-service/internal/application/dto"
	"github.com/rusgainew/tunduck-app-mk/auth-service/internal/domain/entity"
)

// MockUserRepository - Mock для тестирования
type MockUserRepository struct {
	users map[string]*entity.User
}

func NewMockUserRepository() *MockUserRepository {
	return &MockUserRepository{
		users: make(map[string]*entity.User),
	}
}

func (m *MockUserRepository) CreateUser(ctx context.Context, user *entity.User) error {
	if _, exists := m.users[user.Email]; exists {
		return errors.New("user already exists")
	}
	m.users[user.Email] = user
	return nil
}

func (m *MockUserRepository) GetUserByEmail(ctx context.Context, email string) (*entity.User, error) {
	user, exists := m.users[email]
	if !exists {
		return nil, errors.New("user not found")
	}
	return user, nil
}

func (m *MockUserRepository) GetUserByID(ctx context.Context, id string) (*entity.User, error) {
	for _, user := range m.users {
		if user.ID == id {
			return user, nil
		}
	}
	return nil, errors.New("user not found")
}

func (m *MockUserRepository) UpdateUser(ctx context.Context, user *entity.User) error {
	if _, exists := m.users[user.Email]; !exists {
		return errors.New("user not found")
	}
	m.users[user.Email] = user
	return nil
}

func (m *MockUserRepository) DeleteUser(ctx context.Context, id string) error {
	for email, user := range m.users {
		if user.ID == id {
			delete(m.users, email)
			return nil
		}
	}
	return errors.New("user not found")
}

func (m *MockUserRepository) UserExists(ctx context.Context, email string) (bool, error) {
	_, exists := m.users[email]
	return exists, nil
}

// MockEventPublisher - Mock для RabbitMQ
type MockEventPublisher struct {
	events []entity.DomainEvent
}

func NewMockEventPublisher() *MockEventPublisher {
	return &MockEventPublisher{
		events: []entity.DomainEvent{},
	}
}

func (m *MockEventPublisher) Publish(ctx context.Context, event entity.DomainEvent) error {
	m.events = append(m.events, event)
	return nil
}

// TESTS

// TestRegisterUserService_Success
func TestRegisterUserService_Success(t *testing.T) {
	mockRepo := NewMockUserRepository()
	mockPublisher := NewMockEventPublisher()
	svc := NewRegisterUserService(mockRepo, mockPublisher)

	req := &dto.RegisterRequest{
		Email:    "test@example.com",
		Name:     "Test User",
		Password: "ValidPass123",
	}

	resp, err := svc.Execute(context.Background(), req)

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if resp == nil {
		t.Errorf("Expected response, got nil")
	}
	if resp.Email != "test@example.com" {
		t.Errorf("Expected email test@example.com, got %s", resp.Email)
	}

	// Проверить, что событие было опубликовано
	if len(mockPublisher.events) != 1 {
		t.Errorf("Expected 1 event, got %d", len(mockPublisher.events))
	}
}

// TestRegisterUserService_DuplicateEmail
func TestRegisterUserService_DuplicateEmail(t *testing.T) {
	mockRepo := NewMockUserRepository()
	mockPublisher := NewMockEventPublisher()
	svc := NewRegisterUserService(mockRepo, mockPublisher)

	// Создать первого пользователя
	req := &dto.RegisterRequest{
		Email:    "test@example.com",
		Name:     "Test User",
		Password: "ValidPass123",
	}
	svc.Execute(context.Background(), req)

	// Попытка создать с тем же email
	req2 := &dto.RegisterRequest{
		Email:    "test@example.com",
		Name:     "Another User",
		Password: "ValidPass456",
	}
	_, err := svc.Execute(context.Background(), req2)

	if err == nil {
		t.Errorf("Expected error for duplicate email, got nil")
	}
}
