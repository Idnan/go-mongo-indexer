package storage

import (
	"fmt"
	"github.com/globalsign/mgo"
)

func Connect(config *Config) (*mgo.Session, error) {

	usernamePassword := ""
	if config.Username != "" {
		usernamePassword = fmt.Sprintf("%s:%s@", config.Username, config.Password)
	}

	uri := fmt.Sprintf("mongodb://%s%s", usernamePassword, config.Host)

	if config.Port != "" {
		uri += fmt.Sprintf(":%s", config.Port)
	}

	uri += fmt.Sprintf("/%s", config.Database)

	return mgo.Dial(uri)
}
