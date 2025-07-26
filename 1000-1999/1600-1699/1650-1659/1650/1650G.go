package main

import (
	"bufio"
	"fmt"
	"os"
)

const mod int64 = 1000000007

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var T int
	fmt.Fscan(reader, &T)
	for ; T > 0; T-- {
		var n, m int
		fmt.Fscan(reader, &n, &m)
		var s, t int
		fmt.Fscan(reader, &s, &t)
		s--
		t--
		g := make([][]int, n)
		for i := 0; i < m; i++ {
			var u, v int
			fmt.Fscan(reader, &u, &v)
			u--
			v--
			g[u] = append(g[u], v)
			g[v] = append(g[v], u)
		}
		ds := make([]int, n)
		for i := range ds {
			ds[i] = -1
		}
		q := make([]int, 0, n)
		q = append(q, s)
		ds[s] = 0
		for head := 0; head < len(q); head++ {
			v := q[head]
			for _, u := range g[v] {
				if ds[u] == -1 {
					ds[u] = ds[v] + 1
					q = append(q, u)
				}
			}
		}
		maxd := 0
		for _, d := range ds {
			if d > maxd {
				maxd = d
			}
		}
		levels := make([][]int, maxd+1)
		for i, d := range ds {
			if d >= 0 {
				levels[d] = append(levels[d], i)
			}
		}
		dp0 := make([]int64, n)
		dp1 := make([]int64, n)
		dp0[s] = 1
		for L := 0; L <= maxd; L++ {
			// edges within the same level add exactly one extra step
			for _, v := range levels[L] {
				for _, u := range g[v] {
					if ds[u] == L {
						dp1[u] = (dp1[u] + dp0[v]) % mod
					}
				}
			}
			// propagate to next level
			if L < maxd {
				for _, v := range levels[L] {
					for _, u := range g[v] {
						if ds[u] == L+1 {
							dp0[u] = (dp0[u] + dp0[v]) % mod
							dp1[u] = (dp1[u] + dp1[v]) % mod
						}
					}
				}
			}
		}
		ans := (dp0[t] + dp1[t]) % mod
		fmt.Fprintln(writer, ans)
	}
}
