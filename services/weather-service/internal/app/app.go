package app

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"weather-service/config"
	"weather-service/internal/api/middleware"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func Run(logger *slog.Logger, cfg *config.Config) {
	container := NewContainer(logger, cfg)
	logger.Info("application initialized successfully")

	mux := http.NewServeMux()
	mux.HandleFunc("GET /api/weather/current", container.Handler.CurrentWeatherHandler)
	mux.HandleFunc("GET /api/weather/forecast", container.Handler.ForecastHandler)
	mux.HandleFunc("GET /api/auth/health", container.Handler.HealthHandler)
	mux.Handle("/metrics", promhttp.Handler())

	handler := middleware.RequestLogger(container.Logger)(mux)

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
