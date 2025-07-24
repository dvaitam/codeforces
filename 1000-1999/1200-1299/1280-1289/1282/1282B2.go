package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var T int
	if _, err := fmt.Fscan(reader, &T); err != nil {
		return
	}
	for ; T > 0; T-- {
		var n, p, k int
		fmt.Fscan(reader, &n, &p, &k)
		prices := make([]int, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(reader, &prices[i])
		}
		sort.Ints(prices)
		dp := make([]int64, n+1)
		ans := 0
		for i := 1; i <= n; i++ {
			dp[i] = dp[i-1] + int64(prices[i-1])
			if i >= k {
				// using offer buying last k items at cost of max price
				if alt := dp[i-k] + int64(prices[i-1]); alt < dp[i] {
					dp[i] = alt
				}
			}
			if dp[i] <= int64(p) {
				ans = i
			}
		}
		fmt.Fprintln(writer, ans)
	}
}
