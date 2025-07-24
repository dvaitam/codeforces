package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	in := bufio.NewReader(os.Stdin)
	var n int
	if _, err := fmt.Fscan(in, &n); err != nil {
		return
	}
	a := make([]int, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &a[i])
	}

	// DP arrays
	dp0 := make([][]int, n)
	dp1 := make([][]int, n)
	for i := 0; i < n; i++ {
		dp0[i] = make([]int, n)
		dp1[i] = make([]int, n)
	}

	for i := 0; i < n; i++ {
		dp0[i][i] = a[i]
		dp1[i][i] = a[i]
	}

	for len := 2; len <= n; len++ {
		for l := 0; l+len-1 < n; l++ {
			r := l + len - 1
			dp0[l][r] = max(dp1[l+1][r], dp1[l][r-1])
			dp1[l][r] = min(dp0[l+1][r], dp0[l][r-1])
		}
	}

	out := make([]int, n)
	for k := 0; k < n; k++ {
		best := 0
		for x := 0; x <= k; x++ {
			l := x
			r := n - 1 - (k - x)
			if l <= r {
				v := dp0[l][r]
				if v > best {
					best = v
				}
			}
		}
		out[k] = best
	}

	for i := 0; i < n; i++ {
		if i > 0 {
			fmt.Print(" ")
		}
		fmt.Print(out[i])
	}
	fmt.Println()
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
