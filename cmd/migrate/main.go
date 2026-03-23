package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"
	"www.github.com/maxbrt/colibri/internal/utils"
	schema "www.github.com/maxbrt/colibri/sql"

	"github.com/pressly/goose/v3"
)

func main() {
	driver := os.Getenv("DB_DRIVER")
	dbString, err := utils.GetSecret(os.Getenv("DB_STRING_FILE"))
	if err != nil {
		log.Printf("error loading secret: %s", err)
		os.Exit(1)
	}
	goose.SetBaseFS(schema.MigrationFiles)

	db, err := sql.Open(driver, dbString)
	if err != nil {
		log.Fatalf("Error connecting to db | %s", err)
	}
	switch os.Args[1] {
	case "up":
		if err := goose.Up(db, "schema"); err != nil {
			log.Printf("Error applying migrations: %s\n", err)
		}
	case "down":
		if err := goose.Down(db, "schema"); err != nil {
			log.Printf("Error applying migrations: %s\n", err)
		}
	default:
		fmt.Println("Invalid argument: <up | down>")
	}
}
