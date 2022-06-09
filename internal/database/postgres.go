package database

import (
	"fmt"

	"github.com/Pythonyan3/payment-service/config"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type PostgresDB struct {
	*sqlx.DB
}

func NewPostgresDB(cfg *config.Config) (*PostgresDB, error) {
	// create connection to DB
	db, err := sqlx.Open("postgres", fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=%s",
		cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBName, cfg.DBPassword, cfg.DBSSLMode))
	// check successfully connection
	if err != nil {
		return nil, err
	}
	// check #2 trying to ping
	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return &PostgresDB{db}, nil
}
