package main

import (
	"database/sql"
	"fmt"
	"log/slog"
	"net/http"
	"os"

	"github.com/falasefemi2/goreact-boilerplate/internal/config"
	"github.com/falasefemi2/goreact-boilerplate/internal/server"
	_ "github.com/jackc/pgx/v5/stdlib"
)

func main() {
	cfg := config.Load()

	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	// Connect to database
	db, err := sql.Open("pgx", cfg.DatabaseURL)
	if err != nil {
		slog.Error("failed to open database", "error", err)
		os.Exit(1)
	}
	defer db.Close()

	// Verify connection is actually alive
	if err := db.Ping(); err != nil {
		slog.Error("failed to connect to database", "error", err)
		os.Exit(1)
	}

	slog.Info("database connected")
	slog.Info("starting server", "port", cfg.Port, "env", cfg.Env)

	srv := server.New(db, cfg.JWTSecret)

	if err := http.ListenAndServe(fmt.Sprintf(":%s", cfg.Port), srv); err != nil {
		slog.Error("server failed", "error", err)
		os.Exit(1)
	}
}
