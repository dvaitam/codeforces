package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var n int
	if _, err := fmt.Fscan(reader, &n); err != nil {
		return
	}
	const MAXV = 5000
	const LIMIT = 13 // at most 13 numbers needed for any XOR up to 2^13-1
	// Collect at most LIMIT occurrences of each value
	counts := make([]int, MAXV+1)
	vals := make([]int, 0)
	for i := 0; i < n; i++ {
		var x int
		fmt.Fscan(reader, &x)
		if counts[x] < LIMIT {
			counts[x]++
			vals = append(vals, x)
		}
	}

	const MAXX = 1 << 13 // 8192
	const INF = MAXV + 1
	dp := make([]int, MAXX)
	for i := range dp {
		dp[i] = INF
	}
	dp[0] = 0 // empty subsequence

	for _, v := range vals {
		for x := 0; x < MAXX; x++ {
			if dp[x] < v {
				if v < dp[x^v] {
					dp[x^v] = v
				}
			}
		}
	}

	res := []int{}
	for x := 0; x < MAXX; x++ {
		if dp[x] < INF {
			res = append(res, x)
		}
	}
	fmt.Fprintln(writer, len(res))
	for i, v := range res {
		if i > 0 {
			fmt.Fprint(writer, " ")
		}
		fmt.Fprint(writer, v)
	}
	fmt.Fprintln(writer)
}
