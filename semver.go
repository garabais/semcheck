package main

import (
	"bufio"
	"fmt"
	"strings"
)

// valid semver->{version core}|{version core}-{pre-release}|{version core}+{build}|{version core}-{pre-release}+{build}
// version core->{major}.{minor}.{patch}
// major->{numeric identifier}
// minor->{numeric identifier}
// patch->{numeric identifier}
// pre-release->{dot-separated pre-release identifiers}
// dot-separated pre-release identifiers->{pre-release identifier}|{pre-release identifier}.{dot-separated pre-release identifiers}
// build->{dot-separated build identifiers}
// dot-separated build identifiers->{build identifier}|{build identifier}.{dot-separated build identifiers}
// pre-release identifier->{alphanumeric identifier}|{numeric identifier}
// build identifier->{alphanumeric identifier}|{digits}
// alphanumeric identifier->{non-digit}|{non-digit}{identifier characters}|{identifier characters}{non-digit}|{identifier characters}{non-digit}{identifier characters}
// numeric identifier->0|{positive digit}|{positive digit}{digits}
// identifier characters->{identifier character}|{identifier character}{identifier characters}
// identifier character->{digit}|{non-digit}
// non-digit->{letter}|-
// digits->{digit}|{digit}{digits}
// digit->0|{positive digit}
// positive digit->1|2|3|4|5|6|7|8|9
// letter->A|B|C|D|E|F|G|H|I|J|K|L|M|N|O|P|Q|R|S|T|U|V|W|X|Y|Z|a|b|c|d|e|f|g|h|i|j|k|l|m|n|o|p|q|r|s|t|u|v|w|x|y|z
const semver = `valid semver->{major}{C0}|{version core}{C1}|{version core}{C2}|{version core}{C3}
major->9|0|8|{positive digit}{digits}|1|2|3|4|5|6|7
minor->9|0|8|{positive digit}{digits}|1|2|3|4|5|6|7
patch->9|0|8|{positive digit}{digits}|1|2|3|4|5|6|7
version core->{major}{C0}
pre-release->z|y|{pre-release identifier}{C4}|x|w|v|{non-digit}{identifier characters}|{identifier characters}{non-digit}|{identifier characters}{C5}|0|u|{positive digit}{digits}|1|2|3|4|5|6|7|8|9|t|-|A|B|C|D|E|F|G|H|I|J|K|L|M|N|O|P|Q|R|S|T|U|V|W|X|Y|Z|a|b|c|d|e|f|g|h|i|j|k|l|m|n|o|p|q|r|s
build->9|8|{build identifier}{C6}|7|6|5|{non-digit}{identifier characters}|{identifier characters}{non-digit}|{identifier characters}{C5}|4|{digit}{digits}|3|-|0|2|A|B|C|D|E|F|G|H|I|J|K|L|M|N|O|P|Q|R|S|T|U|V|W|X|Y|Z|a|b|c|d|e|f|g|h|i|j|k|l|m|n|o|p|q|r|s|t|u|v|w|x|y|z|1
positive digit->1|2|3|4|5|6|7|8|9
digits->9|{digit}{digits}|0|8|1|2|3|4|5|6|7
pre-release identifier->z|y|x|{non-digit}{identifier characters}|{identifier characters}{non-digit}|{identifier characters}{C5}|0|w|{positive digit}{digits}|1|2|3|4|5|6|7|8|9|v|-|A|B|C|D|E|F|G|H|I|J|K|L|M|N|O|P|Q|R|S|T|U|V|W|X|Y|Z|a|b|c|d|e|f|g|h|i|j|k|l|m|n|o|p|q|r|s|t|u
dot-separated pre-release identifiers->z|{pre-release identifier}{C4}|y|x|w|{non-digit}{identifier characters}|{identifier characters}{non-digit}|{identifier characters}{C5}|0|v|{positive digit}{digits}|1|2|3|4|5|6|7|8|9|u|-|A|B|C|D|E|F|G|H|I|J|K|L|M|N|O|P|Q|R|S|T|U|V|W|X|Y|Z|a|b|c|d|e|f|g|h|i|j|k|l|m|n|o|p|q|r|s|t
non-digit->z|-|A|B|C|D|E|F|G|H|I|J|K|L|M|N|O|P|Q|R|S|T|U|V|W|X|Y|Z|a|b|c|d|e|f|g|h|i|j|k|l|m|n|o|p|q|r|s|t|u|v|w|x|y
identifier characters->9|{identifier character}{identifier characters}|8|7|0|6|5|-|A|B|C|D|E|F|G|H|I|J|K|L|M|N|O|P|Q|R|S|T|U|V|W|X|Y|Z|a|b|c|d|e|f|g|h|i|j|k|l|m|n|o|p|q|r|s|t|u|v|w|x|y|z|1|2|3|4
build identifier->9|8|7|{non-digit}{identifier characters}|{identifier characters}{non-digit}|{identifier characters}{C5}|6|{digit}{digits}|5|-|0|4|A|B|C|D|E|F|G|H|I|J|K|L|M|N|O|P|Q|R|S|T|U|V|W|X|Y|Z|a|b|c|d|e|f|g|h|i|j|k|l|m|n|o|p|q|r|s|t|u|v|w|x|y|z|1|2|3
dot-separated build identifiers->9|{build identifier}{C6}|8|7|6|{non-digit}{identifier characters}|{identifier characters}{non-digit}|{identifier characters}{C5}|5|{digit}{digits}|4|-|0|3|A|B|C|D|E|F|G|H|I|J|K|L|M|N|O|P|Q|R|S|T|U|V|W|X|Y|Z|a|b|c|d|e|f|g|h|i|j|k|l|m|n|o|p|q|r|s|t|u|v|w|x|y|z|1|2
digit->0|9|1|2|3|4|5|6|7|8
identifier character->9|8|0|7|6|-|A|B|C|D|E|F|G|H|I|J|K|L|M|N|O|P|Q|R|S|T|U|V|W|X|Y|Z|a|b|c|d|e|f|g|h|i|j|k|l|m|n|o|p|q|r|s|t|u|v|w|x|y|z|1|2|3|4|5
L0->.
L1->-
L2->+
C0->{L0}{C7}
C1->{L1}{pre-release}
C2->{L2}{build}
C3->{L1}{C8}
C4->{L0}{dot-separated pre-release identifiers}
C5->{non-digit}{identifier characters}
C6->{L0}{dot-separated build identifiers}
C7->{minor}{C9}
C8->{pre-release}{C10}
C9->{L0}{patch}
C10->{L2}{build}`

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
