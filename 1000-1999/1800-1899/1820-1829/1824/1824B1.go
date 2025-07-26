package main

import (
	"bufio"
	"fmt"
	"os"
)

const MOD int64 = 1_000_000_007

var (
	n       int
	k       int
	adj     [][]int
	sz      []int
	sumDist int64
)

func dfs(v, p int) {
	sz[v] = 1
	for _, to := range adj[v] {
		if to == p {
			continue
		}
		dfs(to, v)
		s := sz[to]
		sumDist += int64(s) * int64(n-s)
		sz[v] += s
	}
}

func powmod(a, e int64) int64 {
	a %= MOD
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
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	if _, err := fmt.Fscan(in, &n, &k); err != nil {
		return
	}
	adj = make([][]int, n+1)
	for i := 0; i < n-1; i++ {
		var u, v int
		fmt.Fscan(in, &u, &v)
		adj[u] = append(adj[u], v)
		adj[v] = append(adj[v], u)
	}

	if k == 1 || k == 3 {
		fmt.Fprintln(out, 1)
		return
	}

	sz = make([]int, n+1)
	sumDist = 0
	dfs(1, 0)

	numerator := (2 * (sumDist % MOD)) % MOD
	denom := int64(n) * int64(n-1) % MOD
	invDenom := powmod(denom, MOD-2)
	ans := (1 + numerator*invDenom%MOD) % MOD
	fmt.Fprintln(out, ans)
}
