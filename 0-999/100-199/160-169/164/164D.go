package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"strconv"
)

type Point struct {
	x, y int64
}

func dist2(p1, p2 Point) int64 {
	dx := p1.x - p2.x
	dy := p1.y - p2.y
	return dx*dx + dy*dy
}

type int64Slice []int64

func (x int64Slice) Len() int           { return len(x) }
func (x int64Slice) Less(i, j int) bool { return x[i] < x[j] }
func (x int64Slice) Swap(i, j int)      { x[i], x[j] = x[j], x[i] }

func solveVC(adj [][]int, active []bool, k int) []int {
	n := len(adj)
	deg := make([]int, n)
	maxDeg := 0
	maxU := -1
	edgeCount := 0

	for i := 0; i < n; i++ {
		if !active[i] || len(adj[i]) == 0 {
			continue
		}
		d := 0
		for _, v := range adj[i] {
			if active[v] {
				d++
			}
		}
		deg[i] = d
		edgeCount += d
		if d > maxDeg {
			maxDeg = d
			maxU = i
		}
	}
	edgeCount /= 2

	if edgeCount == 0 {
		return []int{}
	}
	if k == 0 {
		return nil
	}
	if edgeCount > k*maxDeg {
		return nil
	}

	if maxDeg <= 2 {
		leaf := -1
		for i := 0; i < n; i++ {
			if active[i] && deg[i] == 1 {
				leaf = i
				break
			}
		}
		if leaf != -1 {
			parent := -1
			for _, v := range adj[leaf] {
				if active[v] {
					parent = v
					break
				}
			}
			active[parent] = false
			res := solveVC(adj, active, k-1)
			active[parent] = true
			if res != nil {
				return append(res, parent)
			}
			return nil
		}

		u := -1
		for i := 0; i < n; i++ {
			if active[i] && deg[i] == 2 {
				u = i
				break
			}
		}
		if u != -1 {
			active[u] = false
			res := solveVC(adj, active, k-1)
			active[u] = true
			if res != nil {
				return append(res, u)
			}
		}
		return nil
	}

	if maxDeg > k {
		active[maxU] = false
		res := solveVC(adj, active, k-1)
		active[maxU] = true
		if res != nil {
			return append(res, maxU)
		}
		return nil
	}

	active[maxU] = false
	res1 := solveVC(adj, active, k-1)
	active[maxU] = true
	if res1 != nil {
		return append(res1, maxU)
	}

	neighbors := make([]int, 0, deg[maxU])
	for _, v := range adj[maxU] {
		if active[v] {
			neighbors = append(neighbors, v)
		}
	}
	if k >= len(neighbors) {
		for _, v := range neighbors {
			active[v] = false
		}
		res2 := solveVC(adj, active, k-len(neighbors))
		for _, v := range neighbors {
			active[v] = true
		}
		if res2 != nil {
			return append(res2, neighbors...)
		}
	}

	return nil
}

func check(D2 int64, pts []Point, k int) []int {
	n := len(pts)
	edgeCount := 0
	for i := 0; i < n; i++ {
		for j := i + 1; j < n; j++ {
			if dist2(pts[i], pts[j]) > D2 {
				edgeCount++
			}
		}
	}
	if edgeCount > k*(n-1) {
		return nil
	}

	adj := make([][]int, n)
	for i := 0; i < n; i++ {
		for j := i + 1; j < n; j++ {
			if dist2(pts[i], pts[j]) > D2 {
				adj[i] = append(adj[i], j)
				adj[j] = append(adj[j], i)
			}
		}
	}

	active := make([]bool, n)
	for i := range active {
		active[i] = true
	}

	return solveVC(adj, active, k)
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Split(bufio.ScanWords)

	nextInt := func() int {
		if scanner.Scan() {
			res, _ := strconv.Atoi(scanner.Text())
			return res
		}
		return 0
	}

	n := nextInt()
	k := nextInt()
	if n == 0 {
		return
	}

	pts := make([]Point, n)
	for i := 0; i < n; i++ {
		pts[i].x = int64(nextInt())
		pts[i].y = int64(nextInt())
	}

	dists := make([]int64, 0, n*(n-1)/2)
	for i := 0; i < n; i++ {
		for j := i + 1; j < n; j++ {
			dists = append(dists, dist2(pts[i], pts[j]))
		}
	}
	sort.Sort(int64Slice(dists))

	uniqueDists := make([]int64, 0, len(dists)+1)
	uniqueDists = append(uniqueDists, 0)
	if len(dists) > 0 {
		if dists[0] != 0 {
			uniqueDists = append(uniqueDists, dists[0])
		}
		for i := 1; i < len(dists); i++ {
			if dists[i] != dists[i-1] {
				uniqueDists = append(uniqueDists, dists[i])
			}
		}
	}

	low := 0
	high := len(uniqueDists) - 1
	var bestRes []int

	for low <= high {
		mid := low + (high-low)/2
		D2 := uniqueDists[mid]

		res := check(D2, pts, k)
		if res != nil {
			bestRes = res
			high = mid - 1
		} else {
			low = mid + 1
		}
	}

	deleted := make([]bool, n)
	for _, idx := range bestRes {
		deleted[idx] = true
	}

	var finalAns []int
	for _, idx := range bestRes {
		finalAns = append(finalAns, idx+1)
	}

	for i := 0; i < n && len(finalAns) < k; i++ {
		if !deleted[i] {
			finalAns = append(finalAns, i+1)
			deleted[i] = true
		}
	}

	for i, val := range finalAns {
		if i > 0 {
			fmt.Print(" ")
		}
		fmt.Print(val)
	}
	fmt.Println()
}
