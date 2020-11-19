package main

import (
	"fmt"
	"log"
	"os"
)

func main() {
	if len(os.Args) != 2 {
		panic("Invalid number of arguments")
	}

	word := os.Args[1]

	// Generate grammatic
	semChecker, err := newSemChecker()
	if err != nil {
		log.Fatal(err)
	}

	// Use CYK to check if the word matches
	valid := semChecker.check(word)

	if valid {
		fmt.Printf("\x1b[1;32m%s is a valid grammatic\x1b[0m\n", word)
	} else {
		fmt.Printf("\x1b[1;31m%s is an invalid grammatic\x1b[0m\n", word)
	}

}
