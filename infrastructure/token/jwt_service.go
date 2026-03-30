package token

import (
	"sso/internal/auth"
)

// JWTTokenService implements TokenService using JWT
type JWTTokenService struct {
	jwtManager *auth.JWTManager
}

// NewJWTTokenService creates a new JWT token service
func NewJWTTokenService(jwtManager *auth.JWTManager) *JWTTokenService {
	return &JWTTokenService{
		jwtManager: jwtManager,
	}
}

// GenerateToken generates a new JWT token
func (s *JWTTokenService) GenerateToken(userID string) (string, error) {
	return s.jwtManager.GenerateToken(userID)
}

// VerifyToken verifies a JWT token
func (s *JWTTokenService) VerifyToken(token string) (string, error) {
	return s.jwtManager.VerifyToken(token)
}
