package user

import "testing"

func TestNewEmail_Valid(t *testing.T) {
	tests := []struct {
		name     string
		email    string
		wantErr  bool
	}{
		{
			name:    "valid email",
			email:   "user@example.com",
			wantErr: false,
		},
		{
			name:    "valid email with dots",
			email:   "user.name@example.co.uk",
			wantErr: false,
		},
		{
			name:    "empty email",
			email:   "",
			wantErr: true,
		},
		{
			name:    "invalid email format",
			email:   "notanemail",
			wantErr: true,
		},
		{
			name:    "invalid email no domain",
			email:   "user@",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			email, err := NewEmail(tt.email)

			if (err != nil) != tt.wantErr {
				t.Errorf("NewEmail() error = %v, wantErr %v", err, tt.wantErr)
			}

			if !tt.wantErr && email.Value() != tt.email {
				t.Errorf("Email value = %v, want %v", email.Value(), tt.email)
			}
		})
	}
}

func TestEmail_Equals(t *testing.T) {
	email1, _ := NewEmail("user@example.com")
	email2, _ := NewEmail("user@example.com")
	email3, _ := NewEmail("other@example.com")

	if !email1.Equals(email2) {
		t.Error("Same emails should be equal")
	}

	if email1.Equals(email3) {
		t.Error("Different emails should not be equal")
	}
}

func TestNewPassword(t *testing.T) {
	tests := []struct {
		name    string
		hash    string
		wantErr bool
	}{
		{
			name:    "valid hash",
			hash:    "$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcg7b3XeKeUxWdeS86E36jbMfsm",
			wantErr: false,
		},
		{
			name:    "empty hash",
			hash:    "",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pwd, err := NewPassword(tt.hash)

			if (err != nil) != tt.wantErr {
				t.Errorf("NewPassword() error = %v, wantErr %v", err, tt.wantErr)
			}

			if !tt.wantErr && pwd.Hash() != tt.hash {
				t.Errorf("Password hash = %v, want %v", pwd.Hash(), tt.hash)
			}
		})
	}
}

func TestNewUserId(t *testing.T) {
	tests := []struct {
		name    string
		id      string
		wantErr bool
	}{
		{
			name:    "valid id",
			id:      "user@example.com",
			wantErr: false,
		},
		{
			name:    "empty id",
			id:      "",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			userId, err := NewUserId(tt.id)

			if (err != nil) != tt.wantErr {
				t.Errorf("NewUserId() error = %v, wantErr %v", err, tt.wantErr)
			}

			if !tt.wantErr && userId.Value() != tt.id {
				t.Errorf("UserId value = %v, want %v", userId.Value(), tt.id)
			}
		})
	}
}

func TestNewToken(t *testing.T) {
	tests := []struct {
		name    string
		token   string
		wantErr bool
	}{
		{
			name:    "valid token",
			token:   "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIn0.dozjgNryP4J3jVmNHl0w5N_XgL0n3I9PlFUP0THsR8U",
			wantErr: false,
		},
		{
			name:    "empty token",
			token:   "",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			token, err := NewToken(tt.token)

			if (err != nil) != tt.wantErr {
				t.Errorf("NewToken() error = %v, wantErr %v", err, tt.wantErr)
			}

			if !tt.wantErr && token.Value() != tt.token {
				t.Errorf("Token value = %v, want %v", token.Value(), tt.token)
			}
		})
	}
}
