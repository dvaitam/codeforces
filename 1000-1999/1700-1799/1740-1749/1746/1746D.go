package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

var (
	g [][]int
	s []int64
)

func dfs(u int, k int64) (int64, int64) {
	val := k * s[u]
	if len(g[u]) == 0 {
		return val, s[u]
	}
	m := int64(len(g[u]))
	base := k / m
	rem := int(k % m)
	diffs := make([]int64, len(g[u]))
	for i, v := range g[u] {
		childVal, childDiff := dfs(v, base)
		val += childVal
		diffs[i] = childDiff
	}
	sort.Slice(diffs, func(i, j int) bool { return diffs[i] > diffs[j] })
	for i := 0; i < rem; i++ {
		val += diffs[i]
	}
	delta := s[u] + diffs[rem]
	return val, delta
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var T int
	fmt.Fscan(in, &T)
	for ; T > 0; T-- {
		var n int
		var k int64
		fmt.Fscan(in, &n, &k)
		g = make([][]int, n)
		p := make([]int, n)
		for i := 1; i < n; i++ {
			fmt.Fscan(in, &p[i])
			p[i]--
			g[p[i]] = append(g[p[i]], i)
		}
		s = make([]int64, n)
		for i := range s {
			fmt.Fscan(in, &s[i])
		}
		ans, _ := dfs(0, k)
		fmt.Fprintln(out, ans)
	}
}
