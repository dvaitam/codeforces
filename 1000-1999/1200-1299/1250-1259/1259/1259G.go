package main

import (
	"bufio"
	"fmt"
	"os"
)

const mod int64 = 998244353

func dfs(v, p int, g [][]int) (int64, int64) {
	type pair struct{ a0, a1 int64 }
	children := make([]pair, 0, len(g[v]))
	for _, to := range g[v] {
		if to == p {
			continue
		}
		a0, a1 := dfs(to, v, g)
		children = append(children, pair{a0, a1})
	}
	m := len(children)
	pref := make([]int64, m+1)
	pref[0] = 1
	for i := 0; i < m; i++ {
		val := (children[i].a0 + children[i].a1) % mod
		pref[i+1] = pref[i] * val % mod
	}
	suff := make([]int64, m+1)
	suff[m] = 1
	for i := m - 1; i >= 0; i-- {
		val := (children[i].a0 + children[i].a1) % mod
		suff[i] = suff[i+1] * val % mod
	}
	unmatched := pref[m] % mod
	matched := int64(0)
	for i := 0; i < m; i++ {
		prod := pref[i] * suff[i+1] % mod
		prod = prod * children[i].a0 % mod
		matched = (matched + prod) % mod
	}
	return unmatched, matched
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var n int
	fmt.Fscan(reader, &n)
	g := make([][]int, n)
	for i := 0; i < n-1; i++ {
		var u, v int
		fmt.Fscan(reader, &u, &v)
		u--
		v--
		g[u] = append(g[u], v)
		g[v] = append(g[v], u)
	}

	a0, a1 := dfs(0, -1, g)
	ans := (a0 + a1) % mod
	fmt.Fprintln(writer, ans)
}
