//go:build ignore
// +build ignore

package main

import (
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"log"
)

func main() {
	hash, err := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(string(hash))
}
