// Package db handles the postgresql database connection
package pgdb

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/lib/pq"
)

// ConnectDB initializes and tests the connection to a PostgreSQL container.
func ConnectDB() (*sql.DB, error) {
	password := os.Getenv("POSTGRES_PASSWORD")
	if password == "" {
		return nil, fmt.Errorf("POSTGRES_PASSWORD environment variable is not set")
	}

	host := "localhost"
	port := "5432"
	user := "postgres"
	dbname := "postgres"

	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("error opening database configuration: %w", err)
	}

	if err := db.Ping(); err != nil {
		db.Close()
		return nil, fmt.Errorf("error connecting to the database: %w", err)
	}

	fmt.Println("Successfully connected to the PostgreSQL container!")
	return db, nil
}
