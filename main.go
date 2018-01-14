package main

import (
	"flag"
	"fmt"
	"os"
	"strings"
)

func main() {
	loginCmd := flag.NewFlagSet("login", flag.ExitOnError)
	cid := loginCmd.Int("cid", 0, "Customer ID <required>")
	user := loginCmd.String("u", "", "User ID <required>")
	pw := loginCmd.String("p", "", "User password <required>")

	addUserCmd := flag.NewFlagSet("add-user", flag.ExitOnError)
	cid2 := addUserCmd.Int("cid", 0, "Customer ID")
	user2 := addUserCmd.String("u", "", "User ID")
	pw2 := addUserCmd.String("p", "", "User password")

	if len(os.Args) == 1 {
		fmt.Println("Missing subcommand or valid option.")
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

	if loginCmd.Parsed() {
		if *cid == 0 {
			fmt.Println("Subcommand login: Customer ID is required")
			loginCmd.PrintDefaults()
			os.Exit(1)
		}
		if *user == "" {
			fmt.Println("Subcommand login: User ID is required")
			loginCmd.PrintDefaults()
			os.Exit(1)
		}
		if *pw == "" {
			fmt.Println("Subcommand login: Password is required")
			loginCmd.PrintDefaults()
			os.Exit(1)
		}

		fmt.Println(*cid)
		fmt.Println(*user)
		fmt.Println(*pw)
		os.Exit(0)
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
