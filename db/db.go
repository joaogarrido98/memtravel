package db

import (
	"database/sql"

	_ "github.com/lib/pq"

	"memtravel/configs"
)

// DBConnect init the database connection
func DBConnect() (*sql.DB, error) {
	connStr := configs.Envs.DBUser + "://" +
		configs.Envs.DBUser + ":" +
		configs.Envs.DBPassword + "@" +
		configs.Envs.DBAddress + "/" +
		configs.Envs.DBName + "?sslmode=disable"

	database, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}

	err = database.Ping()
	if err != nil {
		return nil, err
	}

	return database, err
}
