package main

import (
	"flag"
	"github.com/davecgh/go-spew/spew"
	"github.com/globalsign/mgo"
	"github.com/idnan/go-mongo-indexer/pkg/storage"
	_ "github.com/joho/godotenv/autoload"
	"log"
	"os"
)

var (
	config *string
	apply  *bool
	db     *mgo.Session
)

func init() {
	config = flag.String("c", "", "index config file")
	apply = flag.Bool("apply", false, "Apply the changes")
	flag.Parse()

	initDb()
}

func main() {
	if *config == "" {
		log.Fatalln("index config file is required")
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
		log.Fatalln(err.Error())
	}

	db = session
}

func dd(data ...interface{}) {
	spew.Dump(data)
	os.Exit(0)
}
