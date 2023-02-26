package db

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
	"github.com/mrbttf/bridge-server/pkg/config"
	"github.com/mrbttf/bridge-server/pkg/log"
)

func New(config *config.Config) (*sql.DB, error) {
	connectionString := fmt.Sprintf("postgres://%s:%s@%s/%s?sslmode=disable", config.DBUser, config.DBPassword, config.DBHost, config.DBName)
	db, err := sql.Open("postgres", connectionString)
	if err != nil {
		return nil, err
	}
	err = db.Ping()
	if err != nil {
		return nil, err
	}
	log.Info("Connected to database")
	return db, nil
}
