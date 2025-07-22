package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var k, n int
	if _, err := fmt.Fscan(in, &k, &n); err != nil {
		return
	}
	times := make([]int64, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &times[i])
	}
	sort.Slice(times, func(i, j int) bool { return times[i] < times[j] })

	const INF int64 = 1 << 60
	prev := make([]int64, n+1)
	curr := make([]int64, n+1)
	for j := 1; j <= k; j++ {
		prefix := INF
		curr[0] = INF
		for i := 1; i <= n; i++ {
			cand := prefix + times[i-1]
			if cand < curr[i-1] {
				curr[i] = cand
			} else {
				curr[i] = curr[i-1]
			}
			val := prev[i-1] - times[i-1]
			if val < prefix {
				prefix = val
			}
		}
		copy(prev, curr)
	}
	fmt.Fprintln(out, prev[n])
}
