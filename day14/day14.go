package day14


import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
	"math"
)

type runeCount map[rune]int

type Grammar struct {
	initialTemplate string
	rules map[string]rune
	table map[string][]runeCount
}

func combineCounts(leftCounts, rightCounts runeCount) runeCount {
	counts := make(runeCount, 0)
	for r, count := range leftCounts {
		counts[r] += count
	}
	for r, count := range rightCounts {
		counts[r] += count
	}
	return counts
}

func (g *Grammar) computeRuneCounts(lhs string, steps int) runeCount {
	insert, ok := g.rules[lhs]
	lhsRunes := []rune(lhs)

	if !ok || steps == 0 {
		counts := make(runeCount, 0)
		for _, r := range(lhsRunes) {
			counts[r]++
		}
		return counts
	} else if counts := g.table[lhs][steps]; counts != nil {
		return counts
	} else {
		leftCounts := g.computeRuneCounts(string([]rune { lhsRunes[0], insert }), steps-1)
		rightCounts := g.computeRuneCounts(string([]rune { insert, lhsRunes[1] }), steps-1)
		counts := combineCounts(leftCounts, rightCounts)
		counts[insert]--
		g.table[lhs][steps] = counts
		return counts
	}
}

func (rc runeCount) printCounts() {
	fmt.Printf("{")
	for r, count := range rc {
		fmt.Printf(" %c: %d", r, count)
	}
	fmt.Printf("}\n")
}

func (g *Grammar) countsAfterSteps(template string, steps int) runeCount {
	counts := make(runeCount, 0)
	templateRunes := []rune(template)
	for i, r := range templateRunes {
		if i > 0 {
			context := string(templateRunes[i-1:i+1])
			contextCounts := g.computeRuneCounts(context, steps)
			counts = combineCounts(counts, contextCounts)
			counts[r]--
		}
	}
	counts[templateRunes[len(templateRunes)-1]]++
	return counts
}

func (g Grammar) performInsertions(template string) string {
	var sb strings.Builder
	templateRunes := []rune(template)
	for i, r := range templateRunes {
		if i > 0 {
			context := string(templateRunes[i-1:i+1])
			if insert, ok := g.rules[context]; ok {
				sb.WriteRune(insert)
			}
		}
		sb.WriteRune(r)
	}
	return sb.String()
}

func (g Grammar) printGrammar() {
	fmt.Println(g.initialTemplate)
	fmt.Println()
	for context, insert := range g.rules {
		fmt.Printf("%s -> %c\n", context, insert)
	}
}

func readGrammar(filename string, steps int) Grammar {
	file, err := os.Open(filename)

	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	scanner.Scan()
	g := Grammar{ initialTemplate: scanner.Text(), rules: make(map[string]rune, 0) }
	scanner.Scan()
	for scanner.Scan() {
		ruleParts := strings.Split(scanner.Text(), " -> ")
		g.rules[ruleParts[0]] = []rune(ruleParts[1])[0]
	}
	g.table = make(map[string][]runeCount, len(g.rules))
	for lhs, _ := range g.rules {
		g.table[lhs] = make([]runeCount, steps+1)
	}
	for lhs, _ := range g.rules {
		for i := 0; i < steps; i++ {
			g.computeRuneCounts(lhs, steps)
		}
	}

	return g
}

func Day14() {
	steps := 40
	grammar := readGrammar(os.Args[1], steps)
	grammar.printGrammar()
	template := grammar.initialTemplate
	var counts runeCount
	for step := 0; step < steps+1; step++ {
		counts = grammar.countsAfterSteps(template, step)
		fmt.Printf("after step %d\n", step)
		counts.printCounts()
	}
	minCount := math.MaxInt
	maxCount := 0
	for r, count := range counts {
		if count > maxCount {
			maxCount = count
		}
		if count < minCount {
			minCount = count
		}
		fmt.Printf("%c occurs %d times\n", r, count)
	}
	fmt.Printf("max %d min %d difference %d\n", maxCount, minCount, maxCount - minCount)
}
