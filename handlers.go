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
	var u userNew
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
	if !validToken(c.Rtoken) {
		fmt.Println("User not logged in.")
		os.Exit(0)
	}
	if !validToken(c.Atoken) {
		newToken, err := refreshToken(c)
		if err != nil {
			log.Fatalln(err)
		}
		c = newToken
		err = writeConfig(c)
		if err != nil {
			log.Fatalln(err)
		}
	}

	// Validate input
	u.CustomerID = *p.CID
	u.Email = *p.Email
	if *p.User == "" {
		u.UserID = u.Email
	} else {
		u.UserID = *p.User
	}
	if *p.Role == "" {
		*p.Role = "user"
	}
	u.Role = *p.Role
	u.Firstname = *p.Firstname
	u.Lastname = *p.Lastname
	err := json.Unmarshal([]byte(*p.Address), &u.Address)
	if err != nil && *p.Address != "" {
		log.Fatalln("Address not valid JSON.")
	}
	err = json.Unmarshal([]byte(*p.GroupID), &u.GroupID)
	if err != nil && *p.GroupID != "" {
		log.Fatalln("Group ID not valid JSON array.")
	}
	err = json.Unmarshal([]byte(*p.CustomAttr), &u.CustomAttr)
	if err != nil && *p.CustomAttr != "" {
		log.Fatalln("Customer Attributes not valid JSON.")
	}

	body, err := json.Marshal(u)
	if err != nil {
		log.Fatalln("Error converting parameters to JOSN.")
	}
	apiURL := url.URL{Scheme: "http", Host: "127.0.0.1:3000", Path: "api/v1/user"}
	client := &http.Client{}
	req, err := http.NewRequest("POST", apiURL.String(), strings.NewReader(string(body)))
	if err != nil {
		log.Fatalf("Error w/ adduser API call: %s\n", err)
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

	var res addUserRetMsg
	err = json.Unmarshal(b, &res)
	if err != nil {
		log.Printf("Error parsing JSON response from API call. %v\n", err)
	}
	// TODO: Return hash going forward or email to user?
	log.Printf("\n\nUser %s created. PW Hash:\n\n%s\n\n", res.Data.UserID, res.Data.Hash)

	return
}

func listUserHandler(l *flag.FlagSet, p userParams) {

	// TODO: sysadmin needs to be able to list users for any customerID. Enable this via the list users endpoint.

	c, _ := readConfig()
	if !validToken(c.Rtoken) {
		fmt.Println("User not logged in.")
		os.Exit(0)
	}
	if !validToken(c.Atoken) {
		newToken, err := refreshToken(c)
		if err != nil {
			log.Fatalln(err)
		}
		c = newToken
		err = writeConfig(c)
		if err != nil {
			log.Fatalln(err)
		}
	}

	// TODO: Is user checked to be admin before calling API?
	apiURL := url.URL{Scheme: "http", Host: "127.0.0.1:3000", Path: "api/v1/user"}
	client := &http.Client{}
	req, err := http.NewRequest("GET", apiURL.String(), nil)
	if err != nil {
		log.Fatalf("Error w/ listuser API call: %s\n", err)
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

	log.Println(string(b))

	// TODO: Maybe add feature to automatically page results to screen. Keep pausing for keypress.

	var res addUserRetMsg
	err = json.Unmarshal(b, &res)
	if err != nil {
		log.Printf("Error parsing JSON response from API call. %v\n", err)
	}
	// TODO: Return hash going forward or email to user?
	log.Printf("\n\nUser %s created. PW Hash:\n\n%s\n\n", res.Data.UserID, res.Data.Hash)

	return
}
