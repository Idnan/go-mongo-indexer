package main

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/globalsign/mgo"
	"github.com/idnan/go-mongo-indexer/pkg/util"
	"log"
	"reflect"
)

type IndexDiff struct {
	Old map[string]map[string][]string
	New map[string]map[string][]string
}

func execute() {

	indexDiff := getIndexesDiff()

	if !*apply {
		showDiff(indexDiff)
	}

	if *apply {
		applyDiff(indexDiff)
	}
}

func applyDiff(indexDiff *IndexDiff) {
	for _, collection := range Collections() {
		indexesToRemove := indexDiff.Old[collection]
		indexesToAdd := indexDiff.New[collection]

		if indexesToRemove == nil && indexesToAdd == nil {
			fmt.Printf("\nNothing to change in %s!\n\n", collection)
			continue
		}

		fmt.Printf("\nApplying Changes: %s\n", collection)

		// @todo cap size

		for indexName, columns := range indexesToRemove {
			util.PrintRed(fmt.Sprintf("- Dropping index %s: %s\n", indexName, util.JsonEncode(columns)))
			dropIndex(collection, indexName)
		}

		for indexName, columns := range indexesToAdd {
			util.PrintGreen(fmt.Sprintf("+ Adding index %s: %s\n", indexName, util.JsonEncode(columns)))
			createIndex(collection, indexName, columns)
		}
	}
}

func createIndex(collection string, indexName string, columns []string) bool {
	index := mgo.Index{
		Key:              columns,
		Background:       true,
		Name:             indexName,
		LanguageOverride: "search_lang",
	}

	err := db.DB("oms_api").C(collection).EnsureIndex(index)

	if err != nil {
		log.Fatalln(err.Error())
	}

	return true
}

func dropIndex(collection string, indexName string) bool {
	err := db.DB("oms_api").C(collection).DropIndexName(indexName)

	if err != nil {
		log.Fatalln(err.Error())
	}

	return true
}

func showDiff(indexDiff *IndexDiff) {

	for _, collection := range Collections() {
		indexesToRemove := indexDiff.Old[collection]
		indexesToAdd := indexDiff.New[collection]

		if indexesToRemove == nil && indexesToAdd == nil {
			fmt.Printf("\nNothing to change in %s!\n\n", collection)
			continue
		}

		fmt.Printf("\n%s\n", collection)

		for indexName, columns := range indexesToRemove {
			util.PrintRed(fmt.Sprintf("- %s: %s\n", indexName, util.JsonEncode(columns)))
		}

		for indexName, columns := range indexesToAdd {
			util.PrintGreen(fmt.Sprintf("+ %s: %s\n", indexName, util.JsonEncode(columns)))
		}
	}
}

func getIndexesDiff() *IndexDiff {

	oldIndexes := make(map[string]map[string][]string)
	newIndexes := make(map[string]map[string][]string)

	for _, collection := range Collections() {

		var alreadyAppliedIndexesColumns []interface{}
		var alreadyAppliedIndexesNames []string
		var givenIndexes [][]string

		configCollection := GetConfigCollection(collection)
		if configCollection != nil {
			givenIndexes = configCollection.Indexes
		}

		// Get current indexes and cap-size
		currentIndexes := DbIndexes(collection)

		// If we don't have the current collection in the index create list then drop all index
		if !IsCollectionToIndex(collection) {
			for indexName, indexDetail := range currentIndexes {
				if oldIndexes[collection] == nil {
					oldIndexes[collection] = make(map[string][]string)
				}
				oldIndexes[collection][indexName] = indexDetail
			}
			continue
		}

		// Prepare the list of indexes that need to be dropped
		for currentIndexName, currentIndexColumns := range currentIndexes {

			isCurrentIndexInConfig := false

			for _, givenIndexColumns := range givenIndexes {

				// If the name of index matches the name of given index
				generatedIndexName := GenerateIndexName(givenIndexColumns)

				if currentIndexName == generatedIndexName {
					isCurrentIndexInConfig = true
					alreadyAppliedIndexesNames = append(alreadyAppliedIndexesNames, generatedIndexName)
					break
				}

				// First check if this column group has the index
				if reflect.DeepEqual(givenIndexColumns, currentIndexColumns) {
					isCurrentIndexInConfig = true
					break
				}
			}

			if !isCurrentIndexInConfig {
				if oldIndexes[collection] == nil {
					oldIndexes[collection] = make(map[string][]string)
				}
				oldIndexes[collection][currentIndexName] = currentIndexColumns
			} else {
				alreadyAppliedIndexesColumns = append(alreadyAppliedIndexesColumns, currentIndexColumns)
			}
		}

		// For each of the given indexes, check if it is already applied or not
		// If not, prepare a list so that those can be applied
		for _, givenIndexColumns := range givenIndexes {

			isAlreadyApplied := false

			// If the name of index matches the name of given index
			generatedIndexName := GenerateIndexName(givenIndexColumns)

			for _, appliedIndexColumns := range alreadyAppliedIndexesColumns {

				if util.StringInSlice(generatedIndexName, alreadyAppliedIndexesNames) {
					isAlreadyApplied = true
					break
				}

				if reflect.DeepEqual(givenIndexColumns, appliedIndexColumns) {
					isAlreadyApplied = true
					break
				}
			}

			if !isAlreadyApplied {
				if newIndexes[collection] == nil {
					newIndexes[collection] = make(map[string][]string)
				}
				newIndexes[collection][generatedIndexName] = givenIndexColumns
			}
		}
	}

	return &IndexDiff{oldIndexes, newIndexes}
}

func GenerateIndexName(indexColumns interface{}) string {
	content, _ := json.Marshal(indexColumns)
	algorithm := md5.New()
	algorithm.Write(content)
	return hex.EncodeToString(algorithm.Sum(nil))
}

func Collections() []string {
	collections, err := db.DB("oms_api").CollectionNames()

	if err != nil {
		log.Fatalln(err.Error())
	}

	return collections
}

func DbIndexes(collection string) map[string][]string {
	indexes, err := db.DB("oms_api").C(collection).Indexes()

	if err != nil {
		log.Fatalln(err.Error())
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
	return GetConfigCollection(collection) != nil
}
