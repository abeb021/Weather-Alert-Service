package app

import (
	"auth-service/config"
	"auth-service/internal/logger"
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

func Run(logger *logger.Log, cfg *config.Config) {
	container := NewContainer(logger, cfg)
	container.Logger.Logger.Info("application initialized successfully")

	mux := http.NewServeMux()
	mux.HandleFunc("/register", container.Handler.RegisterHandler)
	mux.HandleFunc("/login", container.Handler.LoginHandler)

	server := &http.Server{
		Addr:    cfg.Server.Addr,
		Handler: mux,
	}

	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			container.Logger.Logger.Error("Server ListenAndServe error", "error", err)
		}
	}()

	container.Logger.Logger.Info("Server started", "addr", cfg.Server.Addr)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	container.Logger.Logger.Info("Shutdown signal received")

	ctx, cancel := context.WithTimeout(context.Background(), cfg.Server.GracefulShutdownTimeout)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		container.Logger.Logger.Warn("Server forced to shutdown", "error", err)
	} else {
		container.Logger.Logger.Info("Server stoppe gracefully")
	}
}
