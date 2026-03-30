package postgres

import (
	"sso/domain/user"
	"testing"
)

func TestInMemoryUserRepository_Save(t *testing.T) {
	repo := NewInMemoryUserRepository()
	
	userId, _ := user.NewUserId("user@example.com")
	email, _ := user.NewEmail("user@example.com")
	password, _ := user.NewPassword("hashed_pwd")
	usr := user.NewUser(userId, email, password)

	err := repo.Save(usr)
	if err != nil {
		t.Errorf("Save() error = %v, want nil", err)
	}

	// Try to save duplicate
	err = repo.Save(usr)
	if err != user.ErrUserAlreadyExists {
		t.Errorf("Save() duplicate error = %v, want ErrUserAlreadyExists", err)
	}
}

func TestInMemoryUserRepository_GetByEmail(t *testing.T) {
	repo := NewInMemoryUserRepository()
	
	userId, _ := user.NewUserId("user@example.com")
	email, _ := user.NewEmail("user@example.com")
	password, _ := user.NewPassword("hashed_pwd")
	usr := user.NewUser(userId, email, password)

	repo.Save(usr)

	// Get existing user
	retrieved, err := repo.GetByEmail(email)
	if err != nil {
		t.Errorf("GetByEmail() error = %v, want nil", err)
	}

	if !retrieved.Equals(usr) {
		t.Error("Retrieved user should equal saved user")
	}

	// Get non-existing user
	nonExistentEmail, _ := user.NewEmail("nonexistent@example.com")
	_, err = repo.GetByEmail(nonExistentEmail)
	if err != user.ErrUserNotFound {
		t.Errorf("GetByEmail() error = %v, want ErrUserNotFound", err)
	}
}

func TestInMemoryUserRepository_GetByID(t *testing.T) {
	repo := NewInMemoryUserRepository()
	
	userId, _ := user.NewUserId("user@example.com")
	email, _ := user.NewEmail("user@example.com")
	password, _ := user.NewPassword("hashed_pwd")
	usr := user.NewUser(userId, email, password)

	repo.Save(usr)

	// Get existing user
	retrieved, err := repo.GetByID(userId)
	if err != nil {
		t.Errorf("GetByID() error = %v, want nil", err)
	}

	if !retrieved.Equals(usr) {
		t.Error("Retrieved user should equal saved user")
	}

	// Get non-existing user
	nonExistentId, _ := user.NewUserId("nonexistent@example.com")
	_, err = repo.GetByID(nonExistentId)
	if err != user.ErrUserNotFound {
		t.Errorf("GetByID() error = %v, want ErrUserNotFound", err)
	}
}

func TestInMemoryUserRepository_Exists(t *testing.T) {
	repo := NewInMemoryUserRepository()
	
	userId, _ := user.NewUserId("user@example.com")
	email, _ := user.NewEmail("user@example.com")
	password, _ := user.NewPassword("hashed_pwd")
	usr := user.NewUser(userId, email, password)

	repo.Save(usr)

	// Check existing user
	exists, err := repo.Exists(email)
	if err != nil {
		t.Errorf("Exists() error = %v, want nil", err)
	}

	if !exists {
		t.Error("Exists() should return true for existing user")
	}

	// Check non-existing user
	nonExistentEmail, _ := user.NewEmail("nonexistent@example.com")
	exists, err = repo.Exists(nonExistentEmail)
	if err != nil {
		t.Errorf("Exists() error = %v, want nil", err)
	}

	if exists {
		t.Error("Exists() should return false for non-existing user")
	}
}
