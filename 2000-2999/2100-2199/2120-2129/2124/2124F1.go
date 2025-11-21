package main

import (
	"bufio"
	"fmt"
	"os"
)

const MOD int64 = 998244353

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
		for i := range forbidden {
			forbidden[i] = make([]bool, n+1)
		}
		for i := 0; i < m; i++ {
			var pos, val int
			fmt.Fscan(in, &pos, &val)
			forbidden[pos-1][val] = true
		}

		dp := make([]int64, n+1)
		dp[0] = 1
		for pos := 0; pos < n; pos++ {
			if dp[pos] == 0 {
				continue
			}
			for s := 1; pos+s <= n; s++ {
				good := 0
				for r := 1; r <= s; r++ {
					ok := true
					val := r
					for offset := 0; offset < s; offset++ {
						if forbidden[pos+offset][val] {
							ok = false
							break
						}
						val++
						if val > s {
							val = 1
						}
					}
					if ok {
						good++
					}
				}
				if good > 0 {
					dp[pos+s] = (dp[pos+s] + dp[pos]*int64(good)) % MOD
				}
			}
		}

		fmt.Fprintln(out, dp[n]%MOD)
	}
}
