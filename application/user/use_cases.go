package user

import (
"sso/domain/user"
"golang.org/x/crypto/bcrypt"
)

// TokenService defines the interface for token operations
type TokenService interface {
GenerateToken(userID string) (string, error)
VerifyToken(token string) (string, error)
}

// RegisterUserUseCase handles user registration
type RegisterUserUseCase struct {
repository user.Repository
eventPublisher user.EventPublisher
}

// NewRegisterUserUseCase creates a new registration use case
func NewRegisterUserUseCase(repo user.Repository, publisher user.EventPublisher) *RegisterUserUseCase {
return &RegisterUserUseCase{
repository: repo,
eventPublisher: publisher,
}
}

// Input for registration
type RegisterInput struct {
Email    string
Password string
}

// Execute registers a new user
func (uc *RegisterUserUseCase) Execute(input RegisterInput) error {
// Validate input
if input.Email == "" {
return user.ErrEmptyEmail
}
if input.Password == "" {
return user.ErrEmptyPassword
}

// Validate password strength
if err := user.ValidatePassword(input.Password); err != nil {
return err
}

// Create value objects
email, err := user.NewEmail(input.Email)
if err != nil {
return err
}

// Check if user already exists
exists, err := uc.repository.Exists(email)
if err != nil {
return err
}
if exists {
return user.ErrUserAlreadyExists
}

// Hash password
hashed, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
if err != nil {
return err
}

password, err := user.NewPassword(string(hashed))
if err != nil {
return err
}

// Create user ID (using email as ID)
userId, err := user.NewUserId(input.Email)
if err != nil {
return err
}

// Create and save user entity
usr := user.NewUser(userId, email, password)
if err := uc.repository.Save(usr); err != nil {
return err
}

// Publish domain event
event := user.NewUserRegisteredEvent(userId.Value(), email.Value())
if err := uc.eventPublisher.Publish(event); err != nil {
// Log but don't fail - event publishing is non-critical
_ = err
}

return nil
}

// AuthenticateUserUseCase handles user authentication
type AuthenticateUserUseCase struct {
repository user.Repository
tokenService TokenService
eventPublisher user.EventPublisher
}

// NewAuthenticateUserUseCase creates a new authentication use case
func NewAuthenticateUserUseCase(repo user.Repository, ts TokenService, publisher user.EventPublisher) *AuthenticateUserUseCase {
return &AuthenticateUserUseCase{
repository: repo,
tokenService: ts,
eventPublisher: publisher,
}
}

// Input for authentication
type AuthenticateInput struct {
Email    string
Password string
}

// Output of authentication
type AuthenticateOutput struct {
Token string
}

// Execute authenticates a user and returns a token
func (uc *AuthenticateUserUseCase) Execute(input AuthenticateInput) (AuthenticateOutput, error) {
output := AuthenticateOutput{}

// Validate input
if input.Email == "" {
return output, user.ErrEmptyEmail
}
if input.Password == "" {
return output, user.ErrEmptyPassword
}

// Create value object
email, err := user.NewEmail(input.Email)
if err != nil {
return output, err
}

// Get user from repository
usr, err := uc.repository.GetByEmail(email)
if err != nil {
if err == user.ErrUserNotFound {
return output, user.ErrInvalidCredentials
}
return output, err
}

// Verify password
if err := bcrypt.CompareHashAndPassword(
[]byte(usr.GetPasswordHash()),
[]byte(input.Password),
); err != nil {
return output, user.ErrInvalidCredentials
}

// Generate token
token, err := uc.tokenService.GenerateToken(usr.IdentifyBy())
if err != nil {
return output, err
}

output.Token = token

// Publish domain event
event := user.NewUserLoggedInEvent(usr.IdentifyBy())
if err := uc.eventPublisher.Publish(event); err != nil {
// Log but don't fail - event publishing is non-critical
_ = err
}

return output, nil
}
