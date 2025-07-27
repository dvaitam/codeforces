package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

func bfs(start int, adj [][]int) []int {
	n := len(adj) - 1
	dist := make([]int, n+1)
	for i := 1; i <= n; i++ {
		dist[i] = -1
	}
	q := make([]int, 0, n)
	dist[start] = 0
	q = append(q, start)
	for head := 0; head < len(q); head++ {
		u := q[head]
		for _, v := range adj[u] {
			if dist[v] == -1 {
				dist[v] = dist[u] + 1
				q = append(q, v)
			}
		}
	}
	return dist
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var T int
	if _, err := fmt.Fscan(in, &T); err != nil {
		return
	}
	for ; T > 0; T-- {
		var n, m int
		var a, b, c int
		fmt.Fscan(in, &n, &m, &a, &b, &c)
		prices := make([]int, m)
		for i := range prices {
			fmt.Fscan(in, &prices[i])
		}
		sort.Ints(prices)
		pref := make([]int64, m+1)
		for i := 1; i <= m; i++ {
			pref[i] = pref[i-1] + int64(prices[i-1])
		}
		adj := make([][]int, n+1)
		for i := 0; i < m; i++ {
			var u, v int
			fmt.Fscan(in, &u, &v)
			adj[u] = append(adj[u], v)
			adj[v] = append(adj[v], u)
		}

		da := bfs(a, adj)
		db := bfs(b, adj)
		dc := bfs(c, adj)

		ans := int64(1<<62 - 1)
		for i := 1; i <= n; i++ {
			d1 := da[i]
			d2 := db[i]
			d3 := dc[i]
			if d1 == -1 || d2 == -1 || d3 == -1 {
				continue
			}
			total := d1 + d2 + d3
			if total > m {
				continue
			}
			cost := pref[d2] + pref[total]
			if cost < ans {
				ans = cost
			}
		}
		fmt.Fprintln(out, ans)
	}
}
