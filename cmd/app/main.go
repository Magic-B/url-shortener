package main

import (
	"fmt"
	"log"
	"log/slog"
	"os"

	"github.com/Magic-B/url-shortener/internal/config"
	httplogger "github.com/Magic-B/url-shortener/internal/http/middleware/logger"
	"github.com/Magic-B/url-shortener/internal/storage/sqlite"
	"github.com/Magic-B/url-shortener/pkg/logger/slg"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/joho/godotenv"
)

const (
	local = "local"
	dev   = "dev"
	prod  = "prod"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("cannot load envs")
	}

	cfg := config.MustLoad()

	logger := NewLogger(cfg.Env)

	fmt.Println(logger)

	storage, err := sqlite.New(cfg.StoragePath)
	_ = storage

	if err != nil {
		logger.Error("failed to init storage", slg.Error(err))
		os.Exit(1)
	}

	if err != nil {
		logger.Error("failed to save url", slg.Error(err))
	}

	router := chi.NewRouter()

	router.Use(middleware.RequestID)
	router.Use(httplogger.New(logger))
	router.Use(middleware.Recoverer)
	router.Use(middleware.URLFormat)

	//TODO: Run server
}

func NewLogger(env string) *slog.Logger {
	var logger *slog.Logger

	switch env {
	case local:
		logger = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case dev:
		logger = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case prod:
		logger = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	}

	return logger
}
