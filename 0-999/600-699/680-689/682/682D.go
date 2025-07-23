package main

import (
	"bufio"
	"fmt"
	"os"
)

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func main() {
	in := bufio.NewReader(os.Stdin)
	var n, m, K int
	if _, err := fmt.Fscan(in, &n, &m, &K); err != nil {
		return
	}
	var s, t string
	fmt.Fscan(in, &s)
	fmt.Fscan(in, &t)

	f := make([][][]int, K+1)
	for k := range f {
		f[k] = make([][]int, n+1)
		for i := range f[k] {
			f[k][i] = make([]int, m+1)
		}
	}
	g := make([][]int, n+1)
	for i := range g {
		g[i] = make([]int, m+1)
	}

	for k := 1; k <= K; k++ {
		for i := 1; i <= n; i++ {
			for j := 1; j <= m; j++ {
				if s[i-1] == t[j-1] {
					v1 := g[i-1][j-1] + 1
					v2 := f[k-1][i-1][j-1] + 1
					if v2 > v1 {
						v1 = v2
					}
					g[i][j] = v1
				} else {
					g[i][j] = 0
				}
				v := f[k][i-1][j]
				if f[k][i][j-1] > v {
					v = f[k][i][j-1]
				}
				if g[i][j] > v {
					v = g[i][j]
				}
				f[k][i][j] = v
			}
		}
	}

	fmt.Println(f[K][n][m])
}
