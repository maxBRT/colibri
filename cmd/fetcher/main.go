package main

import (
	"log"
	"os"
	"sync"

	amqp "github.com/rabbitmq/amqp091-go"
	"www.github.com/maxbrt/colibri/internal/pubsub"
	ps "www.github.com/maxbrt/colibri/internal/pubsub"
	"www.github.com/maxbrt/colibri/internal/rss"
	s "www.github.com/maxbrt/colibri/internal/sources"
)

func main() {
	conn, err := amqp.Dial(ps.ConnectionString)
	if err != nil {
		log.Printf("%s", err)
		os.Exit(1)
	}
	defer conn.Close()

	ch, _, err := ps.DeclareAndBind(
		conn,
		ps.ColibriExchange,
		ps.ColibriPostsQueue,
		ps.ColibriPostsKey,
		ps.DurableQueue,
	)
	if err != nil {
		log.Printf("%s", err)
		os.Exit(1)
	}
	defer ch.Close()

	sources, err := s.ReadSources("./sources/sources.csv")
	if err != nil {
		log.Fatalf("%s", err)
	}

	var wg sync.WaitGroup

	for _, source := range sources {
		wg.Add(1)

		if err := ps.PublishJSON(
			ch,
			ps.ColibriExchange,
			ps.ColibriSourcesKey,
			source); err != nil {
			log.Printf("%s", err)
			os.Exit(1)
		}

		go func(source s.Source) {
			defer wg.Done()

			posts, err := rss.FetchAndParse(source)
			if err != nil {
				log.Printf("%s", err)
				os.Exit(1)
			}
			for _, post := range posts {
				if err := pubsub.PublishJSON(
					ch,
					ps.ColibriExchange,
					ps.ColibriPostsKey,
					post); err != nil {
					log.Printf("%s", err)
					os.Exit(1)
				}
			}
		}(source)
	}
	wg.Wait()
}
