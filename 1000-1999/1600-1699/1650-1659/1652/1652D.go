package main

import (
	"bufio"
	"fmt"
	"os"
)

const mod int64 = 998244353

var spf []int

func sieve(n int) {
	spf = make([]int, n+1)
	for i := 2; i <= n; i++ {
		if spf[i] == 0 {
			spf[i] = i
			if i*i <= n {
				for j := i * i; j <= n; j += i {
					if spf[j] == 0 {
						spf[j] = i
					}
				}
			}
		}
	}
	for i := 2; i <= n; i++ {
		if spf[i] == 0 {
			spf[i] = i
		}
	}
}

func modPow(a, e int64) int64 {
	res := int64(1)
	a %= mod
	for e > 0 {
		if e&1 != 0 {
			res = res * a % mod
		}
		a = a * a % mod
		e >>= 1
	}
	return res
}

func modInv(a int64) int64 {
	return modPow(a, mod-2)
}

type Edge struct {
	to int
	x  int
	y  int
}

func factorUpdate(x int, sign int, cnt, minCnt map[int]int) {
	for x > 1 {
		p := spf[x]
		c := 0
		for x%p == 0 {
			x /= p
			c++
		}
		cnt[p] += sign * c
		if cnt[p] < minCnt[p] {
			minCnt[p] = cnt[p]
		}
	}
}

func dfs(u, p int, adj [][]Edge, val []int64, inv []int64, cnt, minCnt map[int]int) {
	for _, e := range adj[u] {
		if e.to == p {
			continue
		}
		factorUpdate(e.y, 1, cnt, minCnt)
		factorUpdate(e.x, -1, cnt, minCnt)
		val[e.to] = val[u] * int64(e.y) % mod * inv[e.x] % mod
		dfs(e.to, u, adj, val, inv, cnt, minCnt)
		factorUpdate(e.y, -1, cnt, minCnt)
		factorUpdate(e.x, 1, cnt, minCnt)
	}
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var T int
	fmt.Fscan(in, &T)

	const MAX = 200000
	sieve(MAX)
	inv := make([]int64, MAX+1)
	for i := 1; i <= MAX; i++ {
		inv[i] = modInv(int64(i))
	}

	for ; T > 0; T-- {
		var n int
		fmt.Fscan(in, &n)
		adj := make([][]Edge, n+1)
		for i := 0; i < n-1; i++ {
			var u, v, x, y int
			fmt.Fscan(in, &u, &v, &x, &y)
			adj[u] = append(adj[u], Edge{v, x, y})
			adj[v] = append(adj[v], Edge{u, y, x})
		}
		val := make([]int64, n+1)
		val[1] = 1
		cnt := make(map[int]int)
		minCnt := make(map[int]int)
		dfs(1, 0, adj, val, inv, cnt, minCnt)

		base := int64(1)
		for p, e := range minCnt {
			if e < 0 {
				base = base * modPow(int64(p), int64(-e)) % mod
			}
		}
		ans := int64(0)
		for i := 1; i <= n; i++ {
			ans = (ans + val[i]*base) % mod
		}
		fmt.Fprintln(out, ans)
	}
}
