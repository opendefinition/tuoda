package neo4jdb

import (
	"fmt"
	"log"
	"strings"

	"github.com/neo4j/neo4j-go-driver/v4/neo4j"
)

type Neo4JDB struct {
	Client neo4j.Driver
}

func NewClient(address string, username string, password string) *Neo4JDB {
	if len(address) == 0 || len(username) == 0 || len(password) == 0 {
		log.Fatalf("Empty Neo4J settings observed")
	}

	neo4jdb := new(Neo4JDB)
	driver, err := neo4j.NewDriver(address, neo4j.BasicAuth(username, password, ""))

	if err != nil {
		panic(err)
	}

	neo4jdb.Client = driver

	//defer neo4jdb.Client.Close()

	return neo4jdb
}

func (n4j *Neo4JDB) InsertLogItem(collection string, item map[string]interface{}) {
	session := n4j.Client.NewSession(neo4j.SessionConfig{})
	defer session.Close()

	itemmap := []string{}

	for key, value := range item {
		if key == "_key" {
			item["entryid"] = value
			key = "entryid"
			delete(item, "_key")
		}

		itemmap = append(itemmap, fmt.Sprintf("%s: $%s", key, key))
	}

	createstatement := fmt.Sprintf("CREATE (a:%s { %s })", collection, strings.Join(itemmap, ", "))

	_, err := session.WriteTransaction(func(tx neo4j.Transaction) (interface{}, error) {
		result, err := tx.Run(createstatement, item)

		if err != nil {
			log.Fatalf("Unable to insert document: %v", err)
		} else {
			fmt.Println("[ ", item["docsum"], " ] State: inserted")
		}

		return result.Consume()
	})

	if err != nil {
		panic(err)
	}
}
