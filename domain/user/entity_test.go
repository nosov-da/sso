package user

import "testing"

func TestNewUser(t *testing.T) {
	userId, _ := NewUserId("user@example.com")
	email, _ := NewEmail("user@example.com")
	password, _ := NewPassword("hashed_password")

	user := NewUser(userId, email, password)

	if user == nil {
		t.Error("NewUser() should not return nil")
	}

	if !user.ID().Equals(userId) {
		t.Errorf("User ID = %v, want %v", user.ID(), userId)
	}

	if !user.Email().Equals(email) {
		t.Errorf("User Email = %v, want %v", user.Email(), email)
	}

	if !user.Password().Equals(password) {
		t.Errorf("User Password = %v, want %v", user.Password(), password)
	}
}

func TestUser_Equals(t *testing.T) {
	userId1, _ := NewUserId("user1@example.com")
	userId2, _ := NewUserId("user2@example.com")
	email, _ := NewEmail("user1@example.com")
	password, _ := NewPassword("hashed_password")

	user1 := NewUser(userId1, email, password)
	user2 := NewUser(userId1, email, password)
	user3 := NewUser(userId2, email, password)

	if !user1.Equals(user2) {
		t.Error("Users with same ID should be equal")
	}

	if user1.Equals(user3) {
		t.Error("Users with different ID should not be equal")
	}

	if user1.Equals(nil) {
		t.Error("User should not be equal to nil")
	}
}

func TestUser_GettersAndHelpers(t *testing.T) {
	userId, _ := NewUserId("user@example.com")
	email, _ := NewEmail("user@example.com")
	passwordHash := "bcrypt_hash_value"
	password, _ := NewPassword(passwordHash)

	user := NewUser(userId, email, password)

	if user.IdentifyBy() != "user@example.com" {
		t.Errorf("IdentifyBy() = %v, want 'user@example.com'", user.IdentifyBy())
	}

	if user.GetEmail() != "user@example.com" {
		t.Errorf("GetEmail() = %v, want 'user@example.com'", user.GetEmail())
	}

	if user.GetPasswordHash() != passwordHash {
		t.Errorf("GetPasswordHash() = %v, want %v", user.GetPasswordHash(), passwordHash)
	}
}
