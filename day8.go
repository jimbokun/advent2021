package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
	"sort"
	
)

type sevenSegmentDisplay struct {
	digitsByLen [8]map[string]bool
	digits [10]string
	output []string
}

func parseSevenSegmentDisplay(line string) sevenSegmentDisplay {
	var ssd sevenSegmentDisplay
	for i, _  :=  range ssd.digitsByLen {
		ssd.digitsByLen[i] = make(map[string]bool, 0)
	}

	digitsOutputs := strings.Split(line, "|")

	digits := strings.Fields(digitsOutputs[0])
	for _, digit := range digits {
		ssd.digitsByLen[len(digit)][sortRunes(digit)] = true
	}

	outputs := strings.Fields(digitsOutputs[1])
	ssd.output = make([]string, len(outputs))
	for i, output := range outputs {
		ssd.output[i] = sortRunes(output)
	}
	
	return ssd
}

func readSevenSegmentDisplays(filename string) []sevenSegmentDisplay {
	file, err := os.Open(filename)

	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	ssds := make([]sevenSegmentDisplay, 0)
	for scanner.Scan() {
		line := scanner.Text()
		ssds = append(ssds, parseSevenSegmentDisplay(line))
	}

	return ssds
}

func sortRunes(s string) string {
	runes := []rune(s)
	sort.Slice(runes, func(i, j int) bool {	return runes[i] < runes[j] })
	return string(runes)
}

func ContainsAll(s1 string, s2 string) bool {
	return len(s2) == len(intersection(s1, s2))
}

func intersection(s1 string, s2 string) string {
	s1chars := make(map[rune]bool, len(s1))
	for _, r := range s1 {
		s1chars[r] = true
	}
	var intersection []rune
	for _, r := range s2 {
		if _, ok := s1chars[r]; ok {
			intersection = append(intersection, r)
		}
	}
	return string(intersection)
}

func solveSevenSegmentDisplay(ssd *sevenSegmentDisplay) {
	// 1 is only digit made up of 2 segments
	for k, _ := range ssd.digitsByLen[2] {
		ssd.digits[1] = k
	}
	delete(ssd.digitsByLen[2], ssd.digits[1])
	// 4 is only digit made up of 4 segments
	for k, _ := range ssd.digitsByLen[4] {
		ssd.digits[4] = k
	}
	delete(ssd.digitsByLen[4], ssd.digits[4])
	// 7 is only digit made up of 3 segments
	for k, _ := range ssd.digitsByLen[3] {
		ssd.digits[7] = k
	}
	delete(ssd.digitsByLen[3], ssd.digits[7])
	// 8 is only digit made up of 7 segments
	for k, _ := range ssd.digitsByLen[7] {
		ssd.digits[8] = k
	}
	delete(ssd.digitsByLen[7], ssd.digits[8])
	// 3 is only digit of 5 segments containing all the segements of digit 1
	for k, _ := range ssd.digitsByLen[5] {
		if ContainsAll(k, ssd.digits[1]) {
			ssd.digits[3] = k
			break
		}
	}
	delete(ssd.digitsByLen[5], ssd.digits[3])
	// of digits with 6 segements, only 6 does *not* contain the segements of digit 1
	for k, _ := range ssd.digitsByLen[6] {
		if !ContainsAll(k, ssd.digits[1]) {
			ssd.digits[6] = k
			break
		}
	}
	delete(ssd.digitsByLen[6], ssd.digits[6])
	// 5 and 6 share the bottom right segment
	bottomRight := intersection(ssd.digits[1], ssd.digits[6])
	for k, _ := range ssd.digitsByLen[5] {
		if ContainsAll(k, bottomRight) {
			ssd.digits[5] = k
			break
		}
	}
	delete(ssd.digitsByLen[5], ssd.digits[5])
	// 2 is only remaining digit with 5 segments
	for k, _ := range ssd.digitsByLen[5] {
		ssd.digits[2] = k
	}
	delete(ssd.digitsByLen[5], ssd.digits[2])
	// 9 shares all the segments of 4
	for k, _ := range ssd.digitsByLen[6] {
		if ContainsAll(k, ssd.digits[4]) {
			ssd.digits[9] = k
			break
		}
	}
	delete(ssd.digitsByLen[6], ssd.digits[9])
	// 0 is only remaining digit with 6 segments
	for k, _ := range ssd.digitsByLen[6] {
		ssd.digits[0] = k
	}
	delete(ssd.digitsByLen[6], ssd.digits[0])
}

func (ssd sevenSegmentDisplay) decodeOutput() int {
	digitsMap := make(map[string]int, 10)
	for i, d := range ssd.digits {
		digitsMap[d] = i
	}
	decoded := 0
	for _, d := range ssd.output {
		decoded = 10 * decoded + digitsMap[d]
	}
	return decoded
}

func main() {
	ssds := readSevenSegmentDisplays("day8_input.txt")
	total := 0 
	for _, ssd := range ssds {
		solveSevenSegmentDisplay(&ssd)
		// fmt.Println(ssd)
		// for i, d := range ssd.digits {
		// 	fmt.Printf(" %d: %s", i, d)
		// }
		// fmt.Println()
		decoded := ssd.decodeOutput()
		fmt.Println(decoded)
		total += decoded
	}
	fmt.Println(total)
}
