package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"
	"www.github.com/maxbrt/colibri/internal/database"
	"www.github.com/maxbrt/colibri/internal/server"
	"www.github.com/maxbrt/colibri/internal/utils"
)

func main() {
	driver := os.Getenv("DB_DRIVER")
	dbString, err := utils.GetSecret(os.Getenv("DB_STRING_FILE"))
	if err != nil {
		log.Printf("error loading secret: %s", err)
		os.Exit(1)
	}

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
	if err := s.ListenAndServe(":8080"); err != nil {
		log.Printf("%s", err)
		os.Exit(1)
	}
}
