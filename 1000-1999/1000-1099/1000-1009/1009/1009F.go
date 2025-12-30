package main

import (
	"bufio"
	"fmt"
	"os"
)

const N = 1000010

var (
	adj   [N][]int
	t     [N][]int
	mx    [N]int
	ans   [N]int
	n     int
)

// merge merges two slices using the smaller-to-larger trick.
// p and q represent the indices of the maximum values in slices x and y.
func merge(u, v int) {
	// If the child's slice is larger, swap it with the parent's slice.
	// This is equivalent to x.swap(y) in C++.
	if len(t[u]) < len(t[v]) {
		t[u], t[v] = t[v], t[u]
		mx[u], mx[v] = mx[v], mx[u]
	}

	// Merge the smaller slice into the larger one.
	// Aligns the ends of the slices (depth-based).
	du := len(t[u])
	dv := len(t[v])
	diff := du - dv

	for i := 0; i < dv; i++ {
		targetIdx := i + diff
		t[u][targetIdx] += t[v][i]

		// Update the index of the maximum frequency (mx[u]).
		// Condition: new frequency is strictly greater, 
		// or equal but at a shallower depth (smaller index relative to the end).
		if t[u][targetIdx] > t[u][mx[u]] || (t[u][targetIdx] == t[u][mx[u]] && targetIdx > mx[u]) {
			mx[u] = targetIdx
		}
	}
	// Free memory for the merged-away slice
	t[v] = nil 
}

func dfs(u, fa int) {
	for _, v := range adj[u] {
		if v != fa {
			dfs(v, u)
			merge(u, v)
		}
	}

	// Add the current node itself (depth 0 relative to itself)
	t[u] = append(t[u], 1)
	
	// If the current node (count 1) is the new maximum frequency
	if t[u][mx[u]] == 1 {
		mx[u] = len(t[u]) - 1
	}
	
	// The answer is the depth (distance from u), which is (TotalLength - 1) - mxIdx
	ans[u] = len(t[u]) - 1 - mx[u]
}

func main() {
	// Use fast I/O for 10^6 nodes
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	fmt.Fscan(reader, &n)
	for i := 1; i < n; i++ {
		var u, v int
		fmt.Fscan(reader, &u, &v)
		adj[u] = append(adj[u], v)
		adj[v] = append(adj[v], u)
	}

	dfs(1, 0)

	for i := 1; i <= n; i++ {
		fmt.Fprintln(writer, ans[i])
	}
}
