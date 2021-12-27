package day3

import (
	"fmt"
	// "strconv"
	// "math/bits"
	"bufio"
	"log"
	"os"
)

// const maxBit = 5
// const maxBit = 12

// func parseBits(val string) int64 {
// 	i, _ := strconv.ParseInt(fmt.Sprintf("0b%s", val), 0, 64)
// 	return i
// }

// func checkBit(i int64, mask int) int {
// 	return bits.OnesCount64(uint64(i & (1 << mask)))
// }

func incrementBitCounts(binaries []int64) []int {
	bitCounts := make([]int, maxBit)
	for i := 0; i < len(binaries); i++ {
		for mask := 0; mask < len(bitCounts); mask++ {
			bitCounts[mask] += checkBit(binaries[i], mask)
		}
	}
	return bitCounts
}

// func gamma(bitCounts []int, total int) int {
// 	gamma := 0
// 	for i := 0; i < len(bitCounts); i++ {
// 		if bitCounts[i] * 2 > total {
// 			gamma += 1 << i
// 		}
// 	}

// 	return gamma
// }

// func epsilon(bitCounts []int, total int) int {
// 	epsilon := 0
// 	for i := 0; i < len(bitCounts); i++ {
// 		if bitCounts[i] * 2 < total {
// 			epsilon += 1 << i
// 		}
// 	}

// 	return epsilon
// }

// func powerConsumption(bitCounts []int, total int) int {
// 	gamma := 0
// 	epsilon := 0
// 	for i := 0; i < len(bitCounts); i++ {
// 		if bitCounts[i] * 2 > total {
// 			gamma += 1 << i
// 		} else {
// 			epsilon += 1 << i
// 		}
// 	}
// 	return gamma * epsilon
// }

func oxygenGeneratorCheckBit(binaries []int64, mask int) int {
	bitCounts := incrementBitCounts(binaries)
	if bitCounts[mask] * 2 >= len(binaries) {
		return 1
	} else {
		return 0
	}
}

func co2ScrubberRatingCheckBit(binaries []int64, mask int) int {
	bitCounts := incrementBitCounts(binaries)
	if bitCounts[mask] * 2 >= len(binaries) {
		return 0
	} else {
		return 1
	}
}

func printBinaries(binaries []int64) {
	fmt.Printf("[")
	for i := 0; i < len(binaries); i++ {
		fmt.Printf(" %12b", binaries[i])
	}
	fmt.Printf("]")
}

func filterBinaries(binaries []int64, filterBit func(binaries []int64, mask int) int) int64 {
	for mask := maxBit-1; mask >= 0; mask-- {
		keep := make([]int64, 0)
		check := filterBit(binaries, mask)
		for i := 0; i < len(binaries); i++ {
			if checkBit(binaries[i], mask) == check {
				keep = append(keep, binaries[i])
			}
		}

		// fmt.Printf("keep: ")
		// printBinaries(keep)
		// fmt.Printf("\n")
		
		if len(keep) == 1 {
			return keep[0]
		}
		binaries = keep
	}
	// to shut up compiler, should probably panic or something
	return binaries[0]
}

func lifeSupportRating(binaries []int64) int64 {
	return filterBinaries(binaries, oxygenGeneratorCheckBit) * filterBinaries(binaries, co2ScrubberRatingCheckBit)
}

func readBinaries(fname string) []int64 {
	file, err := os.Open(fname)

	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	binaries := make([]int64, 1)
	for scanner.Scan() {
		binaries = append(binaries, parseBits(scanner.Text()))
	}

	return binaries
}

func Day3Part2() {
	binaries := readBinaries("day3/day3_input.txt")
	fmt.Printf("oxygen generator value: %12b\n", filterBinaries(binaries, oxygenGeneratorCheckBit))
	fmt.Printf("co2 scrubber value: %12b\n", filterBinaries(binaries, co2ScrubberRatingCheckBit))
	fmt.Printf("life support rating: %d\n", lifeSupportRating(binaries))
}
