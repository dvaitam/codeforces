package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

// attempt solution for problem F
// We try to binary search the minimal distance between chosen vertices
// and greedily pick vertices when forced by uncovered edges.

func feasible(n int, pre [][]int, edges [][2]int, D int) bool {
	chosen := make([]bool, n+1)
	last := -D // last chosen index
	for i := 1; i <= n; i++ {
		needed := false
		for _, j := range pre[i] {
			if !chosen[j] {
				needed = true
				break
			}
		}
		if needed {
			if i-last < D {
				return false
			}
			chosen[i] = true
			last = i
		}
	}
	for _, e := range edges {
		if !chosen[e[0]] && !chosen[e[1]] {
			return false
		}
	}
	// check spacing
	prev := -D
	for i := 1; i <= n; i++ {
		if chosen[i] {
			if i-prev < D {
				return false
			}
			prev = i
		}
	}
	return true
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var t int
	fmt.Fscan(in, &t)
	for ; t > 0; t-- {
		var n, m int
		fmt.Fscan(in, &n, &m)
		pre := make([][]int, n+1)
		edges := make([][2]int, m)
		for i := 0; i < m; i++ {
			var u, v int
			fmt.Fscan(in, &u, &v)
			if u > v {
				u, v = v, u
			}
			edges[i] = [2]int{u, v}
			pre[v] = append(pre[v], u)
		}
		sort.Slice(edges, func(i, j int) bool {
			if edges[i][1] == edges[j][1] {
				return edges[i][0] < edges[j][0]
			}
			return edges[i][1] < edges[j][1]
		})
		l, r, ans := 1, n, 1
		for l <= r {
			mid := (l + r) >> 1
			if feasible(n, pre, edges, mid) {
				ans = mid
				l = mid + 1
			} else {
				r = mid - 1
			}
		}
		fmt.Fprintln(out, ans)
	}
}
