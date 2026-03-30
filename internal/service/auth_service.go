package service

import (
    "errors"
    "sso/internal/auth"
    "sso/internal/model"
    "sso/internal/repository"
    "golang.org/x/crypto/bcrypt"
)

type AuthService struct {
    repo       repository.UserRepository
    jwtManager *auth.JWTManager
}

func NewAuthService(repo repository.UserRepository, jwtManager *auth.JWTManager) *AuthService {
    return &AuthService{repo: repo, jwtManager: jwtManager}
}

func (s *AuthService) Register(email, password string) error {
    _, err := s.repo.GetByEmail(email)
    if err == nil {
        return errors.New("user already exists")
    }

    hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
    if err != nil {
        return err
    }

    return s.repo.Create(model.User{ID: email, Email: email, Password: string(hashed)})
}

func (s *AuthService) Login(email, password string) (string, error) {
    user, err := s.repo.GetByEmail(email)
    if err != nil {
        return "", err
    }

    if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
        return "", errors.New("invalid credentials")
    }

    return s.jwtManager.GenerateToken(user.ID)
}
