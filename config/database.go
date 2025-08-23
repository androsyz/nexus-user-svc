package config

import (
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	_ "github.com/pressly/goose/v3"
)

func NewDatabase(cfg Database) (*sqlx.DB, error) {
	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable TimeZone=Asia/Jakarta", cfg.Address, cfg.Port, cfg.Username, cfg.Password, cfg.DBName)
	return connectDB(connStr, cfg)
}

func connectDB(connStr string, cfg Database) (*sqlx.DB, error) {
	db, err := sqlx.Connect("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to postgres database: %w", err)
	}

	if cfg.MaxOpenConn > 0 {
		db.SetMaxOpenConns(cfg.MaxOpenConn)
	}
	if cfg.MaxIdleConn >= 0 {
		db.SetMaxIdleConns(cfg.MaxIdleConn)
	}
	if cfg.ConnMaxLifeTime > 0 {
		db.SetConnMaxLifetime(time.Second * time.Duration(cfg.ConnMaxLifeTime))
	}
	if cfg.ConnMaxIdleTime >= 0 {
		db.SetConnMaxIdleTime(time.Second * time.Duration(cfg.ConnMaxIdleTime))
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping DB: %w", err)
	}

	return db, nil
}
