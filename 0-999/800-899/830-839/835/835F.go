package main

import (
	"bufio"
	"fmt"
	"os"
)

type Edge struct {
	to int
	w  int64
}

var (
	n            int
	adj          [][]Edge
	inCycle      []bool
	branchDepth  []int64
	baseDiameter int64
)

func dfs(u, p int) int64 {
	max1, max2 := int64(0), int64(0)
	for _, e := range adj[u] {
		v := e.to
		if v == p || inCycle[v] {
			continue
		}
		d := dfs(v, u) + e.w
		if d > max1 {
			max2 = max1
			max1 = d
		} else if d > max2 {
			max2 = d
		}
	}
	if max1+max2 > baseDiameter {
		baseDiameter = max1 + max2
	}
	return max1
}

func max64(a, b int64) int64 {
	if a > b {
		return a
	}
	return b
}

func check(m int, a []int64, pos []int64, limit int64) bool {
	type pair struct {
		idx int
		val int64
	}
	dq := make([]pair, 0)
	l := 0
	for r := 0; r < 2*m; r++ {
		val := a[r] - pos[r]
		for len(dq) > 0 && dq[len(dq)-1].val <= val {
			dq = dq[:len(dq)-1]
		}
		dq = append(dq, pair{r, val})
		for {
			for len(dq) > 0 && dq[0].idx < l {
				dq = dq[1:]
			}
			if len(dq) == 0 {
				break
			}
			if a[r]+pos[r]+dq[0].val <= limit && r-l+1 <= m {
				break
			}
			if dq[0].idx == l {
				dq = dq[1:]
			}
			l++
		}
		if r-l+1 == m {
			return true
		}
	}
	return false
}

func main() {
	in := bufio.NewReader(os.Stdin)
	fmt.Fscan(in, &n)
	adj = make([][]Edge, n)
	for i := 0; i < n; i++ {
		var u, v int
		var l int64
		fmt.Fscan(in, &u, &v, &l)
		u--
		v--
		adj[u] = append(adj[u], Edge{v, l})
		adj[v] = append(adj[v], Edge{u, l})
	}

	deg := make([]int, n)
	for i := 0; i < n; i++ {
		deg[i] = len(adj[i])
	}
	queue := make([]int, 0)
	for i := 0; i < n; i++ {
		if deg[i] == 1 {
			queue = append(queue, i)
		}
	}
	removed := make([]bool, n)
	for len(queue) > 0 {
		v := queue[0]
		queue = queue[1:]
		removed[v] = true
		for _, e := range adj[v] {
			to := e.to
			if removed[to] {
				continue
			}
			deg[to]--
			if deg[to] == 1 {
				queue = append(queue, to)
			}
		}
	}
	inCycle = make([]bool, n)
	for i := 0; i < n; i++ {
		if !removed[i] {
			inCycle[i] = true
		}
	}

	// collect cycle nodes in order
	start := -1
	for i := 0; i < n; i++ {
		if inCycle[i] {
			start = i
			break
		}
	}
	cycle := make([]int, 0)
	edges := make([]int64, 0)
	prev := -1
	u := start
	for {
		cycle = append(cycle, u)
		next := -1
		var w int64
		for _, e := range adj[u] {
			if e.to == prev || !inCycle[e.to] {
				continue
			}
			next = e.to
			w = e.w
			break
		}
		if next == -1 {
			break
		}
		prev = u
		u = next
		edges = append(edges, w)
		if u == start {
			break
		}
	}

	m := len(cycle)
	branchDepth = make([]int64, n)
	baseDiameter = 0
	for i := 0; i < m; i++ {
		v := cycle[i]
		for _, e := range adj[v] {
			if inCycle[e.to] {
				continue
			}
			d := dfs(e.to, v) + e.w
			if d > branchDepth[v] {
				branchDepth[v] = d
			}
		}
	}

	// prepare arrays for check
	a := make([]int64, 2*m)
	pos := make([]int64, 2*m)
	for i := 0; i < 2*m; i++ {
		a[i] = branchDepth[cycle[i%m]]
		if i > 0 {
			pos[i] = pos[i-1] + edges[(i-1)%m]
		}
	}

	low := baseDiameter
	high := int64(1)
	for i := 0; i < m; i++ {
		high += a[i] + pos[i+1] - pos[0]
	}
	if high < low {
		high = low
	}

	for low < high {
		mid := (low + high) / 2
		if check(m, a, pos, mid) {
			high = mid
		} else {
			low = mid + 1
		}
	}
	fmt.Println(low)
}
