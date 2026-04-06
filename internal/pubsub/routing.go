package pubsub

import (
	"log"
	"os"

	"www.github.com/maxbrt/colibri/internal/utils"
)

const (
	ColibriExchange     = "colibri_topic"
	ColibriPostsKey     = "posts"
	ColibriSourcesKey   = "sources"
	ColibriLogoKey      = "logo"
	ColibriPostsQueue   = "posts_queue"
	ColibriSourcesQueue = "sources_queue"
	ColibriLogoQueue    = "logo_queue"
	CDNBaseURL          = "d1mc6q6crhpy86.cloudfront.net"
)

var ConnectionString = getConnectionString()

func getConnectionString() string {
	secretFile := os.Getenv("AMQP_URL_FILE")
	if secretFile == "" {
		log.Println("AMQP_URL_FILE not set, using default localhost connection")
		return "amqp://guest:guest@localhost:5672/"
	}
	v, err := utils.GetSecret(secretFile)
	if err != nil {
		log.Fatalf("Failed to read AMQP secret %s", err)
	}
	return v
}
