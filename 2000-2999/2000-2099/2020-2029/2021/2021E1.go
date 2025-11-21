package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

type edge struct {
	u, v int
	w    int64
}

const inf int64 = 1 << 60

func find(parent []int, x int) int {
	if parent[x] != x {
		parent[x] = find(parent, parent[x])
	}
	return parent[x]
}

func dfs(v int, n int, pLimit int, leftChild, rightChild []int, weight []int64, special []bool, capNodes, sizeSpec []int, dp [][]int64) []int64 {
	if dp[v] != nil {
		return dp[v]
	}
	if v <= n {
		capNodes[v] = 1
		if special[v] {
			sizeSpec[v] = 1
		} else {
			sizeSpec[v] = 0
		}
		limit := 1
		if pLimit < limit {
			limit = pLimit
		}
		arr := make([]int64, limit+1)
		for i := range arr {
			arr[i] = inf
		}
		if special[v] {
			if limit >= 1 {
				arr[1] = 0
			}
		} else {
			arr[0] = 0
			if limit >= 1 {
				arr[1] = 0
			}
		}
		dp[v] = arr
		return arr
	}
	lc := leftChild[v]
	rc := rightChild[v]
	dpL := dfs(lc, n, pLimit, leftChild, rightChild, weight, special, capNodes, sizeSpec, dp)
	dpR := dfs(rc, n, pLimit, leftChild, rightChild, weight, special, capNodes, sizeSpec, dp)

	capNodes[v] = capNodes[lc] + capNodes[rc]
	sizeSpec[v] = sizeSpec[lc] + sizeSpec[rc]
	limit := capNodes[v]
	if pLimit < limit {
		limit = pLimit
	}
	arr := make([]int64, limit+1)
	for i := range arr {
		arr[i] = inf
	}
	lenL := len(dpL) - 1
	lenR := len(dpR) - 1

	for i := 0; i <= lenL; i++ {
		if dpL[i] >= inf {
			continue
		}
		for j := 0; j <= lenR && i+j <= limit; j++ {
			if dpR[j] >= inf {
				continue
			}
			t := i + j
			val := dpL[i] + dpR[j]
			if val < arr[t] {
				arr[t] = val
			}
		}
	}

	w := weight[v]
	sizeL := sizeSpec[lc]
	sizeR := sizeSpec[rc]

	if sizeL > 0 {
		for j := 1; j <= lenR && j <= limit; j++ {
			if dpR[j] >= inf {
				continue
			}
			cost := dpR[j] + int64(sizeL)*w
			if cost < arr[j] {
				arr[j] = cost
			}
		}
	}
	if sizeR > 0 {
		for i := 1; i <= lenL && i <= limit; i++ {
			if dpL[i] >= inf {
				continue
			}
			cost := dpL[i] + int64(sizeR)*w
			if cost < arr[i] {
				arr[i] = cost
			}
		}
	}

	if sizeSpec[v] == 0 {
		for i := 0; i <= limit; i++ {
			if arr[i] > 0 {
				arr[i] = 0
			}
		}
	}

	dp[v] = arr
	return arr
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var T int
	fmt.Fscan(in, &T)
	for ; T > 0; T-- {
		var n, m, p int
		fmt.Fscan(in, &n, &m, &p)
		special := make([]bool, 2*n+5)
		need := make([]int, p)
		for i := 0; i < p; i++ {
			fmt.Fscan(in, &need[i])
			special[need[i]] = true
		}
		edges := make([]edge, m)
		for i := 0; i < m; i++ {
			fmt.Fscan(in, &edges[i].u, &edges[i].v, &edges[i].w)
		}
		sort.Slice(edges, func(i, j int) bool { return edges[i].w < edges[j].w })

		total := 2*n + 2
		parent := make([]int, total)
		for i := 1; i < total; i++ {
			parent[i] = i
		}
		leftChild := make([]int, total)
		rightChild := make([]int, total)
		weight := make([]int64, total)
		nextID := n
		for _, e := range edges {
			ru := find(parent, e.u)
			rv := find(parent, e.v)
			if ru == rv {
				continue
			}
			nextID++
			leftChild[nextID] = ru
			rightChild[nextID] = rv
			weight[nextID] = e.w
			parent[ru] = nextID
			parent[rv] = nextID
			parent[nextID] = nextID
		}
		root := find(parent, 1)

		capNodes := make([]int, total)
		sizeSpec := make([]int, total)
		dp := make([][]int64, total)
		dfs(root, n, p, leftChild, rightChild, weight, special, capNodes, sizeSpec, dp)
		dpRoot := dp[root]

		if p == 0 {
			for k := 1; k <= n; k++ {
				if k > 1 {
					fmt.Fprint(out, " ")
				}
				fmt.Fprint(out, 0)
			}
			fmt.Fprintln(out)
			continue
		}

		limit := len(dpRoot) - 1
		prefix := make([]int64, limit+1)
		best := inf
		for t := 0; t <= limit; t++ {
			if dpRoot[t] < best {
				best = dpRoot[t]
			}
			prefix[t] = best
		}

		for k := 1; k <= n; k++ {
			if k > 1 {
				fmt.Fprint(out, " ")
			}
			use := k
			if use > limit {
				use = limit
			}
			if use < 0 {
				fmt.Fprint(out, 0)
				continue
			}
			ans := prefix[use]
			if ans >= inf {
				ans = 0
			}
			fmt.Fprint(out, ans)
		}
		fmt.Fprintln(out)
	}
}
