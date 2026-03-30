package repository

import "sso/internal/model"

type UserRepository interface {
    Create(user model.User) error
    GetByEmail(email string) (model.User, error)
}
