package database

import (
	"database/sql"
	"fmt"
	"itmo-profile/config"
)

func Connect() (*sql.DB, error) {

	var host = config.GetEnv("PG_HOST", "")

	var port = config.GetEnvAsInt("PG_PORT", 0)

	var user = config.GetEnv("PG_USER", "")
	var password = config.GetEnv("PG_PASSWORD", "")
	var dbname = config.GetEnv("PG_DBNAME", "")

	dsn := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return &sql.DB{}, err
	}

	err = db.Ping()
	if err != nil {
		return &sql.DB{}, err
	}

	return db, nil
}
