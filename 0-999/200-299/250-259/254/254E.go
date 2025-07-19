package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

type pair struct {
	fi, idx int
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var n, v int
	fmt.Fscan(reader, &n, &v)
	a := make([]int, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(reader, &a[i])
	}
	var m int
	fmt.Fscan(reader, &m)
	f := make([][]pair, n)
	for i := 0; i < m; i++ {
		var l, r, fi int
		fmt.Fscan(reader, &l, &r, &fi)
		l--
		r--
		for j := l; j <= r; j++ {
			f[j] = append(f[j], pair{fi, i + 1})
		}
	}
	for i := 0; i < n; i++ {
		sort.Slice(f[i], func(p, q int) bool { return f[i][p].fi < f[i][q].fi })
	}

	// dp[i][j]: max count using first i, with carry j
	dp := make([][]int, n+1)
	tp := make([][]int, n+1)
	for i := 0; i <= n; i++ {
		dp[i] = make([]int, 401)
		tp[i] = make([]int, 401)
		for j := 0; j <= 400; j++ {
			dp[i][j] = -1
		}
	}
	dp[0][0] = 0

	for i := 0; i < n; i++ {
		for j := 0; j <= 400; j++ {
			if dp[i][j] < 0 {
				continue
			}
			x := j + a[i] - v
			if x < 0 {
				continue
			}
			// no picks
			nxt := min(a[i], x)
			if dp[i+1][nxt] < dp[i][j] {
				dp[i+1][nxt] = dp[i][j]
				tp[i+1][nxt] = j
			}
			// picks
			t := x
			for k := 0; k < len(f[i]); k++ {
				t -= f[i][k].fi
				if t < 0 {
					break
				}
				nxt2 := min(a[i], t)
				if dp[i+1][nxt2] < dp[i][j]+k+1 {
					dp[i+1][nxt2] = dp[i][j] + k + 1
					tp[i+1][nxt2] = j
				}
			}
		}
	}

	// find best
	ans, fj := 0, 0
	for j := 0; j <= 400; j++ {
		if dp[n][j] > ans {
			ans = dp[n][j]
			fj = j
		}
	}
	fmt.Fprintln(writer, ans)
	// backtrack
	nf := make([]int, 0, n)
	for i := n; i > 0; i-- {
		tf := tp[i][fj]
		nf = append(nf, dp[i][fj]-dp[i-1][tf])
		fj = tf
	}
	// reverse
	for i, j := 0, len(nf)-1; i < j; i, j = i+1, j-1 {
		nf[i], nf[j] = nf[j], nf[i]
	}
	// output picks
	for i := 0; i < n; i++ {
		fmt.Fprint(writer, nf[i])
		for j := 0; j < nf[i]; j++ {
			fmt.Fprint(writer, " ", f[i][j].idx)
		}
		fmt.Fprintln(writer)
	}
}
