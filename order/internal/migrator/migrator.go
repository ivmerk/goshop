package migrator

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/pressly/goose/v3"
)

type Migrator struct {
	db            *sql.DB
	migrationsDir string
}

func NewMigrator(db *sql.DB, migrationsDir string) *Migrator {
	return &Migrator{
		db:            db,
		migrationsDir: migrationsDir,
	}
}

func (m *Migrator) Up() error {

	if _, err := os.Stat(m.migrationsDir); err != nil {
		return fmt.Errorf("migrations dir %q: %w", m.migrationsDir, err)
	}

	err := goose.SetDialect("postgres")
	if err != nil {
		return err
	}
	err = goose.Up(m.db, m.migrationsDir)
	if err != nil {
		return err
	}

	return nil
}
