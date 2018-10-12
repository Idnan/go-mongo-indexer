package main

import (
	"fmt"
	"reflect"
)

func execute() {

	getIndexesDiff()

	if !*apply {
		// todo show diff
	}

	if *apply {
		// todo apply the indexes
	}
}

func getIndexesDiff() {

	var indexDiff map[string]map[string]map[string]map[string]interface{}

	for _, collection := range Collections() {

		// Get current indexes and cap-size
		currentIndexes := DbIndexes(collection)

		// If we don't have the current collection in the index create list then drop all index
		if !IsCollectionToIndex(collection) {
			for indexName, indexDetail := range currentIndexes {

				if indexDiff[collection] == nil {

				}

				indexDiff[collection]["old"]["indexes"][indexName] = indexDetail
			}
			fmt.Println(indexDiff)
		}

	}
}

func Collections() []string {
	collections, err := db.DB("oms_api").CollectionNames()

	if err != nil {
		logger.Error(err.Error())
	}

	return collections
}

func DbIndexes(collection string) map[string][]string {
	indexes, err := db.DB("oms_api").C(collection).Indexes()

	if err != nil {
		logger.Error(err.Error())
	}

	dbIndexes := make(map[string][]string)

	for _, index := range indexes {

		keys := index.Key

		if len(keys) == 0 || reflect.DeepEqual(keys, []string{"_id"}) {
			continue
		}

		dbIndexes[index.Name] = keys
	}

	return dbIndexes
}

// Drop index from collection by index name
func IsCollectionToIndex(collection string) bool {
	indexes := ConfigCollectionIndexes(*config, collection)

	return indexes != nil
}
