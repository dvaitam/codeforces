package main

import (
	"bufio"
	"fmt"
	"os"
)

const inf int64 = 1 << 60

type edge struct {
	to int
	id int
}

type childData struct {
	id int
	c  int64
}

var (
	n, m  int
	g     [][]edge
	edges [][2]int
	isBr  []bool
	tin   []int
	low   []int
	timer int
	brAdj [][]int
	dp0   []int64
	dp1   []int64
	ans   [][2]int
)

func addEdge(u, v int) {
	id := len(edges)
	edges = append(edges, [2]int{u, v})
	g[u] = append(g[u], edge{to: v, id: id})
	g[v] = append(g[v], edge{to: u, id: id})
}

func dfsBridge(v, pEdge int) {
	timer++
	tin[v] = timer
	low[v] = timer
	for _, e := range g[v] {
		if e.id == pEdge {
			continue
		}
		to := e.to
		if tin[to] == 0 {
			dfsBridge(to, e.id)
			if low[to] > tin[v] {
				isBr[e.id] = true
			}
			if low[to] < low[v] {
				low[v] = low[to]
			}
		} else {
			if tin[to] < low[v] {
				low[v] = tin[to]
			}
		}
	}
}

// choose selects a subset of children (described by vals) whose size has required parity.
// It returns the selected set and the additional cost (sum of c values).
func choose(vals []childData, parity int, needOne bool) (map[int]bool, int64, bool) {
	sel := make(map[int]bool)
	cnt := 0
	var sum int64
	bestPosVal := inf
	bestPosID := -1
	bestNegVal := int64(-inf)
	bestNegID := -1
	bestAllVal := inf
	bestAllID := -1

	for _, cd := range vals {
		if cd.c < 0 {
			sel[cd.id] = true
			cnt++
			sum += cd.c
			if cd.c > bestNegVal {
				bestNegVal = cd.c
				bestNegID = cd.id
			}
		} else {
			if cd.c < bestPosVal {
				bestPosVal = cd.c
				bestPosID = cd.id
			}
		}
		if cd.c < bestAllVal {
			bestAllVal = cd.c
			bestAllID = cd.id
		}
	}

	if len(vals) == 0 && needOne {
		return nil, 0, false
	}

	if cnt%2 != parity {
		addCost := bestPosVal
		remCost := int64(inf)
		if bestNegID != -1 {
			remCost = -bestNegVal
		}
		if bestPosID == -1 && remCost >= inf {
			return nil, 0, false
		}
		if addCost <= remCost {
			sel[bestPosID] = true
			sum += bestPosVal
			cnt++
		} else {
			delete(sel, bestNegID)
			sum -= bestNegVal
			cnt--
		}
	}

	if needOne && cnt == 0 {
		if bestAllID == -1 {
			return nil, 0, false
		}
		sel[bestAllID] = true
		sum += bestAllVal
		cnt++
	}

	if cnt%2 != parity {
		return nil, 0, false
	}
	return sel, sum, true
}

// calcCost computes minimal cost (in doubled units) for the subtree rooted at v when its
// edge to parent is treated according to offset/parity/subtract flags.
// offset: 0 normally, 1 when parent edge is also counted as open (dp0 state).
// parityDesired: 0 if total open edges should be even, 1 if odd.
// subtract: 1 only for root1 state, otherwise 0.
// needOne: whether at least one open edge is required (root1 state).
// parent is ignored for root processing.
func calcCost(v, parent int, offset int, parityDesired int, subtract int, needOne bool) (int64, map[int]bool, bool) {
	base := int64(0)
	openForced := 0
	vals := make([]childData, 0)

	for _, to := range brAdj[v] {
		if to == parent {
			continue
		}
		d0 := dp0[to]
		d1 := dp1[to]
		if d0 >= inf && d1 >= inf {
			return inf, nil, false
		}
		if d0 >= inf {
			base += d1
			openForced++
		} else if d1 >= inf {
			base += d0
		} else {
			base += d0
			vals = append(vals, childData{id: to, c: d1 - d0 + 1})
		}
	}

	parityNeeded := (parityDesired ^ (openForced & 1) ^ (offset & 1)) & 1
	needPositive := needOne && openForced == 0
	sel, add, ok := choose(vals, parityNeeded, needPositive)
	if !ok {
		return inf, nil, false
	}

	totalOpen := openForced + offset
	totalCost := base + int64(totalOpen) - int64(subtract) + add
	return totalCost, sel, true
}

func dfsDP(v, parent int) {
	for _, to := range brAdj[v] {
		if to == parent {
			continue
		}
		dfsDP(to, v)
	}
	cost1, _, ok1 := calcCost(v, parent, 0, 0, 0, false)
	cost0, _, ok0 := calcCost(v, parent, 1, 0, 0, false)
	if ok1 {
		dp1[v] = cost1
	} else {
		dp1[v] = inf
	}
	if ok0 {
		dp0[v] = cost0
	} else {
		dp0[v] = inf
	}
}

func build(v, parent int, state int, parentEndpoint int) int {
	// state: 0 parent edge open (dp1), 1 parent edge closed here (dp0), 2 root0, 3 root1
	offset := 0
	parityDesired := 0
	needOne := false
	switch state {
	case 0:
		offset = 0
		parityDesired = 0
	case 1:
		offset = 1
		parityDesired = 0
	case 2:
		offset = 0
		parityDesired = 0
	case 3:
		offset = 0
		parityDesired = 1
		needOne = true
	}

	// Recompute selection with the same rules used in dp.
	base := int64(0)
	openForced := 0
	vals := make([]childData, 0)
	forcedOpen := make(map[int]bool)
	forcedClose := make(map[int]bool)

	for _, to := range brAdj[v] {
		if to == parent {
			continue
		}
		d0 := dp0[to]
		d1 := dp1[to]
		if d0 >= inf && d1 >= inf {
			continue
		}
		if d0 >= inf {
			base += d1
			openForced++
			forcedOpen[to] = true
		} else if d1 >= inf {
			base += d0
			forcedClose[to] = true
		} else {
			base += d0
			vals = append(vals, childData{id: to, c: d1 - d0 + 1})
		}
	}

	parityNeeded := (parityDesired ^ (openForced & 1) ^ (offset & 1)) & 1
	needPositive := needOne && openForced == 0
	sel, _, ok := choose(vals, parityNeeded, needPositive)
	if !ok {
		sel = make(map[int]bool)
	}

	// Determine child states and gather open endpoints.
	opens := make([]int, 0)
	for _, to := range brAdj[v] {
		if to == parent {
			continue
		}
		childState := 0
		if forcedOpen[to] || sel[to] {
			childState = 0 // dp1 for child (edge open)
		} else {
			childState = 1 // dp0 for child (edge closed at child)
		}
		endpoint := build(to, v, childState, v)
		if forcedOpen[to] || sel[to] {
			opens = append(opens, endpoint)
		}
	}

	if state == 1 {
		opens = append(opens, parentEndpoint)
	}

	if state == 3 {
		// leave one open edge unpaired
		if len(opens) > 0 {
			opens = opens[1:]
		}
	}

	for i := 0; i+1 < len(opens); i += 2 {
		ans = append(ans, [2]int{opens[i], opens[i+1]})
	}

	if state == 0 {
		return v
	}
	return -1
}

func processComponent(start int) {
	dfsDP(start, -1)
	costRoot0, _, ok0 := calcCost(start, -1, 0, 0, 0, false)
	costRoot1, _, ok1 := calcCost(start, -1, 0, 1, 1, true)
	chosenState := 2
	if !ok0 || (ok1 && costRoot1 < costRoot0) {
		chosenState = 3
	}
	build(start, -1, chosenState, -1)
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	for {
		if _, err := fmt.Fscan(in, &n, &m); err != nil {
			return
		}
		if n == 0 && m == 0 {
			break
		}

		g = make([][]edge, n+1)
		edges = edges[:0]
		for i := 0; i < m; i++ {
			var s int
			fmt.Fscan(in, &s)
			path := make([]int, s)
			for j := 0; j < s; j++ {
				fmt.Fscan(in, &path[j])
			}
			for j := 0; j+1 < s; j++ {
				addEdge(path[j], path[j+1])
			}
		}

		isBr = make([]bool, len(edges))
		tin = make([]int, n+1)
		low = make([]int, n+1)
		timer = 0
		for v := 1; v <= n; v++ {
			if tin[v] == 0 {
				dfsBridge(v, -1)
			}
		}

		brAdj = make([][]int, n+1)
		for id, e := range edges {
			if isBr[id] {
				u, v := e[0], e[1]
				brAdj[u] = append(brAdj[u], v)
				brAdj[v] = append(brAdj[v], u)
			}
		}

		dp0 = make([]int64, n+1)
		dp1 = make([]int64, n+1)
		ans = ans[:0]

		visited := make([]bool, n+1)
		var stack []int
		for v := 1; v <= n; v++ {
			if visited[v] || len(brAdj[v]) == 0 {
				continue
			}
			// gather component vertices to run dfs safely
			stack = stack[:0]
			stack = append(stack, v)
			visited[v] = true
			for len(stack) > 0 {
				x := stack[len(stack)-1]
				stack = stack[:len(stack)-1]
				for _, to := range brAdj[x] {
					if !visited[to] {
						visited[to] = true
						stack = append(stack, to)
					}
				}
			}
			processComponent(v)
		}

		fmt.Fprintln(out, len(ans))
		for _, e := range ans {
			fmt.Fprintln(out, e[0], e[1])
		}
	}
}
