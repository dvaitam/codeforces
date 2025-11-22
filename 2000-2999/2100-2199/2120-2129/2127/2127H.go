package main

import (
	"bufio"
	"fmt"
	"os"
)

type edge struct {
	u, v int
}

var (
	gEdges []edge
	adj    [][]int
	best   int
	n, m   int
)

// Depth-first branch and bound.
func dfs(cap []int, state []int8, degRem []int, unknown, chosen, sumCap int) {
	// Upper bound pruning.
	ub := chosen
	if unknown < sumCap/2 {
		ub += unknown
	} else {
		ub += sumCap / 2
	}
	if ub <= best {
		return
	}

	// Find a vertex whose remaining incident unknown edges exceed its capacity.
	v := -1
	for i := 0; i < n; i++ {
		if degRem[i] > cap[i] {
			v = i
			break
		}
	}

	if v == -1 {
		// All vertices can accommodate remaining edges; take them all.
		if chosen+unknown > best {
			best = chosen + unknown
		}
		return
	}

	// Collect unknown incident edges of v.
	inc := make([]int, 0, degRem[v])
	for _, idx := range adj[v] {
		if state[idx] == 0 {
			inc = append(inc, idx)
		}
	}

	if cap[v] == 0 {
		// Must reject all incident unknown edges.
		newCap := append([]int(nil), cap...)
		newState := append([]int8(nil), state...)
		newDeg := append([]int(nil), degRem...)
		newUnknown := unknown
		for _, idx := range inc {
			e := gEdges[idx]
			u, w := e.u, e.v
			newState[idx] = -1
			newUnknown--
			newDeg[u]--
			newDeg[w]--
		}
		newDeg[v] = 0
		dfs(newCap, newState, newDeg, newUnknown, chosen, sumCap)
		return
	}

	// Helper to apply a choice of kept edges (subset) among inc.
	apply := func(keep []int) {
		newCap := append([]int(nil), cap...)
		newState := append([]int8(nil), state...)
		newDeg := append([]int(nil), degRem...)
		newUnknown := unknown
		newChosen := chosen
		newSumCap := sumCap

		keepSet := make(map[int]struct{}, len(keep))
		for _, id := range keep {
			keepSet[id] = struct{}{}
		}

		valid := true

		for _, idx := range inc {
			e := gEdges[idx]
			u, w := e.u, e.v
			newDeg[u]--
			newDeg[w]--
			newUnknown--
			if _, ok := keepSet[idx]; ok {
				if newCap[u] == 0 || newCap[w] == 0 {
					valid = false
					break
				}
				newCap[u]--
				newCap[w]--
				newSumCap -= 2
				newChosen++
				newState[idx] = 1
			} else {
				newState[idx] = -1
			}
		}

		if !valid {
			return
		}

		dfs(newCap, newState, newDeg, newUnknown, newChosen, newSumCap)
	}

	lenInc := len(inc)

	// Enumerate subsets of size up to cap[v] (cap <= 2).
	// Size 0
	apply([]int{})

	// Size 1
	if cap[v] >= 1 {
		for i := 0; i < lenInc; i++ {
			apply([]int{inc[i]})
		}
	}

	// Size 2
	if cap[v] >= 2 {
		for i := 0; i < lenInc; i++ {
			for j := i + 1; j < lenInc; j++ {
				apply([]int{inc[i], inc[j]})
			}
		}
	}
}

func solve() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var T int
	fmt.Fscan(in, &T)
	for ; T > 0; T-- {
		fmt.Fscan(in, &n, &m)
		gEdges = make([]edge, m)
		adj = make([][]int, n)
		for i := 0; i < m; i++ {
			var u, v int
			fmt.Fscan(in, &u, &v)
			u--
			v--
			gEdges[i] = edge{u, v}
			adj[u] = append(adj[u], i)
			adj[v] = append(adj[v], i)
		}

		cap := make([]int, n)
		for i := range cap {
			cap[i] = 2
		}
		state := make([]int8, m) // 0 unknown, 1 chosen, -1 rejected
		degRem := make([]int, n)
		for i := 0; i < n; i++ {
			degRem[i] = len(adj[i])
		}
		unknown := m
		chosen := 0
		sumCap := 2 * n
		best = 0

		dfs(cap, state, degRem, unknown, chosen, sumCap)

		fmt.Fprintln(out, best)
	}
}

func main() {
	solve()
}
