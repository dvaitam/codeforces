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
		pre := 0
		dp := make([]int, n+1)
		mp := map[int]int{0: 0}
		for i := 1; i <= n; i++ {
			pre ^= a[i-1]
			dp[i] = dp[i-1]
			if v, ok := mp[pre]; ok {
				if dp[i] < v+1 {
					dp[i] = v + 1
				}
			}
			if mpv, ok := mp[pre]; !ok || mpv < dp[i] {
				mp[pre] = dp[i]
			}
		}
		fmt.Fprintln(out, n-dp[n])
	}
}
