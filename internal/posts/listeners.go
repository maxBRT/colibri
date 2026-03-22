package posts

import (
	"context"
	"database/sql"
	"fmt"
	"log"

	"www.github.com/maxbrt/colibri/internal/database"
	"www.github.com/maxbrt/colibri/internal/pubsub"
)

func Listener(db *database.Queries) func(Post) pubsub.AckType {
	return func(p Post) pubsub.AckType {
		id, err := db.DeduplicatePosts(
			context.Background(),
			database.DeduplicatePostsParams{
				Title:       p.Title,
				Description: sql.NullString{String: p.Description, Valid: true},
				Link:        p.Link,
				Guid:        p.GUID,
				PubDate:     p.PubDate,
				SourceID:    p.SourceID,
			})
		if err != nil {
			log.Printf("%s", err)
			return pubsub.Ack
		}
		if id == "" {
			fmt.Printf("duplicate post: %s", p.GUID)
			return pubsub.Ack
		}

		a, err := NewDescriptionAgent()
		if err != nil {
			fmt.Printf("error creating description agent: %s", err)
			return pubsub.NackRequeue
		}
		err = a.GenerateDescription(&p)
		if err != nil {
			fmt.Printf("error generating description: %s", err)
			return pubsub.NackRequeue
		}

		err = db.UpdatePost(
			context.Background(),
			database.UpdatePostParams{
				Title:       p.Title,
				Description: sql.NullString{String: p.Description, Valid: true},
				Link:        p.Link,
				Guid:        p.GUID,
				PubDate:     p.PubDate,
				SourceID:    p.SourceID,
			})
		if err != nil {
			fmt.Printf("error updating post: %s", err)
			return pubsub.NackRequeue
		}

		return pubsub.Ack
	}
}
