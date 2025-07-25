package main

import (
	"bufio"
	"fmt"
	"os"
)

var (
	g [][]int
	a []int64
)

func dfs(v int) (int64, int64) {
	if len(g[v]) == 0 {
		return a[v], 0
	}
	var total, cost int64
	for _, u := range g[v] {
		s, c := dfs(u)
		total += s
		cost += c
	}
	if total < a[v] {
		diff := a[v] - total
		cost += diff
		total += diff
	}
	return total + a[v], cost
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var T int
	if _, err := fmt.Fscan(in, &T); err != nil {
		return
	}
	for ; T > 0; T-- {
		var n int
		fmt.Fscan(in, &n)
		a = make([]int64, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(in, &a[i])
		}
		g = make([][]int, n)
		for i := 1; i < n; i++ {
			var p int
			fmt.Fscan(in, &p)
			g[p-1] = append(g[p-1], i)
		}
		_, ans := dfs(0)
		fmt.Fprintln(out, ans)
	}
}
