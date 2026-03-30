package user

import (
	"sso/domain/user"
	"sso/infrastructure/events"
	inmemoryrep "sso/infrastructure/persistence"
	"testing"
)

func TestRegisterUserUseCase_Success(t *testing.T) {
	repo := inmemoryrep.NewInMemoryUserRepository()
	publisher := events.NewInMemoryEventPublisher()
	uc := NewRegisterUserUseCase(repo, publisher)

	input := RegisterInput{
		Email:    "newuser@example.com",
		Password: "password123",
	}

	err := uc.Execute(input)
	if err != nil {
		t.Errorf("Execute() error = %v, want nil", err)
	}

	// Verify user was saved
	email, _ := user.NewEmail(input.Email)
	savedUser, err := repo.GetByEmail(email)
	if err != nil {
		t.Errorf("GetByEmail() error = %v, want nil", err)
	}

	if savedUser.GetEmail() != input.Email {
		t.Errorf("Saved user email = %v, want %v", savedUser.GetEmail(), input.Email)
	}
}

func TestRegisterUserUseCase_InvalidEmail(t *testing.T) {
	repo := inmemoryrep.NewInMemoryUserRepository()
	publisher := events.NewInMemoryEventPublisher()
	uc := NewRegisterUserUseCase(repo, publisher)

	input := RegisterInput{
		Email:    "invalid-email",
		Password: "password123",
	}

	err := uc.Execute(input)
	if err == nil {
		t.Error("Execute() should return error for invalid email")
	}
}

func TestRegisterUserUseCase_DuplicateUser(t *testing.T) {
	repo := inmemoryrep.NewInMemoryUserRepository()
	publisher := events.NewInMemoryEventPublisher()
	uc := NewRegisterUserUseCase(repo, publisher)

	input := RegisterInput{
		Email:    "existing@example.com",
		Password: "password123",
	}

	// First registration should succeed
	err := uc.Execute(input)
	if err != nil {
		t.Errorf("First Execute() error = %v, want nil", err)
	}

	// Second registration should fail
	err = uc.Execute(input)
	if err != user.ErrUserAlreadyExists {
		t.Errorf("Second Execute() error = %v, want ErrUserAlreadyExists", err)
	}
}

func TestRegisterUserUseCase_PublishesEvent(t *testing.T) {
	repo := inmemoryrep.NewInMemoryUserRepository()
	publisher := events.NewInMemoryEventPublisher()
	
	eventFired := false
	publisher.Subscribe("user.registered", func(event user.DomainEvent) error {
		eventFired = true
		return nil
	})

	uc := NewRegisterUserUseCase(repo, publisher)
	input := RegisterInput{
		Email:    "eventuser@example.com",
		Password: "password123",
	}

	uc.Execute(input)

	if !eventFired {
		t.Error("UserRegisteredEvent should be published")
	}
}

func TestAuthenticateUserUseCase_Success(t *testing.T) {
	// Setup
	repo := inmemoryrep.NewInMemoryUserRepository()
	publisher := events.NewInMemoryEventPublisher()
	
	// Create a mock token service
	mockTokenService := &mockTokenService{
		token: "test-token-123",
	}

	registerUC := NewRegisterUserUseCase(repo, publisher)
	authenticateUC := NewAuthenticateUserUseCase(repo, mockTokenService, publisher)

	email := "user@example.com"
	password := "password123"

	// Register user first
	registerUC.Execute(RegisterInput{Email: email, Password: password})

	// Try to authenticate
	input := AuthenticateInput{
		Email:    email,
		Password: password,
	}

	output, err := authenticateUC.Execute(input)
	if err != nil {
		t.Errorf("Execute() error = %v, want nil", err)
	}

	if output.Token != "test-token-123" {
		t.Errorf("Token = %v, want 'test-token-123'", output.Token)
	}
}

func TestAuthenticateUserUseCase_WrongPassword(t *testing.T) {
	repo := inmemoryrep.NewInMemoryUserRepository()
	publisher := events.NewInMemoryEventPublisher()
	mockTokenService := &mockTokenService{token: "test-token"}

	registerUC := NewRegisterUserUseCase(repo, publisher)
	authenticateUC := NewAuthenticateUserUseCase(repo, mockTokenService, publisher)

	email := "user@example.com"
	registerUC.Execute(RegisterInput{Email: email, Password: "correct-password"})

	input := AuthenticateInput{
		Email:    email,
		Password: "wrong-password",
	}

	_, err := authenticateUC.Execute(input)
	if err == nil {
		t.Error("Execute() should return error for wrong password")
	}
}

func TestAuthenticateUserUseCase_NonExistentUser(t *testing.T) {
	repo := inmemoryrep.NewInMemoryUserRepository()
	publisher := events.NewInMemoryEventPublisher()
	mockTokenService := &mockTokenService{token: "test-token"}

	authenticateUC := NewAuthenticateUserUseCase(repo, mockTokenService, publisher)

	input := AuthenticateInput{
		Email:    "nonexistent@example.com",
		Password: "password",
	}

	_, err := authenticateUC.Execute(input)
	if err == nil {
		t.Error("Execute() should return error for non-existent user")
	}
}

func TestAuthenticateUserUseCase_PublishesEvent(t *testing.T) {
	repo := inmemoryrep.NewInMemoryUserRepository()
	publisher := events.NewInMemoryEventPublisher()
	mockTokenService := &mockTokenService{token: "test-token"}

	eventFired := false
	publisher.Subscribe("user.logged_in", func(event user.DomainEvent) error {
		eventFired = true
		return nil
	})

	registerUC := NewRegisterUserUseCase(repo, publisher)
	authenticateUC := NewAuthenticateUserUseCase(repo, mockTokenService, publisher)

	email := "user@example.com"
	registerUC.Execute(RegisterInput{Email: email, Password: "password123"})

	authenticateUC.Execute(AuthenticateInput{Email: email, Password: "password123"})

	if !eventFired {
		t.Error("UserLoggedInEvent should be published")
	}
}

// Mock token service for testing
type mockTokenService struct {
	token string
}

func (m *mockTokenService) GenerateToken(userID string) (string, error) {
	return m.token, nil
}

func (m *mockTokenService) VerifyToken(token string) (string, error) {
	return "user-id", nil
}
