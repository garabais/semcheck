package main

import (
	"fmt"
)

const EPSILON = '~'

type simbol struct {
	displayName string
	productions []*production
}

type production struct {
	elements []interface{}
	remove   bool
}

func (s *simbol) add(prod *production) bool {
	// Check if the production to add already exists
	for _, p := range s.productions {
		if len(p.elements) == len(prod.elements) {
			var equal bool = true
			for i := 0; i < len(p.elements); i++ {
				switch v1 := prod.elements[i].(type) {
				case rune:
					v2, ok := p.elements[i].(rune)
					if ok {

						if v1 != v2 {
							equal = false
							break
						}
					} else {
						equal = false
						break
					}
				case *simbol:
					if v2, ok := p.elements[i].(*simbol); ok {
						if v1 != v2 {
							equal = false
							break
						}
					} else {
						equal = false
						break
					}
				}
			}
			if equal {
				return false
			}
		}
	}
	// If the production is new add it
	s.productions = append(s.productions, prod)
	return true
}

func (sim *simbol) String() string {
	s := fmt.Sprintf("%s->", sim.displayName)

	for i, prod := range sim.productions {
		if i != 0 {
			s += "|"
		}
		s += fmt.Sprint(prod)
	}

	return s
}

func (p *production) String() string {
	var s string

	if len(p.elements) == 0 {
		s += string(EPSILON)
	}

	for _, e := range p.elements {
		// Depending on the underling type how to print it
		switch v := e.(type) {
		case rune:
			s += string(v)
		case *simbol:
			s += fmt.Sprintf("{%s}", v.displayName)
		}
	}

	return s
}
