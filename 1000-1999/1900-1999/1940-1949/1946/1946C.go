package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

// This program solves the problem described in problemC.txt.
// We are given a tree and must remove exactly k edges so that
// every resulting connected component contains at least x vertices.
// To find the maximum possible x, we binary search on x and
// check feasibility.
//
// For a fixed x, a DFS greedily cuts a subtree whenever its size
// reaches x. Each cut subtree becomes an independent component.
// All such components are disjoint. To maximize the size of the
// remaining component containing the root while cutting exactly
// k edges, it is optimal to cut the k smallest available subtrees.
// If after cutting them the size of the remaining vertices is at
// least x, then x is feasible.

func can(g [][]int, n, k, x int) bool {
	sizes := make([]int, 0)
	var dfs func(v, p int) int
	dfs = func(v, p int) int {
		sum := 1
		for _, to := range g[v] {
			if to == p {
				continue
			}
			sum += dfs(to, v)
		}
		if sum >= x {
			sizes = append(sizes, sum)
			return 0
		}
		return sum
	}
	root := dfs(0, -1)
	_ = root // root size equals n - sum(sizes), not used directly

	if len(sizes) < k {
		return false
	}
	sort.Ints(sizes)
	total := 0
	for i := 0; i < k; i++ {
		total += sizes[i]
	}
	return n-total >= x
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var t int
	fmt.Fscan(in, &t)
	for ; t > 0; t-- {
		var n, k int
		fmt.Fscan(in, &n, &k)
		g := make([][]int, n)
		for i := 0; i < n-1; i++ {
			var u, v int
			fmt.Fscan(in, &u, &v)
			u--
			v--
			g[u] = append(g[u], v)
			g[v] = append(g[v], u)
		}

		low, high, ans := 1, n, 1
		for low <= high {
			mid := (low + high) / 2
			if can(g, n, k, mid) {
				ans = mid
				low = mid + 1
			} else {
				high = mid - 1
			}
		}
		fmt.Fprintln(out, ans)
	}
}
