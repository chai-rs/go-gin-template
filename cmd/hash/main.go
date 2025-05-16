package main

import (
	"flag"
	"fmt"

	"github.com/chai-rs/simple-bookstore/pkg/crypto"
)

func main() {
	v := flag.String("v", "", "value to hash")
	flag.Parse()

	if *v == "" {
		fmt.Println("-v flag is required")
		return
	}

	hashedPassword, err := crypto.HashPassword(*v)
	if err != nil {
		fmt.Println("failed to hash password:", err)
		return
	}

	fmt.Println(hashedPassword)
}
