package main

import (
	"encoding/json"
	"fmt"
	"log" // TODO: remove
	"os"
)

type Config struct {
	ExeFilepath string `json:"exe_filepath"`
	Changes     []struct {
		Description string   `json:"description"`
		Inputs      []string `json:"inputs"`
		Ouput       string   `json:"ouput"`
	} `json:"changes"`
}

func main() {
	config := parseConfig("./config.json")
	fmt.Println(config)
}

func parseConfig(filepath string) Config {
	var config Config

	data, err := os.ReadFile(filepath)
	if err != nil {
		log.Fatalf("Failed to read config: %v", err)
	}

	if err := json.Unmarshal(data, &config); err != nil {
		log.Fatalf("parse: %v", err)
	}

	return config
}
