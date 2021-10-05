package main

import (
	"database/sql"
	"embed"
	"fmt"
	"net/http"
	"os"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/sqlite"
	"github.com/golang-migrate/migrate/v4/source/httpfs"
	_ "modernc.org/sqlite"
)

func main() {
	db, err := sql.Open("sqlite", "file:example.db")
	check(err, "failed to open sqlite connection")
	check(ensureSchema(db), "failed to migrate database")
}

//go:embed migrations
var migrations embed.FS

const schemaVersion = 1

func ensureSchema(db *sql.DB) error {
	sourceInstance, err := httpfs.New(http.FS(migrations), "migrations")
	if err != nil {
		return fmt.Errorf("invalid source instance, %w", err)
	}
	targetInstance, err := sqlite.WithInstance(db, new(sqlite.Config))
	if err != nil {
		return fmt.Errorf("invalid target sqlite instance err, %w", err)
	}
	m, err := migrate.NewWithInstance("httpfs", sourceInstance,
		"sqlite", targetInstance)
	if err != nil {
		return fmt.Errorf("failed to initialize migrate instance, %w", err)
	}
	err = m.Migrate(schemaVersion)
	if err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("migration failed, %w", err)
	}
	return sourceInstance.Close()
}

func check(err error, msg string) {
	if err != nil {
		fmt.Printf("%v, err: %v\n", msg, err)
		os.Exit(1)
	}
}
