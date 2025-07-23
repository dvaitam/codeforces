package main

import (
	"bufio"
	"fmt"
	"os"
)

const inf = int(1e9)

var (
	n    int
	adj  [][]int
	root int
)

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func dfs(v, p int) ([]int, []int, int) {
	isLeaf := len(adj[v]) == 1 && v != root
	if isLeaf {
		dp0 := []int{0, inf}
		dp1 := []int{inf, 0}
		return dp0, dp1, 1
	}
	dp0 := []int{0}
	dp1 := []int{0}
	leaves := 0
	for _, u := range adj[v] {
		if u == p {
			continue
		}
		c0, c1, sz := dfs(u, v)
		best0 := make([]int, sz+1)
		best1 := make([]int, sz+1)
		for t := 0; t <= sz; t++ {
			c0t := inf
			if t < len(c0) {
				c0t = c0[t]
			}
			c1t := inf
			if t < len(c1) {
				c1t = c1[t]
			}
			best0[t] = min(c0t, c1t+1)
			best1[t] = min(c0t+1, c1t)
		}
		newSize := leaves + sz
		new0 := make([]int, newSize+1)
		new1 := make([]int, newSize+1)
		for i := 0; i <= newSize; i++ {
			new0[i] = inf
			new1[i] = inf
		}
		for i := 0; i <= leaves; i++ {
			if dp0[i] < inf {
				for t := 0; t <= sz; t++ {
					val := dp0[i] + best0[t]
					if val < new0[i+t] {
						new0[i+t] = val
					}
				}
			}
			if dp1[i] < inf {
				for t := 0; t <= sz; t++ {
					val := dp1[i] + best1[t]
					if val < new1[i+t] {
						new1[i+t] = val
					}
				}
			}
		}
		dp0 = new0
		dp1 = new1
		leaves += sz
	}
	return dp0, dp1, leaves
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	if _, err := fmt.Fscan(in, &n); err != nil {
		return
	}
	adj = make([][]int, n+1)
	for i := 0; i < n-1; i++ {
		var x, y int
		fmt.Fscan(in, &x, &y)
		adj[x] = append(adj[x], y)
		adj[y] = append(adj[y], x)
	}
	leafCount := 0
	for i := 1; i <= n; i++ {
		if len(adj[i]) == 1 {
			leafCount++
		}
	}
	root = 1
	for i := 1; i <= n; i++ {
		if len(adj[i]) > 1 {
			root = i
			break
		}
	}
	dp0, dp1, _ := dfs(root, -1)
	if len(adj[root]) == 1 {
		ext := make([]int, len(dp1)+1)
		for i := 0; i < len(ext); i++ {
			ext[i] = inf
		}
		for i := 0; i < len(dp1); i++ {
			if dp1[i] < ext[i+1] {
				ext[i+1] = dp1[i]
			}
		}
		dp1 = ext
		if len(dp0) < len(dp1) {
			tmp := make([]int, len(dp1))
			for i := range tmp {
				tmp[i] = inf
			}
			copy(tmp, dp0)
			dp0 = tmp
		}
	}
	need := leafCount / 2
	ans := inf
	if need < len(dp0) && dp0[need] < ans {
		ans = dp0[need]
	}
	if need < len(dp1) && dp1[need] < ans {
		ans = dp1[need]
	}
	fmt.Fprintln(out, ans)
}
