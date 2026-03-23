package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	_ "github.com/lib/pq"
	amqp "github.com/rabbitmq/amqp091-go"
	"www.github.com/maxbrt/colibri/internal/database"
	p "www.github.com/maxbrt/colibri/internal/posts"
	ps "www.github.com/maxbrt/colibri/internal/pubsub"
	s "www.github.com/maxbrt/colibri/internal/sources"
	"www.github.com/maxbrt/colibri/internal/utils"
)

func main() {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

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

	conn, err := amqp.Dial(ps.ConnectionString)
	if err != nil {
		log.Printf("%s", err)
		os.Exit(1)
	}
	defer conn.Close()

	err = ps.SubscribeJSON(
		conn,
		ps.ColibriExchange,
		ps.ColibriPostsQueue,
		ps.ColibriPostsKey,
		ps.DurableQueue,
		p.Listener(db),
	)
	if err != nil {
		log.Printf("%s", err)
	}

	err = ps.SubscribeJSON(
		conn,
		ps.ColibriExchange,
		ps.ColibriSourcesQueue,
		ps.ColibriSourcesKey,
		ps.DurableQueue,
		s.Listener(db),
	)
	if err != nil {
		log.Printf("%s", err)
	}

	fmt.Println("Listening for messages...")
	<-sigs
	fmt.Println("")
	fmt.Println("Program killed")
}
