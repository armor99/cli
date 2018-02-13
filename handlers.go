package main

import (
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
	"strings"
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

	// TODO: have valid access token? Yes, then done, no,continue.

	// TODO: have valid refresh token? Yes, then get new access, no, continue

	// TODO: call API and print result.
	//u, _ := url.ParseRequestURI("127.0.0.1:3000/api/v1/auth/token")
	var u url.URL
	u.Host = "127.0.0.1:3000"
	u.Path = "api/v1/auth/token"
	var client *http.Client
	req, err := http.NewRequest("POST", u.String(), strings.NewReader("test"))
	if err != nil {
		log.Fatal(err)
	}
	//req.Header.Add("If-None-Match", `W/"wyzzy"`)
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalln("Error making http request: %v\n", err)
	}
	defer resp.Body.Close()
	log.Println(resp.Status)

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
