package user

// Repository defines the interface for user persistence
type Repository interface {
	// Save stores a new user in the repository
	Save(user *User) error

	// GetByEmail retrieves a user by their email
	GetByEmail(email Email) (*User, error)

	// GetByID retrieves a user by their ID
	GetByID(id UserId) (*User, error)

	// Exists checks if a user with given email exists
	Exists(email Email) (bool, error)
}
