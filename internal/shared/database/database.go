package database

import (
	"database/sql"

	_ "modernc.org/sqlite"
)

type DatabaseConnection struct {
	DatabaseName string
}

func GenerateDatabaseConnection(params DatabaseConnection, initFunction func(db *sql.DB) error) (*sql.DB, error) {
	db, err := sql.Open("sqlite", params.DatabaseName)
	if err != nil {
		return nil, err
	}

	if initFunction != nil {
		if err := initFunction(db); err != nil {
			return nil, err
		}
	}

	return db, nil
}
