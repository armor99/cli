package main

import (
	"flag"
	"fmt"
	"log"
	"os"
)

func loginHandler(l *flag.FlagSet, p params) {
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
