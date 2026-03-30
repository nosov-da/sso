package postgres

import (
	"sync"
	"sso/domain/user"
)

// InMemoryUserRepository is an in-memory implementation of the user repository
type InMemoryUserRepository struct {
	mu    sync.RWMutex
	users map[string]*user.User
}

// NewInMemoryUserRepository creates a new in-memory repository
func NewInMemoryUserRepository() *InMemoryUserRepository {
	return &InMemoryUserRepository{
		users: make(map[string]*user.User),
	}
}

// Save stores a user in memory
func (r *InMemoryUserRepository) Save(usr *user.User) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	key := usr.GetEmail()
	if _, exists := r.users[key]; exists {
		return user.ErrUserAlreadyExists
	}

	r.users[key] = usr
	return nil
}

// GetByEmail retrieves a user by email
func (r *InMemoryUserRepository) GetByEmail(email user.Email) (*user.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	usr, exists := r.users[email.Value()]
	if !exists {
		return nil, user.ErrUserNotFound
	}

	return usr, nil
}

// GetByID retrieves a user by ID
func (r *InMemoryUserRepository) GetByID(id user.UserId) (*user.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	// For in-memory implementation, we search by ID value
	for _, u := range r.users {
		if u.ID().Equals(id) {
			return u, nil
		}
	}

	return nil, user.ErrUserNotFound
}

// Exists checks if a user with given email exists
func (r *InMemoryUserRepository) Exists(email user.Email) (bool, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	_, exists := r.users[email.Value()]
	return exists, nil
}
