package main

import (
	"bufio"
	"fmt"
	"os"
)

const LOG = 20

func lca(u, v int, depth []int, up [][]int) int {
	if depth[u] < depth[v] {
		u, v = v, u
	}
	diff := depth[u] - depth[v]
	for k := 0; k < LOG; k++ {
		if diff>>k&1 == 1 {
			u = up[k][u]
		}
	}
	if u == v {
		return u
	}
	for k := LOG - 1; k >= 0; k-- {
		if up[k][u] != up[k][v] {
			u = up[k][u]
			v = up[k][v]
		}
	}
	return up[0][u]
}

func dist(u, v int, depth []int, up [][]int) int {
	l := lca(u, v, depth, up)
	return depth[u] + depth[v] - 2*depth[l]
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
		p := make([]int, n+1)
		for i := 2; i <= n; i++ {
			fmt.Fscan(in, &p[i])
		}

		up := make([][]int, LOG)
		for k := 0; k < LOG; k++ {
			up[k] = make([]int, n+1)
		}
		depth := make([]int, n+1)
		deg := make([]int, n+1)
		ans := make([]int, n+1)

		for k := 0; k < LOG; k++ {
			up[k][1] = 1
		}
		depth[1] = 0
		deg[1] = 0
		ans[1] = 1

		a, b := 1, 1
		length := 0
		invalid := false

		for i := 2; i <= n; i++ {
			parent := p[i]
			deg[parent]++
			if deg[parent] > 3 {
				invalid = true
			}
			deg[i] = 1
			depth[i] = depth[parent] + 1
			up[0][i] = parent
			for k := 1; k < LOG; k++ {
				up[k][i] = up[k-1][up[k-1][i]]
			}
			if invalid {
				ans[i] = -1
				continue
			}
			d := dist(i, a, depth, up)
			if d > length {
				length = d
				b = i
			} else {
				d = dist(i, b, depth, up)
				if d > length {
					length = d
					a = i
				}
			}
			radius := (length + 1) / 2
			ans[i] = radius + 1
		}

		for i := 1; i <= n; i++ {
			if i > 1 {
				fmt.Fprint(out, " ")
			}
			fmt.Fprint(out, ans[i])
		}
		fmt.Fprintln(out)
	}
}
