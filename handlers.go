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
	"strconv"
	"strings"
)

func logoutHandler(l *flag.FlagSet, p params) {
	c, _ := readConfig()

	if *p.CID == 0 && c.CustomerID == 0 {
		log.Println("Subcommand login: Customer ID is required")
		l.PrintDefaults()
		os.Exit(1)
	}

	if *p.CID != 0 {
		c.CustomerID = *p.CID
	}
	if *p.User != "" {
		c.UserID = *p.User
	}
	if !validToken(c.Rtoken) {
		log.Fatalln("User not logged in")
	}
	if !validToken(c.Atoken) {
		res, err := refreshToken(c)
		if err != nil {
			log.Fatalln(err)
		}
		err = writeConfig(res)
		if err != nil {
			log.Fatalln(err)
		}
		c.Atoken = res.Atoken
		c.Rtoken = res.Rtoken
	}

	tid, err := getTID(c.Rtoken)
	if err != nil {
		log.Fatalln(err)
	}
	path := "api/v1/auth/token/" + strconv.Itoa(tid)
	log.Println(path)
	u := url.URL{Scheme: "http", Host: "127.0.0.1:3000", Path: path}
	client := &http.Client{}
	req, err := http.NewRequest("DELETE", u.String(), nil)
	if err != nil {
		log.Fatalf("Error w/ DELETE token query: %s\n", err)
	}
	auth := "Bearer " + c.Atoken
	req.Header.Add("Authorization", auth)
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalln(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		log.Fatalln(resp.Status)
	}
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}

	var res returnMsg
	log.Println(string(b))
	err = json.Unmarshal(b, &res)
	if err != nil {
		log.Println(err)
	}

	c.Atoken = ""
	c.Rtoken = ""
	err = writeConfig(c)
	if err != nil {
		log.Fatalln(err)
	}

	return
}

func loginHandler(l *flag.FlagSet, p params) {

	// TODO: Change all fmt to log and exit(1) to log.fatal
	// Get config values from file if exists
	c, _ := readConfig()

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
	if validToken(c.Atoken) {
		fmt.Printf("%s is logged in\n", c.UserID)
		os.Exit(0)
	}
	fmt.Println("Access token not valid")
	if validToken(c.Rtoken) {
		c, err := refreshToken(c)
		if err != nil {
			goto End
		}
		err = writeConfig(c)
		if err != nil {
			log.Println(err)
			os.Exit(1)
		}
		fmt.Printf("%s is logged in\n", c.UserID)
		os.Exit(0)
	End:
	}
	fmt.Println("Refresh token not valid")
	if *p.Passwd == "" {
		fmt.Println("Subcommand login: Password is required")
		l.PrintDefaults()
		os.Exit(1)
	}

	err := userLogin(&c, *p.Passwd)
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
	err = writeConfig(c)
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
	fmt.Printf("%s is logged in\n", c.UserID)

	return
}

func addUserHandler(l *flag.FlagSet, p userParams) {
	c, _ := readConfig()
	if *p.CID == 0 {
		fmt.Println("Subcommand login: Customer ID is required")
		l.PrintDefaults()
		os.Exit(1)
	}
	if *p.Email == "" {
		fmt.Println("Subcommand login: User email is required")
		l.PrintDefaults()
		os.Exit(1)
	}
	if *p.Role == "" {
		*p.Role = "user"
	}

	if !validToken(c.Rtoken) {
		fmt.Println("User not logged in.")
		os.Exit(0)
	}
	if !validToken(c.Atoken) {
		c, err := refreshToken(c)
		if err != nil {
			log.Fatalln(err)
		}
		err = writeConfig(c)
		if err != nil {
			log.Fatalln(err)
		}
	}

	// TODO: validate address is valid JSON

	// TODO: validate group ID is valid JSON

	// TODO: validate custom attributes is valid JSON

	// TODO: make API call using real params

	apiBody := `
	{
		"customer_id": 1,
		"user_id": "mike98",
		"email": "mmillsap98@cox.net",
		"firstname": "Mike",
		"lastname": "Miller",
		"address": {
		  "address1": "3103 E Killarney St",
		  "address2": "",
		  "city": "Gilbert",
		  "state": "AZ",
		  "zip": 85298
		},
		"group_id": [
		  12345,
		  12346,
		  12347,
		  12348
		],
		"role": "user",
		"custom_attr": {"name":"sap"}
	}`

	u := url.URL{Scheme: "http", Host: "127.0.0.1:3000", Path: "api/v1/user"}
	client := &http.Client{}
	req, err := http.NewRequest("POST", u.String(), strings.NewReader(apiBody))
	if err != nil {
		log.Fatalf("Error w/ adduser query: %s\n", err)
	}
	auth := "Bearer " + c.Atoken
	req.Header.Add("Authorization", auth)

	resp, err := client.Do(req)
	if err != nil {
		log.Fatalln(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		log.Fatalln(resp.Status)
	}
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}

	var res returnMsg
	log.Println(string(b))
	err = json.Unmarshal(b, &res)
	if err != nil {
		log.Println(err)
	}

	return
}
