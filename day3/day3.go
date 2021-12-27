package day3

import (
	"fmt"
	"strconv"
	"math/bits"
	"bufio"
	"log"
	"os"
)

// const maxBit = 5
const maxBit = 12

func parseBits(val string) int64 {
	i, _ := strconv.ParseInt(fmt.Sprintf("0b%s", val), 0, 64)
	return i
}

func checkBit(i int64, mask int) int {
	return bits.OnesCount64(uint64(i & (1 << mask)))
}

func incrementBits(bitCounts []int, val string) {
	i := parseBits(val)
	for mask := 0; mask < len(bitCounts); mask++ {
		bitCounts[mask] += checkBit(i, mask)
	}
}

func gamma(bitCounts []int, total int) int {
	gamma := 0
	for i := 0; i < len(bitCounts); i++ {
		if bitCounts[i] * 2 > total {
			gamma += 1 << i
		}
	}

	return gamma
}

func epsilon(bitCounts []int, total int) int {
	epsilon := 0
	for i := 0; i < len(bitCounts); i++ {
		if bitCounts[i] * 2 < total {
			epsilon += 1 << i
		}
	}

	return epsilon
}

func powerConsumption(bitCounts []int, total int) int {
	gamma := 0
	epsilon := 0
	for i := 0; i < len(bitCounts); i++ {
		if bitCounts[i] * 2 > total {
			gamma += 1 << i
		} else {
			epsilon += 1 << i
		}
	}
	return gamma * epsilon
}

func Day3() {
	file, err := os.Open("day3/day3_input.txt")

	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	bitCounts := make([]int, maxBit)
	total := 0
	for scanner.Scan() {
		val := scanner.Text()
		incrementBits(bitCounts, val)
		fmt.Println(bitCounts)
		total++
	}
	fmt.Println(total)
	fmt.Printf("%12b\n", gamma(bitCounts, total))
	fmt.Printf("%12b\n", epsilon(bitCounts, total))
	fmt.Println(powerConsumption(bitCounts, total))
}
