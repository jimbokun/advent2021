package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"math"
)

func absdiff(i, j int) int {
	if i > j {
		return i - j
	} else {
		return j - i
	}
}

func expdiff(i, j int) int {
	abs := absdiff(i, j)
	cost := 0
	for k := 1; k<=abs; k++ {
		cost += k
	}
	return cost
}

func cost(crabs []int, position int, costfn func(int, int) int) int {
	cost := 0
	for _, crab := range crabs {
		cost += costfn(crab, position)
	}
	return cost
}

func minCost(crabs []int, costfn func(int, int) int) (int, int) {
	minCrab := math.MaxInt
	maxCrab := math.MinInt
	for _, crab := range crabs {
		if crab < minCrab {
			minCrab = crab
		}
		if crab > maxCrab {
			maxCrab = crab
		}
	}
	
	minCost := math.MaxInt
	var minPosition int
	
	for position := minCrab; position <= maxCrab; position++ {
		cost := cost(crabs, position, costfn)
		fmt.Printf("Cost at position %d is %d\n", position, cost)
		if cost < minCost {
			minCost = cost
			minPosition = position
		}
	}

	return minPosition, minCost
}

func readInts(filename string) []int {
	file, err := os.Open(filename)

	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	scanner.Scan()
	intVals := strings.Split(scanner.Text(), ",")
	ints := make([]int, len(intVals))
	for i := range intVals {
		j, _ := strconv.Atoi(intVals[i])
		ints[i] = j
	}

	return ints
}

func main() {
	// crabs := []int {16,1,2,0,4,2,7,1,2,14}
	crabs := readInts("day7_input.txt")
	fmt.Printf("aligning %d crabs\n", len(crabs))
	position, cost := minCost(crabs, expdiff)
	fmt.Printf("Minimum cost is %d at position %d\n", cost, position)
}
