package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/user"
	"path/filepath"
	"time"
)

func main() {
	t := time.Now()
	tStr := t.Format("2006-01-02 15:04:05")
	fmt.Printf("Timestamp: %s\n", tStr)

	usr, err := user.Current()
	if err != nil {
		log.Fatal(err)
	}
	homePath := usr.HomeDir
	dirPath := filepath.Join(homePath, ".idaas")
	configPath := filepath.Join(dirPath, "config.json")
	fmt.Printf("Path: %s\n", configPath)

	// Create directory if doesn't exist
	_, err = os.Stat(dirPath)
	if os.IsNotExist(err) {
		err = os.MkdirAll(dirPath, 0755)
		if err != nil {
			log.Fatalf("Error creating directory: %s\n", err)
		}
	}

	// Read file contents
	content, err := ioutil.ReadFile(configPath)
	if !os.IsNotExist(err) && err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Old file content: %s", content)

	// Overwrite file whether it exists or not.
	content = []byte(tStr)
	err = ioutil.WriteFile(configPath, content, 777)
	if err != nil {
		log.Fatal(err)
	}
}
