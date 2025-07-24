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
		ans := 0
		for start := 0; start < n; start += 256 {
			end := start + 256
			if end > n {
				end = n
			}
			m := end - start
			dp := make([]int, m)
			best := 1
			for i := 0; i < m; i++ {
				dp[i] = 1
				for j := 0; j < i; j++ {
					if (a[start+j] ^ (start + i)) < (a[start+i] ^ (start + j)) {
						if dp[j]+1 > dp[i] {
							dp[i] = dp[j] + 1
						}
					}
				}
				if dp[i] > best {
					best = dp[i]
				}
			}
			if best > ans {
				ans = best
			}
		}
		fmt.Fprintln(out, ans)
	}
}
