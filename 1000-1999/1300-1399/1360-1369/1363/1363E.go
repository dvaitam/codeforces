package main

import (
	"bufio"
	"fmt"
	"os"
)

const N = 200005

var (
	n   int
	a   [N]int64
	b   [N]int
	c   [N]int
	g   [N][]int
	ans int64
)

func min64(x, y int64) int64 {
	if x < y {
		return x
	}
	return y
}

func dfs(v, p int, mn int64) (int, int) {
	cnt01, cnt10 := 0, 0
	if b[v] != c[v] {
		if b[v] == 0 {
			cnt01 = 1
		} else {
			cnt10 = 1
		}
	}
	for _, to := range g[v] {
		if to == p {
			continue
		}
		x01, x10 := dfs(to, v, min64(mn, a[to]))
		cnt01 += x01
		cnt10 += x10
	}
	t := int(min64(int64(cnt01), int64(cnt10)))
	ans += 2 * int64(t) * mn
	cnt01 -= t
	cnt10 -= t
	return cnt01, cnt10
}

func main() {
	in := bufio.NewReader(os.Stdin)
	fmt.Fscan(in, &n)
	for i := 1; i <= n; i++ {
		fmt.Fscan(in, &a[i], &b[i], &c[i])
		g[i] = g[i][:0]
	}
	for i := 0; i < n-1; i++ {
		var u, v int
		fmt.Fscan(in, &u, &v)
		g[u] = append(g[u], v)
		g[v] = append(g[v], u)
	}
	dfs01, dfs10 := dfs(1, 0, a[1])
	if dfs01 != 0 || dfs10 != 0 {
		fmt.Println(-1)
	} else {
		fmt.Println(ans)
	}
}
