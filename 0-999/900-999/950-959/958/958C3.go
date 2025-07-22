package main

import (
	"bufio"
	"fmt"
	"os"
)

const INF = int(1e9)

func update(dp, pre, suf [][]int, j, r, val, p int) {
	if val >= dp[j][r] {
		return
	}
	dp[j][r] = val
	if val < pre[j][r] {
		pre[j][r] = val
		for t := r + 1; t < p && pre[j][t] > val; t++ {
			pre[j][t] = val
		}
	}
	if val < suf[j][r] {
		suf[j][r] = val
		for t := r - 1; t >= 0 && suf[j][t] > val; t-- {
			suf[j][t] = val
		}
	}
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n, k, p int
	if _, err := fmt.Fscan(in, &n, &k, &p); err != nil {
		return
	}

	a := make([]int, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &a[i])
	}

	dp := make([][]int, k+1)
	pre := make([][]int, k+1)
	suf := make([][]int, k+1)
	for i := 0; i <= k; i++ {
		dp[i] = make([]int, p)
		pre[i] = make([]int, p)
		suf[i] = make([]int, p)
		for j := 0; j < p; j++ {
			dp[i][j] = INF
			pre[i][j] = INF
			suf[i][j] = INF
		}
	}
	dp[0][0] = 0
	// initialize pre and suf for j=0
	minv := INF
	for r := 0; r < p; r++ {
		if dp[0][r] < minv {
			minv = dp[0][r]
		}
		pre[0][r] = minv
	}
	minv = INF
	for r := p - 1; r >= 0; r-- {
		if dp[0][r] < minv {
			minv = dp[0][r]
		}
		suf[0][r] = minv
	}

	rem := 0
	total := 0
	for idx, val := range a {
		total += val
		rem = (rem + val) % p
		maxj := k
		if idx+1 < k {
			maxj = idx + 1
		}
		for j := maxj; j >= 1; j-- {
			best := pre[j-1][rem]
			if rem+1 < p {
				alt := suf[j-1][rem+1]
				if alt+1 < best {
					best = alt + 1
				}
			}
			update(dp, pre, suf, j, rem, best, p)
		}
	}
	ans := total%p + p*dp[k][rem]
	fmt.Fprintln(out, ans)
}
