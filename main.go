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

	// Put the grammatic in FNC
	gramatic.normalize()

	// Save the grammatic
	fmt.Println("Saving Chomsky grammatic in fnc-")
	f, err := os.Create("fnc-")
	if err != nil {
		log.Println("Failed creating file", err)
	}
	fmt.Fprint(f, gramatic)

	// Use CYK to check if the word matches
	matches := gramatic.match(word)

	if matches {
		fmt.Printf("The word \x1b[1;31m%s\x1b[0m is accepted by the grammatic\n", word)
	} else {
		fmt.Printf("The word \x1b[1;31m%s\x1b[0m isn't accepted by the grammatic\n", word)
		os.Exit(1)
	}

}
