package pubsub

import (
	"os"

	"www.github.com/maxbrt/colibri/internal/utils"
)

const (
	ColibriExchange     = "colibri_topic"
	ColibriPostsKey     = "posts"
	ColibriSourcesKey   = "sources"
	ColibriPostsQueue   = "posts_queue"
	ColibriSourcesQueue = "sources_queue"
)

var ConnectionString = getConnectionString()

func getConnectionString() string {
	v, _ := utils.GetSecret(os.Getenv("AMQP_URL_FILE"))
	if v != "" {
		return v
	}
	return "amqp://guest:guest@localhost:5672/"
}
