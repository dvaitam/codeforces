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

	var T int
	fmt.Fscan(in, &T)
	for ; T > 0; T-- {
		var n, k int
		fmt.Fscan(in, &n, &k)
		parentArr := make([]int, n+1)
		adj := make([][]int, n+1)
		for i := 2; i <= n; i++ {
			fmt.Fscan(in, &parentArr[i])
			p := parentArr[i]
			adj[p] = append(adj[p], i)
			adj[i] = append(adj[i], p)
		}
		LOG := 0
		for (1 << LOG) <= n {
			LOG++
		}
		up := make([][]int, LOG)
		for i := range up {
			up[i] = make([]int, n+1)
		}
		depth := make([]int, n+1)
		tin := make([]int, n+1)
		tout := make([]int, n+1)
		order := make([]int, 0, n)
		var timer int
		var dfs func(int, int)
		dfs = func(v, p int) {
			tin[v] = timer
			order = append(order, v)
			up[0][v] = p
			if v == 1 {
				depth[v] = 0
			} else {
				depth[v] = depth[p] + 1
			}
			timer++
			for _, to := range adj[v] {
				if to == p {
					continue
				}
				dfs(to, v)
			}
			tout[v] = timer - 1
		}
		dfs(1, 0)
		for j := 1; j < LOG; j++ {
			for i := 1; i <= n; i++ {
				prev := up[j-1][i]
				if prev != 0 {
					up[j][i] = up[j-1][prev]
				}
			}
		}
		nodes := make([]int, n)
		for i := 1; i <= n; i++ {
			nodes[i-1] = i
		}
		sort.Slice(nodes, func(i, j int) bool {
			return depth[nodes[i]] > depth[nodes[j]]
		})

		jump := func(v, d int) int {
			for i := 0; d > 0 && v != 0; i++ {
				if d&1 != 0 {
					v = up[i][v]
				}
				d >>= 1
			}
			if v == 0 {
				return 1
			}
			return v
		}

		can := func(h int) bool {
			parentDSU := make([]int, n+1)
			for i := range parentDSU {
				parentDSU[i] = i
			}
			var find func(int) int
			find = func(x int) int {
				if parentDSU[x] != x {
					parentDSU[x] = find(parentDSU[x])
				}
				return parentDSU[x]
			}
			ops := 0
			for _, v := range nodes {
				if depth[v] <= h {
					break
				}
				if find(tin[v]) > tout[v] {
					continue
				}
				ops++
				if ops > k {
					return false
				}
				u := jump(v, h-1)
				l, r := tin[u], tout[u]
				for x := find(l); x <= r; x = find(x) {
					parentDSU[x] = x + 1
				}
			}
			return true
		}

		left, right := 1, n
		for left < right {
			mid := (left + right) >> 1
			if can(mid) {
				right = mid
			} else {
				left = mid + 1
			}
		}
		fmt.Fprintln(out, left)
	}
}
