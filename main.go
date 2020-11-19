package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
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
	matches, graph := gramatic.match(word)

	if matches {

		// Save the graph
		fmt.Printf("The word \x1b[1;31m%s\x1b[0m is accepted by the grammatic\n", word)
		fmt.Println("Saving derivation tree dot file in graph.dot")
		f, err = os.Create("graph.dot")
		if err != nil {
			log.Fatal("Failed creating file", err)
		}
		fmt.Fprint(f, graph)

		// Make an image of the graph
		fmt.Println("Creating image of the derivation tree graph.png")
		dotPath, err := exec.LookPath("dot")
		if err != nil {
			log.Fatal("graphviz not found")
		}
		cmd := exec.Command(dotPath, "-Tpng", "graph.dot", "-o", "graph.png")
		err = cmd.Run()
		if err != nil {
			log.Fatal("Error runing graphviz", err)
		}
	} else {
		fmt.Printf("The word \x1b[1;31m%s\x1b[0m isn't accepted by the grammatic\n", word)
		os.Exit(1)
	}

}
