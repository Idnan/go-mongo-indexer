package main

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"reflect"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/idnan/go-mongo-indexer/pkg/util"
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

		util.PrintBold(fmt.Sprintf("\n%s.%s\n", db.Name(), collection))

		if indexesToRemove == nil && indexesToAdd == nil && capToAdd == 0 {
			util.PrintGreen(fmt.Sprintln("No index changes"))
			continue
		}

		if capToAdd != 0 {
			util.PrintGreen(fmt.Sprintf("+ Adding cap of %d\n", capToAdd))
			SetCapSize(collection, capToAdd)
		}

		for _, index := range indexesToRemove {
			util.PrintRed(fmt.Sprintf("- Dropping index %s: %s\n", index.Name, util.JsonEncode(index.Keys)))
			DropIndex(collection, index.Name)
		}

		for _, index := range indexesToAdd {
			util.PrintGreen(fmt.Sprintf("+ Adding index %s: %s\n", index.Name, util.JsonEncode(index.Keys)))
			CreateIndex(collection, index.Name, index)
		}
	}
}

// Create index of on the given collection with index Name and columns
func CreateIndex(collection string, indexName string, indexModel IndexModel) bool {

	keys := indexModel.Keys
	background := true
	languageOverride := "search_lang"

	_unique, exists := keys["_unique"]
	if !exists {
		_unique = 0
	}

	unique := false
	if _unique == 1 {
		unique = true
	}

	_expireAfterSeconds, exists := keys["_expireAfterSeconds"]
	if !exists {
		_expireAfterSeconds = 0
	}

	// setting options
	opts := &options.IndexOptions{
		Unique:           &unique,
		Background:       &background,
		Name:             &indexName,
		LanguageOverride: &languageOverride,
	}
	expireAfterSeconds := int32(_expireAfterSeconds)
	if expireAfterSeconds > 0 {
		opts.ExpireAfterSeconds = &expireAfterSeconds
	}

	// remove the non index fields
	delete(keys, "_unique")
	delete(keys, "_expireAfterSeconds")

	index := mongo.IndexModel{
		Keys:    keys,
		Options: opts,
	}

	indexView := db.Collection(collection).Indexes()

	_, err := indexView.CreateOne(context.TODO(), index)

	if err != nil {
		log.Fatalln(err.Error())
	}

	return true
}

// Drop an index by Name from given collection
func DropIndex(collection string, indexName string) bool {
	indexes := db.Collection(collection).Indexes()
	_, err := indexes.DropOne(context.TODO(), indexName)

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

		util.PrintBold(fmt.Sprintf("\n%s.%s\n", db.Name(), collection))

		if indexesToRemove == nil && indexesToAdd == nil && capToAdd == 0 {
			util.PrintGreen(fmt.Sprintln("No index changes"))
			continue
		}

		if capToAdd != 0 {
			util.PrintGreen(fmt.Sprintf("+ Capsize to set: %d\n", capToAdd))
		}

		for _, index := range indexesToRemove {
			util.PrintRed(fmt.Sprintf("- %s: %s\n", index.Name, util.JsonEncode(index.Keys)))
		}

		for _, index := range indexesToAdd {
			util.PrintGreen(fmt.Sprintf("+ %s: %s\n", index.Name, util.JsonEncode(index.Keys)))
		}
	}
}

// Match existing indexes with the given config file and match and find the diff
// the indexes that are not inside the config will be deleted, only the indexes in
// the config file will be created
func getIndexesDiff() *IndexDiff {

	oldIndexes := make(map[string]map[string]IndexModel)
	newIndexes := make(map[string]map[string]IndexModel)
	capSize := make(map[string]int)

	for _, collection := range Collections() {

		var alreadyAppliedIndexesColumns []interface{}
		var alreadyAppliedIndexesNames []string
		var givenIndexes []map[string]int

		configCollection := GetConfigCollection(collection)

		if configCollection != nil {
			givenIndexes = configCollection.Indexes
		}

		// Get current database collection indexes
		currentIndexes := DbIndexes(collection)

		// If we don't have the current collection in the index create list then drop all index
		if !IsCollectionToIndex(collection) {
			for _, dbIndex := range currentIndexes {
				if oldIndexes[collection] == nil {
					oldIndexes[collection] = make(map[string]IndexModel)
				}
				oldIndexes[collection][dbIndex.Name] = dbIndex
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
		for _, dbIndex := range currentIndexes {

			isCurrentIndexInConfig := false

			for _, givenIndexColumns := range givenIndexes {

				// If the Name of index matches the Name of given index
				generatedIndexName := GenerateIndexName(givenIndexColumns)

				if dbIndex.Name == generatedIndexName {
					isCurrentIndexInConfig = true
					alreadyAppliedIndexesNames = append(alreadyAppliedIndexesNames, generatedIndexName)
					break
				}

				// First check if this column group has the index
				if reflect.DeepEqual(givenIndexColumns, dbIndex.Keys) {
					isCurrentIndexInConfig = true
					break
				}
			}

			if !isCurrentIndexInConfig {
				if oldIndexes[collection] == nil {
					oldIndexes[collection] = make(map[string]IndexModel)
				}
				oldIndexes[collection][dbIndex.Name] = dbIndex
			} else {
				alreadyAppliedIndexesColumns = append(alreadyAppliedIndexesColumns, dbIndex.Keys)
			}
		}

		// For each of the given indexes, check if it is already applied or not
		// If not, prepare a list so that those can be applied
		for _, givenIndexColumns := range givenIndexes {

			isAlreadyApplied := false

			// If the Name of index matches the Name of given index
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
					newIndexes[collection] = make(map[string]IndexModel)
				}
				newIndexes[collection][generatedIndexName] = IndexModel{generatedIndexName, givenIndexColumns}
			}
		}
	}

	return &IndexDiff{oldIndexes, newIndexes, capSize}
}

// Generate index Name by doing md5 of indexes json
func GenerateIndexName(indexColumns interface{}) string {
	content, _ := json.Marshal(indexColumns)
	algorithm := md5.New()
	algorithm.Write(content)

	return hex.EncodeToString(algorithm.Sum(nil))
}

// Return list of database collections
func Collections() []string {
	collections, err := db.ListCollectionNames(context.TODO(), bson.M{})

	if err != nil {
		log.Fatalln(err.Error())
	}

	return collections
}

// Return database collection indexes
func DbIndexes(collection string) []IndexModel {
	cursor, err := db.Collection(collection).Indexes().List(context.TODO())

	if err != nil {
		log.Fatalln(err.Error())
	}

	dbIndexes := make([]IndexModel, 0)

	for cursor.Next(context.TODO()) {

		index := bson.M{}
		if err := cursor.Decode(&index); err != nil {
			log.Fatalln(err.Error())
		}

		keys := map[string]int{}
		keysByte, _ := bson.Marshal(index["key"])
		if err := bson.Unmarshal(keysByte, &keys); err != nil {
			log.Fatalln(err)
		}

		// ignore the _id index as it's the default index
		if len(keys) == 0 || reflect.DeepEqual(keys, map[string]int{"_id": 1}) {
			continue
		}

		// check if there's a unique index or not
		_, exists := index["unique"]
		if exists {
			keys["_unique"] = 1
		}

		// check if there's a unique index or not
		expireAfterSeconds, exists := index["expireAfterSeconds"]
		if exists {
			keys["_expireAfterSeconds"] = int(expireAfterSeconds.(int32))
		}

		name := index["name"].(string)

		dbIndexes = append(dbIndexes, IndexModel{name, keys})
	}

	return dbIndexes
}

// Drop index from collection by index Name
func IsCollectionToIndex(collection string) bool {
	return GetConfigCollection(collection) != nil
}

// Check if the collection is already capped
func IsCollectionCaped(collection string) bool {

	command := map[string]string{"collStats": collection}
	result := db.RunCommand(context.TODO(), command)

	var doc bson.M
	if err := result.Decode(&doc); err != nil {
		log.Fatalln(err.Error())
	}

	return doc["capped"].(bool)
}

func SetCapSize(collection string, size int) bool {

	command := map[string]interface{}{"convertToCapped": collection, "size": size}
	result := db.RunCommand(context.TODO(), command)

	var doc bson.M
	if err := result.Decode(&doc); err != nil {
		log.Fatalln(err.Error())
	}

	return true
}
