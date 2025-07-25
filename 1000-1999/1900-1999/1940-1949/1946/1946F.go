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
	if _, err := fmt.Fscan(in, &t); err != nil {
		return
	}
	for ; t > 0; t-- {
		var n, q int
		fmt.Fscan(in, &n, &q)
		a := make([]int, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(in, &a[i])
		}
		for ; q > 0; q-- {
			var l, r int
			fmt.Fscan(in, &l, &r)
			l--
			r--
			dp := make([]int64, r-l+1)
			var ans int64
			for i := l; i <= r; i++ {
				dp[i-l] = 1
				for j := l; j < i; j++ {
					if a[i]%a[j] == 0 {
						dp[i-l] += dp[j-l]
					}
				}
				ans += dp[i-l]
			}
			fmt.Fprintln(out, ans)
		}
	}
}
