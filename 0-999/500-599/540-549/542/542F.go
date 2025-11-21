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

	var n, T int
	fmt.Fscan(in, &n, &T)
	items := make([]struct {
		t int
		q int
	}, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &items[i].t, &items[i].q)
	}

	// dp[c] = max total interest using total time c
	dp := make([]int, T+1)
	for _, it := range items {
		for c := T; c >= it.t; c-- {
			if dp[c-it.t]+it.q > dp[c] {
				dp[c] = dp[c-it.t] + it.q
			}
		}
	}

	ans := 0
	for c := 0; c <= T; c++ {
		if dp[c] > ans {
			ans = dp[c]
		}
	}
	fmt.Fprintln(out, ans)
}
