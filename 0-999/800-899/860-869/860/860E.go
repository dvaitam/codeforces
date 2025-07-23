package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	in := bufio.NewReader(os.Stdin)
	var n int
	if _, err := fmt.Fscan(in, &n); err != nil {
		return
	}
	parent := make([]int, n+1)
	root := 0
	for i := 1; i <= n; i++ {
		fmt.Fscan(in, &parent[i])
		if parent[i] == 0 {
			root = i
		}
	}
	g := make([][]int, n+1)
	for i := 1; i <= n; i++ {
		if parent[i] != 0 {
			g[parent[i]] = append(g[parent[i]], i)
		}
	}
	depth := make([]int, n+1)
	var dfsDepth func(int, int)
	dfsDepth = func(v, d int) {
		depth[v] = d
		for _, to := range g[v] {
			dfsDepth(to, d+1)
		}
	}
	dfsDepth(root, 1)

	// collect subtree nodes for each node (inefficient)
	sub := make([][]int, n+1)
	var dfsSub func(int) []int
	dfsSub = func(v int) []int {
		nodes := []int{v}
		for _, to := range g[v] {
			nodes = append(nodes, dfsSub(to)...)
		}
		sub[v] = nodes
		return nodes
	}
	dfsSub(root)

	// compute ancestors for each node
	ancestors := make([][]int, n+1)
	var dfsAnc func(int, []int)
	dfsAnc = func(v int, path []int) {
		ancestors[v] = append([]int(nil), path...)
		path = append(path, v)
		for _, to := range g[v] {
			dfsAnc(to, path)
		}
	}
	dfsAnc(root, nil)

	z := make([]int, n+1)
	for a := 1; a <= n; a++ {
		for _, b := range ancestors[a] {
			if b == a {
				continue
			}
			cnt := 0
			for _, x := range sub[b] {
				if x != b && depth[x] <= depth[a] {
					cnt++
				}
			}
			z[a] += cnt
		}
	}

	out := bufio.NewWriter(os.Stdout)
	for i := 1; i <= n; i++ {
		if i > 1 {
			out.WriteByte(' ')
		}
		fmt.Fprint(out, z[i])
	}
	out.WriteByte('\n')
	out.Flush()
}
