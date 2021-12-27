package matrix

import (
	"bufio"
	"log"
	"os"
)

type Matrix interface {
	get(i, j int) int
	set(i, j, v int)
	Increment(i, j int)
}

type ArrayMatrix struct {
	Rows, Cols int
	Values []int
}

func (m ArrayMatrix) get(i, j int) int {
	// fmt.Printf("get(%d, %d)\n", i, j)
	return m.Values[m.Cols * i + j]
}

func (m ArrayMatrix) set(i, j, v int) {
	// fmt.Printf("set(%d, %d, %d)\n", i, j, v)
	m.Values[m.Cols * i + j] = v
}

func (m ArrayMatrix) Increment(i, j int) {
	m.Values[m.Cols * j + i]++
}

type point struct {
	x, y int
}

func MakeArrayMatrix(rows, cols int) ArrayMatrix {
	return ArrayMatrix{Rows: rows, Cols: cols, Values: make([]int, rows * cols)}
}

func makePoint(x, y int) point {
	return point{ x: x, y: y }
}

func (m ArrayMatrix) adjacent(p point) []point {
	i := p.x
	j := p.y
	if i == 0 && j == 0 {
		return []point { makePoint(i, j+1), makePoint(i+1, j) }
	} else if i == 0 && j == m.Cols-1 {
		return []point { makePoint(i, j-1), makePoint(i+1, j) }
	} else if i == 0 {
		return []point { makePoint(i, j-1), makePoint(i, j+1), makePoint(i+1, j) }
	} else if i == m.Rows-1 && j == 0 {
		return []point { makePoint(i-1, j), makePoint(i, j+1) }
	} else if j == 0 {
		return []point { makePoint(i-1, j), makePoint(i, j+1), makePoint(i+1, j) }
	} else if i == m.Rows-1 && j == m.Cols-1 {
		return []point { makePoint(i-1, j), makePoint(i, j-1) }
	} else if i == m.Rows-1 {
		return []point { makePoint(i-1, j), makePoint(i, j-1), makePoint(i, j+1) }
	} else if j == m.Cols-1 {
		return []point { makePoint(i-1, j), makePoint(i, j-1), makePoint(i+1, j) }
	} else {
		return []point { makePoint(i-1, j), makePoint(i, j-1), makePoint(i, j+1), makePoint(i+1, j) }
	}
}

func (m ArrayMatrix) adjacentValues(i, j int) []int {
	adj := m.adjacent(point{ x: i, y: j})
	adjValues := make([]int, len(adj))
	for i, p := range adj {
		adjValues[i] = m.get(p.x, p.y)
	}
	return adjValues
}

func readTopographicMap(filename string, rows, cols int) ArrayMatrix {
	file, err := os.Open(filename)

	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	matrix := MakeArrayMatrix(rows, cols)
	row := 0
	for scanner.Scan() {
		for col, r := range scanner.Text() {
			matrix.set(row, col, int(r - '0'))
		}
		row++
	}

	return matrix
}
