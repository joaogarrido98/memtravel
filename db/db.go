package db

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"

	"memtravel/configs"
)

type Database struct {
	*sql.DB
}

// DBConnect init the database connection
func DBConnect() (Database, error) {
	connStr := configs.Envs.DBUser + "://" +
		configs.Envs.DBUser + ":" +
		configs.Envs.DBPassword + "@" +
		configs.Envs.DBAddress + "/" +
		configs.Envs.DBName + "?sslmode=disable"

	database, err := sql.Open("postgres", connStr)
	if err != nil {
		return Database{}, err
	}

	err = database.Ping()
	if err != nil {
		return Database{}, err
	}

	return Database{database}, err
}

// ExecQuery executes a non-select query
func (database Database) ExecQuery(query string, data ...any) error {
	result, err := database.Exec(query, data...)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return fmt.Errorf("expected at least one row to be affected but rececvied: %d", rows)
	}

	return nil
}
