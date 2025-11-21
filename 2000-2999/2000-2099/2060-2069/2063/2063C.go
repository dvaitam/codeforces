package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

func edgeKey(u, v int) uint64 {
	if u > v {
		u, v = v, u
	}
	return (uint64(u) << 32) | uint64(v)
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var t int
	fmt.Fscan(in, &t)
	for ; t > 0; t-- {
		var n int
		fmt.Fscan(in, &n)

		deg := make([]int, n+1)
		edges := make(map[uint64]struct{}, n-1)
		for i := 0; i < n-1; i++ {
			var u, v int
			fmt.Fscan(in, &u, &v)
			deg[u]++
			deg[v]++
			edges[edgeKey(u, v)] = struct{}{}
		}

		nodes := make([]int, n)
		for i := 0; i < n; i++ {
			nodes[i] = i + 1
		}
		sort.Slice(nodes, func(i, j int) bool {
			if deg[nodes[i]] == deg[nodes[j]] {
				return nodes[i] < nodes[j]
			}
			return deg[nodes[i]] > deg[nodes[j]]
		})

		k := n
		if k > 2000 {
			k = 2000
		}

		best := 0
		for i := 0; i < k; i++ {
			u := nodes[i]
			for j := i + 1; j < k; j++ {
				v := nodes[j]
				adj := 0
				if _, ok := edges[edgeKey(u, v)]; ok {
					adj = 1
				}
				val := deg[u] + deg[v] - 1 - adj
				if val > best {
					best = val
				}
			}
		}

		fmt.Fprintln(out, best)
	}
}
