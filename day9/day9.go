package day9

import (
	"fmt"
	"sort"
	"jimbokun/advent/matrix"
)

func All(vs []matrix.Point, f func(matrix.Point) bool) bool {
	for _, v := range vs {
		if !f(v) {
			return false
		}
	}
	return true
}

func lowPoints(m matrix.ArrayMatrix) []matrix.Point {
	low := make([]matrix.Point, 0)
	for row := 0; row < m.Rows; row++ {
		for col := 0; col < m.Cols; col++ {
			depth := m.Get(row, col)
			p := matrix.MakePoint(row, col)
			if All(m.Adjacent(p), func(p matrix.Point) bool { return m.Get(p.X, p.Y) > depth }) {
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

func basin(m matrix.ArrayMatrix, low matrix.Point) int {
	q := make([]matrix.Point, 0)
	for _, p := range m.Adjacent(low) {
		q = append(q, p)
	}
	basin := matrix.MakeArrayMatrix(m.Rows, m.Cols)
	basinSize := 0
	for len(q) > 0 {
		p := q[0]
		if basin.Get(p.X, p.Y) == 0 && m.Get(p.X, p.Y) < 9 {
			basinSize++
			for _, p := range m.Adjacent(p) {
				q = append(q, p)
			}
			basin.Set(p.X, p.Y, 1)
		}
		q = q[1:]
	}
	fmt.Println("basin (%d, %d): %v\n", low.X, low.Y, basin)
	return basinSize
}

func Day9() {
	topo := matrix.ReadTopographicMap("day9/day9_input.txt", 100, 100)
	// topo := readTopographicMap("day9_sample_input.txt", 5, 10)
	for row := 0; row < topo.Rows; row++ {
		for col := 0; col < topo.Cols; col++ {
			fmt.Printf("adj %d, %d: %v\n", row, col, topo.AdjacentValues(row, col))
		}
	}
	top3 := make([]int, 3)
	for _, low := range lowPoints(topo) {
		top3 = append(top3, basin(topo, low))
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
