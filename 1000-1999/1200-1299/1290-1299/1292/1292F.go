package main

import (
	"bufio"
	"fmt"
	"os"
)

const MOD int = 1000000007

type Component struct {
	roots  []int
	others []int
}

func cntOrder(s []int, t []int) int {
	p := len(s)
	m := len(t)
	inMask := make([]int, m)
	for x := 0; x < p; x++ {
		for i := 0; i < m; i++ {
			if t[i]%s[x] == 0 {
				inMask[i] |= 1 << x
			}
		}
	}
	cnt := make([]int, 1<<p)
	for mask := 0; mask < (1 << p); mask++ {
		for i := 0; i < m; i++ {
			if inMask[i]&mask == inMask[i] {
				cnt[mask]++
			}
		}
	}
	dp := make([][]int, m+1)
	for i := range dp {
		dp[i] = make([]int, 1<<p)
	}
	for i := 0; i < m; i++ {
		dp[1][inMask[i]] = (dp[1][inMask[i]] + 1) % MOD
	}
	for k := 1; k < m; k++ {
		for mask := 0; mask < (1 << p); mask++ {
			val := dp[k][mask]
			if val == 0 {
				continue
			}
			for i := 0; i < m; i++ {
				if inMask[i]&mask != 0 && inMask[i]&^mask != 0 {
					dp[k+1][mask|inMask[i]] = (dp[k+1][mask|inMask[i]] + val) % MOD
				}
			}
			remain := cnt[mask] - k
			if remain > 0 {
				dp[k+1][mask] = (dp[k+1][mask] + val*remain) % MOD
			}
		}
	}
	return dp[m][(1<<p)-1]
}

func main() {
	in := bufio.NewReader(os.Stdin)
	var n int
	fmt.Fscan(in, &n)
	a := make([]int, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &a[i])
	}

	C := make([][]int, n+1)
	for i := 0; i <= n; i++ {
		C[i] = make([]int, n+1)
		C[i][0] = 1
		for j := 1; j <= i; j++ {
			C[i][j] = (C[i-1][j-1] + C[i-1][j]) % MOD
		}
	}

	graph := make([][]int, n)
	degIn := make([]int, n)
	for u := 0; u < n; u++ {
		for v := 0; v < n; v++ {
			if u != v && a[v]%a[u] == 0 {
				graph[u] = append(graph[u], v)
				graph[v] = append(graph[v], u)
				degIn[v]++
			}
		}
	}

	visited := make([]bool, n)
	var ans int = 1
	curLen := 0
	var dfs func(int, *Component)
	dfs = func(u int, comp *Component) {
		visited[u] = true
		if degIn[u] == 0 {
			comp.roots = append(comp.roots, a[u])
		} else {
			comp.others = append(comp.others, a[u])
		}
		for _, v := range graph[u] {
			if !visited[v] {
				dfs(v, comp)
			}
		}
	}

	for u := 0; u < n; u++ {
		if !visited[u] {
			comp := &Component{}
			dfs(u, comp)
			if len(comp.others) > 0 {
				sz := len(comp.others) - 1
				cnt := cntOrder(comp.roots, comp.others)
				ans = (ans * cnt) % MOD
				ans = (ans * C[curLen+sz][sz]) % MOD
				curLen += sz
			}
		}
	}
	fmt.Println(ans)
}
