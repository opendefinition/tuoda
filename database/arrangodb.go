package database

import (
	"fmt"
	"log"

	"github.com/arangodb/go-driver"
	"github.com/arangodb/go-driver/http"
)

type ArrangoDB struct {
	Connection driver.Connection
	Client     driver.Client
	Database   driver.Database
}

func ArrangoDBClient(address string, database string, username string, password string) *ArrangoDB {
	arrango := new(ArrangoDB)
	var err error

	arrango.Connection, err = http.NewConnection(http.ConnectionConfig{
		Endpoints: []string{address},
	})

	if err != nil {
		log.Fatalf("Failed to create HTTP connection: %v", err)
	}

	arrango.Client, err = driver.NewClient(driver.ClientConfig{
		Connection: arrango.Connection,
		Authentication: driver.BasicAuthentication(
			username,
			password,
		),
	})

	db_exists, err := arrango.Client.DatabaseExists(nil, database)

	if db_exists {
		arrango.Database, err = arrango.Client.Database(nil, database)

		if err != nil {
			log.Fatalf("Failed to open database: %v", err)
		}
	} else {
		arrango.Database, err = arrango.Client.CreateDatabase(nil, database, nil)

		if err != nil {
			log.Fatalf("Failed to create database: %v", err)
		}
	}

	return arrango
}

func (ac *ArrangoDB) InsertLogItem(collection string, item map[string]interface{}) {
	_, err := ac.Database.CollectionExists(nil, collection)

	if err != nil {
		fmt.Println("Collection already exists: %v", err)
	} else {
		var col driver.Collection
		col, err := ac.Database.CreateCollection(nil, collection, nil)

		if err != nil {
			col, err = ac.Database.Collection(nil, collection)

			if err != nil {
				log.Fatalf("Error loading collection: %v", err)
			}
		}

		// Insert document
		_, err = col.CreateDocument(nil, item)

		if err != nil {
			log.Fatalf("Unable to insert document: %v", err)
		} else {
			// fmt.Println("Document inserted!")
		}
	}
}
