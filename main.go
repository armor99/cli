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

	// add-user subcommand setup
	addUserCmd := flag.NewFlagSet("add-user", flag.ExitOnError)
	cid2 := addUserCmd.Int("cid", 0, "Customer ID")
	user2 := addUserCmd.String("u", "", "User ID")
	pw2 := addUserCmd.String("p", "", "User password")

	// TODO: Is this the right way to handler missing or invalid commands.
	if len(os.Args) == 1 {
		fmt.Println(`Subcommand required:
 - login
 - add-user`)
		os.Exit(0)
	}

	cmd := strings.ToLower(os.Args[1])
	//os.Args = append(os.Args[:1], os.Args[2:]...)
	switch cmd {
	case "login":
		loginCmd.Parse(os.Args[2:])
	case "-v", "--v", "-version", "--version":
		fmt.Println("CLI version: 1.0")
	default:
		fmt.Println("Invalid subcommand or option.")
	}

	// Execute subcommands
	if loginCmd.Parsed() {
		loginHandler(loginCmd, loginParams)
	}

	if addUserCmd.Parsed() {
		fmt.Println(*cid2)
		fmt.Println(*user2)
		fmt.Println(*pw2)
		os.Exit(0)
	}

	fmt.Println(cmd)
	fmt.Println(len(os.Args))
}
