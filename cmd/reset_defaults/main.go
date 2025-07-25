package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	resources_dir := filepath.Join(exeDir(), "default_resources")

	entries, err := os.ReadDir(resources_dir)
	if err != nil {
		log.Fatalf("Failed to find default resources directory: %v", err)
	}

	for _, entry := range entries {
		err = os.RemoveAll(filepath.Join(resources_dir, entry.Name()))
		if err != nil {
			log.Fatalf("Failed to remove file: %v", err)
		}
	}
}

func exeDir() string {
	exePath, err := os.Executable()
	if err != nil {
		panic(err)
	}

	exeDir := filepath.Dir(exePath)

	if strings.Contains(exeDir, "go-build") {
		fmt.Println("using dev fallback")
		wd, _ := os.Getwd()
		return wd
	}

	return exeDir
}
