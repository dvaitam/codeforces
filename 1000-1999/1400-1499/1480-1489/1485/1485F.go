package main

import (
	"bufio"
	"fmt"
	"os"
)

const mod int64 = 1000000007

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var t int
	fmt.Fscan(in, &t)
	for ; t > 0; t-- {
		var n int
		fmt.Fscan(in, &n)
		b := make([]int64, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(in, &b[i])
		}
		dp := map[int64]int64{0: 1}
		offset := int64(0)
		ans := int64(1)
		for _, x := range b {
			oldAns := ans
			offset += x
			key := x - offset
			val := dp[key]
			ans = (ans*2 - val) % mod
			if ans < 0 {
				ans += mod
			}
			dp[key] = (dp[key] + oldAns - val) % mod
			if dp[key] < 0 {
				dp[key] += mod
			}
		}
		fmt.Fprintln(out, ans%mod)
	}
}
