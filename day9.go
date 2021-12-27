package main


import (
	"bufio"
	"fmt"
	"log"
	"os"
	"sort"
	// "matrix"
)

// type matrix interface {
// 	get(i, j int) int
// 	set(i, j, v int)
// }

// type arrayMatrix struct {
// 	rows, cols int
// 	values []int
// }

// func (m arrayMatrix) get(i, j int) int {
// 	// fmt.Printf("get(%d, %d)\n", i, j)
// 	return m.values[m.cols * i + j]
// }

// func (m arrayMatrix) set(i, j, v int) {
// 	// fmt.Printf("set(%d, %d, %d)\n", i, j, v)
// 	m.values[m.cols * i + j] = v
// }

// type point struct {
// 	x, y int
// }

func All(vs []point, f func(point) bool) bool {
	for _, v := range vs {
		if !f(v) {
			return false
		}
	}
	return true
}

func lowPoints(m arrayMatrix) []point {
	low := make([]point, 0)
	for row := 0; row < m.rows; row++ {
		for col := 0; col < m.cols; col++ {
			depth := m.get(row, col)
			p := point{ x: row, y: col }
			if All(m.adjacent(p), func(p point) bool { return m.get(p.x, p.y) > depth }) {
				low = append(low, p)
			}
		}
	}
	return low
}

// find basins:
// for each basin:
//  construct matrix to mark points in basin
//  for each point in adj:
//    if already marked, ignore
//    if value of point < 9, mark as in the basin and add to queue
//    continue until no more points in queue
//  count up points marked in basin

func (m arrayMatrix) basin(low point) int {
	q := make([]point, 0)
	for _, p := range m.adjacent(low) {
		q = append(q, p)
	}
	basin := makeArrayMatrix(m.rows, m.cols)
	basinSize := 0
	for len(q) > 0 {
		p := q[0]
		if basin.get(p.x, p.y) == 0 && m.get(p.x, p.y) < 9 {
			basinSize++
			for _, p := range m.adjacent(p) {
				q = append(q, p)
			}
			basin.set(p.x, p.y, 1)
		}
		q = q[1:]
	}
	fmt.Println("basin (%d, %d): %v\n", low.x, low.y, basin)
	return basinSize
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

func main() {
	topo := readTopographicMap("day9_input.txt", 100, 100)
	// topo := readTopographicMap("day9_sample_input.txt", 5, 10)
	for row := 0; row < topo.rows; row++ {
		for col := 0; col < topo.cols; col++ {
			fmt.Printf("adj %d, %d: %v\n", row, col, topo.adjacentValues(row, col))
		}
	}
	top3 := make([]int, 3)
	for _, low := range lowPoints(topo) {
		top3 = append(top3, topo.basin(low))
		sort.Slice(top3, func(i, j int) bool { return top3[i] > top3[j] })
		top3 = top3[0:3]
	}
	total := 1
	fmt.Println(top3)
	for _, v := range top3 {
		total *= v
	}
	fmt.Println(total)
}
