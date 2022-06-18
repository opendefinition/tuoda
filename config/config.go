package config

import (
	"encoding/json"
	"errors"
	"io/ioutil"
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

func LoadConfiguration() (Configuration, error) {
	user_home, home_err := os.UserHomeDir()

	if home_err != nil {
		return Configuration{}, errors.New("Unable to retrive user home directory")
	}

	tuoda_folder := path.Join(user_home, "tuoda")

	_, folder_err := os.Stat(tuoda_folder)

	if folder_err != nil {
		create_folder_err := os.Mkdir(tuoda_folder, 0755)

		if create_folder_err != nil {
			return Configuration{}, errors.New("Unable to create Tuoda configuration folder")
		}
	}

	config_path := path.Join(tuoda_folder, "config.json")

	_, config_err := os.Stat(config_path)

	if config_err != nil {
		empty_config, _ := json.MarshalIndent(Configuration{}, "", "\t")
		os.WriteFile(config_path, empty_config, 0755)

		return Configuration{}, errors.New("No application configuration file found. I have created a new empty one. Fill it out and play again!")
	}

	json_data, err := ioutil.ReadFile(config_path)

	if err != nil {
		return Configuration{}, errors.New("Unable to open configuration file. Please check if Tuoda folder and config.json exists!")
	}

	config := Configuration{}

	json.Unmarshal(json_data, &config)

	return config, nil
}
