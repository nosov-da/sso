package api

import (
	"encoding/json"
	"log"
	"net/http"
	"sso/domain/user"
	"sso/internal/middleware"
	appsvc "sso/application/user"
)

type AuthHandler struct {
	registerUC     *appsvc.RegisterUserUseCase
	authenticateUC *appsvc.AuthenticateUserUseCase
	logger         *log.Logger
}

func NewAuthHandler(registerUC *appsvc.RegisterUserUseCase, authenticateUC *appsvc.AuthenticateUserUseCase) *AuthHandler {
	return &AuthHandler{
		registerUC:     registerUC,
		authenticateUC: authenticateUC,
		logger:         log.Default(),
	}
}

func NewAuthHandlerWithLogger(registerUC *appsvc.RegisterUserUseCase, authenticateUC *appsvc.AuthenticateUserUseCase, logger *log.Logger) *AuthHandler {
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
		h.logger.Printf("failed to decode register request: %v", err)
		middleware.HandleError(w, http.StatusBadRequest, err)
		return
	}

	// Validate input
	if req.Email == "" || req.Password == "" {
		h.logger.Printf("invalid register request: email_empty=%v, password_empty=%v", req.Email == "", req.Password == "")
		middleware.HandleError(w, http.StatusBadRequest, user.ErrEmptyEmail)
		return
	}

	input := appsvc.RegisterInput{
		Email:    req.Email,
		Password: req.Password,
	}

	if err := h.registerUC.Execute(input); err != nil {
		h.logger.Printf("registration failed: email=%s, error=%v", req.Email, err)

		statusCode := http.StatusBadRequest
		if err == user.ErrUserAlreadyExists {
			statusCode = http.StatusConflict
		}

		middleware.HandleError(w, statusCode, err)
		return
	}

	h.logger.Printf("user registered: email=%s", req.Email)
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
		h.logger.Printf("failed to decode login request: %v", err)
		middleware.HandleError(w, http.StatusBadRequest, err)
		return
	}

	// Validate input
	if req.Email == "" || req.Password == "" {
		h.logger.Printf("invalid login request: email_empty=%v, password_empty=%v", req.Email == "", req.Password == "")
		middleware.HandleError(w, http.StatusBadRequest, user.ErrEmptyEmail)
		return
	}

	input := appsvc.AuthenticateInput{
		Email:    req.Email,
		Password: req.Password,
	}

	output, err := h.authenticateUC.Execute(input)
	if err != nil {
		h.logger.Printf("authentication failed: email=%s, error=%v", req.Email, err)
		middleware.HandleError(w, http.StatusUnauthorized, err)
		return
	}

	h.logger.Printf("user authenticated: email=%s", req.Email)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(map[string]string{"token": output.Token})
}
