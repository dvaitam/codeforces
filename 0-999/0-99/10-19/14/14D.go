package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
)

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Split(bufio.ScanWords)

	readInt := func() int {
		scanner.Scan()
		val, _ := strconv.Atoi(scanner.Text())
		return val
	}

	if !scanner.Scan() {
		return
	}
	n, _ := strconv.Atoi(scanner.Text())

	adj := make([][]int, n+1)
	type Edge struct {
		u, v int
	}
	edges := make([]Edge, 0, n-1)

	for i := 0; i < n-1; i++ {
		u := readInt()
		v := readInt()
		adj[u] = append(adj[u], v)
		adj[v] = append(adj[v], u)
		edges = append(edges, Edge{u, v})
	}

	// BFS function to find farthest node and its distance
	// restrictedEdgeU, restrictedEdgeV define the edge to ignore
	bfs := func(start int, restrictedU, restrictedV int) (int, int) {
		dist := make([]int, n+1)
		for i := range dist {
			dist[i] = -1
		}
		dist[start] = 0
		queue := []int{start}
		farthestNode := start
		maxDist := 0

		for len(queue) > 0 {
			u := queue[0]
			queue = queue[1:]

			if dist[u] > maxDist {
				maxDist = dist[u]
				farthestNode = u
			}

			for _, v := range adj[u] {
				// Check if this is the restricted edge
				if (u == restrictedU && v == restrictedV) || (u == restrictedV && v == restrictedU) {
					continue
				}
				if dist[v] == -1 {
					dist[v] = dist[u] + 1
					queue = append(queue, v)
				}
			}
		}
		return farthestNode, maxDist
	}

	getDiameter := func(start int, rU, rV int) int {
		// First BFS to find one endpoint of the diameter
		farthest, _ := bfs(start, rU, rV)
		// Second BFS to find the actual diameter length
		_, diam := bfs(farthest, rU, rV)
		return diam
	}

	maxProd := 0
	for _, e := range edges {
		d1 := getDiameter(e.u, e.u, e.v)
		d2 := getDiameter(e.v, e.u, e.v)
		prod := d1 * d2
		if prod > maxProd {
			maxProd = prod
		}
	}

	fmt.Println(maxProd)
}
