package repository

import (
    "errors"
    "sync"
    "sso/internal/model"
)

type memoryUserRepo struct {
    mu    sync.RWMutex
    users map[string]model.User
}

func NewMemoryUserRepository() UserRepository {
    return &memoryUserRepo{users: make(map[string]model.User)}
}

func (r *memoryUserRepo) Create(user model.User) error {
    r.mu.Lock()
    defer r.mu.Unlock()
    if _, ok := r.users[user.Email]; ok {
        return errors.New("user exists")
    }
    r.users[user.Email] = user
    return nil
}

func (r *memoryUserRepo) GetByEmail(email string) (model.User, error) {
    r.mu.RLock()
    defer r.mu.RUnlock()
    u, ok := r.users[email]
    if !ok {
        return model.User{}, errors.New("user not found")
    }
    return u, nil
}
