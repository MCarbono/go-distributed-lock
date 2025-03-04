package database

import (
	"context"
	"database/sql"
	"distributed-lock/config"
	"embed"
	"errors"
	"fmt"
	"time"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

func OpenDB(cfg config.RelationalDatabase, migrations embed.FS) (*sqlx.DB, error) {
	db, err := sql.Open("postgres", fmt.Sprintf(
		"host=%s port=%s user=%s dbname=%s password=%s sslmode=disable connect_timeout=10",
		cfg.Host,
		cfg.Port,
		cfg.User,
		cfg.Name,
		cfg.Password),
	)
	if err != nil {
		return nil, fmt.Errorf("opening connection: %w", err)
	}

	err = Migrate(db, migrations)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = db.PingContext(ctx)
	if err != nil {
		return nil, fmt.Errorf("pinging: %w", err)
	}
	return sqlx.NewDb(db, "postgres"), nil
}

func Migrate(db *sql.DB, migrations embed.FS) error {
	fs, err := iofs.New(migrations, "migrations")
	if err != nil {
		return fmt.Errorf("creating iofs driver: %w", err)
	}
	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		return fmt.Errorf("creating driver: %w", err)
	}
	m, err := migrate.NewWithInstance("iofs", fs, "postgres", driver)
	if err != nil {
		return err
	}
	err = m.Up()
	if err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return fmt.Errorf("migrating up the database: %w", err)
	}
	return nil
}
