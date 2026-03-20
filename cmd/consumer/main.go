package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	_ "github.com/lib/pq"
	amqp "github.com/rabbitmq/amqp091-go"
	"www.github.com/maxbrt/colibri/internal/database"
	"www.github.com/maxbrt/colibri/internal/pubsub"
	r "www.github.com/maxbrt/colibri/internal/routing"
	"www.github.com/maxbrt/colibri/internal/rss"
)

func main() {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	driver := os.Getenv("DB_DRIVER")
	dbString := os.Getenv("DB_STRING")

	dbConn, err := sql.Open(driver, dbString)
	if err != nil {
		log.Printf("%s", err)
		os.Exit(1)
	}
	defer dbConn.Close()

	db := database.New(dbConn)

	conn, err := amqp.Dial(r.ConnectionString)
	if err != nil {
		log.Printf("%s", err)
		os.Exit(1)
	}
	defer conn.Close()

	err = pubsub.SubscribeJSON(
		conn,
		r.ColibriExchange,
		r.ColibriPostsQueue,
		r.ColibriPostsKey,
		pubsub.DurableQueue,
		handlerPost(db),
	)
	if err != nil {
		log.Printf("%s", err)
	}

	err = pubsub.SubscribeJSON(
		conn,
		r.ColibriExchange,
		r.ColibriSourcesQueue,
		r.ColibriSourcesKey,
		pubsub.DurableQueue,
		handlerSources(db),
	)
	if err != nil {
		log.Printf("%s", err)
	}

	fmt.Println("Listening for messages...")
	<-sigs
	fmt.Println("")
	fmt.Println("Program killed")
}

func handlerPost(db *database.Queries) func(rss.Post) pubsub.AckType {
	return func(p rss.Post) pubsub.AckType {
		post, err := db.CreatePost(
			context.Background(),
			database.CreatePostParams{
				Title: p.Title,
				Description: sql.NullString{
					String: p.Description,
					Valid:  p.Description != "",
				},
				Link:     p.Link,
				Guid:     p.GUID,
				PubDate:  p.PubDate,
				SourceID: p.SourceID,
			})
		if err != nil {
			fmt.Printf("error creating post: %s", err)
		}
		fmt.Printf("Title: %s\nLink: %s\n", post.Title, post.Link)
		return pubsub.Ack
	}
}

func handlerSources(db *database.Queries) func(rss.Source) pubsub.AckType {
	return func(s rss.Source) pubsub.AckType {
		source, err := db.CreateSource(
			context.Background(),
			database.CreateSourceParams{
				ID:       s.ID,
				Name:     s.Name,
				Url:      s.URL,
				Category: s.Category,
			})
		if err != nil {
			log.Printf("error creating source: %s", err)
		}

		fmt.Printf(
			"id: %s\nname: %s\nurl: %s\ncategory: %s\n",
			source.ID,
			source.Name,
			source.Url,
			source.Category,
		)
		return pubsub.Ack
	}
}
