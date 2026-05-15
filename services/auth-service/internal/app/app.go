package app

import (
	"auth-service/config"
	"auth-service/internal/api/middleware"
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func Run(logger *slog.Logger, cfg *config.Config) {
	container := NewContainer(logger, cfg)
	logger.Info("application initialized successfully")

	mux := http.NewServeMux()
	mux.HandleFunc("/api/auth/register", container.Handler.RegisterHandler)
	mux.HandleFunc("/api/auth/login", container.Handler.LoginHandler)
	mux.HandleFunc("/api/auth/refresh", container.Handler.RefreshHandler)
	mux.HandleFunc("/api/auth/validate", container.Handler.ValidateHandler)
	mux.HandleFunc("/api/auth/health", container.Handler.HealthHandler)
	mux.Handle("/metrics", promhttp.Handler())

	publicPaths := map[string]struct{}{
		"/api/auth/register": {},
		"/api/auth/login":    {},
		"/api/auth/refresh":  {},
		"/api/auth/validate": {},
		"/api/auth/health":   {},
	}
	handler := middleware.RequestLogger(container.Logger)(
		middleware.Auth(container.JWTService, publicPaths)(mux),
	)

	server := &http.Server{
		Addr:    cfg.Server.Addr,
		Handler: handler,
	}

	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			container.Logger.Error("Server ListenAndServe error", "error", err)
		}
	}()

	container.Logger.Info("Server started", "addr", cfg.Server.Addr)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	container.Logger.Info("Shutdown signal received")

	ctx, cancel := context.WithTimeout(context.Background(), cfg.Server.GracefulShutdownTimeout)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		container.Logger.Warn("Server forced to shutdown", "error", err)
	} else {
		container.Logger.Info("Server stoppe gracefully")
	}
}
