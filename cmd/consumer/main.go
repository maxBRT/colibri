package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	amqp "github.com/rabbitmq/amqp091-go"
	"www.github.com/maxbrt/colibri/internal/pubsub"
	r "www.github.com/maxbrt/colibri/internal/routing"
	"www.github.com/maxbrt/colibri/internal/rss"
)

func main() {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	conn, err := amqp.Dial(r.ConnectionString)
	if err != nil {
		log.Printf("%s", err)
		os.Exit(1)
	}
	defer conn.Close()

	err = pubsub.SubscribeJSON(
		conn,
		r.ColibriExchange,
		r.ColibriFeedQueue,
		r.ColibriFeedKey,
		pubsub.DurableQueue,
		handlerPost(),
	)
	if err != nil {
		log.Printf("%s", err)
	}

	fmt.Println("Listening for messages...")
	<-sigs
	fmt.Println("")
	fmt.Println("Program killed")
}

func handlerPost() func(rss.Post) pubsub.AckType {
	return func(p rss.Post) pubsub.AckType {
		fmt.Printf("Title: %s\nLink: %s\n", p.Title, p.Link)
		return pubsub.Ack
	}
}
