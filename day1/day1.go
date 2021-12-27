package day1

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"math"
)


func Day1() {
	file, err := os.Open("input.txt")

	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	prev := math.MaxUint32
	incrs := 0
	for scanner.Scan() {
		val, _ := strconv.Atoi(scanner.Text())
		if (prev < val) {
			incrs++
		}
		prev = val
	}
	fmt.Printf("%d values are greater than the line before", incrs)
}
