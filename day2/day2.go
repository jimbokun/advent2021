package day2

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

type position struct {
	horizontal int
	depth int
	aim int
}

func (p position) move(direction string, distance int) position {
	switch direction {
	case "forward":
		return position{horizontal: p.horizontal + distance, depth: p.depth}
	case "up":
		return position{horizontal: p.horizontal, depth: p.depth - distance}
	case "down":
		return position{horizontal: p.horizontal, depth: p.depth + distance}
	default:
		return p
	}
}

func Day2Part1() {
	file, err := os.Open("day2/day2_input.txt")

	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	p := position{horizontal: 0, depth:0}
	for scanner.Scan() {
		val := scanner.Text()
		cmdVal := strings.Split(val, " ")
		fmt.Printf("cmd: %s, val: %s\n", cmdVal[0], cmdVal[1])
		distance, _ := strconv.Atoi(cmdVal[1])
		p = p.move(cmdVal[0], distance)
		fmt.Println(p)
	}
	fmt.Println(p.horizontal * p.depth)
}
