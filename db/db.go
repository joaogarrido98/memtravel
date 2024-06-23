package db

import (
	"database/sql"
	"memtravel/configs"

	_ "github.com/lib/pq"
)

var Database *sql.DB

func DBConnect() error {
	connStr := configs.Envs.DBUser + "://" +
		configs.Envs.DBUser + ":" +
		configs.Envs.DBPassword + "@" +
		configs.Envs.DBAddress + "/" +
		configs.Envs.DBName + "?sslmode=disable"

	Database, err := sql.Open("postgres", connStr)
	if err != nil {
		return err
	}

	defer Database.Close()

	return Database.Ping()
}
