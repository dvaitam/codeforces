package main

import (
	"bufio"
	"fmt"
	"os"
)

const inf int64 = 1 << 60

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n int
	if _, err := fmt.Fscan(in, &n); err != nil {
		return
	}

	c := make([][]int64, n+1)
	for i := range c {
		c[i] = make([]int64, n+1)
	}
	for i := 1; i <= n; i++ {
		for j := 1; j <= n; j++ {
			fmt.Fscan(in, &c[i][j])
		}
	}

	pref := make([][]int64, n+1)
	for i := range pref {
		pref[i] = make([]int64, n+1)
	}
	for i := 1; i <= n; i++ {
		rowSum := int64(0)
		for j := 1; j <= n; j++ {
			rowSum += c[i][j]
			pref[i][j] = pref[i-1][j] + rowSum
		}
	}

	sumRect := func(x1, x2, y1, y2 int) int64 {
		if x1 > x2 || y1 > y2 {
			return 0
		}
		return pref[x2][y2] - pref[x1-1][y2] - pref[x2][y1-1] + pref[x1-1][y1-1]
	}

	dp := make([][]int64, n+2)
	root := make([][]int, n+2)
	for i := range dp {
		dp[i] = make([]int64, n+2)
		root[i] = make([]int, n+2)
		for j := range dp[i] {
			dp[i][j] = inf
		}
	}

	for i := 1; i <= n+1; i++ {
		dp[i][i-1] = 0
	}

	for length := 1; length <= n; length++ {
		for l := 1; l+length-1 <= n; l++ {
			r := l + length - 1
			dp[l][r] = inf
			for k := l; k <= r; k++ {
				left := dp[l][k-1]
				right := dp[k+1][r]
				edgeLeft := sumRect(l, k-1, k, r)
				edgeRight := sumRect(l, k, k+1, r)
				total := left + right + edgeLeft + edgeRight
				if total < dp[l][r] {
					dp[l][r] = total
					root[l][r] = k
				}
			}
		}
	}

	parent := make([]int, n+1)

	var build func(int, int, int)
	build = func(l, r, par int) {
		if l > r {
			return
		}
		k := root[l][r]
		parent[k] = par
		build(l, k-1, k)
		build(k+1, r, k)
	}

	build(1, n, 0)

	for i := 1; i <= n; i++ {
		if i > 1 {
			fmt.Fprint(out, " ")
		}
		fmt.Fprint(out, parent[i])
	}
	fmt.Fprintln(out)
}
