package main

import (
	"bufio"
	"fmt"
	"os"
)

type pair struct {
	c int
	d int
}

const inf = int(1e9)

var (
	adj     [][]int
	paths   [][]pair
	sizeArr []int
	dead    []bool
	best    []int
)

func dfsSize(v, p int) int {
	sizeArr[v] = 1
	for _, to := range adj[v] {
		if to != p && !dead[to] {
			sizeArr[v] += dfsSize(to, v)
		}
	}
	return sizeArr[v]
}

func dfsCentroid(v, p, total int) int {
	for _, to := range adj[v] {
		if to != p && !dead[to] && sizeArr[to] > total/2 {
			return dfsCentroid(to, v, total)
		}
	}
	return v
}

func dfsAdd(v, p, d, cent int) {
	paths[v] = append(paths[v], pair{cent, d})
	for _, to := range adj[v] {
		if to != p && !dead[to] {
			dfsAdd(to, v, d+1, cent)
		}
	}
}

func build(v, p int) {
	total := dfsSize(v, -1)
	c := dfsCentroid(v, -1, total)
	dead[c] = true
	dfsAdd(c, -1, 0, c)
	for _, to := range adj[c] {
		if !dead[to] {
			build(to, c)
		}
	}
}

func update(x int) {
	for _, p := range paths[x] {
		if p.d < best[p.c] {
			best[p.c] = p.d
		}
	}
}

func query(x int) int {
	res := inf
	for _, p := range paths[x] {
		val := best[p.c] + p.d
		if val < res {
			res = val
		}
	}
	return res
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var t int
	fmt.Fscan(reader, &t)
	for ; t > 0; t-- {
		var n, c0 int
		fmt.Fscan(reader, &n, &c0)
		c0--
		order := make([]int, n-1)
		for i := 0; i < n-1; i++ {
			fmt.Fscan(reader, &order[i])
			order[i]--
		}
		adj = make([][]int, n)
		for i := 0; i < n-1; i++ {
			var u, v int
			fmt.Fscan(reader, &u, &v)
			u--
			v--
			adj[u] = append(adj[u], v)
			adj[v] = append(adj[v], u)
		}
		// initialize globals
		paths = make([][]pair, n)
		sizeArr = make([]int, n)
		dead = make([]bool, n)
		build(0, -1)
		best = make([]int, n)
		for i := range best {
			best[i] = inf
		}
		update(c0)
		global := inf
		ans := make([]int, n-1)
		for i, x := range order {
			dist := query(x)
			if dist < global {
				global = dist
			}
			ans[i] = global
			update(x)
		}
		for i, v := range ans {
			if i > 0 {
				fmt.Fprint(writer, " ")
			}
			fmt.Fprint(writer, v)
		}
		fmt.Fprintln(writer)
	}
}
