package config

import (
	"encoding/json"
	"fmt"
	"os"
)

var CONF *Configuration

type Configuration struct {
	ActiveProfile int
	Profiles      Profiles
}

func CreateConfiguration() *Configuration {
	return &Configuration{
		ActiveProfile: 0,
		Profiles:      *createProfileList(),
	}
}

func createConfigFile() {
	WriteConfig(CreateConfiguration())
}

func doesConfigExist() bool {
	file, err := os.Open("config.json")
	if err != nil {
		return false
	}
	file.Close()
	return true
}

func readConfig() *Configuration {
	if !doesConfigExist() {
		createConfigFile()
	}
	file, _ := os.Open("config.json")
	decoder := json.NewDecoder(file)
	configuration := Configuration{}
	err := decoder.Decode(&configuration)
	if err != nil {
		fmt.Println("failed to read file:", err)
	}

	return &configuration
}

func InitConfig() {
	c := readConfig()
	CONF = c
}

func WriteConfig(c *Configuration) {
	// Open the file for writing and truncate it so we always rewrite the full contents.
	file, err := os.OpenFile("config.json", os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0600)
	if err != nil {
		fmt.Println("failed to open config file for write:", err)
		return
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "    ")
	if err := encoder.Encode(c); err != nil {
		fmt.Println("failed to write file:", err)
		return
	}

	// Sync
	if err := file.Sync(); err != nil {
		fmt.Println("failed to sync config file:", err)
	}
}
