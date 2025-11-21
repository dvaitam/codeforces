package main

import (
	"bufio"
	"fmt"
	"os"
)

const LOG = 20

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func solve(in *bufio.Reader, out *bufio.Writer) {
	var n int
	if _, err := fmt.Fscan(in, &n); err != nil {
		return
	}

	adj := make([][]int, n+1)
	for i := 0; i < n-1; i++ {
		var x, y int
		fmt.Fscan(in, &x, &y)
		adj[x] = append(adj[x], y)
		adj[y] = append(adj[y], x)
	}

	depth := make([]int, n+1)
	parent := make([][LOG]int, n+1)
	edgeMax := make([][LOG]int, n+1) // max val along 2^p edges upward starting from the current node

	order := make([]int, 0, n)
	stack := []int{1}
	parent[1][0] = 1
	for len(stack) > 0 {
		u := stack[len(stack)-1]
		stack = stack[:len(stack)-1]
		order = append(order, u)
		for _, v := range adj[u] {
			if v == parent[u][0] {
				continue
			}
			parent[v][0] = u
			depth[v] = depth[u] + 1
			stack = append(stack, v)
		}
	}

	down := make([]int, n+1)
	for i := len(order) - 1; i >= 0; i-- {
		u := order[i]
		for _, v := range adj[u] {
			if v == parent[u][0] {
				continue
			}
			down[u] = max(down[u], 1+down[v])
		}
	}

	// Prepare best other depths for each child.
	top1 := make([]int, n+1)
	top2 := make([]int, n+1)
	bestChild := make([]int, n+1) // index of child giving top1
	for u := 1; u <= n; u++ {
		for _, v := range adj[u] {
			if v == parent[u][0] {
				continue
			}
			val := down[v] + 1
			if val > top1[u] {
				top2[u] = top1[u]
				top1[u] = val
				bestChild[u] = v
			} else if val > top2[u] {
				top2[u] = val
			}
		}
	}

	for v := 1; v <= n; v++ {
		p := parent[v][0]
		bestOther := 0
		if p != v {
			if bestChild[p] == v {
				bestOther = top2[p]
			} else {
				bestOther = top1[p]
			}
		}
		edgeMax[v][0] = bestOther - depth[p]
	}

	for p := 1; p < LOG; p++ {
		for v := 1; v <= n; v++ {
			ancestor := parent[v][p-1]
			parent[v][p] = parent[ancestor][p-1]
			edgeMax[v][p] = max(edgeMax[v][p-1], edgeMax[ancestor][p-1])
		}
	}

	var q int
	fmt.Fscan(in, &q)
	for ; q > 0; q-- {
		var v, k int
		fmt.Fscan(in, &v, &k)
		if k > depth[v] {
			k = depth[v]
		}
		ansMax := down[v] // lca = v case
		cur := v
		rem := k
		for p := LOG - 1; p >= 0; p-- {
			if (1 << p) <= rem {
				ansMax = max(ansMax, depth[v]+edgeMax[cur][p])
				cur = parent[cur][p]
				rem -= 1 << p
			}
		}
		fmt.Fprint(out, ansMax, " ")
	}
	fmt.Fprintln(out)
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var T int
	fmt.Fscan(in, &T)
	for ; T > 0; T-- {
		solve(in, out)
	}
}
