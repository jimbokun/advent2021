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

func (g *Graph) addEdgeToNode(from, to int) {
	fromNode := g.nodes[from]
	g.nodes[from] = Node{
		label: fromNode.label,
		isBig: fromNode.isBig,
		edges: append(fromNode.edges, to) }
}

func (g *Graph) addEdge(from, to string) {
	fromIndex := g.ensureNode(from)
	toIndex := g.ensureNode(to)
	g.addEdgeToNode(fromIndex, toIndex)
	g.addEdgeToNode(toIndex, fromIndex)
}

func (g *Graph) printPath(path []int) {
	for i, node := range path {
		if i > 0 {
			fmt.Print(",")
		}
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
		
		for _, to := range g.nodes[from].edges {
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
	g := readGraph("day12/sample_input.txt")
	// g.print()
	visited := make([]int, len(g.nodes))
	count := 0
	g.paths([]int {}, visited, g.byLabel["start"], g.byLabel["end"], func(path []int) {
		g.printPath(path);
		count++
	})
	fmt.Printf("Found %d paths\n", count)
}
