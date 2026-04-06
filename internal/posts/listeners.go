package posts

import (
	"context"
	"database/sql"
	"fmt"

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
			fmt.Printf("error deduplicating post: %s", err)
			return pubsub.Ack
		}
		if id == "" {
			return pubsub.Ack
		}

		a, err := NewDescriptionAgent()
		if err != nil {
			fmt.Printf("error creating description agent: %s", err)
			return pubsub.Ack
		}
		err = a.GenerateDescription(&p)
		if err != nil {
			fmt.Printf("error generating description: %s", err)
			return pubsub.Ack
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
			return pubsub.Ack
		}

		return pubsub.Ack
	}
}
