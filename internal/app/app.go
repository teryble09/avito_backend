package app

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
	"github.com/teryble09/avito_backend/tree/dev/internal/config"
	"github.com/teryble09/avito_backend/tree/dev/migrations"
)

type App struct {
	Server *http.Server

	DB     *pgxpool.Pool
	Logger *slog.Logger
	Config *config.Config
}

func New(cfg *config.Config, logger *slog.Logger) (*App, error) {
	// Если не можем подключиться к бд достаточно быстро, значит проблемы
	shortCtx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	db, err := pgxpool.New(shortCtx, cfg.DatabaseURL)
	if err != nil {
		return nil, fmt.Errorf("db init: %w", err)
	}
    
    goose.SetBaseFS(migrations.FS)
    
    if err := goose.SetDialect("postgres"); err != nil {
		return nil, fmt.Errorf("goose set dialect: %w", err)
    }
	
    dbconn := stdlib.OpenDBFromPool(db)
    if err := goose.Up(dbconn, "."); err != nil {
		return nil, fmt.Errorf("run migrations: %w", err)
    }
    
    if err := dbconn.Close(); err != nil {
		return nil, fmt.Errof("close connection after migrations: %w", err)
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

	err := app.Server.Shutdown(ctx)
	if err != nil {
		app.Logger.ErrorContext(ctx, "http server shutdown",
			slog.String("error", err.Error()),
		)
	} else {
		app.Logger.InfoContext(ctx, "http server stopped")
	}

	app.DB.Close()
	app.Logger.InfoContext(ctx, "database pool closed")

	app.Logger.InfoContext(ctx, "shutdown completed successfully")

	return nil

