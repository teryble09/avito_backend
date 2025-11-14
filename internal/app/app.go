package app

import (
	"context"
	"log/slog"
	"net/http"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/teryble09/avito_backend/tree/dev/internal/config"
)

type App struct {
	Server *http.Server

	DB     *pgxpool.Pool
	Logger *slog.Logger
	Config *config.Config
}

func New(cfg *config.Config, logger *slog.Logger) (*App, error) {
	logger.Info(
		"Start app assembly",
	)

	// Если не можем подключиться к бд достаточно быстро, значит проблемы
	shortCtx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	db, err := pgxpool.New(shortCtx, cfg.DatabaseURL)
	if err != nil {
		logger.Error(
			"connect to the db",
			slog.String("error", err.Error()),
		)
	}

	return &App{
		Logger: logger,
		DB:     db,
		Config: cfg,
		// Server: ,
	}, nil
}

func (app *App) Run() error {
	return nil
}

func (app *App) Shutdown(ctx context.Context) error {
	app.Logger.InfoContext(ctx, "starting graceful shutdown")

	if err := app.Server.Shutdown(ctx); err != nil {
		app.Logger.ErrorContext(ctx, "http server shutdown error",
			slog.String("error", err.Error()),
		)
	} else {
		app.Logger.InfoContext(ctx, "http server stopped")
	}

	app.DB.Close()
	app.Logger.InfoContext(ctx, "database pool closed")

	app.Logger.InfoContext(ctx, "shutdown completed successfully")
	return nil
}
