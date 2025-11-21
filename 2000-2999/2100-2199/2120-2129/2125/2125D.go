package main

import (
	"bufio"
	"fmt"
	"os"
)

const mod = 998244353

type edge struct {
	to int
	w  int
}

func modPow(a, e int) int {
	result := 1
	base := a % mod
	for e > 0 {
		if e&1 == 1 {
			result = int(int64(result) * int64(base) % mod)
		}
		base = int(int64(base) * int64(base) % mod)
		e >>= 1
	}
	return result
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n, m int
	if _, err := fmt.Fscan(in, &n, &m); err != nil {
		return
	}

	adj := make([][]edge, m+2)
	base := 1
	for i := 0; i < n; i++ {
		var l, r, p, q int
		fmt.Fscan(in, &l, &r, &p, &q)
		prob := int(int64(p) * int64(modPow(q, mod-2)) % mod)
		oneMinus := (1 - prob + mod) % mod
		base = int(int64(base) * int64(oneMinus) % mod)
		weight := int(int64(prob) * int64(modPow(oneMinus, mod-2)) % mod)
		adj[l] = append(adj[l], edge{to: r + 1, w: weight})
	}

	ways := make([]int, m+2)
	ways[1] = 1
	for pos := 1; pos <= m; pos++ {
		cur := ways[pos]
		if cur == 0 {
			continue
		}
		for _, e := range adj[pos] {
			ways[e.to] = (ways[e.to] + int(int64(cur)*int64(e.w)%mod)) % mod
		}
	}

	ans := int(int64(base) * int64(ways[m+1]) % mod)
	fmt.Fprintln(out, ans)
}
