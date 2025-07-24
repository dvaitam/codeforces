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

	var t int
	if _, err := fmt.Fscan(reader, &t); err != nil {
		return
	}
	for ; t > 0; t-- {
		var n, p, k int
		fmt.Fscan(reader, &n, &p, &k)
		prices := make([]int, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(reader, &prices[i])
		}
		sort.Ints(prices)
		cost := make([]int, n+1)
		for i := 1; i <= n; i++ {
			cost[i] = cost[i-1] + prices[i-1]
			if i >= k {
				if c := cost[i-k] + prices[i-1]; c < cost[i] {
					cost[i] = c
				}
			}
		}
		ans := 0
		for i := 1; i <= n; i++ {
			if cost[i] <= p {
				ans = i
			}
		}
		fmt.Fprintln(writer, ans)
	}
}
