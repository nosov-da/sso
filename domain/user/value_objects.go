package user

import (
"regexp"
)

// Email - Value Object for user email
type Email struct {
value string
}

// NewEmail creates a new Email value object
func NewEmail(email string) (Email, error) {
if email == "" {
return Email{}, ErrEmptyEmail
}

// Simple email validation
if !isValidEmail(email) {
return Email{}, ErrInvalidEmail
}

return Email{value: email}, nil
}

// Value returns the string representation of email
func (e Email) Value() string {
return e.value
}

// Equals compares two Email value objects
func (e Email) Equals(other Email) bool {
return e.value == other.value
}

// String implements Stringer interface
func (e Email) String() string {
return e.value
}

// Password - Value Object for user password
type Password struct {
hash string
}

// NewPassword creates a new Password value object with hashed value
func NewPassword(hash string) (Password, error) {
if hash == "" {
return Password{}, ErrInvalidPassword
}
return Password{hash: hash}, nil
}

// Hash returns the hashed password
func (p Password) Hash() string {
return p.hash
}

// Equals compares two Password value objects
func (p Password) Equals(other Password) bool {
return p.hash == other.hash
}

// ValidatePassword validates password strength
func ValidatePassword(password string) error {
if password == "" {
return ErrEmptyPassword
}
if len(password) < 8 {
return ErrWeakPassword
}
return nil
}

// UserId - Value Object for user identifier
type UserId struct {
value string
}

// NewUserId creates a new UserId value object
func NewUserId(id string) (UserId, error) {
if id == "" {
return UserId{}, ErrEmptyUserId
}
return UserId{value: id}, nil
}

// Value returns the string representation of user id
func (u UserId) Value() string {
return u.value
}

// Equals compares two UserId value objects
func (u UserId) Equals(other UserId) bool {
return u.value == other.value
}

// String implements Stringer interface
func (u UserId) String() string {
return u.value
}

// Token - Value Object for authentication token
type Token struct {
value string
}

// NewToken creates a new Token value object
func NewToken(token string) (Token, error) {
if token == "" {
return Token{}, ErrInvalidToken
}
return Token{value: token}, nil
}

// Value returns the string representation of token
func (t Token) Value() string {
return t.value
}

// Equals compares two Token value objects
func (t Token) Equals(other Token) bool {
return t.value == other.value
}

// String implements Stringer interface
func (t Token) String() string {
return t.value
}

// Helper function to validate email format
func isValidEmail(email string) bool {
// Simple email regex pattern
pattern := `^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`
re := regexp.MustCompile(pattern)
return re.MatchString(email)
}
