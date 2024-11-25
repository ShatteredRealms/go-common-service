package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/ShatteredRealms/go-common-service/pkg/config"
	"github.com/ShatteredRealms/go-common-service/pkg/log"
	"github.com/ShatteredRealms/go-common-service/pkg/repository/cacher"
	"github.com/go-gorm/caches/v4"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/stdlib"
	_ "github.com/lib/pq"
	"github.com/uptrace/opentelemetry-go-extra/otelgorm"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/plugin/dbresolver"
)

// ConnectDB Initializes the connection to a Postgres database
func ConnectDB(ctx context.Context, pgPool config.DBPoolConfig, redisPool config.DBPoolConfig) (*gorm.DB, error) {
	database, err := gorm.Open(postgres.Open(pgPool.Master.PostgresDSNWithoutName()), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	if database.Exec(fmt.Sprintf("SELECT 1 FROM pg_database WHERE datname = '%s';", pgPool.Master.Name)).RowsAffected == 0 {
		atabase.Exec(fmt.Sprintf("CREATE DATABASE %s;", pgPool.Master.Name))
	}

	conf, err := pgx.ParseConfig(pgPool.Master.PostgresDSN())
	if err != nil {
		return nil, fmt.Errorf("parse config: %w", err)
	}
	sqlDB := stdlib.OpenDB(*conf)

	sqlDB.SetConnMaxLifetime(time.Second)
	sqlDB.SetMaxOpenConns(0)
	sqlDB.SetMaxIdleConns(10)
	db, err := gorm.Open(postgres.New(postgres.Config{
		Conn: sqlDB,
	}), &gorm.Config{
		Logger: logger.New(
			log.Logger,
			logger.Config{
				SlowThreshold:             time.Millisecond * 500,
				Colorful:                  true,
				IgnoreRecordNotFoundError: true,
				ParameterizedQueries:      true,
			},
		),
	})

	if err != nil {
		return nil, fmt.Errorf("gorm: %w", err)
	}

	if len(pgPool.Slaves) > 0 {
		replicas := make([]gorm.Dialector, len(pgPool.Slaves))
		for _, slave := range pgPool.Slaves {
			replicas = append(replicas, postgres.Open(slave.PostgresDSN()))
		}

		err = db.Use(dbresolver.Register(dbresolver.Config{
			Replicas: replicas,
			Policy:   dbresolver.RandomPolicy{},
		}))

		if err != nil {
			return nil, fmt.Errorf("db replica resolver: %w", err)
		}
	}

	if err := db.Use(otelgorm.NewPlugin(otelgorm.WithDBName(pgPool.Master.Name))); err != nil {
		return nil, fmt.Errorf("opentelemetry: %w", err)
	}

	c, err := cacher.NewRedisCache(ctx, redisPool)
	if err != nil {
		return nil, fmt.Errorf("redis cache: %w", err)
	}
	cachesPlugin := caches.Caches{
		Conf: &caches.Config{
			Easer:  true,
			Cacher: c,
		},
	}
	if err = db.Use(&cachesPlugin); err != nil {
		return nil, fmt.Errorf("cacher: %w", err)
	}

	return db, err
}

