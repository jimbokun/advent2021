// lineSegment struct
// from, to point
// method points return []point containing points between from and to
// determing if from.x == to.x, or from.y == to.y
// and then if from.y > from.y or from.y < from.y, to iterate correctly
// build n*n [][]int to represent the plane
// iterate through line segments and add all the corresponding points by incrementing plane matrix value
// alternatively map[int]map[int]int for sparse matrix representation

package day5

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"jimbokun/advent/matrix"
)

// type matrix interface {
// 	get(i, j int) int
// 	increment(i, j int)
// }

// type arrayMatrix struct {
// 	rows, cols int
// 	values []int
// }

// func (m arrayMatrix) get(i, j int) int {
// 	return m.values[m.cols * j + i]
// }

// func (m matrix.ArrayMatrix) increment(i, j int) {
// 	m.values[m.Cols * j + i]++
// }

// func makeArrayMatrix(rows, cols int) ArrayMatrix {
// 	return ArrayMatrix{rows: rows, cols: cols, values: make([]int, rows * cols)}
// }

type point struct {
	x, y int
}

type lineSegment struct {
	from, to point
}

func (p1 point) add(p2 point) point {
	return point{x: p1.x + p2.x, y: p1.y + p2.y}
}

func (ls lineSegment) points() []point {
	points := make([]point, 0)

	// figure out direction between from and to by comparing points
	var delta point

	if ls.from == ls.to {
		delta = point{0, 0}
	} else if ls.from.x == ls.to.x && ls.from.y > ls.to.y {
		delta = point{0, -1}
	} else if ls.from.x == ls.to.x && ls.from.y < ls.to.y {
		delta = point{0, 1}
	} else if ls.from.x > ls.to.x && ls.from.y == ls.to.y {
		delta = point{-1, 0}
	} else if ls.from.x < ls.to.x && ls.from.y == ls.to.y {
		delta = point{1, 0}
	} else if ls.from.x > ls.to.x && ls.from.y > ls.to.y {
		delta = point{-1, -1}
	} else if ls.from.x < ls.to.x && ls.from.y < ls.to.y {
		delta = point{1, 1}
	} else if ls.from.x < ls.to.x && ls.from.y > ls.to.y {
		delta = point{1, -1}
	} else if ls.from.x > ls.to.x && ls.from.y < ls.to.y {
		delta = point{-1, 1}
	}

	fmt.Printf("from %v\n", ls.from)
	fmt.Printf("to %v\n", ls.to)
	fmt.Printf("delta %v\n", delta)
	for p := ls.from; p != ls.to; p = p.add(delta) {
		fmt.Printf("adding point %v\n", p)
		points = append(points, p)
	}
	points = append(points, ls.to)

	return points
}

func addSegment(m matrix.Matrix, ls lineSegment) {
	points := ls.points()
	for p := range points {
		m.Increment(points[p].x, points[p].y)
	}
}

func makeSegment(x1, y1, x2, y2 int) lineSegment {
	return lineSegment{from: point{x: x1, y: y1}, to: point{x: x2, y: y2}}
}

func parsePoint(pointVal string) point {
	xyVals := strings.Split(pointVal, ",")
	x, _ := strconv.Atoi(xyVals[0])
	y, _ := strconv.Atoi(xyVals[1])
	return point{x: x, y: y}
}

func parseLineSegment(line string) lineSegment {
	pointVals := strings.Split(line, " -> ")
	return lineSegment{from: parsePoint(pointVals[0]), to: parsePoint(pointVals[1])}
}

func readLineSegments(filename string) []lineSegment {
	file, err := os.Open(filename)

	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	lineSegments := make([]lineSegment, 0)
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lineSegments = append(lineSegments, parseLineSegment(scanner.Text()))
	}

	return lineSegments
}


func testSegments() []lineSegment {
	return []lineSegment { makeSegment(0,9, 5,9),
		makeSegment(8,0, 0,8),
		makeSegment(9,4, 3,4),
		makeSegment(2,2, 2,1),
		makeSegment(7,0, 7,4),
		makeSegment(6,4, 2,0),
		makeSegment(0,9, 2,9),
		makeSegment(3,4, 1,4),
		makeSegment(0,0, 8,8),
		makeSegment(5,5, 8,2) }
}

func Day5() {
	m := matrix.MakeArrayMatrix(1000, 1000)
	// segments := testSegments()
	segments := readLineSegments("day5/day5_input.txt")
	for s := range segments {
		addSegment(m, segments[s])
	}

	// for i := 0; i < 10; i++ {
	// 	fmt.Printf("%v\n", m.values[i * 10:(i+1) * 10])
	// }

	intersects := 0
	for v := range m.Values {
		if m.Values[v] > 1 {
			intersects++
		}
	}
	fmt.Printf("found %d intersections\n", intersects)
}
