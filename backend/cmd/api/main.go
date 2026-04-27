package main

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/LeksinMaksim/haven/internal/config"
	"github.com/LeksinMaksim/haven/internal/server"
	"github.com/LeksinMaksim/haven/internal/storage"
	"github.com/LeksinMaksim/haven/internal/todo"
)

func main() {
	slog.SetDefault(slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	})))

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	if err := run(ctx); err != nil {
		slog.Error("startup failed", "error", err)
		os.Exit(1)
	}
}

func run(ctx context.Context) error {
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("load config: %w", err)
	}

	db, err := storage.NewPostgres(ctx, cfg.Postgres.DSN())
	if err != nil {
		return fmt.Errorf("connect postgres: %w", err)
	}
	defer db.Close()

	rdb, err := storage.NewRedis(cfg.Redis.Addr(), cfg.Redis.Password)
	if err != nil {
		return fmt.Errorf("connect redis: %w", err)
	}
	defer rdb.Close()

	slog.Info("connected to databases")

	todoRepo := todo.NewPostgresRepo(db)
	todoSvc := todo.NewService(todoRepo)
	todoHandler := todo.NewHandler(todoSvc)

	srv := server.New(cfg.App.Port)
	todoHandler.RegisterRoutes(srv.Mux())

	srv.Mux().HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"status":"ok"}`))
	})

	errCh := make(chan error, 1)

	go func() {
		slog.Info("server started", "port", cfg.App.Port)
		if err := srv.Start(); !errors.Is(err, http.ErrServerClosed) {
			errCh <- err
		}
	}()

	select {
	case err := <-errCh:
		return fmt.Errorf("server error: %w", err)
	case <-ctx.Done():
		slog.Info("shutting down...")
	}

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	return srv.Shutdown(shutdownCtx)
}
