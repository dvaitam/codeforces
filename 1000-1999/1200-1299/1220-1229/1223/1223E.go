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

	var q int
	if _, err := fmt.Fscan(in, &q); err != nil {
		return
	}
	for ; q > 0; q-- {
		var n, k int
		fmt.Fscan(in, &n, &k)
		adj := make([][]struct{ to, w int }, n+1)
		for i := 0; i < n-1; i++ {
			var u, v, w int
			fmt.Fscan(in, &u, &v, &w)
			adj[u] = append(adj[u], struct{ to, w int }{v, w})
			adj[v] = append(adj[v], struct{ to, w int }{u, w})
		}
		dp0 := make([]int64, n+1)
		dp1 := make([]int64, n+1)
		var dfs func(int, int, int)
		dfs = func(v, p, w int) {
			base := int64(0)
			diffs := make([]int64, 0, len(adj[v]))
			for _, e := range adj[v] {
				if e.to == p {
					continue
				}
				dfs(e.to, v, e.w)
				base += dp0[e.to]
				diffs = append(diffs, dp1[e.to]-dp0[e.to])
			}
			sort.Slice(diffs, func(i, j int) bool { return diffs[i] > diffs[j] })
			sum0 := base
			for i := 0; i < k && i < len(diffs) && diffs[i] > 0; i++ {
				sum0 += diffs[i]
			}
			dp0[v] = sum0
			sum1 := base + int64(w)
			for i := 0; i < k-1 && i < len(diffs) && diffs[i] > 0; i++ {
				sum1 += diffs[i]
			}
			dp1[v] = sum1
		}
		dfs(1, 0, 0)
		fmt.Fprintln(out, dp0[1])
	}
}
