package matrix

import (
	"bufio"
	"log"
	"os"
)

type Matrix interface {
	Get(i, j int) int
	Set(i, j, v int)
	Increment(i, j int)
}

type ArrayMatrix struct {
	Rows, Cols int
	Values []int
}

func (m ArrayMatrix) Get(i, j int) int {
	// fmt.Printf("get(%d, %d)\n", i, j)
	return m.Values[m.Cols * i + j]
}

func (m ArrayMatrix) Set(i, j, v int) {
	// fmt.Printf("set(%d, %d, %d)\n", i, j, v)
	m.Values[m.Cols * i + j] = v
}

func (m ArrayMatrix) Increment(i, j int) {
	m.Values[m.Cols * j + i]++
}

type Point struct {
	X, Y int
}

func (p1 Point) Add(p2 Point) Point {
	return Point{X: p1.X + p2.X, Y: p1.Y + p2.Y}
}

func MakeArrayMatrix(rows, cols int) ArrayMatrix {
	return ArrayMatrix{Rows: rows, Cols: cols, Values: make([]int, rows * cols)}
}

func MakePoint(x, y int) Point {
	return Point{ X: x, Y: y }
}

func (m ArrayMatrix) adjacent(p Point) []Point {
	i := p.X
	j := p.Y
	if i == 0 && j == 0 {
		return []Point { MakePoint(i, j+1), MakePoint(i+1, j) }
	} else if i == 0 && j == m.Cols-1 {
		return []Point { MakePoint(i, j-1), MakePoint(i+1, j) }
	} else if i == 0 {
		return []Point { MakePoint(i, j-1), MakePoint(i, j+1), MakePoint(i+1, j) }
	} else if i == m.Rows-1 && j == 0 {
		return []Point { MakePoint(i-1, j), MakePoint(i, j+1) }
	} else if j == 0 {
		return []Point { MakePoint(i-1, j), MakePoint(i, j+1), MakePoint(i+1, j) }
	} else if i == m.Rows-1 && j == m.Cols-1 {
		return []Point { MakePoint(i-1, j), MakePoint(i, j-1) }
	} else if i == m.Rows-1 {
		return []Point { MakePoint(i-1, j), MakePoint(i, j-1), MakePoint(i, j+1) }
	} else if j == m.Cols-1 {
		return []Point { MakePoint(i-1, j), MakePoint(i, j-1), MakePoint(i+1, j) }
	} else {
		return []Point { MakePoint(i-1, j), MakePoint(i, j-1), MakePoint(i, j+1), MakePoint(i+1, j) }
	}
}

func (m ArrayMatrix) adjacentValues(i, j int) []int {
	adj := m.adjacent(Point{ X: i, Y: j})
	adjValues := make([]int, len(adj))
	for i, p := range adj {
		adjValues[i] = m.Get(p.X, p.Y)
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
			matrix.Set(row, col, int(r - '0'))
		}
		row++
	}

	return matrix
}
