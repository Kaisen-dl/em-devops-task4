package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"

	"github.com/Kaisen-dl/em-devops-task4/internal/config"
	"github.com/Kaisen-dl/em-devops-task4/internal/db"
	"github.com/Kaisen-dl/em-devops-task4/internal/handlers"
	"github.com/Kaisen-dl/em-devops-task4/internal/logger"
	"github.com/Kaisen-dl/em-devops-task4/internal/metrics"
)

func main() {
	log := logger.New()
	if err := run(log); err != nil {
		log.Error("app failed", slog.String("err", err.Error()))
		os.Exit(1)
	}
}

func run(log *slog.Logger) error {
	cfg := config.Load()
	ctx := context.Background()

	pg, err := db.NewPostgres(ctx, cfg.PostgresDSN)
	if err != nil {
		return fmt.Errorf("postgres: %w", err)
	}
	defer pg.Close()

	rdb, err := db.NewRedis(ctx, cfg.RedisAddr, cfg.RedisDB)
	if err != nil {
		return fmt.Errorf("redis: %w", err)
	}
	defer func() {
		if err := rdb.Close(); err != nil {
			log.Error("redis close failed", slog.String("err", err.Error()))
		}
	}()

	h := handlers.New(pg, rdb)

	mux := http.NewServeMux()
	mux.HandleFunc("/cache/set", h.CacheSet)
	mux.HandleFunc("/cache/get", h.CacheGet)
	mux.HandleFunc("/users", h.Users)
	mux.HandleFunc("/healthz", h.Healthz)
	mux.HandleFunc("/readyz", h.Readyz)
	mux.Handle("/metrics", metrics.Handler())

	handler := metrics.Instrument(handlers.Logging(log)(mux))

	addr := ":" + cfg.HTTPPort
	log.Info("server starting", slog.String("addr", addr))
	if err := http.ListenAndServe(addr, handler); err != nil {
		return fmt.Errorf("listen: %w", err)
	}
	return nil
}