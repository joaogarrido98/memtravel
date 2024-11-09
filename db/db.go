package db

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"

	"memtravel/configs"
)

type (
	// Transaction is the type of sql transactions to be executed
	Transaction struct {
		Query  string
		Params []any
	}

	// Database is the blueprint for the db
	Database struct {
		*sql.DB
	}
)

// Connect init the database connection
func Connect() (Database, error) {
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

// ExecQuery executes a non-select query and checks for affected rows, if none, then return error
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

// ExecTransaction executes the given queries inside a transaction block, if any fail, roll all previous ones back
// if they all pass, commit it
func (database Database) ExecTransaction(transactions []Transaction) error {
	tx, err := database.Begin()
	if err != nil {
		return err
	}

	for _, params := range transactions {
		_, err := tx.Exec(params.Query, params.Params...)
		if err != nil {
			rollbackErr := tx.Rollback()
			if rollbackErr != nil {
				return rollbackErr
			}
			return err
		}
	}

	return tx.Commit()
}
