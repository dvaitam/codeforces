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

	var n, m int
	fmt.Fscan(reader, &n, &m)

	g := make([][]int, n)
	rg := make([][]int, n)
	indeg := make([]int, n)
	for i := 0; i < m; i++ {
		var u, v int
		fmt.Fscan(reader, &u, &v)
		u--
		v--
		g[u] = append(g[u], v)
		rg[v] = append(rg[v], u)
		indeg[v]++
	}

	// topological sort to detect cycles
	q := make([]int, 0, n)
	for i := 0; i < n; i++ {
		if indeg[i] == 0 {
			q = append(q, i)
		}
	}
	order := make([]int, 0, n)
	for len(q) > 0 {
		v := q[0]
		q = q[1:]
		order = append(order, v)
		for _, to := range g[v] {
			indeg[to]--
			if indeg[to] == 0 {
				q = append(q, to)
			}
		}
	}
	if len(order) < n {
		fmt.Fprintln(writer, -1)
		return
	}

	down := make([]int, n)
	for i := 0; i < n; i++ {
		down[i] = i
	}
	for i := n - 1; i >= 0; i-- {
		v := order[i]
		for _, to := range g[v] {
			if down[to] < down[v] {
				down[v] = down[to]
			}
		}
	}

	up := make([]int, n)
	for i := 0; i < n; i++ {
		up[i] = i
	}
	for _, v := range order {
		for _, from := range rg[v] {
			if up[from] < up[v] {
				up[v] = up[from]
			}
		}
	}

	res := make([]byte, n)
	cnt := 0
	for i := 0; i < n; i++ {
		if down[i] >= i && up[i] >= i {
			res[i] = 'A'
			cnt++
		} else {
			res[i] = 'E'
		}
	}
	fmt.Fprintln(writer, cnt)
	fmt.Fprintln(writer, string(res))
}
