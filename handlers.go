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
	log.Println(tid)
	path := "api/v1/auth/token/" + strconv.Itoa(tid)
	log.Println(path)
	u := url.URL{Scheme: "http", Host: "127.0.0.1:3000", Path: path}
	client := &http.Client{}
	req, err := http.NewRequest("DELETE", u.String(), nil)
	if err != nil {
		log.Fatalln(err)
	}
	auth := "Bearer " + c.Atoken
	req.Header.Add("Authorization", auth)
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalln(err)
	}
	defer resp.Body.Close()

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}

	var res returnMsg
	err = json.Unmarshal(b, &res)
	if err != nil {
		log.Println(err)
	}
	if res.Status.Code != 200 {
		log.Println(res.Status.Message)
	}

	// TODO: need to delete tokens from config file on success

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
