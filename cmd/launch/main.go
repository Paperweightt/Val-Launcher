package main

import (
	"encoding/json"
	"fmt" // TODO: remove
	"io"
	"log" // TODO: remove
	"math/rand"
	"os"
	"os/exec"
	"time"
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
	// TODO: change in final build
	configFilePath := "C:/Users/henry/Projects/Val_Launcher/config/default.json"
	config := parseConfig(configFilePath)
	s := rand.NewSource(time.Now().Unix())
	r := rand.New(s) // initialize local pseudorandom generator

	cmd := exec.Command(config.ExeFilepath)
	_, err := cmd.Output()
	if err != nil {
		log.Fatalf("Failed to run Valorant: %v", err)
	}

	time.Sleep(5 * time.Second) // Pause for 5 seconds

	for _, change := range config.Changes {
		source := change.Inputs[r.Intn(len(change.Inputs))]
		copyFile(source, change.Ouput)
	}
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

func copyFile(source, destination string) error {
	sourceFile, err := os.Open(source)
	if err != nil {
		return fmt.Errorf("failed to open source file: %w", err)
	}
	defer sourceFile.Close() // Ensure source file is closed

	destinationFile, err := os.Create(destination)
	if err != nil {
		return fmt.Errorf("failed to create destination file: %w", err)
	}
	defer destinationFile.Close() // Ensure destination file is closed

	_, err = io.Copy(destinationFile, sourceFile)
	if err != nil {
		return fmt.Errorf("failed to copy file contents: %w", err)
	}

	return nil
}
