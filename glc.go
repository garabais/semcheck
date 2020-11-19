package main

import (
	"fmt"

	"github.com/awalterschulze/gographviz"
)

type glc struct {
	start   *simbol
	simbols []*simbol
	chomsky bool
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

func (g *glc) cleanTransitions() {
	var added bool = true

	// Check for epsilon/unitary transition until no new transitions are added
	for added {
		added = false
		// For every simbol
		for _, sim := range g.simbols {
			// Check every production
			for _, prod := range sim.productions {

				if len(prod.elements) == 0 {
					// This production is epsilon
					prod.remove = true

				} else if len(prod.elements) == 1 {
					// Unitary transition
					if v, ok := prod.elements[0].(*simbol); ok {
						prod.remove = true
						if v != sim {
							// Add every transition of the unitary transition
							for _, targetP := range v.productions {
								a := sim.add(targetP)
								// If the production is added succesfully change the state of added
								if a {
									added = true
								}
							}
						}
					}
				}

				// For every element of the production
				for i, element := range prod.elements {

					// Check if it's a simbol
					if inSim, ok := element.(*simbol); ok {

						// Check if one of the transitions of that simbol is epsilon
						for _, inProd := range inSim.productions {
							if len(inProd.elements) == 0 {
								// Create a copy of the transition skiping the simbol that has epsilon
								temp := &production{}
								temp.elements = append(temp.elements, prod.elements[:i]...)
								temp.elements = append(temp.elements, prod.elements[i+1:]...)
								a := sim.add(temp)
								// If the production is added succesfully change the state of added
								if a {
									added = true
								}
							}
						}

					}
				}
			}
		}

	}

	// Check every symbol
	for _, s := range g.simbols {
		// And every production of the symbol
		for i := 0; i < len(s.productions); i++ {
			// Erase if the production is epsilon, or is marked to remove
			var erase bool = false
			if len(s.productions[i].elements) == 0 {
				erase = true
			} else if len(s.productions[i].elements) == 1 {
				if _, ok := s.productions[i].elements[0].(*simbol); ok {
					erase = true
				}
			}

			if erase {
				s.productions[i] = s.productions[len(s.productions)-1]
				s.productions = s.productions[:len(s.productions)-1]
				i--
			}
		}
	}
}

func (g *glc) cleanUseless() {

	// Set of useful symbol
	usefull := make(map[*simbol]struct{})

	var added bool = true

	// Check for useful transition until no new transitions are added
	for added {
		added = false

		for _, s := range g.simbols {
			// Skip symbol if is already marked as useful
			if _, ok := usefull[s]; ok {
				continue
			}

			// Check for every production
			for _, prod := range s.productions {
				var add bool = true
				// Check every element if is has a useful symbol or the productions only haver terminal symbols
				for _, e := range prod.elements {

					if v, ok := e.(*simbol); ok {
						add = false

						if _, ok = usefull[v]; ok {
							add = true
							break
						}

					}
				}
				if add {
					usefull[s] = struct{}{}
					added = true
					break
				}
			}
		}
	}

	// Erase useless symbols
	for i := 0; i < len(g.simbols); i++ {
		if _, ok := usefull[g.simbols[i]]; !ok {
			g.simbols[i] = g.simbols[len(g.simbols)-1]
			g.simbols = g.simbols[:len(g.simbols)-1]
			i--
		}
	}

	// Create a new array of symbols that will be populated with reachable symbols
	newSimbols := make([]*simbol, 0)

	// Check if the start symbol is useful
	if _, ok := usefull[g.start]; ok {
		// Set of useful symbol
		reached := make(map[*simbol]struct{})
		// Add the starting symbol
		newSimbols = append(newSimbols, g.start)
		reached[g.start] = struct{}{}

		// Add every symbol you can reach
		for i := 0; i < len(newSimbols); i++ {
			for _, prod := range newSimbols[i].productions {
				for _, e := range prod.elements {
					if v, ok := e.(*simbol); ok {
						if _, added := reached[v]; !added {
							newSimbols = append(newSimbols, v)
							reached[v] = struct{}{}
						}
					}
				}
			}
		}
	} else {
		// If the start symbol es not useful sit it to nil
		g.start = nil
	}
	// Replace the array of symbols
	g.simbols = newSimbols
}

func (g *glc) normalize() {

	g.cleanTransitions()
	g.cleanUseless()

	// Set of reference to the new symbol
	letters := make(map[rune]*simbol)

	var n int

	// Replace every terminal symbol so the production only has generators
	for _, sim := range g.simbols {
		for _, prod := range sim.productions {
			if len(prod.elements) > 1 {
				for i := 0; i < len(prod.elements); i++ {
					if letter, ok := prod.elements[i].(rune); ok {
						lSim, ok := letters[letter]
						if !ok {
							lSim = &simbol{displayName: fmt.Sprintf("L%d", n)}
							n++
							p := &production{}
							p.elements = append(p.elements, letter)
							lSim.add(p)
							letters[letter] = lSim
							g.simbols = append(g.simbols, lSim)
						}
						prod.elements[i] = lSim
					}
				}
			}
		}
	}

	var newPN int

	// Check for every production of every symbol if then length is bigger than 2
	for i := 0; i < len(g.simbols); i++ {
		for j := 0; j < len(g.simbols[i].productions); j++ {
			if len(g.simbols[i].productions[j].elements) > 2 {
				// If the length is bigger than 2 create a new symbol and change the production
				nSim := &simbol{displayName: fmt.Sprintf("C%d", newPN)}
				newPN++
				p := &production{}
				p.elements = append(p.elements, g.simbols[i].productions[j].elements[1:]...)
				nSim.add(p)
				g.simbols[i].productions[j].elements[1] = nSim
				g.simbols[i].productions[j].elements = g.simbols[i].productions[j].elements[:2]
				g.simbols = append(g.simbols, nSim)
			}
		}
	}

	g.chomsky = true
}

func (g *glc) match(word string) (bool, *gographviz.Graph) {
	if !g.chomsky {
		g.normalize()
	}
	// Special case, useless start symbol
	if g.start == nil {
		return false, nil
	}

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
		return false, nil
	}
	// Erase the other simbols and just keep the start simbol
	mat[0][len(mat)-1] = map[*simbol]struct{}{
		g.start: {},
	}

	// Create a graph
	graph := gographviz.NewGraph()
	if err := graph.SetName("G"); err != nil {
		panic(err)
	}
	if err := graph.SetDir(true); err != nil {
		panic(err)
	}
	graph.AddAttr("G", "label", fmt.Sprintf("\"%s\"", word))
	graph.AddAttr("G", "labelloc", "t")
	graph.AddAttr("G", "fontsize", "20")

	// Add starting node
	if err := graph.AddNode("G", fmt.Sprintf("q0%d", len(mat)-1), map[string]string{"label": fmt.Sprintf("\"%s\"", g.start.displayName)}); err != nil {
		panic(err)
	}

	backTrack(word, graph, mat, 0, len(mat)-1)

	return true, graph
}

func backTrack(word string, graph *gographviz.Graph, mat [][]map[*simbol]struct{}, targetI, targetJ int) {
	// If the position is in the diagonal add the leave of the graph
	if targetI == targetJ {
		// Create leave node
		if err := graph.AddNode("G", fmt.Sprintf("f%d", targetI), map[string]string{"label": fmt.Sprintf("\"%s\"", string(word[targetI])), "shape": "doublecircle"}); err != nil {
			panic(err)
		}
		// Connect with parent
		if err := graph.AddEdge(fmt.Sprintf("q%d%d", targetI, targetJ), fmt.Sprintf("f%d", targetI), true, nil); err != nil {
			panic(err)
		}
		return
	}

	// Search for the origin of the simbols in mat[targetI][targetJ]
	for i, j := targetI+1, targetI; j < targetJ && i <= targetJ; i, j = i+1, j+1 {
		for sim := range mat[targetI][targetJ] {
			for _, prod := range sim.productions {
				if len(prod.elements) == 2 {
					for s1 := range mat[targetI][j] {
						for s2 := range mat[i][targetJ] {
							c1 := prod.elements[0].(*simbol)
							c2 := prod.elements[1].(*simbol)

							if c1 == s1 && c2 == s2 {
								// When found, add the 2 children
								if err := graph.AddNode("G", fmt.Sprintf("q%d%d", targetI, j), map[string]string{"label": fmt.Sprintf("\"%s\"", s1.displayName)}); err != nil {
									panic(err)
								}
								if err := graph.AddNode("G", fmt.Sprintf("q%d%d", i, targetJ), map[string]string{"label": fmt.Sprintf("\"%s\"", s2.displayName)}); err != nil {
									panic(err)
								}
								// Connect the children with the parent(current)
								if err := graph.AddEdge(fmt.Sprintf("q%d%d", targetI, targetJ), fmt.Sprintf("q%d%d", targetI, j), true, nil); err != nil {
									panic(err)
								}
								if err := graph.AddEdge(fmt.Sprintf("q%d%d", targetI, targetJ), fmt.Sprintf("q%d%d", i, targetJ), true, nil); err != nil {
									panic(err)
								}
								// Call backTrack in each of the clildren
								backTrack(word, graph, mat, targetI, j)
								backTrack(word, graph, mat, i, targetJ)

								// And stop search
								return

							}
						}
					}
				}
			}
		}
	}
}
