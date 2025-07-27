package main

import (
	"bufio"
	"fmt"
	"os"
)

type edge struct {
	to   int
	cost int64
}

const inf int64 = 1 << 60

func main() {
	in := bufio.NewReader(os.Stdin)
	var n, m int
	fmt.Fscan(in, &n, &m)
	g := make([][]edge, n)
	for i := 0; i < m; i++ {
		var a, b int
		var c int64
		fmt.Fscan(in, &a, &b, &c)
		g[a] = append(g[a], edge{to: b, cost: c})
	}

	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	for s := 0; s < n; s++ {
		dist := make([]int64, n)
		used := make([]bool, n)
		for i := range dist {
			dist[i] = inf
		}
		dist[s] = 0
		for iter := 0; iter < n; iter++ {
			v := -1
			for i := 0; i < n; i++ {
				if !used[i] && (v == -1 || dist[i] < dist[v]) {
					v = i
				}
			}
			used[v] = true
			base := int(dist[v] % int64(n))
			tmp := make([]int64, n)
			for i := range tmp {
				tmp[i] = inf
			}
			for _, e := range g[v] {
				idx := (e.to + base) % n
				if val := dist[v] + e.cost; val < tmp[idx] {
					tmp[idx] = val
				}
			}
			for i := 0; i < 2*n; i++ {
				j := (i + 1) % n
				if tmp[i%n]+1 < tmp[j] {
					tmp[j] = tmp[i%n] + 1
				}
			}
			for i := 0; i < n; i++ {
				if tmp[i] < dist[i] {
					dist[i] = tmp[i]
				}
			}
		}
		for i := 0; i < n; i++ {
			if i > 0 {
				fmt.Fprint(out, " ")
			}
			fmt.Fprint(out, dist[i])
		}
		fmt.Fprintln(out)
	}
}
