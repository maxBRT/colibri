package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"

	"github.com/joho/godotenv"
	"github.com/pressly/goose/v3"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	driver := os.Getenv("DB_DRIVER")
	dbString := os.Getenv("DB_STRING")
	migrationDir := os.Getenv("GOOSE_MIGRATION_DIR")

	db, err := sql.Open(driver, dbString)
	if err != nil {
		log.Fatalf("Error connecting to db | %s", err)
	}
	switch os.Args[1] {
	case "up":
		if err := goose.Up(db, migrationDir); err != nil {
			log.Printf("Error applying migrations: %s\n", err)
		}
	case "down":
		if err := goose.Down(db, migrationDir); err != nil {
			log.Printf("Error applying migrations: %s\n", err)
		}
	default:
		fmt.Println("Invalid argument: <up | down>")
	}
}
