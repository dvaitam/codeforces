package main

import (
	"bufio"
	"fmt"
	"os"
)

type edge struct{ to, id int }

var (
	adj   [][]edge
	disc  []int
	low   []int
	sz    []int
	timer int
	n     int
	ans   int64
)

func dfs(v, pe int) {
	timer++
	disc[v] = timer
	low[v] = timer
	sz[v] = 1
	for _, e := range adj[v] {
		if e.id == pe {
			continue
		}
		to := e.to
		if disc[to] == 0 {
			dfs(to, e.id)
			sz[v] += sz[to]
			if low[to] > disc[v] {
				s := sz[to]
				val := int64(s*(s-1)/2 + (n-s)*(n-s-1)/2)
				if val < ans {
					ans = val
				}
			}
			if low[to] < low[v] {
				low[v] = low[to]
			}
		} else {
			if disc[to] < low[v] {
				low[v] = disc[to]
			}
		}
	}
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var t int
	fmt.Fscan(reader, &t)
	for ; t > 0; t-- {
		var m int
		fmt.Fscan(reader, &n, &m)
		adj = make([][]edge, n)
		id := 0
		for i := 0; i < m; i++ {
			var u, v int
			fmt.Fscan(reader, &u, &v)
			u--
			v--
			adj[u] = append(adj[u], edge{v, id})
			adj[v] = append(adj[v], edge{u, id})
			id++
		}
		disc = make([]int, n)
		low = make([]int, n)
		sz = make([]int, n)
		timer = 0
		ans = int64(n * (n - 1) / 2)
		dfs(0, -1)
		fmt.Fprintln(writer, ans)
	}
}
