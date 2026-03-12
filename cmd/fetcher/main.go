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

func main() {
	conn, err := amqp.Dial(r.ConnectionString)
	if err != nil {
		log.Printf("%s", err)
		os.Exit(1)
	}
	defer conn.Close()

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
	defer ch.Close()

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
}
