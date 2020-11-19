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
	gramatic, err := newSemChecker()
	if err != nil {
		log.Fatal(err)
	}

	// Use CYK to check if the word matches
	matches := gramatic.match(word)

	if matches {
		fmt.Printf("The word \x1b[1;31m%s\x1b[0m is accepted by the grammatic\n", word)
	} else {
		fmt.Printf("The word \x1b[1;31m%s\x1b[0m isn't accepted by the grammatic\n", word)
		os.Exit(1)
	}

}
