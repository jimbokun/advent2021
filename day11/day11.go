package day11

import (
	"fmt"
	"jimbokun/advent/matrix"
)

func Adjacent(m matrix.ArrayMatrix, p matrix.Point) []matrix.Point {
	i := p.X
	j := p.Y
	
	ij1 := matrix.MakePoint(i, j+1)
	i1j := matrix.MakePoint(i+1, j)
	i1j1 := matrix.MakePoint(i+1, j+1)
	ij_1 := matrix.MakePoint(i, j-1)
	i1j_1 := matrix.MakePoint(i+1, j-1)
	i_1j := matrix.MakePoint(i-1, j)
	i_1j1 := matrix.MakePoint(i-1, j+1)
	i_1j_1 := matrix.MakePoint(i-1, j-1)
	
	if i == 0 && j == 0 {
		return []matrix.Point { ij1, i1j, i1j1 }
	} else if i == 0 && j == m.Cols-1 {
		return []matrix.Point { ij_1, i1j, i1j_1 }
	} else if i == 0 {
		return []matrix.Point { ij_1, ij1, i1j_1, i1j1, i1j }
	} else if i == m.Rows-1 && j == 0 {
		return []matrix.Point { i_1j, i_1j1, ij1 }
	} else if j == 0 {
		return []matrix.Point { i_1j, i_1j1, ij1, i1j, i1j1 }
	} else if i == m.Rows-1 && j == m.Cols-1 {
		return []matrix.Point { i_1j, ij_1, i_1j_1 }
	} else if i == m.Rows-1 {
		return []matrix.Point { i_1j_1, i_1j, i_1j1, ij_1, ij1 }
	} else if j == m.Cols-1 {
		return []matrix.Point { i_1j_1, i_1j, ij_1, i1j, i1j_1 }
	} else {
		return []matrix.Point { i_1j_1, i_1j, i_1j1, ij_1, ij1, i1j, i1j1, i1j_1 }
	}
}

func step(m matrix.ArrayMatrix) int {
	for row := 0; row < m.Rows; row++ {
		for col := 0; col < m.Cols; col++ {
			m.Increment(row, col)
		}
	}
	didFlash := matrix.MakeArrayMatrix(m.Rows, m.Cols)
	atLeastOneFlash := true
	for atLeastOneFlash {
		atLeastOneFlash = false
		for row := 0; row < m.Rows; row++ {
			for col := 0; col < m.Cols; col++ {
				if m.Get(row, col) > 9 {
					adj := Adjacent(m, matrix.MakePoint(row, col))
					for _, p := range adj {
						if didFlash.Get(p.X, p.Y) == 0 {
							m.Increment(p.X, p.Y)
						}
					}
					m.Set(row, col, 0)
					didFlash.Set(row, col, 1)
					atLeastOneFlash = true
				}
			}
		}
	}
	flashCount := 0
	for row := 0; row < m.Rows; row++ {
		for col := 0; col < m.Cols; col++ {
			if didFlash.Get(row, col) == 1 {
				flashCount++
			}
		}
	}
	return flashCount
}

func allFlashed(m matrix.ArrayMatrix) bool {
	for row := 0; row < m.Rows; row++ {
		for col := 0; col < m.Cols; col++ {
			if m.Get(row, col) > 0 {
				return false
			}
		}
	}
	return true
}

func Day11() {
	topo := matrix.ReadTopographicMap("day11/input.txt", 10, 10)
	fmt.Println(topo)
	// for row := 0; row < topo.Rows; row++ {
	// 	for col := 0; col < topo.Cols; col++ {
	// 		adj := Adjacent(topo, matrix.MakePoint(row, col))
	// 		fmt.Printf("adj(%d, %d) = %v\n", row, col, adj)
	// 		vals := make([]int, len(adj))
	// 		for i, p := range adj {
	// 			vals[i] = topo.Get(p.X, p.Y)
	// 		}
	// 		fmt.Printf("vals(%d, %d) = %v\n", row, col, vals)
	// 	}
	// }
	flashTotal := 0
	var i int
	for i = 0; !allFlashed(topo); i++ {
		flashTotal += step(topo)
		fmt.Println(topo)
		fmt.Printf("%d flashes after %d steps\n", flashTotal, i)
	}
	fmt.Printf("all octopi flashed on step %d\n", i)
}
