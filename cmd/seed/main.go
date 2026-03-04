package main

import (
	"log/slog"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/joho/godotenv"

	"github.com/mytheresa/go-hiring-challenge/app/database"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))

	// load environment variables from .env file
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
		Logger:   logger,
	})
	if err != nil {
		logger.Error("failed to open database connection", slog.Any("error", err))
		os.Exit(1)
	}

	defer db.Close()

	dir := os.Getenv("POSTGRES_SQL_DIR")
	files, err := os.ReadDir(dir)
	if err != nil {
		logger.Error("failed to read directory", slog.Any("error", err))
		os.Exit(1)
	}

	// Filter and sort .sql files
	var sqlFiles []os.DirEntry
	for _, file := range files {
		if !file.IsDir() && strings.HasSuffix(file.Name(), ".sql") {
			sqlFiles = append(sqlFiles, file)
		}
	}
	sort.Slice(sqlFiles, func(i, j int) bool {
		return sqlFiles[i].Name() < sqlFiles[j].Name()
	})

	for _, file := range sqlFiles {
		path := filepath.Join(dir, file.Name())

		// skip over specific files
		if strings.HasPrefix(file.Name(), "#") {
			logger.Info("skipping " + file.Name() + " successfully")
			continue
		}

		content, err := os.ReadFile(path)
		if err != nil {
			logger.Error("failed to read .sql file", slog.String("file", file.Name()), slog.Any("error", err))
		}

		sql := string(content)
		if err := db.GORM().Exec(sql).Error; err != nil {
			logger.Error("failed to execute SQL instruction", slog.String("file", file.Name()), slog.Any("error", err))
			return
		}

		logger.Info("SQL file execute successfully", slog.String("file", file.Name()))
	}
}
