package day12

import (
	"fmt"
	"strings"
	"os"
	"log"
	"bufio"
)

type Node struct {
	label string
	isBig bool
	edges []int
}

type Graph struct {
	nodes []Node
	byLabel map[string]int
}

func newNode(label string) Node {
	return Node{
		label: label,
		isBig: strings.ToUpper(label) == label,
		edges: make([]int, 0) }
}

func (g *Graph) ensureNode(label string) int {
	if _, ok := g.byLabel[label]; !ok {
		g.nodes = append(g.nodes, newNode(label))
		g.byLabel[label] = len(g.nodes)-1
	}
	return g.byLabel[label]
}

func (g *Graph) addEdge(from, to string) {
	fromIndex := g.ensureNode(from)
	toIndex := g.ensureNode(to)
	fromNode := g.nodes[fromIndex]
	g.nodes[fromIndex] = Node{
		label: fromNode.label,
		isBig: fromNode.isBig,
		edges: append(fromNode.edges, toIndex) }
	toNode := g.nodes[toIndex]
	g.nodes[toIndex] = Node{
		label: toNode.label,
		isBig: toNode.isBig,
		edges: append(toNode.edges, fromIndex) }
}

func (g *Graph) printPath(path []int) {
	first := true
	for _, node := range path {
		if !first {
			fmt.Print(",")
		}
		first = false
		fmt.Print(g.nodes[node].label)
	}
	fmt.Println()
}

func (g *Graph) willVisit(visited []int, to int) bool {
	if g.nodes[to].isBig || visited[to] == 0 {
		return true
	}
	if visited[to] == 1 && (g.byLabel["start"] == to || g.byLabel["end"] == to) {
		return false
	}
	for i, visitCount := range visited {
		if !g.nodes[i].isBig && visitCount > 1 {
			return false
		}
	}
	return true
}

func (g *Graph) paths(soFar []int, visited []int, from, end int, handlePath func(path []int)) {
	if from == end {
		handlePath(append(soFar, from))
	} else {
		nextSoFar := append(soFar, from)
		
		nextVisited := make([]int, len(visited))
		copy(nextVisited[:], visited)
		nextVisited[from]++
		
		fromNode := g.nodes[from]
		for _, to := range fromNode.edges {
			if g.willVisit(nextVisited, to) {
				g.paths(nextSoFar, nextVisited, to, end, handlePath)
			}
		}
	}
}

func (g *Graph) parseEdge(serializedEdge string) {
	labels := strings.Split(serializedEdge, "-")
	g.addEdge(labels[0], labels[1])
}

func (g Graph) print() {
	for i, node := range g.nodes {
		fmt.Printf("node %d label %s\n", i, node.label)
	}
	for _, fromNode := range g.nodes {
		for _, to := range fromNode.edges {
			fmt.Printf("%s-%s\n", fromNode.label, g.nodes[to].label)
		}
	}
}

func readGraph(filename string) *Graph {
	file, err := os.Open(filename)

	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()


	g := &Graph{ nodes: make([]Node, 0), byLabel: make(map[string]int, 0) }
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		g.parseEdge(scanner.Text())
	}

	return g
}

func Day12() {
	g := readGraph("day12/input.txt")
	// g.print()
	visited := make([]int, len(g.nodes))
	count := 0
	g.paths([]int {}, visited, g.byLabel["start"], g.byLabel["end"], func(path []int) {
		g.printPath(path);
		count++
	})
	fmt.Printf("Found %d paths\n", count)
}
