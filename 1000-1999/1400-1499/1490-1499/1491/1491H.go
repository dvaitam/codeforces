package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n, q int
	if _, err := fmt.Fscan(in, &n, &q); err != nil {
		return
	}
	parent := make([]int, n+1)
	for i := 2; i <= n; i++ {
		fmt.Fscan(in, &parent[i])
	}
	parent[1] = 1

	// DSU to skip indices whose parent is already 1
	dsu := make([]int, n+2)
	for i := 1; i <= n+1; i++ {
		dsu[i] = i
	}
	var find func(int) int
	find = func(x int) int {
		if dsu[x] != x {
			dsu[x] = find(dsu[x])
		}
		return dsu[x]
	}
	var unite = func(x, y int) {
		x = find(x)
		y = find(y)
		if x != y {
			dsu[x] = y
		}
	}

	for ; q > 0; q-- {
		var t int
		fmt.Fscan(in, &t)
		if t == 1 {
			var l, r, x int
			fmt.Fscan(in, &l, &r, &x)
			if l < 2 {
				l = 2
			}
			for i := find(l); i <= r; i = find(i + 1) {
				if parent[i] > 1 {
					parent[i] -= x
					if parent[i] < 1 {
						parent[i] = 1
					}
					if parent[i] == 1 {
						unite(i, i+1)
					}
				}
			}
		} else if t == 2 {
			var u, v int
			fmt.Fscan(in, &u, &v)
			fmt.Fprintln(out, lca(u, v, parent))
		}
	}
}

func lca(u, v int, parent []int) int {
	for u != v {
		if u > v {
			u = parent[u]
		} else {
			v = parent[v]
		}
	}
	return u
}
