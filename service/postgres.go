package service

import (
	"database/sql"
	"fmt"
	"role-helper/cfg"
	"time"

	_ "github.com/lib/pq"
)

func InitPostgres(cfg *cfg.Config) (*sql.DB, error) {
	dbConfig := cfg.Postgres
	dataConnection := fmt.Sprintf("host=%s port=%s dbname=%s user=%s password=%s sslmode=disable",
		dbConfig.IP, dbConfig.Port, dbConfig.DBname, dbConfig.User, dbConfig.Password)
	db, err := sql.Open("postgres", dataConnection)
	if err != nil {
		return nil, err
	}

	backoff := time.Second
	deadline := time.Now().Add(30 * time.Second)
	for {
		if err := db.Ping(); err == nil {
			return db, nil
		}
		if time.Now().After(deadline) {
			return nil, fmt.Errorf("postgres is not reachable within timeout")
		}
		time.Sleep(backoff)
		if backoff < 5*time.Second {
			backoff *= 2
		}
	}
}
