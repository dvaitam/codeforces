package main

import (
	"bufio"
	"fmt"
	"os"
)

const MOD int64 = 1_000_000_007

func modPow(a, e int64) int64 {
	res := int64(1)
	for e > 0 {
		if e&1 == 1 {
			res = res * a % MOD
		}
		a = a * a % MOD
		e >>= 1
	}
	return res
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	var n, k int
	if _, err := fmt.Fscan(reader, &n, &k); err != nil {
		return
	}
	adj := make([]map[int]int64, n+1)
	deg := make([]int64, n+1)
	for i := 1; i <= n; i++ {
		adj[i] = make(map[int]int64)
	}
	// initial k-clique
	for i := 1; i <= k; i++ {
		for j := 1; j < i; j++ {
			adj[i][j]++
			adj[j][i]++
			deg[i]++
			deg[j]++
		}
	}
	// read vertices k+1..n
	for v := k + 1; v <= n; v++ {
		for t := 0; t < k; t++ {
			var u int
			fmt.Fscan(reader, &u)
			adj[v][u]++
			adj[u][v]++
			deg[v]++
			deg[u]++
		}
	}

	ans := int64(1)
	// eliminate vertices from n down to 2 (1 is removed for Matrix-Tree)
	for v := n; v >= 2; v-- {
		dv := deg[v] % MOD
		ans = ans * dv % MOD
		inv := modPow(dv, MOD-2)
		// gather neighbors
		neigh := make([]int, 0, len(adj[v]))
		for u := range adj[v] {
			neigh = append(neigh, u)
		}
		// update edges between neighbors
		for i := 0; i < len(neigh); i++ {
			u := neigh[i]
			wvu := adj[v][u] % MOD
			for j := i + 1; j < len(neigh); j++ {
				t := neigh[j]
				wvt := adj[v][t] % MOD
				val := wvu * wvt % MOD * inv % MOD
				adj[u][t] = (adj[u][t] + val) % MOD
				adj[t][u] = adj[u][t]
			}
		}
		// update degrees and remove edges to v
		for _, u := range neigh {
			w := adj[v][u] % MOD
			val := w * w % MOD * inv % MOD
			deg[u] = (deg[u] - val) % MOD
			if deg[u] < 0 {
				deg[u] += MOD
			}
			delete(adj[u], v)
		}
		adj[v] = nil
	}
	if ans < 0 {
		ans += MOD
	}
	fmt.Println(ans % MOD)
}
