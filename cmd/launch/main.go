package main

import (
	"encoding/json"
	"fmt" // TODO: remove
	"io"
	"log" // TODO: remove
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
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
	config := parseConfig("./config.json")
	s := rand.NewSource(time.Now().Unix())
	r := rand.New(s)

	resetValorantState(config)
	time.Sleep(3 * time.Second) // Pause for 5 seconds

	// run valorant
	cmd := exec.Command(config.ExeFilepath)
	_, err := cmd.Output()
	if err != nil {
		log.Fatalf("Failed to run Valorant: %v", err)
	}

	time.Sleep(3 * time.Second) // Pause for 5 seconds

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

func resetValorantState(config Config) error {
	var m = map[string]bool{}

	entries, err := os.ReadDir("./default_resources/")
	if err != nil {
		return err
	}

	for _, entry := range entries {
		if !entry.IsDir() {
			m[entry.Name()] = true
		}
	}

	// check if default resource exists and if not save it
	for _, change := range config.Changes {
		if _, found := m[filepath.Base(change.Ouput)]; !found {
			destination := filepath.Join("./default_resources/", filepath.Base(change.Ouput))

			fmt.Println("added default: ", destination)
			err := deepCopy(change.Ouput, destination)
			if err != nil {
				return err
			}
		}
	}

	// set default resources into game files
	for _, change := range config.Changes {
		base := filepath.Base(change.Ouput)
		err := deepCopy(filepath.Join("./default_resources/", base), change.Ouput)
		// err := deepCopy("C:/Users/henry/Projects/Val_Launcher/resources/red_dress_1.mp4", change.Ouput)
		if err != nil {
			return err
		}
	}

	return nil
}

// CopyFileWithMetadata copies a file from src to dst, preserving content, permissions, and timestamps.
func deepCopy(src, dst string) error {
	// Open source file
	srcFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	// Get source file info
	srcInfo, err := srcFile.Stat()
	if err != nil {
		return err
	}

	// Create destination file (overwrite if it exists)
	dstFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer dstFile.Close()

	// Copy contents
	if _, err := io.Copy(dstFile, srcFile); err != nil {
		return err
	}

	// Copy file permissions
	// if err := os.Chmod(dst, srcInfo.Mode()); err != nil {
	// 	return err
	// }

	// required
	// Copy timestamps (atime is not available via os.Stat, so we reuse mod time for both)
	modTime := srcInfo.ModTime()
	if err := os.Chtimes(dst, modTime, modTime); err != nil {
		return err
	}

	return nil
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
