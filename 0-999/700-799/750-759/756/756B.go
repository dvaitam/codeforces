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

	var n int
	if _, err := fmt.Fscan(reader, &n); err != nil {
		return
	}
	t := make([]int, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(reader, &t[i])
	}

	dp := make([]int, n+1)
	for i := 1; i <= n; i++ {
		time := t[i-1]
		cost := dp[i-1] + 20

		j := sort.Search(len(t), func(k int) bool { return t[k] >= time-89 })
		if c := dp[j] + 50; c < cost {
			cost = c
		}

		j = sort.Search(len(t), func(k int) bool { return t[k] >= time-1439 })
		if c := dp[j] + 120; c < cost {
			cost = c
		}

		dp[i] = cost
		fmt.Fprintln(writer, dp[i]-dp[i-1])
	}
}
