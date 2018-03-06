package main

import (
	"flag"
	"fmt"
	"os"
	"strings"
)

func main() {
	// login subcommand setup
	var loginParams params
	loginCmd := flag.NewFlagSet("login", flag.ExitOnError)
	loginParams.CID = loginCmd.Int("cid", 0, "Customer ID <required>")
	loginParams.User = loginCmd.String("u", "", "User ID <required>")
	loginParams.Passwd = loginCmd.String("p", "", "User password <required>")
	loginParams.URL = loginCmd.String("h", "", "API URL <required>")

	// logout subcommand setup
	var logoutParams params
	logoutCmd := flag.NewFlagSet("logout", flag.ExitOnError)
	logoutParams.CID = logoutCmd.Int("cid", 0, "Customer ID")
	logoutParams.User = logoutCmd.String("u", "", "User ID")

	// add-user subcommand setup
	var addUserParams userParams
	addUserCmd := flag.NewFlagSet("add-user", flag.ExitOnError)
	addUserParams.CID = addUserCmd.Int("cid", 0, "Customer ID")
	addUserParams.User = addUserCmd.String("u", "", "User ID")
	addUserParams.Email = addUserCmd.String("e", "", "User email")
	addUserParams.Role = addUserCmd.String("r", "", "User role")
	addUserParams.Firstname = addUserCmd.String("f", "", "User first name")
	addUserParams.Lastname = addUserCmd.String("l", "", "User last name")
	addUserParams.Address = addUserCmd.String("a", "", "User address")
	addUserParams.GroupID = addUserCmd.String("g", "", "User's group IDs")
	addUserParams.CustomAttr = addUserCmd.String("c", "", "User's custom attributes")

	if len(os.Args) == 1 {
		fmt.Println(`Subcommand required:
 - login
 - add-user`)
		os.Exit(0)
	}

	cmd := strings.ToLower(os.Args[1])
	switch cmd {
	case "login":
		loginCmd.Parse(os.Args[2:])
	case "logout":
		logoutCmd.Parse(os.Args[2:])
	case "add-user":
		addUserCmd.Parse(os.Args[2:])
	case "-v", "--v", "-version", "--version":
		fmt.Println("CLI version: 1.0")
	default:
		fmt.Println("Invalid subcommand or option.")
	}

	// Execute subcommand
	if loginCmd.Parsed() {
		loginHandler(loginCmd, loginParams)
	}
	if logoutCmd.Parsed() {
		logoutHandler(logoutCmd, logoutParams)
	}
	if addUserCmd.Parsed() {
		addUserHandler(addUserCmd, addUserParams)
	}
}
