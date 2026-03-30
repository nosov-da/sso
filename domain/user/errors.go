package user

import "errors"

var (
// User-related errors
ErrUserAlreadyExists  = errors.New("user already exists")
ErrUserNotFound       = errors.New("user not found")
ErrInvalidCredentials = errors.New("invalid credentials")

// Validation errors
ErrInvalidEmail    = errors.New("invalid email format")
ErrEmptyEmail      = errors.New("email cannot be empty")
ErrWeakPassword    = errors.New("password must be at least 8 characters long")
ErrEmptyPassword   = errors.New("password cannot be empty")
ErrInvalidPassword = errors.New("invalid password hash")

// ID errors
ErrEmptyUserId = errors.New("user id cannot be empty")
ErrInvalidToken = errors.New("invalid token")
)
