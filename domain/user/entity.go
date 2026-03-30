package user

// User - Domain Entity representing a user
type User struct {
	id    UserId
	email Email
	// password is stored as a hash only
	password Password
}

// NewUser creates a new User entity
func NewUser(id UserId, email Email, passwordHash Password) *User {
	return &User{
		id:       id,
		email:    email,
		password: passwordHash,
	}
}

// ID returns the user's unique identifier
func (u *User) ID() UserId {
	return u.id
}

// Email returns the user's email
func (u *User) Email() Email {
	return u.email
}

// Password returns the user's hashed password
func (u *User) Password() Password {
	return u.password
}

// Equals checks if two users are the same entity
func (u *User) Equals(other *User) bool {
	if other == nil {
		return false
	}
	return u.id.Equals(other.id)
}

// IdentifyBy returns a string identifier for the user (for JSON/persistence)
func (u *User) IdentifyBy() string {
	return u.id.Value()
}

// GetEmail returns the email string (convenience method)
func (u *User) GetEmail() string {
	return u.email.Value()
}

// GetPasswordHash returns the password hash string (convenience method)
func (u *User) GetPasswordHash() string {
	return u.password.Hash()
}
