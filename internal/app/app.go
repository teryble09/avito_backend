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
	db, err := pgxpool.New(context.Background(), cfg.DatabaseURL)
	if err != nil {
		return nil, fmt.Errorf("db init: %w", err)
	}
    
    goose.SetBaseFS(migrations.FS)
    
    if err := goose.SetDialect("postgres"); err != nil {
		return nil, fmt.Errorf("goose set dialect: %w", err)
    }

	// goose требует стандартный интерфейс
    dbconn := stdlib.OpenDBFromPool(db)
    if err := goose.Up(dbconn, "."); err != nil {
		return nil, fmt.Errorf("run migrations: %w", err)
    }
    
    if err := dbconn.Close(); err != nil {
		return nil, fmt.Errorf("close connection after migrations: %w", err)
    }

	handler, err := assemblyLayers(cfg, pgxpool, logger)
	if err != nil {
		return nil, fmt.Errorf("assembly layers: %w", err)
	}

	server := &http.Server{
		Handler: handler,

		Addr: cfg.ServerAddr,
		// опционально можем добавить другие параметры в конфиг и подставить тут
	}

	return &App{
		Logger: logger,
		DB:     db,
		Config: cfg,
		Server: server,
	}, nil
}

func (app *App) Run() error {
	if err := app.Server.ListenAndServe(); err != nil {
		return fmt.Errorf("run server: %w", err)
	}
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

