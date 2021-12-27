package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"sort"
)

type matrix interface {
	get(i, j int) int
	set(i, j, v int)
}

type arrayMatrix struct {
	rows, cols int
	values []int
}

func (m arrayMatrix) get(i, j int) int {
	// fmt.Printf("get(%d, %d)\n", i, j)
	return m.values[m.cols * i + j]
}

func (m arrayMatrix) set(i, j, v int) {
	// fmt.Printf("set(%d, %d, %d)\n", i, j, v)
	m.values[m.cols * i + j] = v
}

type point struct {
	x, y int
}

func makeArrayMatrix(rows, cols int) arrayMatrix {
	return arrayMatrix{rows: rows, cols: cols, values: make([]int, rows * cols)}
}

func makePoint(x, y int) point {
	return point{ x: x, y: y }
}

func (m arrayMatrix) adjacent(p point) []point {
	i := p.x
	j := p.y
	if i == 0 && j == 0 {
		return []point { makePoint(i, j+1), makePoint(i+1, j) }
	} else if i == 0 && j == m.cols-1 {
		return []point { makePoint(i, j-1), makePoint(i+1, j) }
	} else if i == 0 {
		return []point { makePoint(i, j-1), makePoint(i, j+1), makePoint(i+1, j) }
	} else if i == m.rows-1 && j == 0 {
		return []point { makePoint(i-1, j), makePoint(i, j+1) }
	} else if j == 0 {
		return []point { makePoint(i-1, j), makePoint(i, j+1), makePoint(i+1, j) }
	} else if i == m.rows-1 && j == m.cols-1 {
		return []point { makePoint(i-1, j), makePoint(i, j-1) }
	} else if i == m.rows-1 {
		return []point { makePoint(i-1, j), makePoint(i, j-1), makePoint(i, j+1) }
	} else if j == m.cols-1 {
		return []point { makePoint(i-1, j), makePoint(i, j-1), makePoint(i+1, j) }
	} else {
		return []point { makePoint(i-1, j), makePoint(i, j-1), makePoint(i, j+1), makePoint(i+1, j) }
	}
}

func (m arrayMatrix) adjacentValues(i, j int) []int {
	adj := m.adjacent(point{ x: i, y: j})
	adjValues := make([]int, len(adj))
	for i, p := range adj {
		adjValues[i] = m.get(p.x, p.y)
	}
	return adjValues
}

func readTopographicMap(filename string, rows, cols int) arrayMatrix {
	file, err := os.Open(filename)

	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	matrix := makeArrayMatrix(rows, cols)
	row := 0
	for scanner.Scan() {
		for col, r := range scanner.Text() {
			matrix.set(row, col, int(r - '0'))
		}
		row++
	}

	return matrix
}
