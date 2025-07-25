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

	var n int
	// read number of nodes
	if _, err := fmt.Fscan(in, &n); err != nil {
		return
	}
	// build tree
	adj := make([][]int, n+1)
	for i := 0; i < n-1; i++ {
		var u, v int
		fmt.Fscan(in, &u, &v)
		adj[u] = append(adj[u], v)
		adj[v] = append(adj[v], u)
	}
	// read your subtree
	var k1 int
	fmt.Fscan(in, &k1)
	x := make([]int, k1)
	for i := 0; i < k1; i++ {
		fmt.Fscan(in, &x[i])
	}
	// read Li Chen's subtree
	var k2 int
	fmt.Fscan(in, &k2)
	y := make([]int, k2)
	inY := make(map[int]bool, k2)
	for i := 0; i < k2; i++ {
		fmt.Fscan(in, &y[i])
		inY[y[i]] = true
	}
	// ask B y[0]
	fmt.Fprintf(out, "B %d\n", y[0])
	out.Flush()
	// receive the corresponding your-label
	var u int
	fmt.Fscan(in, &u)
	// BFS from u to find nearest in x
	dist := make([]int, n+1)
	for i := 1; i <= n; i++ {
		dist[i] = -1
	}
	queue := make([]int, 0, n)
	queue = append(queue, u)
	dist[u] = 0
	for idx := 0; idx < len(queue); idx++ {
		v := queue[idx]
		for _, w := range adj[v] {
			if dist[w] == -1 {
				dist[w] = dist[v] + 1
				queue = append(queue, w)
			}
		}
	}
	// find closest x
	best := x[0]
	for _, v := range x {
		if dist[v] >= 0 && dist[v] < dist[best] {
			best = v
		}
	}
	// ask A best
	fmt.Fprintf(out, "A %d\n", best)
	out.Flush()
	var pbest int
	fmt.Fscan(in, &pbest)
	// output result
	if inY[pbest] {
		fmt.Fprintf(out, "C %d\n", best)
	} else {
		fmt.Fprint(out, "C -1\n")
	}
}
