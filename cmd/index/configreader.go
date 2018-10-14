package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
)

func ConfigCollections() []ConfigCollection {

	path, _ := filepath.Abs(fmt.Sprintf("config/%s.json", *config))

	jsonFile, err := os.Open(path)

	if err != nil {
		log.Fatalln(err.Error())
	}

	defer jsonFile.Close()

	content, _ := ioutil.ReadAll(jsonFile)

	var collections []ConfigCollection

	json.Unmarshal(content, &collections)

	return collections
}

func GetConfigCollection(collection string) *ConfigCollection {
	collections := ConfigCollections()

	for _, c := range collections {
		if c.Collection == collection {
			return &c
		}
	}

	return nil
}
