package main

import (
	"bufio"
	"fmt"
	"os"
)

func min64(a, b int64) int64 {
	if a < b {
		return a
	}
	return b
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n int
	if _, err := fmt.Fscan(in, &n); err != nil {
		return
	}
	goals := make([]int64, n+1)
	var total int64
	for i := 1; i <= n; i++ {
		fmt.Fscan(in, &goals[i])
		total += goals[i]
	}

	freebies := make([]int64, n+1)
	contrib := make([]int64, n+1)
	var contribSum int64

	var q int
	fmt.Fscan(in, &q)

	type key struct{ s, t int }
	milestones := make(map[key]int)

	for ; q > 0; q-- {
		var s, t, u int
		fmt.Fscan(in, &s, &t, &u)
		k := key{s, t}
		if oldU, ok := milestones[k]; ok {
			// remove existing milestone
			oldC := contrib[oldU]
			freebies[oldU]--
			newC := min64(goals[oldU], freebies[oldU])
			contrib[oldU] = newC
			contribSum += newC - oldC
			delete(milestones, k)
		}
		if u != 0 {
			milestones[k] = u
			oldC := contrib[u]
			freebies[u]++
			newC := min64(goals[u], freebies[u])
			contrib[u] = newC
			contribSum += newC - oldC
		}
		ans := total - contribSum
		fmt.Fprintln(out, ans)
	}
}
