package matrix

import (
	"bufio"
	"log"
	"os"
	"strings"
	"strconv"
	"fmt"
)

type Matrix interface {
	RowCount() int
	ColCount() int
	Get(i, j int) int
	Set(i, j, v int)
	Increment(i, j int)
	Adjacent(p Point) []Point
	Print()
	PrintTopo()
	All(f func(i, j int))
	SubMatrix(origin Point, rows, cols int) Matrix
}

type ArrayMatrix struct {
	Rows, Cols int
	Values []int
}

func (m ArrayMatrix) RowCount() int {
	return m.Rows
}

func (m ArrayMatrix) ColCount() int {
	return m.Cols
}

func (m ArrayMatrix) Get(i, j int) int {
	// fmt.Printf("get(%d, %d)\n", i, j)
	return m.Values[m.Cols * i + j]
}

func (m ArrayMatrix) Set(i, j, v int) {
	// fmt.Printf("Set(%d, %d, %d)\n", i, j, v)
	m.Values[m.Cols * i + j] = v
}

func (m ArrayMatrix) Increment(i, j int) {
	// fmt.Printf("Increment(%d, %d)\n", i, j)
	m.Values[m.Cols * i + j]++
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

func MakeArrayMatrixWithInitialValue(rows, cols, value int) ArrayMatrix {
	m := ArrayMatrix{Rows: rows, Cols: cols, Values: make([]int, rows * cols)}
	All(m, func(i, j int) { m.Set(i, j, value) })
	return m
}

func MakePoint(x, y int) Point {
	return Point{ X: x, Y: y }
}

func ParsePoint(pointVal string) Point {
	xyVals := strings.Split(pointVal, ",")
	x, _ := strconv.Atoi(xyVals[0])
	y, _ := strconv.Atoi(xyVals[1])
	return MakePoint(x, y)
}

func Adjacent(m Matrix, p Point) []Point {
	i := p.X
	j := p.Y
	if i == 0 && j == 0 {
		return []Point { MakePoint(i, j+1), MakePoint(i+1, j) }
	} else if i == 0 && j == m.ColCount()-1 {
		return []Point { MakePoint(i, j-1), MakePoint(i+1, j) }
	} else if i == 0 {
		return []Point { MakePoint(i, j-1), MakePoint(i, j+1), MakePoint(i+1, j) }
	} else if i == m.RowCount()-1 && j == 0 {
		return []Point { MakePoint(i-1, j), MakePoint(i, j+1) }
	} else if j == 0 {
		return []Point { MakePoint(i-1, j), MakePoint(i, j+1), MakePoint(i+1, j) }
	} else if i == m.RowCount()-1 && j == m.ColCount()-1 {
		return []Point { MakePoint(i-1, j), MakePoint(i, j-1) }
	} else if i == m.RowCount()-1 {
		return []Point { MakePoint(i-1, j), MakePoint(i, j-1), MakePoint(i, j+1) }
	} else if j == m.ColCount()-1 {
		return []Point { MakePoint(i-1, j), MakePoint(i, j-1), MakePoint(i+1, j) }
	} else {
		return []Point { MakePoint(i-1, j), MakePoint(i, j-1), MakePoint(i, j+1), MakePoint(i+1, j) }
	}
}

func (m ArrayMatrix) Adjacent(p Point) []Point {
	return Adjacent(m, p)
}

func (m ArrayMatrix) AdjacentValues(i, j int) []int {
	adj := Adjacent(m, Point{ X: i, Y: j})
	adjValues := make([]int, len(adj))
	for i, p := range adj {
		adjValues[i] = m.Get(p.X, p.Y)
	}
	return adjValues
}

func ReadTopographicMap(filename string, rows, cols int) ArrayMatrix {
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

func All(m Matrix, f func(i, j int)) {
	for row := 0; row < m.RowCount(); row++ {
		for col := 0; col < m.ColCount(); col++ {
			f(row, col)
		}
	}
}

func (m ArrayMatrix) All(f func(i, j int)) {
	All(m, f)
}


func printMatrix(m Matrix) {
	for row := 0; row < m.RowCount(); row++ {
		for col := 0; col < m.ColCount(); col++ {
			if col > 0 {
				fmt.Print(", ")
			}
			fmt.Printf("%d", m.Get(row, col))
		}
		fmt.Println()
	}
}

func printTopo(m Matrix) {
	for row := 0; row < m.RowCount(); row++ {
		for col := 0; col < m.ColCount(); col++ {
			fmt.Printf("%d", m.Get(row, col))
		}
		fmt.Println()
	}
}

func (m ArrayMatrix) Print() {
	printMatrix(m)
}

func (m ArrayMatrix) PrintTopo() {
	printTopo(m)
}

func SubMatrix(m Matrix, origin Point, rows, cols int) Matrix {
	sub := MakeArrayMatrix(rows, cols)
	sub.All(func(i, j int) {
		sub.Set(i, j, m.Get(origin.X + i, origin.Y + j))
	})
	return sub
}

func (m ArrayMatrix) SubMatrix(origin Point, rows, cols int) Matrix {
	return SubMatrix(m, origin, rows, cols)
}
