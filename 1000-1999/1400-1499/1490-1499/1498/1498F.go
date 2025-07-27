package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	in := bufio.NewReader(os.Stdin)
	var n, k int
	if _, err := fmt.Fscan(in, &n, &k); err != nil {
		return
	}
	adj := make([][]int, n)
	for i := 0; i < n-1; i++ {
		var x, y int
		fmt.Fscan(in, &x, &y)
		x--
		y--
		adj[x] = append(adj[x], y)
		adj[y] = append(adj[y], x)
	}
	a := make([]int, n)
	for i := range a {
		fmt.Fscan(in, &a[i])
	}

	K := 2 * k
	sub := make([][]int, n)
	for i := 0; i < n; i++ {
		sub[i] = make([]int, K)
	}
	parent := make([]int, n)
	for i := range parent {
		parent[i] = -1
	}
	order := make([]int, 0, n)
	stack := []int{0}
	parent[0] = 0
	for len(stack) > 0 {
		v := stack[len(stack)-1]
		stack = stack[:len(stack)-1]
		order = append(order, v)
		for _, u := range adj[v] {
			if u == parent[v] {
				continue
			}
			parent[u] = v
			stack = append(stack, u)
		}
	}
	for i := len(order) - 1; i >= 0; i-- {
		v := order[i]
		sub[v][0] ^= a[v]
		for _, u := range adj[v] {
			if u == parent[v] {
				continue
			}
			for j := 0; j < K; j++ {
				sub[v][(j+1)%K] ^= sub[u][j]
			}
		}
	}

	all := make([][]int, n)
	for i := 0; i < n; i++ {
		all[i] = make([]int, K)
	}
	copy(all[0], sub[0])
	stack = []int{0}
	for len(stack) > 0 {
		v := stack[len(stack)-1]
		stack = stack[:len(stack)-1]
		for _, u := range adj[v] {
			if u == parent[v] {
				continue
			}
			rest := make([]int, K)
			for i := 0; i < K; i++ {
				rest[i] = all[v][i] ^ sub[u][(i-1+K)%K]
			}
			up := make([]int, K)
			for i := 0; i < K; i++ {
				up[(i+1)%K] = rest[i]
			}
			for i := 0; i < K; i++ {
				all[u][i] = sub[u][i] ^ up[i]
			}
			stack = append(stack, u)
		}
	}

	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()
	for i := 0; i < n; i++ {
		x := 0
		for j := k; j < K; j++ {
			x ^= all[i][j]
		}
		if x != 0 {
			fmt.Fprint(out, 1)
		} else {
			fmt.Fprint(out, 0)
		}
		if i+1 < n {
			fmt.Fprint(out, " ")
		}
	}
	fmt.Fprintln(out)
}
