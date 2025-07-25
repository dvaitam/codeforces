package main

import (
	"bufio"
	"fmt"
	"os"
)

const mod int64 = 998244353

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n int
	fmt.Fscan(in, &n)
	a := make([]int, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &a[i])
	}

	inc := make([]int, n)
	if n > 0 {
		inc[n-1] = 1
		for i := n - 2; i >= 0; i-- {
			if a[i] < a[i+1] {
				inc[i] = inc[i+1] + 1
			} else {
				inc[i] = 1
			}
		}
	}

	maxExp := n * (n - 1) / 2
	pow2 := make([]int64, maxExp+1)
	pow2[0] = 1
	for i := 1; i <= maxExp; i++ {
		pow2[i] = pow2[i-1] * 2 % mod
	}

	dp := make([]int64, n+1)
	dp[1] = 1
	for i := 0; i < n; i++ {
		newDp := make([]int64, n+1)
		for q := 1; q <= n; q++ {
			if dp[q] == 0 {
				continue
			}
			pos := i + q
			if pos > n {
				continue
			}
			maxlen := 0
			if pos < n {
				maxlen = inc[pos]
			}
			minK := 0
			if q == 1 && pos < n {
				minK = 1
			}
			for k := minK; k <= maxlen; k++ {
				qNew := q - 1 + k
				exp := k*(pos-i-1) + k*(k-1)/2
				newDp[qNew] = (newDp[qNew] + dp[q]*pow2[exp]) % mod
			}
		}
		dp = newDp
	}

	fmt.Fprintln(out, dp[0])
}
