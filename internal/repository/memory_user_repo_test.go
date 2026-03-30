package repository

import (
	"sso/internal/model"
	"testing"
)

func TestMemoryUserRepository_Create(t *testing.T) {
	tests := []struct {
		name    string
		user    model.User
		wantErr bool
		errMsg  string
	}{
		{
			name:    "create new user",
			user:    model.User{ID: "test@example.com", Email: "test@example.com", Password: "hashed_pwd"},
			wantErr: false,
		},
		{
			name:    "create duplicate user",
			user:    model.User{ID: "dup@example.com", Email: "dup@example.com", Password: "hashed_pwd1"},
			wantErr: true,
			errMsg:  "user exists",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := NewMemoryUserRepository()

			// First creation
			err := repo.Create(tt.user)
			if tt.name == "create duplicate user" {
				// Create first, then test duplicate
				err = repo.Create(tt.user)
			}

			if (err != nil) != tt.wantErr {
				t.Errorf("Create() error = %v, wantErr %v", err, tt.wantErr)
			}

			if tt.wantErr && err != nil && err.Error() != tt.errMsg {
				t.Errorf("Create() error message = %v, want %v", err.Error(), tt.errMsg)
			}
		})
	}
}

func TestMemoryUserRepository_GetByEmail(t *testing.T) {
	user := model.User{
		ID:       "123",
		Email:    "test@example.com",
		Password: "hashed_password",
	}

	tests := []struct {
		name    string
		email   string
		setup   func() (UserRepository, error)
		wantErr bool
	}{
		{
			name:  "get existing user",
			email: "test@example.com",
			setup: func() (UserRepository, error) {
				repo := NewMemoryUserRepository()
				err := repo.Create(user)
				return repo, err
			},
			wantErr: false,
		},
		{
			name:  "get non-existing user",
			email: "notfound@example.com",
			setup: func() (UserRepository, error) {
				return NewMemoryUserRepository(), nil
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo, err := tt.setup()
			if err != nil {
				t.Fatalf("setup failed: %v", err)
			}

			got, err := repo.GetByEmail(tt.email)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetByEmail() error = %v, wantErr %v", err, tt.wantErr)
			}

			if !tt.wantErr {
				if got.Email != tt.email {
					t.Errorf("GetByEmail() got email = %v, want %v", got.Email, tt.email)
				}
				if got.Password != user.Password {
					t.Errorf("GetByEmail() got password = %v, want %v", got.Password, user.Password)
				}
			}
		})
	}
}

func TestMemoryUserRepository_ConcurrentAccess(t *testing.T) {
	repo := NewMemoryUserRepository()
	numGoroutines := 10

	// Concurrent writes
	done := make(chan bool, numGoroutines)
	for i := 0; i < numGoroutines; i++ {
		go func(idx int) {
			email := "user" + string(rune(idx)) + "@example.com"
			user := model.User{ID: email, Email: email, Password: "pwd"}
			_ = repo.Create(user)
			done <- true
		}(i)
	}

	// Wait for all goroutines
	for i := 0; i < numGoroutines; i++ {
		<-done
	}

	// Verify all users exist
	for i := 0; i < 10; i++ {
		email := "user" + string(rune(i)) + "@example.com"
		_, err := repo.GetByEmail(email)
		if err != nil {
			t.Errorf("Expected to find user %s, but got error: %v", email, err)
		}
	}
}
