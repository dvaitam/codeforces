package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

type Edge struct{ v, w int }

var adj [][]Edge

// Find path from u to target in tree, return edge weights on path
func findPath(u, target, parent int) ([]int, bool) {
	if u == target {
		return nil, true
	}
	for _, e := range adj[u] {
		if e.v == parent {
			continue
		}
		if path, ok := findPath(e.v, target, u); ok {
			return append([]int{e.w}, path...), true
		}
	}
	return nil, false
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var n, l, r int
	fmt.Fscan(reader, &n, &l, &r)

	adj = make([][]Edge, n+1)
	for i := 0; i < n-1; i++ {
		var a, b, c int
		fmt.Fscan(reader, &a, &b, &c)
		adj[a] = append(adj[a], Edge{b, c})
		adj[b] = append(adj[b], Edge{a, c})
	}

	bestMed := -1
	bestU, bestV := 1, 1

	for u := 1; u <= n; u++ {
		for v := u + 1; v <= n; v++ {
			path, ok := findPath(u, v, -1)
			if !ok {
				continue
			}
			length := len(path)
			if length < l || length > r {
				continue
			}
			sorted := make([]int, length)
			copy(sorted, path)
			sort.Ints(sorted)
			med := sorted[length/2]
			if med > bestMed {
				bestMed = med
				bestU = u
				bestV = v
			}
		}
	}

	fmt.Fprintln(writer, bestU, bestV)
}
