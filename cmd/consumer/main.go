package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	_ "github.com/lib/pq"
	amqp "github.com/rabbitmq/amqp091-go"
	"www.github.com/maxbrt/colibri/internal/database"
	"www.github.com/maxbrt/colibri/internal/logo"
	p "www.github.com/maxbrt/colibri/internal/posts"
	ps "www.github.com/maxbrt/colibri/internal/pubsub"
	s "www.github.com/maxbrt/colibri/internal/sources"
	"www.github.com/maxbrt/colibri/internal/utils"
)

func main() {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	driver := os.Getenv("DB_DRIVER")
	dbString, err := utils.GetSecret(os.Getenv("DB_STRING_FILE"))
	if err != nil {
		log.Printf("error loading secret: %s", err)
		os.Exit(1)
	}

	dbConn, err := sql.Open(driver, dbString)
	if err != nil {
		log.Printf("%s", err)
		os.Exit(1)
	}
	defer dbConn.Close()

	db := database.New(dbConn)

	bucketFile := os.Getenv("LOGO_S3_BUCKET_FILE")
	if bucketFile == "" {
		log.Println("LOGO_S3_BUCKET_FILE is required")
		os.Exit(1)
	}
	bucket, err := utils.GetSecret(bucketFile)
	if err != nil {
		log.Printf("failed to read bucket secret: %s", err)
		os.Exit(1)
	}
	regionFile := os.Getenv("LOGO_S3_REGION_FILE")
	if regionFile == "" {
		log.Println("LOGO_S3_REGION_FILE is required")
		os.Exit(1)
	}
	region, err := utils.GetSecret(regionFile)
	if err != nil {
		log.Printf("failed to read region secret: %s", err)
		os.Exit(1)
	}

	awsAccessKeyFile := os.Getenv("AWS_ACCESS_KEY_ID_FILE")
	if awsAccessKeyFile == "" {
		log.Println("AWS_ACCESS_KEY_ID_FILE is required")
		os.Exit(1)
	}
	awsAccessKey, err := utils.GetSecret(awsAccessKeyFile)
	if err != nil {
		log.Printf("failed to read AWS access key secret: %s", err)
		os.Exit(1)
	}
	awsSecretKeyFile := os.Getenv("AWS_SECRET_ACCESS_KEY_FILE")
	if awsSecretKeyFile == "" {
		log.Println("AWS_SECRET_ACCESS_KEY_FILE is required")
		os.Exit(1)
	}
	awsSecretKey, err := utils.GetSecret(awsSecretKeyFile)
	if err != nil {
		log.Printf("failed to read AWS secret key: %s", err)
		os.Exit(1)
	}
	awsSessionToken := ""
	if tokenFile := os.Getenv("AWS_SESSION_TOKEN_FILE"); tokenFile != "" {
		awsSessionToken, err = utils.GetSecret(tokenFile)
		if err != nil {
			log.Printf("failed to read AWS session token: %s", err)
			os.Exit(1)
		}
	}

	awsCfg, err := config.LoadDefaultConfig(
		context.Background(),
		config.WithRegion(region),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(awsAccessKey, awsSecretKey, awsSessionToken)),
	)
	if err != nil {
		log.Printf("failed to load aws config: %s", err)
		os.Exit(1)
	}

	s3Client := s3.NewFromConfig(awsCfg)

	logoService, err := logo.NewService(logo.Config{
		DB:       db,
		S3Client: s3Client,
		Bucket:   bucket,
	})
	if err != nil {
		log.Printf("failed to bootstrap logo service: %s", err)
		os.Exit(1)
	}

	conn, err := amqp.Dial(ps.ConnectionString)
	if err != nil {
		log.Printf("%s", err)
		os.Exit(1)
	}
	defer conn.Close()

	err = ps.SubscribeJSON(
		conn,
		ps.ColibriExchange,
		ps.ColibriPostsQueue,
		ps.ColibriPostsKey,
		ps.DurableQueue,
		p.Listener(db),
	)
	if err != nil {
		log.Printf("%s", err)
	}

	err = ps.SubscribeJSON(
		conn,
		ps.ColibriExchange,
		ps.ColibriLogoQueue,
		ps.ColibriLogoKey,
		ps.DurableQueue,
		logo.Listener(logoService),
	)
	if err != nil {
		log.Printf("%s", err)
	}

	err = ps.SubscribeJSON(
		conn,
		ps.ColibriExchange,
		ps.ColibriSourcesQueue,
		ps.ColibriSourcesKey,
		ps.DurableQueue,
		s.Listener(db),
	)
	if err != nil {
		log.Printf("%s", err)
	}

	fmt.Println("Listening for messages...")
	<-sigs
	fmt.Println("")
	fmt.Println("Program killed")
}
