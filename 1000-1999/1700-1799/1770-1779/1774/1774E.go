package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var n, d int
	if _, err := fmt.Fscan(reader, &n, &d); err != nil {
		return
	}
	adj := make([][]int, n+1)
	for i := 0; i < n-1; i++ {
		var u, v int
		fmt.Fscan(reader, &u, &v)
		adj[u] = append(adj[u], v)
		adj[v] = append(adj[v], u)
	}
	var m1 int
	fmt.Fscan(reader, &m1)
	a := make([]int, m1)
	for i := 0; i < m1; i++ {
		fmt.Fscan(reader, &a[i])
	}
	var m2 int
	fmt.Fscan(reader, &m2)
	b := make([]int, m2)
	for i := 0; i < m2; i++ {
		fmt.Fscan(reader, &b[i])
	}

	log := 0
	for (1 << log) <= n {
		log++
	}
	up := make([][]int, log)
	for i := range up {
		up[i] = make([]int, n+1)
	}
	depth := make([]int, n+1)
	parent := make([]int, n+1)
	parent[1] = 0
	up[0][1] = 1
	stack := []int{1}
	order := make([]int, 0, n)

	// Build parent/depth arrays with iterative DFS to avoid recursion limits.
	for len(stack) > 0 {
		u := stack[len(stack)-1]
		stack = stack[:len(stack)-1]
		order = append(order, u)
		for _, v := range adj[u] {
			if v == parent[u] {
				continue
			}
			parent[v] = u
			depth[v] = depth[u] + 1
			up[0][v] = u
			stack = append(stack, v)
		}
	}

	for i := 1; i < log; i++ {
		for v := 1; v <= n; v++ {
			up[i][v] = up[i-1][up[i-1][v]]
		}
	}

	kthAncestor := func(u, k int) int {
		for i := 0; k > 0 && u != 1; i++ {
			if k&1 == 1 {
				u = up[i][u]
			}
			k >>= 1
		}
		return u
	}

	need1 := make([]bool, n+1)
	need2 := make([]bool, n+1)
	for _, v := range a {
		need1[v] = true
		if depth[v] > d {
			anc := kthAncestor(v, d)
			need2[anc] = true
		}
	}
	for _, v := range b {
		need2[v] = true
		if depth[v] > d {
			anc := kthAncestor(v, d)
			need1[anc] = true
		}
	}

	has1 := make([]bool, n+1)
	has2 := make([]bool, n+1)
	var edges1, edges2 int
	for i := len(order) - 1; i >= 0; i-- {
		u := order[i]
		total1 := need1[u]
		total2 := need2[u]
		for _, v := range adj[u] {
			if parent[v] != u {
				continue
			}
			if has1[v] {
				edges1++
				total1 = true
			}
			if has2[v] {
				edges2++
				total2 = true
			}
		}
		has1[u] = total1
		has2[u] = total2
	}

	ans := int64(edges1+edges2) * 2
	fmt.Fprintln(writer, ans)
}
