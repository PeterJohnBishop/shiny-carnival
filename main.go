package main

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/joho/godotenv"
	pgdb "github.com/peterjohnbishop/shiny-carnival/db"
)

var db *sql.DB

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println("Warning: No .env file found or error loading it")
	}

	db, err = pgdb.ConnectDB()
	if err != nil {
		fmt.Printf("error connecting to the postgres container: %s", err)
	}
	_ = db
}
