package api

import (
"bytes"
"encoding/json"
"log/slog"
"net/http"
"net/http/httptest"
"os"
appsvc "sso/application/user"
"sso/infrastructure/events"
inmemoryrep "sso/infrastructure/persistence"
"testing"
)

func createTestHandler() (*AuthHandler, *appsvc.RegisterUserUseCase, *appsvc.AuthenticateUserUseCase) {
repo := inmemoryrep.NewInMemoryUserRepository()
publisher := events.NewInMemoryEventPublisher()

mockTokenService := &mockTokenService{token: "test-token-123"}

registerUC := appsvc.NewRegisterUserUseCase(repo, publisher)
authenticateUC := appsvc.NewAuthenticateUserUseCase(repo, mockTokenService, publisher)

logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelError}))
handler := NewAuthHandlerWithLogger(registerUC, authenticateUC, logger)
return handler, registerUC, authenticateUC
}

func TestAuthHandler_Register_Success(t *testing.T) {
handler, _, _ := createTestHandler()

body, _ := json.Marshal(map[string]string{
"email":    "user@example.com",
"password": "password123",
})

req := httptest.NewRequest("POST", "/register", bytes.NewReader(body))
w := httptest.NewRecorder()

handler.Register(w, req)

if w.Code != http.StatusCreated {
t.Errorf("Status code = %d, want %d", w.Code, http.StatusCreated)
}
}

func TestAuthHandler_Register_Duplicate(t *testing.T) {
handler, registerUC, _ := createTestHandler()

email := "duplicate@example.com"
password := "password123"

// Register first user
registerUC.Execute(appsvc.RegisterInput{Email: email, Password: password})

// Try to register again
body, _ := json.Marshal(map[string]string{
"email":    email,
"password": password,
})

req := httptest.NewRequest("POST", "/register", bytes.NewReader(body))
w := httptest.NewRecorder()

handler.Register(w, req)

if w.Code != http.StatusConflict {
t.Errorf("Status code = %d, want %d", w.Code, http.StatusConflict)
}
}

func TestAuthHandler_Register_InvalidJSON(t *testing.T) {
handler, _, _ := createTestHandler()

req := httptest.NewRequest("POST", "/register", bytes.NewReader([]byte("invalid json")))
w := httptest.NewRecorder()

handler.Register(w, req)

if w.Code != http.StatusBadRequest {
t.Errorf("Status code = %d, want %d", w.Code, http.StatusBadRequest)
}
}

func TestAuthHandler_Login_Success(t *testing.T) {
handler, registerUC, _ := createTestHandler()

email := "login@example.com"
password := "password123"

// Register user first
registerUC.Execute(appsvc.RegisterInput{Email: email, Password: password})

body, _ := json.Marshal(map[string]string{
"email":    email,
"password": password,
})

req := httptest.NewRequest("POST", "/login", bytes.NewReader(body))
w := httptest.NewRecorder()

handler.Login(w, req)

if w.Code != http.StatusOK {
t.Errorf("Status code = %d, want %d", w.Code, http.StatusOK)
}

var response map[string]string
json.NewDecoder(w.Body).Decode(&response)
if response["token"] == "" {
t.Error("Token should not be empty")
}
}

func TestAuthHandler_Login_WrongPassword(t *testing.T) {
handler, registerUC, _ := createTestHandler()

email := "user@example.com"
password := "password123"

// Register user first
registerUC.Execute(appsvc.RegisterInput{Email: email, Password: password})

body, _ := json.Marshal(map[string]string{
"email":    email,
"password": "wrongpassword",
})

req := httptest.NewRequest("POST", "/login", bytes.NewReader(body))
w := httptest.NewRecorder()

handler.Login(w, req)

if w.Code != http.StatusUnauthorized {
t.Errorf("Status code = %d, want %d", w.Code, http.StatusUnauthorized)
}
}

func TestAuthHandler_Login_NonExistentUser(t *testing.T) {
handler, _, _ := createTestHandler()

body, _ := json.Marshal(map[string]string{
"email":    "nonexistent@example.com",
"password": "password123",
})

req := httptest.NewRequest("POST", "/login", bytes.NewReader(body))
w := httptest.NewRecorder()

handler.Login(w, req)

if w.Code != http.StatusUnauthorized {
t.Errorf("Status code = %d, want %d", w.Code, http.StatusUnauthorized)
}
}

func TestAuthHandler_Login_InvalidJSON(t *testing.T) {
handler, _, _ := createTestHandler()

req := httptest.NewRequest("POST", "/login", bytes.NewReader([]byte("invalid json")))
w := httptest.NewRecorder()

handler.Login(w, req)

if w.Code != http.StatusBadRequest {
t.Errorf("Status code = %d, want %d", w.Code, http.StatusBadRequest)
}
}

func TestAuthHandler_LoginResponseContentType(t *testing.T) {
handler, registerUC, _ := createTestHandler()

email := "user@example.com"
password := "password123"

registerUC.Execute(appsvc.RegisterInput{Email: email, Password: password})

body, _ := json.Marshal(map[string]string{
"email":    email,
"password": password,
})

req := httptest.NewRequest("POST", "/login", bytes.NewReader(body))
w := httptest.NewRecorder()

handler.Login(w, req)

contentType := w.Header().Get("Content-Type")
if contentType != "application/json" {
t.Errorf("Content-Type = %s, want application/json", contentType)
}
}

func TestAuthHandler_RegisterAndLoginFlow(t *testing.T) {
handler, _, _ := createTestHandler()

email := "testflow@example.com"
password := "password123"

// Register
body, _ := json.Marshal(map[string]string{
"email":    email,
"password": password,
})
req := httptest.NewRequest("POST", "/register", bytes.NewReader(body))
w := httptest.NewRecorder()
handler.Register(w, req)

if w.Code != http.StatusCreated {
t.Errorf("Register status = %d, want %d", w.Code, http.StatusCreated)
}

// Login
body, _ = json.Marshal(map[string]string{
"email":    email,
"password": password,
})
req = httptest.NewRequest("POST", "/login", bytes.NewReader(body))
w = httptest.NewRecorder()
handler.Login(w, req)

if w.Code != http.StatusOK {
t.Errorf("Login status = %d, want %d", w.Code, http.StatusOK)
}

var response map[string]string
json.NewDecoder(w.Body).Decode(&response)
if response["token"] == "" {
t.Error("Token should not be empty in login response")
}
}

type mockTokenService struct {
token string
err   error
}

func (m *mockTokenService) GenerateToken(userID string) (string, error) {
if m.err != nil {
return "", m.err
}
return m.token, nil
}

func (m *mockTokenService) VerifyToken(token string) (string, error) {
return "test-user-id", nil
}
