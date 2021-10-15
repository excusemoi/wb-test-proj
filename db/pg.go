package db

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"goProj/config"
)

type PgDb struct {
	port 				string
	hostname 			string
	name 				string
	user 				string
	password 			string
}

func Dial(cfg *config.Config) (*sqlx.DB, error) {
	dbInfo := fmt.Sprintf("host=%s port=%s user=%s "+
		"password=%s dbname=%s sslmode=disable",
		cfg.PgHostName, cfg.PgPort, cfg.PgUser, cfg.PgPassword,	cfg.PgDb)
	db, err := sqlx.Open("postgres", dbInfo)
	if err != nil {
		return nil, err
	}
	err = db.Ping()
	if err != nil {
		return nil, err
	}
	return db, nil

}
