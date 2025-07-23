package main

import (
	"bufio"
	"fmt"
	"os"
)

// This program solves the problem described in problemA.txt.
// Given the probabilities of requesting each of the n videos and
// the cache size k using an LRU policy, it outputs the probability
// that each video is present in the cache after a long sequence of
// requests. Under the independent reference model, this is equal to
// the probability that a video's last request is among the k most
// recent distinct requests.
func main() {
	in := bufio.NewReader(os.Stdin)
	var n, k int
	if _, err := fmt.Fscan(in, &n, &k); err != nil {
		return
	}
	p := make([]float64, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &p[i])
	}

	// dp[mask] holds the probability that the set of the first t
	// distinct items (in order) equals 'mask'. We iterate for t = 0..k.
	dp := map[int]float64{0: 1.0}
	for step := 0; step < k; step++ {
		ndp := make(map[int]float64)
		for mask, prob := range dp {
			var used float64
			for i := 0; i < n; i++ {
				if mask>>i&1 != 0 {
					used += p[i]
				}
			}
			remain := 1 - used
			if remain <= 0 {
				// No other items can appear.
				ndp[mask] += prob
				continue
			}
			for i := 0; i < n; i++ {
				if mask>>i&1 == 0 {
					ndp[mask|1<<i] += prob * p[i] / remain
				}
			}
		}
		dp = ndp
	}

	ans := make([]float64, n)
	for mask, prob := range dp {
		for i := 0; i < n; i++ {
			if mask>>i&1 != 0 {
				ans[i] += prob
			}
		}
	}

	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()
	for i := 0; i < n; i++ {
		if i > 0 {
			fmt.Fprint(out, " ")
		}
		fmt.Fprintf(out, "%.9f", ans[i])
	}
	fmt.Fprintln(out)
}
