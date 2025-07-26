package main

import (
	"bufio"
	"fmt"
	"os"
)

const MOD int64 = 998244353

type Edge struct {
	to int
	w  int
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var n int
	if _, err := fmt.Fscan(reader, &n); err != nil {
		return
	}

	g := make([][]Edge, n+1)
	indeg := make([]int, n+1)
	for i := 1; i <= n; i++ {
		var s int
		fmt.Fscan(reader, &s)
		for j := 0; j < s; j++ {
			var v, w int
			fmt.Fscan(reader, &v, &w)
			g[i] = append(g[i], Edge{to: v, w: w})
			indeg[v]++
		}
	}

	// topological order
	order := make([]int, 0, n)
	queue := make([]int, 0)
	for i := 1; i <= n; i++ {
		if indeg[i] == 0 {
			queue = append(queue, i)
		}
	}
	for len(queue) > 0 {
		u := queue[0]
		queue = queue[1:]
		order = append(order, u)
		for _, e := range g[u] {
			indeg[e.to]--
			if indeg[e.to] == 0 {
				queue = append(queue, e.to)
			}
		}
	}

	cnt0 := make([]int64, n+1)
	cnt1 := make([]int64, n+1)
	inv := make([]int64, n+1)

	for i := len(order) - 1; i >= 0; i-- {
		u := order[i]
		var p0, p1, invU int64
		for _, e := range g[u] {
			if e.w == 1 {
				p1 = (p1 + 1) % MOD
				invU = (invU + p1*cnt0[e.to] + inv[e.to]) % MOD
				p1 = (p1 + cnt1[e.to]) % MOD
				p0 = (p0 + cnt0[e.to]) % MOD
			} else {
				invU = (invU + p1*cnt0[e.to] + inv[e.to]) % MOD
				p1 = (p1 + cnt1[e.to]) % MOD
				p0 = (p0 + cnt0[e.to]) % MOD
				invU = (invU + p1) % MOD
				p0 = (p0 + 1) % MOD
			}
		}
		cnt0[u] = p0 % MOD
		cnt1[u] = p1 % MOD
		inv[u] = invU % MOD
	}

	fmt.Fprintln(writer, inv[1]%MOD)
}
