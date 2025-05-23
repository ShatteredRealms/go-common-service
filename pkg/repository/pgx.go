package repository

import (
	"context"
	"fmt"

	"github.com/ShatteredRealms/go-common-service/pkg/log"
	"github.com/cenkalti/backoff/v4"
	"github.com/exaring/otelpgx"
	"github.com/golang-migrate/migrate/v4"
	migratepgx "github.com/golang-migrate/migrate/v4/database/pgx/v5"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
	pgxuuid "github.com/vgarvardt/pgx-google-uuid/v5"
)

type PgxMigrater struct {
	migrate *migrate.Migrate
	Conn    *pgxpool.Pool
}

func NewPgxMigrater(ctx context.Context, pgpoolUrl string, migrationPath string, autoMigrate bool) (*PgxMigrater, error) {
	migrater := &PgxMigrater{
		migrate: &migrate.Migrate{
			Log: &MigrateLogger{},
		},
	}
	pgConfig, err := pgxpool.ParseConfig(pgpoolUrl)
	if err != nil {
		return nil, fmt.Errorf("parsing pool: %w", err)
	}

	pgConfig.ConnConfig.Tracer = otelpgx.NewTracer()
	pgConfig.AfterConnect = func(ctx context.Context, conn *pgx.Conn) error {
		pgxuuid.Register(conn.TypeMap())
		return nil
	}

	migrater.Conn, err = pgxpool.NewWithConfig(context.Background(), pgConfig)
	if err != nil {
		return nil, fmt.Errorf("connecting pool: %w", err)
	}

	err = backoff.Retry(func() error {
		return migrater.Conn.Ping(ctx)
	}, backoff.NewExponentialBackOff())
	if err != nil {
		return nil, fmt.Errorf("pg pool not available: %w", err)
	}

	driver, err := migratepgx.WithInstance(stdlib.OpenDBFromPool(migrater.Conn), &migratepgx.Config{})
	if err != nil {
		return nil, fmt.Errorf("creating migrate driver: %w", err)
	}

	migrater.migrate, err = migrate.NewWithDatabaseInstance(fmt.Sprintf("file://%s", migrationPath), "postgres", driver)
	if autoMigrate {
		if err != nil {
			return nil, fmt.Errorf("creating migrate instance: %w", err)
		}

		err = migrater.migrate.Up()
		if err != nil && err != migrate.ErrNoChange {
			return nil, fmt.Errorf("migrating: %w", err)
		}
	}

	return migrater, nil
}

type MigrateLogger struct{}

func (m *MigrateLogger) Printf(format string, v ...interface{}) {
	log.Logger.Infof(format, v...)
}

func (m *MigrateLogger) Verbose() bool {
	return false
}
