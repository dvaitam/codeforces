package main

import (
	"bufio"
	"fmt"
	"os"
)

const MOD int = 998244353

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var t int
	fmt.Fscan(in, &t)
	for ; t > 0; t-- {
		var n, m int
		fmt.Fscan(in, &n, &m)

		forbidden := make([][]bool, n)
		for i := 0; i < n; i++ {
			forbidden[i] = make([]bool, n+1)
		}
		for i := 0; i < m; i++ {
			var pos, val int
			fmt.Fscan(in, &pos, &val)
			forbidden[pos-1][val] = true
		}

		dp := make([]int, n+1)
		dp[0] = 1
		for pos := 0; pos < n; pos++ {
			for s := 1; s <= n-pos; s++ {
				seq := make([]int, s)
				for i := 0; i < s; i++ {
					seq[i] = i + 1
				}
				for shift := 0; shift < s; shift++ {
					valid := true
					for j := 0; j < s; j++ {
						val := seq[(j+shift)%s]
						if forbidden[pos+j][val] {
							valid = false
							break
						}
					}
					if valid {
						dp[pos+s] = (dp[pos+s] + dp[pos]) % MOD
					}
				}
			}
		}

		fmt.Fprintln(out, dp[n]%MOD)
	}
}
