package main

import (
	"bufio"
	"fmt"
	"os"
)

const N = 2005

var (
	a, b, up, dn [N][N]bool
	n, m         int
	ans          int64
)

func solve(grid *[N][N]bool, n, m, w int) {
	for j := 1; j <= m; j++ {
		up[0][j] = false
	}
	for i := 1; i <= n; i++ {
		for j := 1; j <= m; j++ {
			up[i][j] = up[i-1][j] || (*grid)[i][j]
		}
	}
	for i := n; i >= 1; i-- {
		for j := 1; j <= m; j++ {
			if i == n {
				dn[i][j] = (*grid)[i][j]
			} else {
				dn[i][j] = dn[i+1][j] || (*grid)[i][j]
			}
		}
	}
	for i := 2; i < n; i++ {
		v := 0
		if !(*grid)[i][1] && w == 1 {
			v = 1
		}
		up[i][1] = true
		dn[i][1] = true
		for j := 2; j < m; j++ {
			if (*grid)[i][j] {
				v = 0
				continue
			}
			if !up[i][j] {
				ans += int64(v)
			}
			if !dn[i][j] {
				ans += int64(v)
			}
			if !up[i][j] && !dn[i][j-1] {
				ans++
			}
			if !dn[i][j] && !up[i][j-1] {
				ans++
			}
			if !up[i][j-1] {
				v++
			}
			if !dn[i][j-1] {
				v++
			}
		}
		if !(*grid)[i][m] && w == 1 {
			ans += int64(v)
			if !up[i][m-1] {
				ans++
			}
			if !dn[i][m-1] {
				ans++
			}
		}
	}
	if w == 1 {
		for j := 2; j < m; j++ {
			if !up[n][j] {
				ans++
			}
		}
	}
}

func main() {
	in := bufio.NewReader(os.Stdin)
	fmt.Fscan(in, &n, &m)
	for i := 1; i <= n; i++ {
		var s string
		fmt.Fscan(in, &s)
		for j := 1; j <= m; j++ {
			if s[j-1] == '#' {
				a[i][j] = true
				b[j][i] = true
			}
		}
	}
	solve(&a, n, m, 1)
	solve(&b, m, n, 0)
	out := bufio.NewWriter(os.Stdout)
	fmt.Fprintln(out, ans)
	out.Flush()
}
