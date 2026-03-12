// Package routing regroup most of the constant required for RabbitMQ
package routing

import "os"

const (
	ColibriExchange     = "colibri_topic"
	ColibriPostsKey     = "posts"
	ColibriSourcesKey   = "sources"
	ColibriPostsQueue   = "posts_queue"
	ColibriSourcesQueue = "sources_queue"
)

var ConnectionString = getConnectionString()

func getConnectionString() string {
	if v := os.Getenv("AMQP_URL"); v != "" {
		return v
	}
	return "amqp://guest:guest@localhost:5672/"
}
