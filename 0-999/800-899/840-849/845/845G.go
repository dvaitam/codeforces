package main

import (
	"bufio"
	"fmt"
	"os"
)

type Edge struct {
	to int
	w  int64
}

const maxBits = 60

func insertBasis(basis []int64, x int64) {
	for i := maxBits - 1; i >= 0; i-- {
		if (x>>i)&1 == 0 {
			continue
		}
		if basis[i] != 0 {
			x ^= basis[i]
		} else {
			basis[i] = x
			break
		}
	}
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var n, m int
	if _, err := fmt.Fscan(reader, &n, &m); err != nil {
		return
	}

	g := make([][]Edge, n+1)
	for i := 0; i < m; i++ {
		var x, y int
		var w int64
		fmt.Fscan(reader, &x, &y, &w)
		g[x] = append(g[x], Edge{y, w})
		g[y] = append(g[y], Edge{x, w})
	}

	dist := make([]int64, n+1)
	vis := make([]bool, n+1)
	parent := make([]int, n+1)
	queue := make([]int, 0, n)

	basis := make([]int64, maxBits)

	vis[1] = true
	queue = append(queue, 1)

	for len(queue) > 0 {
		u := queue[0]
		queue = queue[1:]
		for _, e := range g[u] {
			v := e.to
			if !vis[v] {
				vis[v] = true
				parent[v] = u
				dist[v] = dist[u] ^ e.w
				queue = append(queue, v)
			} else if parent[u] != v {
				val := dist[u] ^ dist[v] ^ e.w
				insertBasis(basis, val)
			}
		}
	}

	x := dist[n]
	for i := maxBits - 1; i >= 0; i-- {
		if (x ^ basis[i]) < x {
			x ^= basis[i]
		}
	}

	fmt.Fprintln(writer, x)
}
