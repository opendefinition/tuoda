package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"
)

type ArangoDBConfig struct {
	Address  string
	Username string
	Password string
	Database string
}

type Configuration struct {
	ArangoDB ArangoDBConfig
}

func LoadConfiguration() Configuration {
	user_home, err := os.UserHomeDir()

	if err != nil {
		log.Fatalf("Unable to retrive user home directory")
	}

	config_path := path.Join(user_home, "tuoda", "config.json")
	fmt.Println(config_path)

	json_data, err := ioutil.ReadFile(config_path)

	if err != nil {
		log.Fatalf("Unable to open configuration file")
	}

	config := Configuration{}
	json.Unmarshal(json_data, &config)

	return config
}
