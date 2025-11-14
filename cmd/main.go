package cmd

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/joho/godotenv"
	"github.com/teryble09/avito_backend/tree/dev/internal/app"
	"github.com/teryble09/avito_backend/tree/dev/internal/config"
	"github.com/teryble09/avito_backend/tree/dev/internal/logger"
)

func main() {
	godotenv.Load()

	cfg := config.Load()

	logger, err := logger.Setup(cfg)
	if err != nil {
		slog.Error("logger setup",
			slog.String("error", err.Error()),
		)
		os.Exit(1)
	}

	app, err := app.New(cfg, logger)
	if err != nil {
		logger.Error("Fail to assembly app",
			slog.String("error", err.Error()),
		)
		os.Exit(1)
	}

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)

	serverErr := make(chan error, 1)
	go func() {
		app.Logger.Info("starting server",
			slog.String("addr", cfg.ServerAddr),
		)
		if err := app.Run(); err != nil {
			serverErr <- err
		}
	}()

	select {
	case err := <-serverErr:
		app.Logger.Error("server error",
			slog.String("error", err.Error()),
		)
		os.Exit(1)

	case sig := <-sigChan:
		app.Logger.Info("shutdown signal received",
			slog.String("signal", sig.String()),
		)
	}

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	app.Logger.Info("shutting down server gracefully")

	if err := app.Shutdown(shutdownCtx); err != nil {
		app.Logger.Error("forced shutdown",
			slog.String("error", err.Error()),
		)
		os.Exit(1)
	}

	app.Logger.Info("server stopped successfully")
}
