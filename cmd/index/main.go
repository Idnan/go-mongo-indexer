package main

import (
	"flag"
	"github.com/davecgh/go-spew/spew"
	"github.com/globalsign/mgo"
	"log"
	"os"
	"time"
)

var (
	config   *string
	apply    *bool
	mongoUri *string
	session  *mgo.Session
	db       *mgo.Database
)

func init() {
	config = flag.String("config", "", "[REQUIRED] index config file")
	apply = flag.Bool("apply", false, "apply the changes")
	mongoUri = flag.String("uri", "", "[REQUIRED] mongo uri path")
	flag.Parse()

	initDb()
}

func main() {
	if *config == "" || *mongoUri == "" {
		usage()
	}

	defer session.Close()
	execute()
}

func initDb() {
	var dialInfo *mgo.DialInfo
	var err error

	dialInfo, err = mgo.ParseURL(*mongoUri)
	if err != nil {
		log.Fatalln(err.Error())
	}

	session, err = mgo.DialWithTimeout(*mongoUri, time.Second*3)
	if err != nil {
		log.Fatalln(err.Error())
	}

	db = session.DB(dialInfo.Database)
}

func usage() {
	flag.Usage()
	os.Exit(1)
}

func dd(data ...interface{}) {
	spew.Dump(data)
	os.Exit(0)
}
