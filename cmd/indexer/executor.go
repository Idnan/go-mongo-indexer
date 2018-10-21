package main

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/globalsign/mgo"
	"github.com/idnan/go-mongo-indexer/pkg/util"
	"gopkg.in/mgo.v2/bson"
	"log"
	"reflect"
)

const GB1 = 1000000000

// Execute the command
func execute() {

	indexDiff := getIndexesDiff()

	if !*apply {
		showDiff(indexDiff)
	}

	if *apply {
		applyDiff(indexDiff)
	}
}

// Drop and apply the indexes
func applyDiff(indexDiff *IndexDiff) {
	for _, collection := range Collections() {
		indexesToRemove := indexDiff.Old[collection]
		indexesToAdd := indexDiff.New[collection]
		capToAdd := indexDiff.Cap[collection]

		util.PrintBold(fmt.Sprintf("\n%s.%s\n", db.Name, collection))

		if indexesToRemove == nil && indexesToAdd == nil && capToAdd == 0 {
			util.PrintGreen(fmt.Sprintln("No index changes"))
			continue
		}

		if capToAdd != 0 {
			util.PrintGreen(fmt.Sprintf("+ Adding cap of %d\n", capToAdd))
			SetCapSize(collection, capToAdd)
		}

		for indexName, columns := range indexesToRemove {
			util.PrintRed(fmt.Sprintf("- Dropping index %s: %s\n", indexName, util.JsonEncode(columns)))
			DropIndex(collection, indexName)
		}

		for indexName, columns := range indexesToAdd {
			util.PrintGreen(fmt.Sprintf("+ Adding index %s: %s\n", indexName, util.JsonEncode(columns)))
			CreateIndex(collection, indexName, columns)
		}
	}
}

// Create index of on the given collection with index name and columns
func CreateIndex(collection string, indexName string, columns []string) bool {
	index := mgo.Index{
		Key:              columns,
		Background:       true,
		Name:             indexName,
		LanguageOverride: "search_lang",
	}

	err := db.C(collection).EnsureIndex(index)

	if err != nil {
		log.Fatalln(err.Error())
	}

	return true
}

// Drop an index by name from given collection
func DropIndex(collection string, indexName string) bool {
	err := db.C(collection).DropIndexName(indexName)

	if err != nil {
		log.Fatalln(err.Error())
	}

	return true
}

// Show the index difference, the indexes with `-` will be deleted only
// the ones with the `+` will be created
func showDiff(indexDiff *IndexDiff) {

	for _, collection := range Collections() {
		indexesToRemove := indexDiff.Old[collection]
		indexesToAdd := indexDiff.New[collection]
		capToAdd := indexDiff.Cap[collection]

		util.PrintBold(fmt.Sprintf("\n%s.%s\n", db.Name, collection))

		if indexesToRemove == nil && indexesToAdd == nil && capToAdd == 0 {
			util.PrintGreen(fmt.Sprintln("No index changes"))
			continue
		}

		if capToAdd != 0 {
			util.PrintGreen(fmt.Sprintf("+ Capsize to set: %d\n", capToAdd))
		}

		for indexName, columns := range indexesToRemove {
			util.PrintRed(fmt.Sprintf("- %s: %s\n", indexName, util.JsonEncode(columns)))
		}

		for indexName, columns := range indexesToAdd {
			util.PrintGreen(fmt.Sprintf("+ %s: %s\n", indexName, util.JsonEncode(columns)))
		}
	}
}

// Match existing indexes with the given config file and match and find the diff
// the indexes that are not inside the config will be deleted, only the indexes in
// the config file will be created
func getIndexesDiff() *IndexDiff {

	oldIndexes := make(map[string]map[string][]string)
	newIndexes := make(map[string]map[string][]string)
	capSize := make(map[string]int)

	for _, collection := range Collections() {

		var alreadyAppliedIndexesColumns []interface{}
		var alreadyAppliedIndexesNames []string
		var givenIndexes [][]string

		configCollection := GetConfigCollection(collection)
		if configCollection != nil {
			givenIndexes = configCollection.Indexes
		}

		// Get current database collection indexes
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

		// Get the config cap size
		givenCapSize := GetConfigCollectionCapSize(collection)
		isAlreadyCapped := IsCollectionCaped(collection)

		minAllowedCapSize := GB1 / 2

		// Add the cap size
		if givenCapSize > 0 && (givenCapSize >= minAllowedCapSize) && !isAlreadyCapped {
			capSize[collection] = givenCapSize
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

	return &IndexDiff{oldIndexes, newIndexes, capSize}
}

// Generate index name by doing md5 of indexes json
func GenerateIndexName(indexColumns interface{}) string {
	content, _ := json.Marshal(indexColumns)
	algorithm := md5.New()
	algorithm.Write(content)

	return hex.EncodeToString(algorithm.Sum(nil))
}

// Return list of database collections
func Collections() []string {
	collections, err := db.CollectionNames()

	if err != nil {
		log.Fatalln(err.Error())
	}

	return collections
}

// Return database collection indexes
func DbIndexes(collection string) map[string][]string {
	indexes, err := db.C(collection).Indexes()

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

// Check if the collection is already capped
func IsCollectionCaped(collection string) bool {
	var doc bson.M
	err := db.Run(map[string]string{"collStats": collection}, &doc)
	if err != nil {
		log.Fatalln(err.Error())
	}

	return doc["capped"].(bool)
}

func SetCapSize(collection string, size int) bool {
	var doc bson.M
	err := db.Run(map[string]interface{}{"convertToCapped": collection, "size": size}, &doc)

	if err != nil {
		log.Fatalln(err.Error())
	}

	return true
}
