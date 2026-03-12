package main

import (
	"log"
	"os"
	"sync"

	amqp "github.com/rabbitmq/amqp091-go"
	"www.github.com/maxbrt/colibri/internal/pubsub"
	r "www.github.com/maxbrt/colibri/internal/routing"
	"www.github.com/maxbrt/colibri/internal/rss"
)

const (
	ConnectionString = "amqp://guest:guest@localhost:5672/"
)

func main() {
	conn, err := amqp.Dial(ConnectionString)
	if err != nil {
		log.Printf("%s", err)
		os.Exit(1)
	}

	ch, _, err := pubsub.DeclareAndBind(
		conn,
		r.ColibriExchange,
		r.ColibriFeedQueue,
		r.ColibriFeedKey,
		pubsub.DurableQueue,
	)
	if err != nil {
		log.Printf("%s", err)
		os.Exit(1)
	}

	sources, err := rss.ReadSources("./sources")
	if err != nil {
		log.Fatalf("%s", err)
	}

	var wg sync.WaitGroup

	for _, s := range sources {
		wg.Add(1)

		go func(s rss.Source) {
			defer wg.Done()

			posts, err := rss.FetchAndParse(s)
			if err != nil {
				log.Printf("%s", err)
				os.Exit(1)
			}
			for _, p := range posts {
				if err := pubsub.PublishJSON(
					ch,
					r.ColibriExchange,
					r.ColibriFeedKey,
					p); err != nil {
					log.Printf("%s", err)
					os.Exit(1)
				}
			}
		}(s)
	}
	wg.Wait()
	conn.Close()
	ch.Close()
}
