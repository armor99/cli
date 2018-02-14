package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/user"
	"path/filepath"
)

func loginHandler(l *flag.FlagSet, p params) {
	// Get config values from file if exists
	usr, err := user.Current()
	if err != nil {
		log.Fatal(err)
	}
	dirPath := filepath.Join(usr.HomeDir, ".idaas")
	configPath := filepath.Join(dirPath, "config.json")
	content, _ := ioutil.ReadFile(configPath)

	var c config
	err = json.Unmarshal(content, &c)
	if err != nil {
		log.Printf("Error parsing configuration file: %v", err)
	}

	if *p.CID == 0 && c.CustomerID == 0 {
		fmt.Println("Subcommand login: Customer ID is required")
		l.PrintDefaults()
		os.Exit(1)
	}
	if *p.CID != 0 {
		c.CustomerID = *p.CID
	}
	if *p.User == "" && c.UserID == "" {
		fmt.Println("Subcommand login: User ID is required")
		l.PrintDefaults()
		os.Exit(1)
	}
	if *p.User != "" {
		c.UserID = *p.User
	}
	if *p.URL == "" && c.URL == "" {
		fmt.Println("Subcommand login: API URL is required")
		l.PrintDefaults()
		os.Exit(1)
	}
	if *p.URL != "" {
		c.URL = *p.URL
	}

	// TODO: Password should only be required if tokens expired
	if *p.Passwd == "" {
		fmt.Println("Subcommand login: Password is required")
		l.PrintDefaults()
		os.Exit(1)
	}

	// TODO: Exit if access token valid

	// TODO: Exit if refresh token valid

	body := login{CustomerID: c.CustomerID, IP: "127.0.0.1"}
	bodyJSON, err := json.Marshal(body)
	if err != nil {
		log.Fatal(err)
	}
	u := url.URL{Scheme: "http", Host: "127.0.0.1:3000", Path: "api/v1/auth/token"}
	client := &http.Client{}
	req, err := http.NewRequest("POST", u.String(), bytes.NewReader(bodyJSON))
	if err != nil {
		log.Fatal(err)
	}
	auth := encodeBasicAuth(c.UserID, *p.Passwd)
	// TODO: remove
	log.Printf("Authorization: %s\n", auth)
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", auth)
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("Error making http request: %v\n", err)
	}
	defer resp.Body.Close()
	// TODO: remove below
	log.Printf("Response status: %v\n", resp.Status)
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("Error reading response body: %v\n", err)
	}
	log.Println(string(b))

	// TODO: Unmarshal body (b) to obtain access & refresh tokens
	// Create directory if doesn't exist
	_, err = os.Stat(dirPath)
	if os.IsNotExist(err) {
		err = os.MkdirAll(dirPath, 0755)
		if err != nil {
			log.Fatalf("Error creating directory: %s\n", err)
		}
	}
	// TODO: Should only save values after successful login
	content, err = json.MarshalIndent(c, "", "   ")
	err = ioutil.WriteFile(configPath, content, 0644)
	if err != nil {
		log.Printf("Error writing config file: %v\n", err)
	}

	return
}
