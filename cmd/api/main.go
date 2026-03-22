package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	_ "github.com/lib/pq"
	"www.github.com/maxbrt/colibri/internal/database"
	"www.github.com/maxbrt/colibri/internal/server"
)

func main() {
	driver := os.Getenv("DB_DRIVER")
	dbString := os.Getenv("DB_STRING")

	dbConn, err := sql.Open(driver, dbString)
	if err != nil {
		log.Printf("%s", err)
		os.Exit(1)
	}
	defer dbConn.Close()

	db := database.New(dbConn)
	s := server.NewServer()
	s.MountHandlers(db)

	fmt.Println("Listening on port 8080")
	if err := http.ListenAndServe(":8080", s.Router); err != nil {
		log.Printf("%s", err)
		os.Exit(1)
	}
}
