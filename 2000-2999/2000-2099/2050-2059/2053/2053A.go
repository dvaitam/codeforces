package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var t int
	fmt.Fscan(in, &t)

	for ; t > 0; t-- {
		var n int
		fmt.Fscan(in, &n)

		a := make([]int, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(in, &a[i])
		}

		stable := make([][]bool, n)
		for i := 0; i < n; i++ {
			stable[i] = make([]bool, n)
			minVal, maxVal := a[i], a[i]
			for j := i; j < n; j++ {
				if a[j] < minVal {
					minVal = a[j]
				}
				if a[j] > maxVal {
					maxVal = a[j]
				}
				if 2*minVal > maxVal {
					stable[i][j] = true
				}
			}
		}

		dp := make([]int, n+1)
		dp[0] = 1
		for i := 0; i < n; i++ {
			if dp[i] == 0 {
				continue
			}
			for j := i; j < n; j++ {
				if stable[i][j] {
					dp[j+1] += dp[i]
					if dp[j+1] > 2 {
						dp[j+1] = 2
					}
				}
			}
		}

		if dp[n] >= 2 {
			fmt.Fprintln(out, "YES")
		} else {
			fmt.Fprintln(out, "NO")
		}
	}
}
