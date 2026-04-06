package sources

import (
	"context"
	"log"

	"www.github.com/maxbrt/colibri/internal/database"
	"www.github.com/maxbrt/colibri/internal/pubsub"
)

func Listener(db *database.Queries) func(Source) pubsub.AckType {
	return func(s Source) pubsub.AckType {
		_, err := db.CreateSource(
			context.Background(),
			database.CreateSourceParams{
				ID:       s.ID,
				Name:     s.Name,
				Url:      s.URL,
				Category: s.Category,
			})
		if err != nil {
			log.Printf("error creating source: %s", err)
			return pubsub.Ack
		}

		return pubsub.Ack
	}
}
