package main

import (
	"bufio"
	"fmt"
	"strings"
)

const semver = `valid semver->{version core}|{version core}-{pre-release}|{version core}+{build}|{version core}-{pre-release}+{build}
version core->{major}.{minor}.{patch}
major->{numeric identifier}
minor->{numeric identifier}
patch->{numeric identifier}
pre-release->{dot-separated pre-release identifiers}
dot-separated pre-release identifiers->{pre-release identifier}|{pre-release identifier}.{dot-separated pre-release identifiers}
build->{dot-separated build identifiers}
dot-separated build identifiers->{build identifier}|{build identifier}.{dot-separated build identifiers}
pre-release identifier->{alphanumeric identifier}|{numeric identifier}
build identifier->{alphanumeric identifier}|{digits}
alphanumeric identifier->{non-digit}|{non-digit}{identifier characters}|{identifier characters}{non-digit}|{identifier characters}{non-digit}{identifier characters}
numeric identifier->0|{positive digit}|{positive digit}{digits}
identifier characters->{identifier character}|{identifier character}{identifier characters}
identifier character->{digit}|{non-digit}
non-digit->{letter}|-
digits->{digit}|{digit}{digits}
digit->0|{positive digit}
positive digit->1|2|3|4|5|6|7|8|9
letter->A|B|C|D|E|F|G|H|I|J|K|L|M|N|O|P|Q|R|S|T|U|V|W|X|Y|Z|a|b|c|d|e|f|g|h|i|j|k|l|m|n|o|p|q|r|s|t|u|v|w|x|y|z`

func newSemChecker() (*glc, error) {
	gramathic := strings.NewReader(semver)
	scanner := bufio.NewScanner(gramathic)

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
