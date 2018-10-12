package main

import (
	"flag"
	"github.com/globalsign/mgo"
	"github.com/idnan/go-mongo-indexer/pkg/storage"
	_ "github.com/joho/godotenv/autoload"
	"github.com/withmandala/go-log"
	"os"
)

var (
	config *string
	apply  *bool
	db     *mgo.Session

	logger = log.New(os.Stderr).WithDebug()
)

func init() {
	config = flag.String("c", "", "index config file")
	apply = flag.Bool("apply", false, "Apply the changes")
	flag.Parse()

	initDb()
}

func main() {
	if *config == "" {
		logger.Error("index config file is required")
	}
	defer db.Close()

	execute()
}

func initDb() {
	session, err := storage.Connect(&storage.Config{
		Host:     os.Getenv("OMS_API.DB.MONGO.HOST"),
		Port:     os.Getenv("OMS_API.DB.MONGO.PORT"),
		Database: os.Getenv("OMS_API.DB.MONGO.DATABASE"),
		Username: os.Getenv("OMS_API.DB.MONGO.USERNAME"),
		Password: os.Getenv("OMS_API.DB.MONGO.PASSWORD"),
		Options:  os.Getenv("OMS_API.DB.MONGO.OPTS"),
	})

	if err != nil {
		logger.Fatal(err)
	}

	db = session
}
