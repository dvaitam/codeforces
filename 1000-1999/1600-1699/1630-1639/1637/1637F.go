package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

func main() {
	in := bufio.NewReader(os.Stdin)
	var n int
	if _, err := fmt.Fscan(in, &n); err != nil {
		return
	}
	h := make([]int, n+1)
	for i := 1; i <= n; i++ {
		fmt.Fscan(in, &h[i])
	}
	adj := make([][]int, n+1)
	for i := 0; i < n-1; i++ {
		var u, v int
		fmt.Fscan(in, &u, &v)
		adj[u] = append(adj[u], v)
		adj[v] = append(adj[v], u)
	}
	if n == 2 {
		// only two vertices, need towers at both
		res := int64(h[1] + h[2])
		if h[1] > h[2] {
			res = int64(2 * h[1])
		} else if h[2] > h[1] {
			res = int64(2 * h[2])
		}
		fmt.Println(res)
		return
	}

	var ans int64
	visited := make([]bool, n+1)
	var dfs func(int, int) int
	dfs = func(v, p int) int {
		visited[v] = true
		childs := make([]int, 0)
		for _, to := range adj[v] {
			if to == p {
				continue
			}
			childs = append(childs, dfs(to, v))
		}
		if len(childs) == 0 {
			ans += int64(h[v])
			return h[v]
		}
		sort.Slice(childs, func(i, j int) bool { return childs[i] > childs[j] })
		if len(childs) == 1 {
			if childs[0] < h[v] {
				ans += int64(h[v] - childs[0])
				childs[0] = h[v]
			}
			return childs[0]
		}
		// len >= 2
		if childs[1] < h[v] {
			ans += int64(h[v] - childs[1])
			childs[1] = h[v]
		}
		if childs[0] < h[v] {
			childs[0] = h[v]
		}
		return childs[0]
	}

	root := 1
	topVals := make([]int, 0)
	for _, to := range adj[root] {
		topVals = append(topVals, dfs(to, root))
	}
	sort.Slice(topVals, func(i, j int) bool { return topVals[i] > topVals[j] })

	if len(topVals) == 1 {
		if topVals[0] < h[root] {
			ans += int64(h[root] - topVals[0])
			topVals[0] = h[root]
		}
		ans += int64(h[root])
		fmt.Println(ans)
		return
	}

	cost1 := ans
	if topVals[0] < h[root] {
		cost1 += int64(h[root] - topVals[0])
	}
	if topVals[1] < h[root] {
		cost1 += int64(h[root] - topVals[1])
	}

	cost2 := ans
	if topVals[0] < h[root] {
		cost2 += int64(h[root] - topVals[0])
		topVals[0] = h[root]
	}
	cost2 += int64(h[root])

	if cost1 < cost2 {
		fmt.Println(cost1)
	} else {
		fmt.Println(cost2)
	}
}
