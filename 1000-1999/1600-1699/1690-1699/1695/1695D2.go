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

	var t int
	fmt.Fscan(in, &t)
	for ; t > 0; t-- {
		var n int
		fmt.Fscan(in, &n)
		g := make([][]int, n)
		for i := 0; i < n-1; i++ {
			var x, y int
			fmt.Fscan(in, &x, &y)
			x--
			y--
			g[x] = append(g[x], y)
			g[y] = append(g[y], x)
		}
		if n == 1 {
			fmt.Fprintln(out, 0)
			continue
		}
		deg := make([]int, n)
		maxDeg := 0
		for i := 0; i < n; i++ {
			deg[i] = len(g[i])
			if deg[i] > maxDeg {
				maxDeg = deg[i]
			}
		}
		if maxDeg <= 2 {
			fmt.Fprintln(out, 1)
			continue
		}
		majors := make([]int, 0)
		for i := 0; i < n; i++ {
			if deg[i] >= 3 {
				majors = append(majors, i)
			}
		}
		if len(majors) == 0 {
			fmt.Fprintln(out, 1)
			continue
		}
		dist := make([]int, n)
		parent := make([]int, n)
		for i := range dist {
			dist[i] = -1
		}
		q := append([]int(nil), majors...)
		for _, v := range majors {
			dist[v] = 0
		}
		for head := 0; head < len(q); head++ {
			u := q[head]
			for _, w := range g[u] {
				if dist[w] == -1 {
					dist[w] = dist[u] + 1
					parent[w] = u
					q = append(q, w)
				}
			}
		}
		cnt := make([]int, n)
		for i := 0; i < n; i++ {
			if deg[i] == 1 {
				v := i
				for dist[v] != 0 {
					v = parent[v]
				}
				cnt[v]++
			}
		}
		ans := 0
		for _, v := range majors {
			if cnt[v] > 0 {
				ans += cnt[v] - 1
			}
		}
		fmt.Fprintln(out, ans)
	}
}
