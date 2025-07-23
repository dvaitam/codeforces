package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

const inf int64 = 1 << 60

func min64(a, b int64) int64 {
	if a < b {
		return a
	}
	return b
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n, k int
	if _, err := fmt.Fscan(in, &n, &k); err != nil {
		return
	}
	arr := make([]int, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &arr[i])
	}
	sort.Ints(arr)
	// convert to 1-indexed int64 slice
	a := make([]int64, n+1)
	for i, v := range arr {
		a[i+1] = int64(v)
	}

	base := n / k
	extra := n % k

	dp := make([][]int64, k+1)
	for i := range dp {
		dp[i] = make([]int64, extra+1)
		for j := range dp[i] {
			dp[i][j] = inf
		}
	}
	dp[0][0] = 0

	for t := 0; t < k; t++ {
		maxj := t
		if maxj > extra {
			maxj = extra
		}
		for j := 0; j <= maxj; j++ {
			cur := dp[t][j]
			if cur == inf {
				continue
			}
			idx := t*base + j
			if t-j < k-extra { // place group of size base
				val := cur + a[idx+base] - a[idx+1]
				if val < dp[t+1][j] {
					dp[t+1][j] = val
				}
			}
			if j < extra { // place group of size base+1
				val := cur + a[idx+base+1] - a[idx+1]
				if val < dp[t+1][j+1] {
					dp[t+1][j+1] = val
				}
			}
		}
	}
	fmt.Fprintln(out, dp[k][extra])
}
