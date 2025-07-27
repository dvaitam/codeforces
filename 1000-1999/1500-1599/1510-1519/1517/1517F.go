package main

import (
	"bufio"
	"fmt"
	"os"
)

const mod = 998244353

func main() {
	in := bufio.NewReader(os.Stdin)
	var n int
	fmt.Fscan(in, &n)
	g := make([][]int, n)
	for i := 0; i < n-1; i++ {
		var a, b int
		fmt.Fscan(in, &a, &b)
		a--
		b--
		g[a] = append(g[a], b)
		g[b] = append(g[b], a)
	}
	// precompute all-pairs shortest paths using BFS from each node
	dist := make([][]int, n)
	for i := 0; i < n; i++ {
		dist[i] = bfs(i, g, n)
	}

	total := 0
	maxMask := 1 << n
	for mask := 0; mask < maxMask; mask++ {
		if mask == 0 {
			total = (total - 1) % mod
			continue
		}
		if mask == maxMask-1 {
			total = (total + n) % mod
			continue
		}
		level := -1
		for v := 0; v < n; v++ {
			r := 0
			for {
				ok := true
				for u := 0; u < n; u++ {
					if dist[v][u] <= r && (mask&(1<<u)) == 0 {
						ok = false
						break
					}
				}
				if !ok {
					break
				}
				r++
			}
			if r-1 > level {
				level = r - 1
			}
		}
		total = (total + level) % mod
	}

	inv2 := powmod(2, mod-1-n)
	ans := total % mod
	if ans < 0 {
		ans += mod
	}
	ans = ans * inv2 % mod
	fmt.Println(ans)
}

func bfs(start int, g [][]int, n int) []int {
	d := make([]int, n)
	for i := range d {
		d[i] = -1
	}
	q := []int{start}
	d[start] = 0
	for len(q) > 0 {
		v := q[0]
		q = q[1:]
		for _, to := range g[v] {
			if d[to] == -1 {
				d[to] = d[v] + 1
				q = append(q, to)
			}
		}
	}
	return d
}

func powmod(a, b int) int {
	res := 1
	base := a % mod
	for b > 0 {
		if b&1 != 0 {
			res = res * base % mod
		}
		base = base * base % mod
		b >>= 1
	}
	return res
}
