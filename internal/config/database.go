package config

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

// -------------------------------------------------------------
// ensureDatabase() → Auto-creates DB if missing (when connected to postgres default DB)
// -------------------------------------------------------------
// func ensureDatabase(cfg *config.Database) error {
// 	connStr := fmt.Sprintf(
// 		"host=%s port=%d user=%s password=%s dbname=postgres-rest-api sslmode=%s",
// 		cfg.Database.Host,
// 		cfg.Database.Port,
// 		cfg.Database.User,
// 		cfg.Database.Pass,
// 		cfg.Database.SSLMode,
// 	)

// 	db, err := sql.Open("pgx", connStr)
// 	if err != nil {
// 		return fmt.Errorf("failed to connect to postgres system DB: %w", err)
// 	}
// 	defer db.Close()

// 	query := fmt.Sprintf("CREATE DATABASE %s;", cfg.Database.Name)
// 	_, err = db.Exec(query)
// 	if err != nil && !strings.Contains(err.Error(), "already exists") {
// 		return fmt.Errorf("failed to create database: %w", err)
// 	}

// 	return nil
// }

func NewDBPool(ctx context.Context, cfg *Config) (*pgxpool.Pool, error) {
	dsn := fmt.Sprintf("postgres://%s:%s@%s:%d/%s",
		cfg.Database.User, cfg.Database.Pass, cfg.Database.Host, cfg.Database.Port, cfg.Database.Name)

	pcfg, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return nil, fmt.Errorf("parse db config: %w", err)
	}
	pcfg.MaxConns = cfg.Database.MaxConns
	pcfg.MinConns = cfg.Database.MinConns
	pcfg.MaxConnLifetime = time.Duration(cfg.Database.MaxConnLifetimeMinutes) * time.Minute

	return pgxpool.NewWithConfig(ctx, pcfg)
}
