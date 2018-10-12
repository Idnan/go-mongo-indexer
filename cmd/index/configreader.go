package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
)

func ConfigCollections(file string) []ConfigCollection {

	path, _ := filepath.Abs(fmt.Sprintf("config/%s.json", file))

	jsonFile, err := os.Open(path)

	if err != nil {
		logger.Error(err.Error())
	}

	defer jsonFile.Close()

	content, _ := ioutil.ReadAll(jsonFile)

	var collections []ConfigCollection

	json.Unmarshal(content, &collections)

	return collections
}

func ConfigCollectionIndexes(file string, collection string) *ConfigCollection {
	collections := ConfigCollections(file)

	for _, c := range collections {
		if c.Collection == collection {
			return &c
		}
	}

	return nil
}
