package pgdb

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/brianvoe/gofakeit/v6"
)

type User struct {
	ID    int
	Name  string
	Email string
}

// CreateUserTable creates the users table if it doesn't already exist.
func CreateUserTable(db *sql.DB) error {
	query := `
	CREATE TABLE IF NOT EXISTS users (
		id SERIAL PRIMARY KEY,
		name VARCHAR(100) NOT NULL,
		email VARCHAR(100) UNIQUE NOT NULL
	);`

	_, err := db.Exec(query)
	if err != nil {
		return fmt.Errorf("failed to create table: %w", err)
	}

	fmt.Println("Table 'users' ready!")
	return nil
}

// SeedUsers creates randomized user entries using gofakeit in a single batched query.
func SeedUsers(db *sql.DB, count int) error {
	valueStrings := make([]string, 0, count)
	valueArgs := make([]any, 0, count*2)

	for i := range count {
		name := gofakeit.Name()
		email := gofakeit.Email()

		valueStrings = append(valueStrings, fmt.Sprintf("($%d, $%d)", i*2+1, i*2+2))
		valueArgs = append(valueArgs, name, email)
	}

	stmt := fmt.Sprintf("INSERT INTO users (name, email) VALUES %s ON CONFLICT DO NOTHING;", strings.Join(valueStrings, ", "))

	tx, err := db.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	result, err := tx.Exec(stmt, valueArgs...)
	if err != nil {
		return fmt.Errorf("failed to batch insert users: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	rowsInserted, _ := result.RowsAffected()
	fmt.Printf("Successfully seeded %d user records with gofakeit!\n", rowsInserted)
	return nil
}

// InsertUser creates a new user and returns the auto-generated ID.
func InsertUser(db *sql.DB, name, email string) (int, error) {
	query := `INSERT INTO users (name, email) VALUES ($1, $2) RETURNING id;`

	var id int
	err := db.QueryRow(query, name, email).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("insert failed: %w", err)
	}

	return id, nil
}

// GetUserByID fetches a single user by their primary key.
func GetUserByID(db *sql.DB, id int) (*User, error) {
	query := `SELECT id, name, email FROM users WHERE id = $1;`

	user := &User{}
	err := db.QueryRow(query, id).Scan(&user.ID, &user.Name, &user.Email)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("user with id %d not found", id)
		}
		return nil, fmt.Errorf("query failed: %w", err)
	}

	return user, nil
}

// UpdateUserEmail updates the email for a given user ID.
func UpdateUserEmail(db *sql.DB, id int, newEmail string) error {
	query := `UPDATE users SET email = $1 WHERE id = $2;`

	result, err := db.Exec(query, newEmail, id)
	if err != nil {
		return fmt.Errorf("update failed: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to check affected rows: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("no user found with id %d", id)
	}

	fmt.Printf("Updated email for user %d\n", id)
	return nil
}

// DeleteUser removes a user record by ID.
func DeleteUser(db *sql.DB, id int) error {
	query := `DELETE FROM users WHERE id = $1;`

	result, err := db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("delete failed: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to check affected rows: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("no user found to delete with id %d", id)
	}

	fmt.Printf("Deleted user %d\n", id)
	return nil
}
