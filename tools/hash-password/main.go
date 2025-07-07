package main

import (
	"flag"
	"fmt"
	"os"

	"golang.org/x/crypto/bcrypt"
)

func main() {
	var password string
	flag.StringVar(&password, "p", "", "Password to hash")
	flag.Parse()

	if password == "" {
		fmt.Println("Usage: go run tools/hash-password/main.go -p <password>")
		os.Exit(1)
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		fmt.Printf("Error generating hash: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Bcrypt hash for password '%s':\n%s\n", password, string(hash))
}
