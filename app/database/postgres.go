package database

import (
	"database/sql"
	"fmt"
	"log/slog"

	_ "github.com/lib/pq"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// Database here is a struct but it's usually an interface where the actual operations are performed by
// whatever implements it.
type Database struct {
	db     *gorm.DB
	idb    *sql.DB
	logger *slog.Logger
}

type DatabaseOptions struct {
	User     string
	Password string
	Name     string // the name of the database
	Port     string
	Logger   *slog.Logger
}

// Creates a new database instance
func New(opts *DatabaseOptions) (*Database, error) {
	dsn := fmt.Sprintf("postgres://%s:%s@localhost:%s/%s?sslmode=disable", opts.User, opts.Password, opts.Port, opts.Name)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{}) //TODO(rafael.nunes): set up the logger for GORM to use the app's logger (though it needs to implement the logger.Interface)
	if err != nil {
		opts.Logger.Error("failed to connect to database", slog.Any("error", err))
		return nil, err
	}

	sqlDB, err := db.DB()
	if err != nil {
		opts.Logger.Error("Failed to acquire database object", slog.Any("error", err))
		return nil, err
	}

	d := &Database{
		db:     db,
		idb:    sqlDB,
		logger: opts.Logger,
	}

	return d, nil
}

func (d *Database) Close()         { d.idb.Close() }
func (d *Database) GORM() *gorm.DB { return d.db }
