package config

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

func LoadDB() (*sql.DB, error) {
	//username = root
	//password = password
	//address = localhost
	//port = 3306
	//db name = world
	db, err := sql.Open("mysql", "root:password@tcp(localhost:3306)/world")
	if err != nil {
		return nil, fmt.Errorf("failed to open database connection: %w", err)
	}

	// Test the connection
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return db, nil
}
