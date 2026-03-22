package sources

import (
	"context"
	"fmt"
	"log"

	"www.github.com/maxbrt/colibri/internal/database"
	"www.github.com/maxbrt/colibri/internal/pubsub"
)

func Listener(db *database.Queries) func(Source) pubsub.AckType {
	return func(s Source) pubsub.AckType {
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
