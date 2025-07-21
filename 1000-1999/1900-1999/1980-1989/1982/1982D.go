package main

import (
	"bufio"
	"fmt"
	"os"
)

func gcd(a, b int64) int64 {
	for b != 0 {
		a, b = b, a%b
	}
	if a < 0 {
		return -a
	}
	return a
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var t int
	if _, err := fmt.Fscan(in, &t); err != nil {
		return
	}
	for ; t > 0; t-- {
		var n, m, k int
		fmt.Fscan(in, &n, &m, &k)
		heights := make([][]int64, n)
		for i := 0; i < n; i++ {
			heights[i] = make([]int64, m)
			for j := 0; j < m; j++ {
				fmt.Fscan(in, &heights[i][j])
			}
		}
		types := make([]string, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(in, &types[i])
		}
		ones := make([][]int, n)
		var sum0, sum1 int64
		for i := 0; i < n; i++ {
			ones[i] = make([]int, m)
			for j := 0; j < m; j++ {
				if types[i][j] == '1' {
					ones[i][j] = 1
					sum1 += heights[i][j]
				} else {
					sum0 += heights[i][j]
				}
			}
		}
		pref := make([][]int, n+1)
		for i := 0; i <= n; i++ {
			pref[i] = make([]int, m+1)
		}
		for i := 0; i < n; i++ {
			rowSum := 0
			for j := 0; j < m; j++ {
				rowSum += ones[i][j]
				pref[i+1][j+1] = pref[i][j+1] + rowSum
			}
		}
		var g int64
		for i := 0; i+k <= n; i++ {
			for j := 0; j+k <= m; j++ {
				onesCount := pref[i+k][j+k] - pref[i][j+k] - pref[i+k][j] + pref[i][j]
				diff := int64(k*k - 2*onesCount)
				if diff < 0 {
					diff = -diff
				}
				g = gcd(g, diff)
			}
		}
		diffTotal := sum1 - sum0
		if g == 0 {
			if diffTotal == 0 {
				fmt.Fprintln(out, "YES")
			} else {
				fmt.Fprintln(out, "NO")
			}
		} else {
			if diffTotal%g == 0 {
				fmt.Fprintln(out, "YES")
			} else {
				fmt.Fprintln(out, "NO")
			}
		}
	}
}
