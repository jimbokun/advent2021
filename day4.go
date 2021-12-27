// card is 5x5 int[] (given card size of 5)
// bingo is sequence of balls drawn plus card[]
// test for win against initial subsequence:
//   - row or column is subset of the subsequence
//   - row is just a range
//   - column calculated by slicing increments of card size
// read bingo:
//  - first line is sequence of balls drawn
//  - read card until input exhausted
//  - read card: read new line, read row and split on space until reaching card size

package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

const cardSize = 5

type card [][]int
type bingo struct {
	balls []int
	cards []card
}

func subset(s1 []int, s2 map[int]bool) bool {
	for i := 0; i < len(s1); i++ {
		if _, ok := s2[s1[i]]; !ok {
			return false
		}
	}
	return true
}

func checkCard(ballSet map[int]bool, c card) bool {
	// check rows
	for row := 0; row < cardSize; row++ {
		if subset(c[row], ballSet) {
			fmt.Printf("WINNER! row %d: %v\n", row, c[row])
			return true
		}
	}
	// check columns
	for col := 0; col < cardSize; col++ {
		column := make([]int, cardSize)
		for i := 0; i < cardSize; i++ {
			column[i] = c[i][col]
		}
		if subset(column, ballSet) {
			fmt.Printf("WINNER! column %d: %v\n", col, column)
			return true
		}
	}
	return false
}

func checkWinner(cards []card, ballSet map[int]bool) card {
	for c := range cards {
		if checkCard(ballSet, cards[c]) {
			fmt.Printf("Found a winner! %v\n", cards[c])
			return cards[c]
		}
	}
	return nil
}

func sliceToSet(slice []int) map[int]bool {
	set := make(map[int]bool, len(slice))
	for i := 0; i < len(slice); i++ {
		set[slice[i]] = true
	}
	return set
}

func (b bingo) play() (map[int]bool, card, int) {
	for round := 0; round < len(b.balls); round++ {
		balls := b.balls[:round]
		ballSet := sliceToSet(balls)
		fmt.Printf("round %d: %v\n", round, balls)
		if c := checkWinner(b.cards, ballSet); c != nil {
			return ballSet, c, balls[len(balls)-1]
		}
	}
	return nil, nil, 0
}

func (b bingo) lastWinner() (map[int]bool, card, int) {
	var lastWinner card
	var lastBallSet map[int]bool
	var lastBall int

	winningCards := make([]bool, len(b.cards))

	fmt.Printf("bingo game has %d cards\n", len(b.cards))
	for round := 0; round < len(b.balls); round++ {
		balls := b.balls[:round]
		ballSet := sliceToSet(balls)
		fmt.Printf("round %d: %v\n", round, balls)
		for c := range b.cards {
			if !winningCards[c] {
				if checkCard(ballSet, b.cards[c]) {
					winningCards[c] = true
					lastBallSet = ballSet
					lastWinner = b.cards[c]
					lastBall = balls[len(balls)-1]
				}
				
			}
		}
	}
	return lastBallSet, lastWinner, lastBall
}

func (b bingo) score(play func() (map[int]bool, card, int)) int {
	ballSet, winningCard, lastBall := play()
	fmt.Printf("Scoring ballSet %v, winningCard %v, lastBall %d\n", ballSet, winningCard, lastBall)
	score := 0
	for i := 0; i < cardSize; i++ {
		for j := 0; j < cardSize; j++ {
			if _, ok := ballSet[winningCard[i][j]]; !ok {
				score += winningCard[i][j]
			}
		}
	}
	return score * lastBall
}

func testGame() bingo {
	balls := []int {7,4,9,5,11,17,23,2,0,14,21,24,10,16,13,6,15,25,12,22,18,20,8,19,3,26,1}
	card1 := [][]int {
		{22, 13, 17, 11, 0},
		{8, 2, 23, 4, 24},
		{21, 9, 14, 16, 7},
		{6, 10, 3, 18, 5},
		{1, 12, 20, 15, 19} }
	card2 := [][]int { {3, 15, 0, 2, 22},
		{9, 18, 13, 17, 5},
		{19, 8, 7, 25, 23},
		{20, 11, 10, 24, 4},
		{14, 21, 16, 12, 6} }
	card3 := [][]int { {14, 21, 17, 24, 4},
		{10, 16, 15, 9, 19},
		{18, 8, 23, 26, 20},
		{22, 11, 13, 6, 5},
		{2, 0, 12, 3, 7} }
	return bingo{balls: balls, cards: []card{card1, card2, card3}}
}

func readGame(filename string) bingo {
	file, err := os.Open(filename)

	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	scanner.Scan()
	ballVals := strings.Split(scanner.Text(), ",")
	balls := make([]int, 0)
	for i := range ballVals {
		ball, _ := strconv.Atoi(ballVals[i])
		balls = append(balls, ball)
	}

	var cards []card
	var c card
	var row int
	for scanner.Scan() {
		line := scanner.Text()
		fmt.Println(line)
		if line == "" {
			if c != nil {
				cards = append(cards, c)
			}
			c = make([][]int, cardSize)
			for i := range c {
				c[i] = make([]int, cardSize)
			}
			row = 0
		} else {
			rowVals := strings.Fields(line)
			for i := range rowVals {
				cardVal, _ := strconv.Atoi(rowVals[i])
				fmt.Printf("row: %d, col: %d, val: %d\n", row, i, cardVal)
				c[row][i] = cardVal
			}
			row++
		}
	}
	if c != nil {
		cards = append(cards, c)
	}

	return bingo{balls: balls, cards: cards}
}

func main() {
	bingoGame := readGame("bingo_input.txt")
	fmt.Println(bingoGame)
	fmt.Printf("Score: %d\n", bingoGame.score(bingoGame.play))
	fmt.Printf("Last Winner Score: %d\n", bingoGame.score(bingoGame.lastWinner))
}
