// @title Mytheresa
// @version 1.0
// @description The API used to interact with the catalog, including products, variants and categories.
// @host localhost:8080
// @BasePath /
package main

import (
	"context"
	"flag"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/joho/godotenv"
	"github.com/mytheresa/go-hiring-challenge/app/database"
	"github.com/mytheresa/go-hiring-challenge/app/logger"
	"github.com/mytheresa/go-hiring-challenge/cmd/server"
)

var (
	level = flag.String("level", "debug", "the level of messages to be logged (debug|info|warning|error)")
)

func main() {
	flag.Parse()

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	// initialize custom logger
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		AddSource: false, // nice to have when looking through logs on a production service
		Level:     logger.From(*level),
	}))

	if err := godotenv.Load(".env"); err != nil {
		logger.Error("failed to load .env file", slog.Any("error", err))

		os.Exit(1)
	}

	// initialize database connection
	db, err := database.New(&database.DatabaseOptions{
		User:     os.Getenv("POSTGRES_USER"),
		Password: os.Getenv("POSTGRES_PASSWORD"),
		Name:     os.Getenv("POSTGRES_DB"),
		Port:     os.Getenv("POSTGRES_PORT"),
	})
	if err != nil {
		os.Exit(1)
	}

	defer db.Close()

	service := server.New(&server.ServerOpts{Host: os.Getenv("MYTHERESA_APP_HOST"), Port: os.Getenv("MYTHERESA_APP_PORT"), Database: db, Logger: logger})
	service.Start()

	select {
	case <-ctx.Done():
		logger.Info("shutting down server...")
		service.Stop()
	case err := <-service.Error():
		logger.Error("failed to start server", slog.Any("error", err))
	}
}
