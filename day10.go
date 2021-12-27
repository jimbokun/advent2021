package main

import (
	"fmt"
	"os"
	"log"
	"bufio"
	"sort"
)

var chunkOpeners = []rune("([{<")
var chunkClosers = []rune(")]}>")
var invalidChunkScores = map[rune]int { ')': 3, ']': 57, '}': 1197, '>': 25137 }
var completeScores = map[rune]int { ')': 1, ']': 2, '}': 3, '>': 4 }

type chunkParser struct {
	openers []rune
}

func pushRune(runes []rune, r rune) []rune {
	updated := make([]rune, len(runes)+1)
	updated[0] = r
	copy(updated[1:], runes)
	return updated
}

func IndexRune(runes []rune, r rune) int {
	for i, r2 := range runes {
		if r2 == r {
			return i
		}
	}
	return -1
}

func ContainsRune(runes []rune, r rune) bool {
	return IndexRune(runes, r) != -1
}

func (parser *chunkParser) pushOpener(r rune) {
	parser.openers = pushRune(parser.openers, r)
}

func (parser *chunkParser) popOpener() {
	parser.openers = parser.openers[1:]
}

func (parser *chunkParser) parse(s string) (int, rune) {
	runes := []rune(s)
	for offset, r := range runes {
		if ContainsRune(chunkOpeners, r) {
			parser.pushOpener(r)
		} else {
			closerIndex := IndexRune(chunkClosers, r)
			if chunkOpeners[closerIndex] == parser.openers[0] {
				parser.popOpener()
			} else {
				openerIndex := IndexRune(chunkOpeners, parser.openers[0])
				fmt.Printf("Expected %c, but found %c instead.\n", chunkClosers[openerIndex], r)
				return offset, r
			}
		}
		
	}
	return len(runes), ' '
}

func (parser *chunkParser) validateChunks(s string) int {
	if _, r := parser.parse(s); r != ' ' {
		return invalidChunkScores[r]
	}
	return 0
}

func (parser *chunkParser) completeChunks(s string) string {
	if _, r := parser.parse(s); r == ' ' {
		complete := make([]rune, len(parser.openers))
		for i, opener := range parser.openers {
			complete[i] = chunkClosers[IndexRune(chunkOpeners, opener)]
		}
		return string(complete)
	}
	return ""
}

func scoresForInput(filename string, score func(string) int) []int {
	file, err := os.Open(filename)

	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scores := make([]int, 0)
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		scores = append(scores, score(scanner.Text()))
	}

	return scores
}

func validateScoreForInput(filename string) int {
	scores := scoresForInput(filename, func (input string) int {
		parser := &chunkParser{}
		return parser.validateChunks(input)
	})
	total := 0
	for _, score := range scores {
		total += score
	}
	return total
}

// iterate through lines
// if validateChunks score > 0, ignore
// otherwise, generate and score completion

func completeScoreForInput(filename string) int {
	scores := scoresForInput(filename, func (input string) int {
		parser := &chunkParser{}
		complete := parser.completeChunks(input)
		total := 0
		for _, c := range complete {
			total = total * 5 + completeScores[c]
		}
		fmt.Printf("score for %s is %d\n", string(complete), total)
		return total
	})

	nonzero := make([]int, 0)
	for _, score := range scores {
		if score > 0 {
			nonzero = append(nonzero, score)
		}
	}
	sort.Slice(nonzero, func(i, j int) bool { return nonzero[i] < nonzero[j] })
	return nonzero[len(nonzero)/2]
}

func main() {
	// fmt.Println(validateScoreForInput("day10_input.txt"))
	fmt.Println(completeScoreForInput("day10_input.txt"))
}
