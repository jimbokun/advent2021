// lineSegment struct
// from, to point
// method points return []point containing points between from and to
// determing if from.X == to.X, or from.Y == to.Y
// and then if from.Y > from.Y or from.Y < from.Y, to iterate correctly
// build n*n [][]int to represent the plane
// iterate through line segments and add all the corresponding points by incrementing plane matrix value
// alternatively map[int]map[int]int for sparse matrix representation

package day5

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
	"jimbokun/advent/matrix"
)

type lineSegment struct {
	from, to matrix.Point
}

func (ls lineSegment) points() []matrix.Point {
	points := make([]matrix.Point, 0)

	// figure out direction between from and to by comparing points
	var delta matrix.Point

	if ls.from == ls.to {
		delta = matrix.MakePoint(0, 0)
	} else if ls.from.X == ls.to.X && ls.from.Y > ls.to.Y {
		delta = matrix.MakePoint(0, -1)
	} else if ls.from.X == ls.to.X && ls.from.Y < ls.to.Y {
		delta = matrix.MakePoint(0, 1)
	} else if ls.from.X > ls.to.X && ls.from.Y == ls.to.Y {
		delta = matrix.MakePoint(-1, 0)
	} else if ls.from.X < ls.to.X && ls.from.Y == ls.to.Y {
		delta = matrix.MakePoint(1, 0)
	} else if ls.from.X > ls.to.X && ls.from.Y > ls.to.Y {
		delta = matrix.MakePoint(-1, -1)
	} else if ls.from.X < ls.to.X && ls.from.Y < ls.to.Y {
		delta = matrix.MakePoint(1, 1)
	} else if ls.from.X < ls.to.X && ls.from.Y > ls.to.Y {
		delta = matrix.MakePoint(1, -1)
	} else if ls.from.X > ls.to.X && ls.from.Y < ls.to.Y {
		delta = matrix.MakePoint(-1, 1)
	}

	fmt.Printf("from %v\n", ls.from)
	fmt.Printf("to %v\n", ls.to)
	fmt.Printf("delta %v\n", delta)
	for p := ls.from; p != ls.to; p = p.Add(delta) {
		fmt.Printf("adding point %v\n", p)
		points = append(points, p)
	}
	points = append(points, ls.to)

	return points
}

func addSegment(m matrix.Matrix, ls lineSegment) {
	points := ls.points()
	for p := range points {
		m.Increment(points[p].X, points[p].Y)
	}
}

func makeSegment(x1, y1, x2, y2 int) lineSegment {
	return lineSegment{from: matrix.MakePoint(x1, y1), to: matrix.MakePoint(x2, y2)}
}

func parseLineSegment(line string) lineSegment {
	pointVals := strings.Split(line, " -> ")
	return lineSegment{from: matrix.ParsePoint(pointVals[0]), to: matrix.ParsePoint(pointVals[1])}
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
