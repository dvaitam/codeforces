package main

import (
	"fmt"
	"math"
)

// This program attempts to solve the following problem:
// "A tourist wants to walk through all the streets in the city and end his
// journey at the same point where he started. What minimum distance will the
// tourist have to walk?" The original statement references a figure listing
// the street lengths. Since that figure is not available in this repository,
// the program uses a small example graph to demonstrate an approach for
// computing the Chinese Postman distance.
//
// The example graph has four intersections (0..3) with the following streets:
// 0-1:1, 1-2:1, 2-3:1, 3-0:1, 0-2:2, 1-3:2.
// The minimal closed walk visiting every street at least once has length 10.
// Replace the edges array with the actual streets from the figure to obtain
// the correct answer for the intended problem.

const n = 4

var edges = []struct{ u, v, w int }{
	{0, 1, 1},
	{1, 2, 1},
	{2, 3, 1},
	{3, 0, 1},
	{0, 2, 2},
	{1, 3, 2},
}

func main() {
	// Sum of all street lengths
	sum := 0
	deg := make([]int, n)
	dist := make([][]int, n)
	for i := range dist {
		dist[i] = make([]int, n)
		for j := range dist[i] {
			if i == j {
				dist[i][j] = 0
			} else {
				dist[i][j] = math.MaxInt32 / 2
			}
		}
	}
	for _, e := range edges {
		sum += e.w
		deg[e.u]++
		deg[e.v]++
		if e.w < dist[e.u][e.v] {
			dist[e.u][e.v] = e.w
			dist[e.v][e.u] = e.w
		}
	}

	// All-pairs shortest paths
	for k := 0; k < n; k++ {
		for i := 0; i < n; i++ {
			for j := 0; j < n; j++ {
				if dist[i][k]+dist[k][j] < dist[i][j] {
					dist[i][j] = dist[i][k] + dist[k][j]
				}
			}
		}
	}

	// Collect odd degree vertices
	var odd []int
	for i, d := range deg {
		if d%2 == 1 {
			odd = append(odd, i)
		}
	}

	// DP over subsets to find minimal pairing cost
	m := len(odd)
	dp := make([]int, 1<<m)
	for i := range dp {
		dp[i] = math.MaxInt32 / 2
	}
	dp[0] = 0
	for mask := 0; mask < (1 << m); mask++ {
		if dp[mask] == math.MaxInt32/2 {
			continue
		}
		// find first unmatched vertex
		var p int
		for p = 0; p < m; p++ {
			if (mask>>p)&1 == 0 {
				break
			}
		}
		if p == m {
			continue
		}
		for q := p + 1; q < m; q++ {
			if (mask>>q)&1 == 0 {
				next := mask | (1 << p) | (1 << q)
				cost := dist[odd[p]][odd[q]]
				if dp[mask]+cost < dp[next] {
					dp[next] = dp[mask] + cost
				}
			}
		}
	}

	// The minimum closed walk length equals the total length of all streets
	// plus the minimal extra distance needed to pair up odd-degree vertices.
	result := sum + dp[(1<<m)-1]
	fmt.Println(result)
}
