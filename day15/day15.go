package day15


import (
	"fmt"
	"os"
	"math"
	"strconv"
	"time"
	"container/heap"
	"jimbokun/advent/matrix"
)

type MinCosts struct {
	costs, minCosts matrix.Matrix
	pathEnd matrix.Point
}

func MakeMinCosts(costs matrix.Matrix, end matrix.Point) MinCosts {
	minCosts := matrix.MakeArrayMatrixWithInitialValue(costs.RowCount(), costs.ColCount(), math.MaxInt)
	minCosts.Set(end.X, end.Y, costs.Get(end.X, end.Y))
	return MinCosts{ costs: costs, minCosts: minCosts, pathEnd: end }
}

func (mc MinCosts) GetMinCost(point matrix.Point) int {
	return mc.minCosts.Get(point.X, point.Y)
}

func (mc MinCosts) SetMinCost(point matrix.Point, value int) {
	mc.minCosts.Set(point.X, point.Y, value)
}

// implements Heap interface methods
type ToVisitQueue struct {
	// priority queue of points to visit
	toVisit []matrix.Point
	// map point to its current index in the queue
	indexes matrix.Matrix
	// use to compute the ordering
	minCosts *MinCosts
}

func (tv *ToVisitQueue) Empty() bool {
	return len(tv.toVisit) == 0
}

func (tv *ToVisitQueue) Len() int { return len(tv.toVisit) }

func (tv *ToVisitQueue) Less(i, j int) bool {
	return tv.minCosts.GetMinCost(tv.toVisit[i]) < tv.minCosts.GetMinCost(tv.toVisit[j])
}

func (tv *ToVisitQueue) Swap(i, j int) {
	toVisit := tv.toVisit
	// swap
	toVisit[i], toVisit[j] = toVisit[j], toVisit[i]
	// update indexes
	tv.indexes.Set(toVisit[i].X, toVisit[i].Y, i)
	tv.indexes.Set(toVisit[j].X, toVisit[j].Y, j)
}

func (tv *ToVisitQueue) Push(x interface{}) {
	point := x.(matrix.Point)
	tv.indexes.Set(point.X, point.Y, len(tv.toVisit))
	tv.toVisit = append(tv.toVisit, point)
}

func (tv *ToVisitQueue) Pop() interface{} {
	n := len(tv.toVisit)
	current := tv.toVisit[n-1]
	tv.toVisit = tv.toVisit[0:n-1]
	tv.indexes.Set(current.X, current.Y, -1)
	return current
}

func (tv *ToVisitQueue) update(point matrix.Point) {
	// lookup index of point and Fix
	heap.Fix(tv, tv.indexes.Get(point.X, point.Y))
}

// print next n items to pop
func (tv ToVisitQueue) PrintHead(n int) {
	toVisit := tv.toVisit
	for i, p := range toVisit[0:n] {
		if i > 0 {
			fmt.Printf(", ")
		}
		fmt.Printf("(%d, %d): %d", p.X, p.Y, tv.minCosts.GetMinCost(p))
	}
	fmt.Println()
}

func (mc *MinCosts) MakeToVisit() *ToVisitQueue {
	indexes := matrix.MakeArrayMatrix(mc.costs.RowCount(), mc.costs.ColCount())
	toVisit := &ToVisitQueue{ toVisit: make([]matrix.Point, 0), minCosts: mc, indexes: indexes }
	mc.minCosts.All(func(i, j int) { heap.Push(toVisit, matrix.MakePoint(i, j)) })
	return toVisit
}

func (mc *MinCosts) computeMinPaths() {
	fmt.Println("computing min paths")
	start := time.Now()
	
	// add all points to toVisit
	toVisit := mc.MakeToVisit()
	iterations := 0
	for !toVisit.Empty() {
		// toVisit.PrintHead(5)
		if iterations % 1000 == 0 {
			fmt.Printf("total time elapsed after %d iterations %v\n", iterations, time.Since(start))
		}

		current := heap.Pop(toVisit).(matrix.Point)
		currentCost := mc.GetMinCost(current)

		for _, neighbor := range mc.minCosts.Adjacent(current) {
			newCost := mc.costs.Get(neighbor.X, neighbor.Y) + currentCost
			if mc.GetMinCost(neighbor) > newCost {
				// fmt.Printf("updating %d, %d from %d to %d + %d = %d\n", neighbor.X, neighbor.Y, mc.GetMinCost(neighbor), mc.costs.Get(neighbor.X, neighbor.Y), currentCost, newCost)
				mc.SetMinCost(neighbor, newCost)
				toVisit.update(neighbor)
			}
		}
		iterations++
	}
}

func incrementedMatrix(m matrix.Matrix) matrix.Matrix {
	incremented := matrix.MakeArrayMatrix(m.RowCount(), m.ColCount())
	m.All(func(i, j int) {
		original := m.Get(i, j)
		if original == 9 {
			incremented.Set(i, j, 1)
		} else {
			incremented.Set(i, j, original + 1)
		}
	})
	return incremented
}

// only call after computeMinPaths
func (mc *MinCosts) minPath(start matrix.Point) []matrix.Point {
	path := make([]matrix.Point, 0)
	current := start
	for {
		path = append(path, current)
		if current.X == mc.pathEnd.X && current.Y == mc.pathEnd.Y {
			return path
		}
		adj := mc.minCosts.Adjacent(current)
		cost := math.MaxInt
		for _, neighbor := range adj {
			neighborCost := mc.GetMinCost(neighbor)
			if neighborCost < cost {
				current = neighbor
				cost = neighborCost
			}
		}
	}
	return path
}

func (mc *MinCosts) pathCost(path []matrix.Point) int {
	cost := 0
	for _, p := range path {
		cost += mc.costs.Get(p.X, p.Y)
	}
	return cost
}

func CopySubMatrix(src, dest matrix.Matrix, destOrigin matrix.Point) {
	src.All(func(i, j int) {
		dest.Set(destOrigin.X + i, destOrigin.Y + j, src.Get(i, j))
	})
}

func expandMatrix(m matrix.Matrix, expandedSize int) matrix.Matrix {
	expanded := matrix.MakeArrayMatrix(m.RowCount() * expandedSize, m.ColCount() * expandedSize)
	// populate first row
	CopySubMatrix(m, expanded, matrix.MakePoint(0, 0))
	prev := m
	for expandedCol := 1; expandedCol < expandedSize; expandedCol++ {
		incremented := incrementedMatrix(prev)
		CopySubMatrix(incremented, expanded, matrix.MakePoint(0, expandedCol * incremented.ColCount()))
		prev = incremented
	}

	// populate remaining rows
	for expandedRow := 1; expandedRow < expandedSize; expandedRow++ {
		// copy first sub matrix from above row
		origin := matrix.MakePoint((expandedRow - 1) * m.RowCount(), 0)
		prev := expanded.SubMatrix(origin, m.RowCount(), m.ColCount())
		
		for expandedCol := 0; expandedCol < 5; expandedCol++ {
			incremented := incrementedMatrix(prev)
			expandedOrigin := matrix.MakePoint(expandedRow * incremented.RowCount(), expandedCol * incremented.ColCount())
			CopySubMatrix(incremented, expanded, expandedOrigin)
			prev = incremented
		}
	}
	
	return expanded
}

func Day15() {
	size, _ := strconv.Atoi(os.Args[2])
	costs := matrix.ReadTopographicMap(os.Args[1], size, size)
	expanded := expandMatrix(costs, 5)
	// expanded := costs
	// expanded.PrintTopo()
	start := matrix.MakePoint(0, 0)
	end := matrix.MakePoint(expanded.RowCount()-1, expanded.ColCount()-1)
	minCosts := MakeMinCosts(expanded, end)
	fmt.Printf("start min cost: %d\n", minCosts.GetMinCost(start))
	minCosts.computeMinPaths()
	// minCosts.minCosts.Print()
	fmt.Printf("min path cost = %d\n", minCosts.GetMinCost(start))
	fmt.Println("computing min path")
	path := minCosts.minPath(start)
	fmt.Println(path)
	fmt.Printf("computed path cost = %d\n", minCosts.pathCost(path[1:]))
	fmt.Printf("lookup min path cost from %d, %d = %d\n", path[1].X, path[1].Y, minCosts.GetMinCost(path[1]))
}
