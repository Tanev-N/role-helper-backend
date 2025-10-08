package service

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"role-helper/cfg"
)

func InitPostgres(cfg *cfg.Config) (*sql.DB, error) {
	dbConfig := cfg.Postgres
	dataConnection := fmt.Sprintf("host=%s port=%s dbname=%s user=%s password=%s sslmode=disable",
		dbConfig.IP, dbConfig.Port, dbConfig.DBname, dbConfig.User, dbConfig.Password)
	db, err := sql.Open("postgres", dataConnection)
	if err != nil {
		return nil, err
	}
	return db, nil
}
