package main

import (
	"bufio"
	"fmt"
	"os"
)

const mod = 998244353

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var T int
	fmt.Fscan(in, &T)
	for ; T > 0; T-- {
		var n, q int
		fmt.Fscan(in, &n, &q)

		adj := make([][]int, n+1)
		for i := 0; i < n-1; i++ {
			var u, v int
			fmt.Fscan(in, &u, &v)
			adj[u] = append(adj[u], v)
			adj[v] = append(adj[v], u)
		}

		queries := make([][2]int, q)
		maxP := 0
		for i := 0; i < q; i++ {
			fmt.Fscan(in, &queries[i][0], &queries[i][1])
			if queries[i][1] > maxP {
				maxP = queries[i][1]
			}
		}

		parent := make([]int, n+1)
		depth := make([]int, n+1)
		order := make([]int, 0, n)
		stack := []int{1}
		parent[1] = 0
		depth[1] = 0
		for len(stack) > 0 {
			v := stack[len(stack)-1]
			stack = stack[:len(stack)-1]
			order = append(order, v)
			for _, to := range adj[v] {
				if to == parent[v] {
					continue
				}
				parent[to] = v
				depth[to] = depth[v] + 1
				stack = append(stack, to)
			}
		}

		f := make([][]int64, n+1)
		g := make([][]int64, n+1)
		for i := 1; i <= n; i++ {
			f[i] = make([]int64, maxP+1)
			g[i] = make([]int64, maxP+1)
		}

		// root values already zero
		for _, u := range order {
			if u == 1 {
				continue
			}
			par := parent[u]
			penalty := int64(2 * (len(adj[u]) - 1))
			for p := 0; p <= maxP; p++ {
				f[u][p] = 1 + g[par][p]
			}
			for p := 0; p <= maxP; p++ {
				noPay := 1 + f[par][p] + penalty
				pay := int64(1 << 60)
				if p > 0 {
					pay = 1 + f[par][p-1]
				}
				if pay < noPay {
					g[u][p] = pay
				} else {
					g[u][p] = noPay
				}
			}
		}

		for _, qu := range queries {
			v, p := qu[0], qu[1]
			if p > maxP {
				p = maxP
			}
			ans := f[v][p] % mod
			fmt.Fprintln(out, ans)
		}
	}
}
