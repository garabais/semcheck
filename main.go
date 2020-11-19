package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"strings"
)

func main() {
	if len(os.Args) != 3 {
		panic("Invalid number of arguments")
	}

	filename := os.Args[1]
	word := os.Args[2]

	// Open file
	f, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}

	// Generate grammatic
	gramatic, err := generateGLC(f)
	f.Close()
	if err != nil {
		log.Fatal(err)
	}

	// Put the grammatic in FNC
	gramatic.normalize()

	// Save the grammatic
	fmt.Println("Saving Chomsky grammatic in fnc-" + filename)
	f, err = os.Create("fnc-" + filename)
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

func generateGLC(f io.Reader) (*glc, error) {

	scanner := bufio.NewScanner(f)

	reference := make(map[string]*simbol)

	simbols := make([]*simbol, 0)

	var lineNumber int = 0

	for scanner.Scan() {
		lineNumber++

		// Separate symbol and productions
		line := scanner.Text()
		lParts := strings.SplitN(line, "->", 2)

		if len(lParts) != 2 {
			return nil, fmt.Errorf("Invalid syntax in line %d", lineNumber)
		}

		// Separate productions
		productions := strings.Split(lParts[1], "|")

		// If the symbol doesn't exist create it
		s, ok := reference[lParts[0]]
		if !ok {
			s = &simbol{displayName: lParts[0]}
			reference[s.displayName] = s
			simbols = append(simbols, s)
		}

		for _, prod := range productions {
			p := &production{}

			var pStart int
			var pGen bool
			for i, c := range prod {
				if pGen {
					if c == '}' {
						genSimStr := prod[pStart:i]
						pGen = false

						// If the symbol doesn't exist create it
						genSim, ok := reference[genSimStr]
						if !ok {
							genSim = &simbol{displayName: genSimStr}
							reference[genSim.displayName] = genSim
							simbols = append(simbols, genSim)
						}

						p.elements = append(p.elements, genSim)
					}
				} else {
					// Al the text until '}' is the name of the symbol
					if c == '{' {
						pGen = true
						pStart = i + 1

					} else if c != EPSILON {
						// Don't add epsilon, in case the whole transition is epsilon is an empty array
						p.elements = append(p.elements, c)
					}
				}
			}
			if pGen {
				return nil, fmt.Errorf("Invalid syntax in line %d", lineNumber)
			}

			// Add the production to the symbol
			s.add(p)
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return &glc{simbols: simbols, start: simbols[0]}, nil
}
