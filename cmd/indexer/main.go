package main

import (
	"context"
	"flag"
	"log"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	config   *string
	apply    *bool
	mongoUri *string
	database *string
	db       *mongo.Database
)

func init() {
	config = flag.String("config", "", "[REQUIRED] index config file")
	apply = flag.Bool("apply", false, "apply the changes")
	mongoUri = flag.String("uri", "", "[REQUIRED] mongo uri path")
	database = flag.String("database", "", "[REQUIRED] database Name")
	flag.Parse()
}

func main() {
	if len(*config) == 0 || len(*mongoUri) == 0 || len(*database) == 0 {
		usage()
	}

	dbClient := initDb()
	db = dbClient.Database(*database)

	defer dbClient.Disconnect(context.TODO())

	execute()
}

func initDb() *mongo.Client {
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	opts := options.Client().ApplyURI(*mongoUri)

	if err := opts.Validate(); err != nil {
		log.Fatalf("Invalid option given: %v", err)
	}

	client, err := mongo.Connect(ctx, opts)
	if err != nil {
		log.Fatalf("Connection failed: %v", err)
	}

	return client
}

func usage() {
	flag.Usage()
	os.Exit(1)
}
