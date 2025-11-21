package main

import (
	"bufio"
	"fmt"
	"os"
)

const MOD int64 = 1_000_000_007

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

		size := n
		dp := make([][]int64, size+1)
		for i := range dp {
			dp[i] = make([]int64, size+1)
		}
		dp[0][0] = 1

		for _, v := range a {
			ndp := make([][]int64, size+1)
			for i := range ndp {
				ndp[i] = make([]int64, size+1)
			}
			for x := 0; x <= size; x++ {
				for y := x; y <= size; y++ {
					cur := dp[x][y]
					if cur == 0 {
						continue
					}
					// Skip current element
					ndp[x][y] = (ndp[x][y] + cur) % MOD
					// Include current element
					if v >= x {
						nx, ny := v, y
						if nx > ny {
							nx, ny = ny, nx
						}
						ndp[nx][ny] = (ndp[nx][ny] + cur) % MOD
					}
				}
			}
			dp = ndp
		}

		var ans int64
		for x := 0; x <= size; x++ {
			for y := x; y <= size; y++ {
				ans = (ans + dp[x][y]) % MOD
			}
		}
		fmt.Fprintln(out, ans)
	}
}
