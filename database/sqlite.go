package database

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

func GetSQLiteInstance() (*sql.DB, error) {
	dbPath := os.Getenv("DB_PATH")
	db, err := sql.Open("sqlite3", dbPath)

	if err != nil {
		return nil, fmt.Errorf("failed to open database: %s", err)
	}

	err = db.Ping()
	if err != nil {
		return nil, fmt.Errorf("failed to ping database: %s", err)
	}

	db.Exec("CREATE TABLE IF NOT EXISTS messages (message TEXT)")

	return db, nil
}
