package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

type lanternFish struct {
	timer int
}

func (lf *lanternFish) tick() *lanternFish {
	if lf.timer == 0 {
		lf.timer = 6
		return &lanternFish{timer: 8}
	} else {
		lf.timer = lf.timer - 1
		return nil
	}
}

type lanternFishes interface {
	print(day int)
	tick() lanternFishes
	count() int
}

type lanternFishSlice []*lanternFish

func makeFish(initState []int) lanternFishes {
	fish := make(lanternFishSlice, len(initState))
	for i := range initState {
		fish[i] = &lanternFish{timer: initState[i]}
	}
	return fish
}

func (fish lanternFishSlice) print(day int) {
	fmt.Printf("After %d days: ", day)
	for f := range fish {
		fmt.Printf(" %d", fish[f].timer)
	}
	fmt.Println()
}

func (fish lanternFishSlice) tick() lanternFishes {
	for f := range fish {
		if baby := fish[f].tick(); baby != nil {
			fish = append(fish, baby)
		}
	}
	return fish
}

func (fish lanternFishSlice) count() int {
	return len(fish)
}

type lanternFishByDay []int

func makeFishByDay(initState []int) lanternFishes {
	fish := make(lanternFishByDay, 9)
	for i := range initState {
		fish[initState[i]]++
	}
	return fish
}

func (fish lanternFishByDay) print(day int) {
	fmt.Printf("After %d days: %v\n", day, fish)
}

func (fish lanternFishByDay) tick() lanternFishes {
	next := make(lanternFishByDay, 9)

	copy(next[0:8], fish[1:9])
	next[8] = fish[0]
	next[6] += fish[0]

	return next
}

func (fish lanternFishByDay) count() int {
	count := 0
	for i := range fish {
		count += fish[i]
	}
	return count
}

func simulate(fish lanternFishes, days int) lanternFishes {
	for day := 0; day < days; day++ {
		fish = fish.tick()
		// fish.print(day)
	}
	return fish
}

func readInitState(filename string) []int {
	file, err := os.Open(filename)

	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	scanner.Scan()
	timerVals := strings.Split(scanner.Text(), ",")
	initState := make([]int, len(timerVals))
	for i := range timerVals {
		timer, _ := strconv.Atoi(timerVals[i])
		initState[i] = timer
	}

	return initState
}

func main() {
	// initState := []int { 3, 4, 3, 1, 2 }
	initState := readInitState("day6_input.txt")
	fish := makeFishByDay(initState)
	fish = simulate(fish, 256)
	fmt.Println(fish.count())
}
