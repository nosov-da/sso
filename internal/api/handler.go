package api

import (
"encoding/json"
"log/slog"
"net/http"
"sso/domain/user"
"sso/internal/middleware"
appsvc "sso/application/user"
)

type AuthHandler struct {
registerUC      *appsvc.RegisterUserUseCase
authenticateUC  *appsvc.AuthenticateUserUseCase
logger          *slog.Logger
}

func NewAuthHandler(registerUC *appsvc.RegisterUserUseCase, authenticateUC *appsvc.AuthenticateUserUseCase) *AuthHandler {
return &AuthHandler{
registerUC:     registerUC,
authenticateUC: authenticateUC,
logger:         slog.Default(),
}
}

func NewAuthHandlerWithLogger(registerUC *appsvc.RegisterUserUseCase, authenticateUC *appsvc.AuthenticateUserUseCase, logger *slog.Logger) *AuthHandler {
return &AuthHandler{
registerUC:     registerUC,
authenticateUC: authenticateUC,
logger:         logger,
}
}

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
var req struct {
Email    string `json:"email"`
Password string `json:"password"`
}

if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
h.logger.Error("failed to decode register request", "error", err)
middleware.HandleError(w, http.StatusBadRequest, err)
return
}

// Validate input
if req.Email == "" || req.Password == "" {
h.logger.Warn("invalid register request", "email_empty", req.Email == "", "password_empty", req.Password == "")
middleware.HandleError(w, http.StatusBadRequest, user.ErrEmptyEmail)
return
}

input := appsvc.RegisterInput{
Email:    req.Email,
Password: req.Password,
}

if err := h.registerUC.Execute(input); err != nil {
h.logger.Warn("registration failed", "email", req.Email, "error", err)

statusCode := http.StatusBadRequest
if err == user.ErrUserAlreadyExists {
statusCode = http.StatusConflict
}

middleware.HandleError(w, statusCode, err)
return
}

h.logger.Info("user registered", "email", req.Email)
w.Header().Set("Content-Type", "application/json")
w.WriteHeader(http.StatusCreated)
_ = json.NewEncoder(w).Encode(map[string]string{"message": "user registered successfully"})
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
var req struct {
Email    string `json:"email"`
Password string `json:"password"`
}

if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
h.logger.Error("failed to decode login request", "error", err)
middleware.HandleError(w, http.StatusBadRequest, err)
return
}

// Validate input
if req.Email == "" || req.Password == "" {
h.logger.Warn("invalid login request", "email_empty", req.Email == "", "password_empty", req.Password == "")
middleware.HandleError(w, http.StatusBadRequest, user.ErrEmptyEmail)
return
}

input := appsvc.AuthenticateInput{
Email:    req.Email,
Password: req.Password,
}

output, err := h.authenticateUC.Execute(input)
if err != nil {
h.logger.Warn("authentication failed", "email", req.Email, "error", err)
middleware.HandleError(w, http.StatusUnauthorized, err)
return
}

h.logger.Info("user authenticated", "email", req.Email)
w.Header().Set("Content-Type", "application/json")
w.WriteHeader(http.StatusOK)
_ = json.NewEncoder(w).Encode(map[string]string{"token": output.Token})
}
