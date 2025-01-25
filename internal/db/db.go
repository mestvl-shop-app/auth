package db

import (
	"fmt"
	"time"

	"github.com/mestvl-shop-app/auth/internal/config"

	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	_ "github.com/lib/pq"
)

func New(cfg config.Database) (*sqlx.DB, error) {
	location, err := time.LoadLocation(cfg.TimeZone)
	if err != nil {
		return nil, fmt.Errorf("time load location failed: %v", err)
	}

	connStr := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s timezone=%s",
		cfg.Host,
		cfg.Port,
		cfg.User,
		cfg.Password,
		cfg.DBName,
		cfg.SSLMode,
		location.String(),
	)

	dbConn, err := sqlx.Connect("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("db connection failed: %v", err)
	}

	dbConn.SetMaxIdleConns(cfg.MaxIdleConnections)
	dbConn.SetMaxOpenConns(cfg.MaxOpenConnections)

	if err := dbConn.Ping(); err != nil {
		return nil, err
	}

	return dbConn, nil
}

func IsDuplicate(err error) bool {
	if pqErr, ok := err.(*pq.Error); ok {
		if pqErr.Code.Name() == "unique_violation" {
			return true
		}
	}

	return false
}
