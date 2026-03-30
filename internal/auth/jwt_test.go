package auth

import (
	"strings"
	"testing"
	"time"
)

func TestJWTManager_GenerateToken(t *testing.T) {
	manager := NewJWTManager("my-secret-key", 3600)

	tests := []struct {
		name   string
		userID string
	}{
		{
			name:   "generate token for valid userID",
			userID: "user123",
		},
		{
			name:   "generate token for email userID",
			userID: "user@example.com",
		},
		{
			name:   "generate token with long userID",
			userID: "very_long_user_id_with_many_characters_and_special-chars_123",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			token, err := manager.GenerateToken(tt.userID)

			if err != nil {
				t.Errorf("GenerateToken() error = %v", err)
			}

			if token == "" {
				t.Errorf("GenerateToken() returned empty token")
			}

			// Token should have two parts separated by dot
			parts := strings.Split(token, ".")
			if len(parts) != 2 {
				t.Errorf("GenerateToken() token parts = %d, want 2", len(parts))
			}
		})
	}
}

func TestJWTManager_VerifyToken(t *testing.T) {
	secret := "test-secret-key"
	ttl := 3600
	manager := NewJWTManager(secret, ttl)

	userID := "testuser"
	token, err := manager.GenerateToken(userID)
	if err != nil {
		t.Fatalf("GenerateToken() failed: %v", err)
	}

	tests := []struct {
		name       string
		token      string
		wantUserID string
		wantErr    bool
		errMsg     string
	}{
		{
			name:       "verify valid token",
			token:      token,
			wantUserID: userID,
			wantErr:    false,
		},
		{
			name:    "verify invalid format",
			token:   "invalid_token_format",
			wantErr: true,
			errMsg:  "invalid token format",
		},
		{
			name:    "verify tampered signature",
			token:   strings.Split(token, ".")[0] + ".invalidsignature",
			wantErr: true,
			errMsg:  "invalid token signature",
		},
		{
			name:    "verify empty token",
			token:   "",
			wantErr: true,
			errMsg:  "invalid token format",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			userID, err := manager.VerifyToken(tt.token)

			if (err != nil) != tt.wantErr {
				t.Errorf("VerifyToken() error = %v, wantErr %v", err, tt.wantErr)
			}

			if tt.wantErr && err.Error() != tt.errMsg {
				t.Errorf("VerifyToken() error message = %v, want %v", err.Error(), tt.errMsg)
			}

			if !tt.wantErr && userID != tt.wantUserID {
				t.Errorf("VerifyToken() userID = %v, want %v", userID, tt.wantUserID)
			}
		})
	}
}

func TestJWTManager_TokenExpiration(t *testing.T) {
	// Create manager with 1 second TTL
	manager := NewJWTManager("secret", 1)

	userID := "user123"
	token, err := manager.GenerateToken(userID)
	if err != nil {
		t.Fatalf("GenerateToken() failed: %v", err)
	}

	// Token should be valid immediately
	returnedID, err := manager.VerifyToken(token)
	if err != nil {
		t.Errorf("VerifyToken() should succeed for fresh token: %v", err)
	}
	if returnedID != userID {
		t.Errorf("VerifyToken() userID = %v, want %v", returnedID, userID)
	}

	// Wait for token to expire
	time.Sleep(2 * time.Second)

	// Token should be expired now
	_, err = manager.VerifyToken(token)
	if err == nil {
		t.Errorf("VerifyToken() should fail for expired token")
	}
	if err.Error() != "token expired" {
		t.Errorf("VerifyToken() error message = %v, want 'token expired'", err.Error())
	}
}

func TestJWTManager_DifferentSecrets(t *testing.T) {
	userID := "user123"

	managerA := NewJWTManager("secret-a", 3600)
	managerB := NewJWTManager("secret-b", 3600)

	token, err := managerA.GenerateToken(userID)
	if err != nil {
		t.Fatalf("GenerateToken() failed: %v", err)
	}

	// ManagerA should verify its own token
	_, err = managerA.VerifyToken(token)
	if err != nil {
		t.Errorf("ManagerA should verify its own token: %v", err)
	}

	// ManagerB should NOT verify token from ManagerA (different secret)
	_, err = managerB.VerifyToken(token)
	if err == nil {
		t.Errorf("ManagerB should not verify token from ManagerA (different secret)")
	}
	if err.Error() != "invalid token signature" {
		t.Errorf("VerifyToken() error = %v, want 'invalid token signature'", err.Error())
	}
}

func TestJWTManager_SpecialCharactersInUserID(t *testing.T) {
	manager := NewJWTManager("secret", 3600)

	specialUserIDs := []string{
		"user@example.com",
		"user-with-dashes",
		"user_with_underscores",
		"user.with.dots",
		"user#with#hashes",
		"user:with:colons",
	}

	for _, uid := range specialUserIDs {
		token, err := manager.GenerateToken(uid)
		if err != nil {
			t.Errorf("GenerateToken() failed for userID %q: %v", uid, err)
			continue
		}

		returnedID, err := manager.VerifyToken(token)
		if err != nil {
			t.Errorf("VerifyToken() failed for userID %q: %v", uid, err)
			continue
		}

		if returnedID != uid {
			t.Errorf("VerifyToken() userID = %v, want %v", returnedID, uid)
		}
	}
}
