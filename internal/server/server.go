package server

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	appsvc "sso/application/user"
	"sso/internal/api"
	"sso/internal/config"
	"sso/internal/middleware"
)

func Run(cfg *config.Config, authSvc interface{}) error {
	// Legacy support - not used anymore
	return nil
}

// RunWithUseCases starts the server with DDD use cases
func RunWithUseCases(cfg *config.Config, registerUC *appsvc.RegisterUserUseCase, authenticateUC *appsvc.AuthenticateUserUseCase) error {
	logger := slog.Default()
	return RunWithUseCasesAndLogger(cfg, registerUC, authenticateUC, logger)
}

// RunWithUseCasesAndLogger starts the server with DDD use cases and custom logger
func RunWithUseCasesAndLogger(cfg *config.Config, registerUC *appsvc.RegisterUserUseCase, authenticateUC *appsvc.AuthenticateUserUseCase, logger *slog.Logger) error {
	h := api.NewAuthHandlerWithLogger(registerUC, authenticateUC, logger)

	mux := http.NewServeMux()

	// Apply middleware
	var handler http.Handler = mux
	handler = middleware.CORSMiddleware()(handler)
	handler = middleware.LoggingMiddleware(logger)(handler)

	// Register handlers
	mux.HandleFunc("POST /register", h.Register)
	mux.HandleFunc("POST /login", h.Login)
	mux.HandleFunc("GET /openapi.yaml", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "docs/openapi.yaml")
	})
	mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"status":"ok"}`))
	})

	srv := &http.Server{
		Addr:         cfg.Addr,
		Handler:      handler,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Start server in goroutine
	go func() {
		logger.Info("server starting", "addr", cfg.Addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Error("server error", "error", err)
		}
	}()

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("shutting down server")

	// Graceful shutdown with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		logger.Error("server shutdown error", "error", err)
		return err
	}

	logger.Info("server stopped")
	return nil
}
