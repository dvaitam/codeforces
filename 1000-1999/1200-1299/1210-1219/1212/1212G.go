package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
)

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n int
	var l int
	if _, err := fmt.Fscan(in, &n, &l); err != nil {
		return
	}
	x := make([]int, n+1)
	b := make([]int, n+1)
	for i := 1; i <= n; i++ {
		fmt.Fscan(in, &x[i], &b[i])
	}

	// helper to run DP for given lambda and optionally keep parents
	dp := make([]float64, n+1)
	parent := make([]int, n+1)
	path := func(lambda float64, keepParent bool) float64 {
		for i := 1; i <= n; i++ {
			dp[i] = math.Inf(1)
			if keepParent {
				parent[i] = -1
			}
		}
		dp[0] = 0
		for i := 1; i <= n; i++ {
			best := math.Inf(1)
			bestParent := -1
			for j := 0; j < i; j++ {
				diff := float64(x[i] - x[j])
				val := dp[j] + math.Sqrt(math.Abs(diff-float64(l))) - lambda*float64(b[i])
				if val < best {
					best = val
					bestParent = j
				}
			}
			dp[i] = best
			if keepParent {
				parent[i] = bestParent
			}
		}
		return dp[n]
	}

	low, high := 0.0, 2000.0
	for iter := 0; iter < 60; iter++ {
		mid := (low + high) / 2
		if path(mid, false) <= 0 {
			high = mid
		} else {
			low = mid
		}
	}
	path(high, true)
	// reconstruct path
	var seq []int
	cur := n
	for cur > 0 {
		seq = append(seq, cur)
		cur = parent[cur]
	}
	// reverse
	for i := len(seq) - 1; i >= 0; i-- {
		if i < len(seq)-1 {
			fmt.Fprint(out, " ")
		}
		fmt.Fprint(out, seq[i])
	}
	if len(seq) > 0 {
		fmt.Fprintln(out)
	}
}
