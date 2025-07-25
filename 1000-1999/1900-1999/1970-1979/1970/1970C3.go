package main

import (
	"bufio"
	"fmt"
	"os"
)

var (
	adj      [][]int
	parent   []int
	dpDown   []bool
	dpUp     []bool
	startWin []bool
)

func dfsDown(v, p int) {
	parent[v] = p
	win := false
	for _, u := range adj[v] {
		if u == p {
			continue
		}
		dfsDown(u, v)
		if !dpDown[u] {
			win = true
		}
	}
	dpDown[v] = win
}

func dfsUp(v int) {
	neighbors := adj[v]
	m := len(neighbors)
	tmp := make([]bool, m)
	for i, u := range neighbors {
		if u == parent[v] {
			tmp[i] = dpUp[v]
		} else {
			tmp[i] = dpDown[u]
		}
	}
	prefix := make([]bool, m+1)
	for i := 0; i < m; i++ {
		prefix[i+1] = prefix[i] || !tmp[i]
	}
	suffix := make([]bool, m+1)
	for i := m - 1; i >= 0; i-- {
		suffix[i] = suffix[i+1] || !tmp[i]
	}
	startWin[v] = prefix[m]
	for i, u := range neighbors {
		if u == parent[v] {
			continue
		}
		dpUp[u] = prefix[i] || suffix[i+1]
	}
	for _, u := range neighbors {
		if u == parent[v] {
			continue
		}
		dfsUp(u)
	}
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n, t int
	if _, err := fmt.Fscan(in, &n, &t); err != nil {
		return
	}

	adj = make([][]int, n)
	for i := 0; i < n-1; i++ {
		var u, v int
		fmt.Fscan(in, &u, &v)
		u--
		v--
		adj[u] = append(adj[u], v)
		adj[v] = append(adj[v], u)
	}

	parent = make([]int, n)
	dpDown = make([]bool, n)
	dpUp = make([]bool, n)
	startWin = make([]bool, n)

	dfsDown(0, -1)
	dpUp[0] = false
	dfsUp(0)

	for i := 0; i < t; i++ {
		var x int
		fmt.Fscan(in, &x)
		if startWin[x-1] {
			fmt.Fprintln(out, "Ron")
		} else {
			fmt.Fprintln(out, "Hermione")
		}
	}
}
