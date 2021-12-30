package day13


import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
	"strconv"
	"jimbokun/advent/matrix"
)

type FoldAxis int

const (
	X FoldAxis = iota
	Y
)

type FoldInstruction struct {
	axis FoldAxis
	foldLine int
}

type Paper struct {
	dots matrix.ArrayMatrix
	folds []FoldInstruction
}

func parseFoldInstruction(line string) FoldInstruction {
	instruction := line[len("fold along "):]
	instructionParts := strings.Split(instruction, "=")
	var axis FoldAxis
	switch instructionParts[0] {
	case "x":
		axis = X
	case "y":
		axis = Y
	}
	foldLine, _ := strconv.Atoi(instructionParts[1])

	return FoldInstruction{axis: axis, foldLine: foldLine}
}

func (p Paper) fold(instruction FoldInstruction) Paper {
	dots := p.dots
	switch instruction.axis {
	case X:
		// calculate dimensions of folded paper
		// make matrix for folded paper
		folded := matrix.MakeArrayMatrix(dots.Rows, instruction.foldLine)
		// copy dots "above the fold"
		for row := 0; row < dots.Rows; row++ {
			for col := 0; col < folded.Cols; col++ {
				folded.Set(row, col, dots.Get(row, col))
			}
		}
		// for dots below the fold, transform coordiantes to folded coordinates
		// (fold X, Y stays the same, new X is max X - X)
		for row := 0; row < dots.Rows; row++ {
			for col := folded.Cols; col < dots.Cols; col++ {
				if dots.Get(row, col) == 1 {
					foldCol := folded.Cols - (col - folded.Cols)
					folded.Set(row, foldCol, dots.Get(row, col))
				}
			}
		}
		return Paper{ dots: folded, folds: p.folds }
	case Y:
		// calculate dimensions of folded paper
		// make matrix for folded paper
		folded := matrix.MakeArrayMatrix(instruction.foldLine, dots.Cols)
		// copy dots "above the fold"
		for row := 0; row < folded.Rows; row++ {
			for col := 0; col < dots.Cols; col++ {
				folded.Set(row, col, dots.Get(row, col))
			}
		}
		// for dots below the fold, transform coordiantes to folded coordinates
		// (fold Y, X stays the same, new Y is max Y - Y)
		for row := folded.Rows; row < dots.Rows; row++ {
			for col := 0; col < dots.Cols; col++ {
				if dots.Get(row, col) == 1 {
					foldRow := folded.Rows - (row - folded.Rows)
					folded.Set(foldRow, col, dots.Get(row, col))
				}
			}
		}
		return Paper{ dots: folded, folds: p.folds }
	default:
		return p
	}
}

func (p Paper) printDots() {
	// print each dot as "#", one row at a time
	dots := p.dots
	for row := 0; row < dots.Rows; row++ {
		line := make([]rune, dots.Cols)
		for col := 0; col < dots.Cols; col++ {
			if dots.Get(row, col) == 1 {
				line[col] = '#'
			} else {
				line[col] = '.'
			}
		}
		fmt.Println(string(line))
	}
}

func (p Paper) countDots() int {
	dots := p.dots
	count := 0
	for row := 0; row < dots.Rows; row++ {
		for col := 0; col < dots.Cols; col++ {
			if dots.Get(row, col) == 1 {
				count++
			}
		}
	}
	return count
}

func readPaper(filename string) Paper {
	file, err := os.Open(filename)

	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	points := readPoints(scanner)

	max := maxPoint(points)
	dots := matrix.MakeArrayMatrix(max.Y+1, max.X+1)
	for _, p := range points {
		dots.Set(p.Y, p.X, 1)
	}

	instructions := make([]FoldInstruction, 0)
	for scanner.Scan() {
		instructions = append(instructions, parseFoldInstruction(scanner.Text()))
	}

	return Paper{dots: dots, folds: instructions}
}

func readPoints(scanner *bufio.Scanner) []matrix.Point {
	points := make([]matrix.Point, 0)
	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			break
		}
		points = append(points, matrix.ParsePoint(line))
	}
	return points
}

func maxPoint(points []matrix.Point) matrix.Point {
	max := matrix.MakePoint(0, 0)
	for _, p := range points {
		if p.X > max.X {
			max = matrix.MakePoint(p.X, max.Y)
		}
		if p.Y > max.Y {
			max = matrix.MakePoint(max.X, p.Y)
		}
	}
	return max
}

func Day13() {
	paper := readPaper(os.Args[1])
	// fmt.Println(paper)
	fmt.Println("original:")
	paper.printDots()
	for i, fold := range paper.folds {
		paper = paper.fold(fold)
		fmt.Printf("after fold %d:\n", i)
		paper.printDots()
		fmt.Printf("found %d dots after fold %d\n", paper.countDots(), i)
	}
}
