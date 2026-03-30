package service

import (
	"sso/internal/auth"
	"sso/internal/repository"
	"testing"
)

func TestAuthService_Register(t *testing.T) {
	tests := []struct {
		name      string
		email     string
		password  string
		wantErr   bool
		wantErrMsg string
	}{
		{
			name:     "successful registration",
			email:    "newuser@example.com",
			password: "password123",
			wantErr:  false,
		},
		{
			name:      "duplicate registration",
			email:     "existing@example.com",
			password:  "password123",
			wantErr:   true,
			wantErrMsg: "user already exists",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := repository.NewMemoryUserRepository()
			jwtMgr := auth.NewJWTManager("secret", 3600)
			svc := NewAuthService(repo, jwtMgr)

			if tt.name == "duplicate registration" {
				// Pre-register the user for duplicate test
				svc.Register("existing@example.com", "password123")
			}

			err := svc.Register(tt.email, tt.password)

			if (err != nil) != tt.wantErr {
				t.Errorf("Register() error = %v, wantErr %v", err, tt.wantErr)
			}

			if tt.wantErr && err.Error() != tt.wantErrMsg {
				t.Errorf("Register() error message = %v, want %v", err.Error(), tt.wantErrMsg)
			}

			if !tt.wantErr {
				// Verify user was created
				user, err := repo.GetByEmail(tt.email)
				if err != nil {
					t.Errorf("Register() should have created user, but GetByEmail failed: %v", err)
				}
				if user.Email != tt.email {
					t.Errorf("Register() user email = %v, want %v", user.Email, tt.email)
				}
			}
		})
	}
}

func TestAuthService_Login(t *testing.T) {
	repo := repository.NewMemoryUserRepository()
	jwtMgr := auth.NewJWTManager("secret", 3600)
	svc := NewAuthService(repo, jwtMgr)

	// Register a test user
	testEmail := "testuser@example.com"
	testPassword := "correctpassword"
	svc.Register(testEmail, testPassword)

	tests := []struct {
		name       string
		email      string
		password   string
		wantErr    bool
		wantErrMsg string
	}{
		{
			name:     "successful login",
			email:    testEmail,
			password: testPassword,
			wantErr:  false,
		},
		{
			name:       "non-existing user",
			email:      "nonexistent@example.com",
			password:   "password",
			wantErr:    true,
			wantErrMsg: "user not found",
		},
		{
			name:       "wrong password",
			email:      testEmail,
			password:   "wrongpassword",
			wantErr:    true,
			wantErrMsg: "invalid credentials",
		},
		{
			name:       "empty email",
			email:      "",
			password:   testPassword,
			wantErr:    true,
			wantErrMsg: "user not found",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			token, err := svc.Login(tt.email, tt.password)

			if (err != nil) != tt.wantErr {
				t.Errorf("Login() error = %v, wantErr %v", err, tt.wantErr)
			}

			if tt.wantErr && err.Error() != tt.wantErrMsg {
				t.Errorf("Login() error message = %v, want %v", err.Error(), tt.wantErrMsg)
			}

			if !tt.wantErr {
				if token == "" {
					t.Errorf("Login() should return non-empty token")
				}

				// Verify token is valid
				userID, err := jwtMgr.VerifyToken(token)
				if err != nil {
					t.Errorf("Login() returned invalid token: %v", err)
				}

				if userID != testEmail {
					t.Errorf("Login() token userID = %v, want %v", userID, testEmail)
				}
			}
		})
	}
}

func TestAuthService_PasswordHashing(t *testing.T) {
	repo := repository.NewMemoryUserRepository()
	jwtMgr := auth.NewJWTManager("secret", 3600)
	svc := NewAuthService(repo, jwtMgr)

	email := "hash@example.com"
	password := "securepassword"

	// Register user
	err := svc.Register(email, password)
	if err != nil {
		t.Fatalf("Register() failed: %v", err)
	}

	// Get user from repo
	user, err := repo.GetByEmail(email)
	if err != nil {
		t.Fatalf("GetByEmail() failed: %v", err)
	}

	// Password should be hashed, not plaintext
	if user.Password == password {
		t.Error("Password should be hashed, not stored as plaintext")
	}

	// Password hash should be bcrypt format
	if len(user.Password) != 60 { // bcrypt hashes are 60 characters
		t.Errorf("Password hash length = %d, want 60 (bcrypt format)", len(user.Password))
	}
}
