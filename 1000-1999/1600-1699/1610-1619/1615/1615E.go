package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n, k int
	if _, err := fmt.Fscan(in, &n, &k); err != nil {
		return
	}

	adj := make([][]int, n+1)
	for i := 0; i < n-1; i++ {
		var u, v int
		fmt.Fscan(in, &u, &v)
		adj[u] = append(adj[u], v)
		adj[v] = append(adj[v], u)
	}

	parent := make([]int, n+1)
	depth := make([]int, n+1)
	vis := make([]bool, n+1)
	q := []int{1}
	vis[1] = true
	for len(q) > 0 {
		u := q[0]
		q = q[1:]
		for _, v := range adj[u] {
			if !vis[v] {
				vis[v] = true
				parent[v] = u
				depth[v] = depth[u] + 1
				q = append(q, v)
			}
		}
	}

	nodes := make([]int, n)
	for i := 1; i <= n; i++ {
		nodes[i-1] = i
	}
	sort.Slice(nodes, func(i, j int) bool {
		return depth[nodes[i]] > depth[nodes[j]]
	})

	up := make([]int, n+1)
	for i := 0; i <= n; i++ {
		up[i] = i
	}
	find := func(x int) int {
		for x != up[x] {
			up[x] = up[up[x]]
			x = up[x]
		}
		return x
	}

	visited := make([]bool, n+1)
	Avals := make([]int, n+1)
	A := 0
	idx := 0
	for idx < n {
		u := nodes[idx]
		idx++
		x := u
		for {
			x = find(x)
			if x == 0 || visited[x] {
				break
			}
			visited[x] = true
			A++
			up[x] = parent[x]
			x = find(x)
		}
		Avals[idx] = A
	}
	for i := idx + 1; i <= n; i++ {
		Avals[i] = A
	}

	best := int64(0)
	if k > n {
		k = n
	}
	for i := 1; i <= k; i++ {
		w := Avals[i] - i
		rb := i + Avals[i] - n
		val := int64(w) * int64(rb)
		if val > best {
			best = val
		}
	}
	fmt.Fprintln(out, best)
}
