package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

type Pair struct {
	val int
	idx int
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var n int
	if _, err := fmt.Fscan(reader, &n); err != nil {
		return
	}
	a := make([]int, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(reader, &a[i])
	}

	pairs := make([]Pair, n)
	for i := 0; i < n; i++ {
		pairs[i] = Pair{val: a[i], idx: i}
	}
	sort.Slice(pairs, func(i, j int) bool { return pairs[i].val < pairs[j].val })

	const INF int64 = 1<<63 - 1
	dp := make([]int64, n+1)
	prev := make([]int, n+1)
	for i := range dp {
		dp[i] = INF
	}
	dp[0] = 0

	for i := 3; i <= n; i++ {
		for s := 3; s <= 5; s++ {
			if i-s < 0 {
				continue
			}
			diff := pairs[i-1].val - pairs[i-s].val
			cost := dp[i-s] + int64(diff)
			if cost < dp[i] {
				dp[i] = cost
				prev[i] = s
			}
		}
	}

	teams := make([]int, n)
	teamCount := 0
	for i := n; i > 0; {
		s := prev[i]
		teamCount++
		for j := i - s; j < i; j++ {
			teams[pairs[j].idx] = teamCount
		}
		i -= s
	}

	fmt.Fprintln(writer, dp[n], teamCount)
	for i := 0; i < n; i++ {
		if i > 0 {
			fmt.Fprint(writer, " ")
		}
		fmt.Fprint(writer, teams[i])
	}
	fmt.Fprintln(writer)
}
