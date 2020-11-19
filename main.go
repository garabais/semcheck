package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"runtime"
	"strings"
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
	matches, graph := gramatic.match(word)

	if matches {

		fmt.Printf("\x1b[1;32m%s is a valid semantic version\x1b[0m\n", word)
		// Save the graph
		dotPath, err := exec.LookPath("dot")
		if err != nil {
			log.Fatal("graphviz not found")
		}
		cmd := exec.Command(dotPath, "-Tpng", "-o", "tree.png")
		cmd.Stdin = strings.NewReader(graph.String())
		err = cmd.Run()
		if err != nil {
			log.Fatal("Error runing graphviz", err)
		}

		var openProg []string

		switch runtime.GOOS {
		case "darwin":
			openProg = []string{"open"}
		case "windows":
			openProg = []string{"cmd", "/c", "start"}
		default:
			openProg = []string{"xdg-open"}
		}

		cmd = exec.Command(openProg[0], append(openProg[1:], "tree.png")...)
		cmd.Start()
	} else {
		fmt.Printf("\x1b[1;31m%s is an invalid semantic version\x1b[0m\n", word)
	}

}
