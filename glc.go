package main

import (
	"fmt"
)

type glc struct {
	start   *simbol
	simbols []*simbol
}

func (g *glc) String() string {
	var s string
	// Special case after cleaning when the starting simbol is useless
	if g.start == nil {
		return "\n"
	}

	s += fmt.Sprintln(g.start)

	for _, sim := range g.simbols {
		if sim != g.start {
			s += fmt.Sprintln(sim)
		}
	}

	return s
}

func (g *glc) check(word string) bool {

	// Matrix of sets of symbols
	mat := make([][]map[*simbol]struct{}, len(word))
	for i := 0; i < len(mat); i++ {
		mat[i] = make([]map[*simbol]struct{}, len(word))
		for j := 0; j < len(mat[i]); j++ {
			mat[i][j] = make(map[*simbol]struct{})
		}
	}

	// Make start diagonal of the matrix
	for i, letter := range word {
		for _, sim := range g.simbols {
			for _, prod := range sim.productions {
				if len(prod.elements) == 1 {
					c := prod.elements[0].(rune)

					if c == letter {
						mat[i][i][sim] = struct{}{}
					}
				}
			}
		}
	}

	// Move through the matrix
	for l := 1; l <= len(mat); l++ {
		for i := 0; i <= len(mat)-l; i++ {
			j := i + l - 1
			for k := i; k <= j-1; k++ {

				// The two sets in (i,k) and (k+1, j)
				for s1 := range mat[i][k] {
					for s2 := range mat[k+1][j] {
						// Check if there's a production that produces the elements of the sets
						for _, sim := range g.simbols {
							for _, prod := range sim.productions {
								if len(prod.elements) == 2 {
									c1 := prod.elements[0].(*simbol)
									c2 := prod.elements[1].(*simbol)

									if c1 == s1 && c2 == s2 {
										mat[i][j][sim] = struct{}{}
									}
								}
							}
						}
					}
				}
			}
		}
	}

	// Check if s is the final part of the matrix
	if _, ok := mat[0][len(mat)-1][g.start]; !ok {
		return false
	}

	return true
}
