package app

import (
	"auth-service/internal/logger"
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func Run(logger *logger.Log) {
	container := NewContainer(logger)
	container.Logger.Logger.Info("application initialized successfully")

	mux := http.NewServeMux()
	mux.HandleFunc("/register", container.Handler.RegistrationHandler)
	mux.HandleFunc("/login", container.Handler.LoginHandler)

	server := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			container.Logger.Logger.Error("Server ListenAndServe error", err)
		}
	}()

	container.Logger.Logger.Info("Server started on :8080")

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	container.Logger.Logger.Info("Shutdown signal received")

	// потом 10сек вынести в конфиг файл
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		container.Logger.Logger.Warn("Server forced to shutdown", err)
	} else {
		container.Logger.Logger.Info("Server stoppe gracefully")
	}
}
