package main

import (
	"bufio"
	"fmt"
	"math/bits"
	"os"
)

// We model the plan as a functional graph with outdegree at most one (each city
// fires at most once). For every passenger (a -> b) the city b just needs to lie on
// the forward path that starts at a.
//
// Key observations
// ----------------
// - If the constraints force two cities to be mutually reachable, they must lie on
//   the same directed cycle. Moreover, such a cycle cannot point outside, so every
//   city reachable from that mutually-reachable set must also belong to the same
//   cycle.
// - After collapsing those forced cycles (and everything downstream of them) into
//   sink components, the remaining constraint graph is a DAG. On a DAG, a functional
//   graph corresponds to choosing at most one outgoing edge per component so that
//   every component’s required destinations lie on a single chain.
//
// Construction outline
// --------------------
// 1) Build the constraint digraph from passengers and compute SCCs.
// 2) In the condensation DAG, any SCC with size>1 forces every component in its
//    reachable closure to join the same cycle. We merge those closures with DSU.
// 3) Compress by these groups; forced-cycle groups become sinks. The remaining graph
//    is a DAG where each node needs to reach certain ancestors. We compute the
//    transitive closure on this DAG and apply the “chain” condition: for each node,
//    all reachable nodes must be pairwise comparable. The nearest ancestor becomes
//    its parent; nodes with empty reach stay as sinks.
// 4) Build the functional graph: each non-cycle node (represented by a single city)
//    points to its parent’s representative when it has one. Each cycle group is
//    wired into a directed cycle internally.
// 5) Verify that every passenger’s destination is reachable along the chosen edges.
//    Build precedence constraints along those paths and topologically sort the
//    launches; if there is a cycle, the plan is impossible.

type dsu struct {
	p   []int
	s   []int
	cyc []bool
}

func newDSU(n int) *dsu {
	p := make([]int, n)
	s := make([]int, n)
	cyc := make([]bool, n)
	for i := 0; i < n; i++ {
		p[i] = i
		s[i] = 1
	}
	return &dsu{p: p, s: s, cyc: cyc}
}

func (d *dsu) find(x int) int {
	if d.p[x] != x {
		d.p[x] = d.find(d.p[x])
	}
	return d.p[x]
}

func (d *dsu) unite(a, b int) {
	a = d.find(a)
	b = d.find(b)
	if a == b {
		return
	}
	if d.s[a] < d.s[b] {
		a, b = b, a
	}
	d.p[b] = a
	d.s[a] += d.s[b]
	d.cyc[a] = d.cyc[a] || d.cyc[b]
}

func main() {
	in := bufio.NewReader(os.Stdin)
	var n, m int
	if _, err := fmt.Fscan(in, &n, &m); err != nil {
		return
	}

	g := make([][]int, n)
	edges := make([][2]int, 0, m)
	for i := 0; i < m; i++ {
		var a, b int
		fmt.Fscan(in, &a, &b)
		a--
		b--
		g[a] = append(g[a], b)
		edges = append(edges, [2]int{a, b})
	}

	// Tarjan SCC.
	idx := 0
	st := make([]int, 0, n)
	on := make([]bool, n)
	tin := make([]int, n)
	comp := make([]int, n)
	for i := range tin {
		tin[i] = -1
	}
	comps := 0

	var dfs func(v int)
	low := make([]int, n)
	dfs = func(v int) {
		tin[v] = idx
		low[v] = idx
		idx++
		st = append(st, v)
		on[v] = true
		for _, to := range g[v] {
			if tin[to] == -1 {
				dfs(to)
				if low[to] < low[v] {
					low[v] = low[to]
				}
			} else if on[to] && tin[to] < low[v] {
				low[v] = tin[to]
			}
		}
		if low[v] == tin[v] {
			for {
				x := st[len(st)-1]
				st = st[:len(st)-1]
				on[x] = false
				comp[x] = comps
				if x == v {
					break
				}
			}
			comps++
		}
	}

	for i := 0; i < n; i++ {
		if tin[i] == -1 {
			dfs(i)
		}
	}

	compSize := make([]int, comps)
	for i := 0; i < n; i++ {
		compSize[comp[i]]++
	}

	compAdj := make([][]int, comps)
	for _, e := range edges {
		a := comp[e[0]]
		b := comp[e[1]]
		if a != b {
			compAdj[a] = append(compAdj[a], b)
		}
	}

	// Remove duplicate edges in compAdj.
	for i := 0; i < comps; i++ {
		if len(compAdj[i]) == 0 {
			continue
		}
		seen := make(map[int]struct{}, len(compAdj[i]))
		uniq := compAdj[i][:0]
		for _, v := range compAdj[i] {
			if _, ok := seen[v]; ok {
				continue
			}
			seen[v] = struct{}{}
			uniq = append(uniq, v)
		}
		compAdj[i] = uniq
	}

	// Topological order of condensation DAG.
	indeg := make([]int, comps)
	for u := 0; u < comps; u++ {
		for _, v := range compAdj[u] {
			indeg[v]++
		}
	}
	topo := make([]int, 0, comps)
	q := make([]int, 0, comps)
	for i := 0; i < comps; i++ {
		if indeg[i] == 0 {
			q = append(q, i)
		}
	}
	for len(q) > 0 {
		v := q[0]
		q = q[1:]
		topo = append(topo, v)
		for _, to := range compAdj[v] {
			indeg[to]--
			if indeg[to] == 0 {
				q = append(q, to)
			}
		}
	}
	if len(topo) != comps {
		fmt.Println(-1)
		return
	}

	// Reachability between components using bitsets.
	blk := (comps + 63) >> 6
	compReach := make([][]uint64, comps)
	for i := 0; i < comps; i++ {
		compReach[i] = make([]uint64, blk)
	}
	for i := comps - 1; i >= 0; i-- {
		v := topo[i]
		rv := compReach[v]
		for _, to := range compAdj[v] {
			rt := compReach[to]
			for b := 0; b < blk; b++ {
				rv[b] |= rt[b]
			}
			rv[to>>6] |= 1 << (uint(to) & 63)
		}
	}

	// DSU merging closures of cyclic components.
	d := newDSU(comps)
	for i := 0; i < comps; i++ {
		if compSize[i] > 1 {
			d.cyc[i] = true
			rv := compReach[i]
			for b := 0; b < blk; b++ {
				val := rv[b]
				if val == 0 {
					continue
				}
				base := b << 6
				for val != 0 {
					bit := bits.TrailingZeros64(val)
					v := base + bit
					d.unite(i, v)
					val &^= 1 << uint(bit)
				}
			}
			d.unite(i, i)
		}
	}

	// Compress components into groups (cycle groups as needed).
	groupID := make([]int, comps)
	idMap := make(map[int]int)
	groupNodes := make([][]int, 0)
	groupCyc := make([]bool, 0)
	for i := 0; i < comps; i++ {
		r := d.find(i)
		if idx, ok := idMap[r]; ok {
			groupID[i] = idx
		} else {
			groupID[i] = len(groupNodes)
			idMap[r] = groupID[i]
			groupNodes = append(groupNodes, nil)
			groupCyc = append(groupCyc, d.cyc[r])
		}
	}

	gCount := len(groupNodes)
	for city := 0; city < n; city++ {
		gIdx := groupID[comp[city]]
		groupNodes[gIdx] = append(groupNodes[gIdx], city)
	}

	// Build group graph.
	gAdj := make([][]int, gCount)
	deg := make([]int, gCount)
	for u := 0; u < comps; u++ {
		gu := groupID[u]
		for _, v := range compAdj[u] {
			gv := groupID[v]
			if gu != gv {
				gAdj[gu] = append(gAdj[gu], gv)
			}
		}
	}

	// Deduplicate edges and compute indegrees.
	for i := 0; i < gCount; i++ {
		if len(gAdj[i]) == 0 {
			continue
		}
		seen := make(map[int]struct{}, len(gAdj[i]))
		uniq := gAdj[i][:0]
		for _, v := range gAdj[i] {
			if _, ok := seen[v]; ok {
				continue
			}
			seen[v] = struct{}{}
			uniq = append(uniq, v)
		}
		gAdj[i] = uniq
	}
	deg = make([]int, gCount)
	for u := 0; u < gCount; u++ {
		for _, v := range gAdj[u] {
			deg[v]++
		}
	}

	// Cyclic groups must be sinks; otherwise impossible.
	for i := 0; i < gCount; i++ {
		if groupCyc[i] && len(gAdj[i]) > 0 {
			fmt.Println(-1)
			return
		}
	}

	// Topological order on group DAG.
	tq := make([]int, 0, gCount)
	topoG := make([]int, 0, gCount)
	for i := 0; i < gCount; i++ {
		if deg[i] == 0 {
			tq = append(tq, i)
		}
	}
	for len(tq) > 0 {
		v := tq[0]
		tq = tq[1:]
		topoG = append(topoG, v)
		for _, to := range gAdj[v] {
			deg[to]--
			if deg[to] == 0 {
				tq = append(tq, to)
			}
		}
	}
	if len(topoG) != gCount {
		fmt.Println(-1)
		return
	}

	// Reachability between groups.
	gBlk := (gCount + 63) >> 6
	gReach := make([][]uint64, gCount)
	for i := 0; i < gCount; i++ {
		gReach[i] = make([]uint64, gBlk)
	}
	for i := gCount - 1; i >= 0; i-- {
		v := topoG[i]
		rv := gReach[v]
		for _, to := range gAdj[v] {
			rt := gReach[to]
			for b := 0; b < gBlk; b++ {
				rv[b] |= rt[b]
			}
			rv[to>>6] |= 1 << (uint(to) & 63)
		}
	}

	gRev := make([][]uint64, gCount)
	for i := 0; i < gCount; i++ {
		gRev[i] = make([]uint64, gBlk)
	}
	for u := 0; u < gCount; u++ {
		ru := gReach[u]
		for b := 0; b < gBlk; b++ {
			val := ru[b]
			if val == 0 {
				continue
			}
			base := b << 6
			for val != 0 {
				bit := bits.TrailingZeros64(val)
				v := base + bit
				gRev[v][u>>6] |= 1 << (uint(u) & 63)
				val &^= 1 << uint(bit)
			}
		}
	}

	parent := make([]int, gCount)
	for i := range parent {
		parent[i] = -1
	}

	// Chain condition and parent selection for non-cyclic groups.
	for u := 0; u < gCount; u++ {
		if groupCyc[u] {
			// Should not need to reach outside its group.
			empty := true
			for _, v := range gReach[u] {
				if v != 0 {
					empty = false
					break
				}
			}
			if !empty {
				fmt.Println(-1)
				return
			}
			continue
		}

		R := gReach[u]
		isEmpty := true
		for _, v := range R {
			if v != 0 {
				isEmpty = false
				break
			}
		}
		if isEmpty {
			continue
		}

		// Check comparability.
		for blkIdx := 0; blkIdx < gBlk; blkIdx++ {
			val := R[blkIdx]
			if val == 0 {
				continue
			}
			base := blkIdx << 6
			for val != 0 {
				bit := bits.TrailingZeros64(val)
				v := base + bit
				val &^= 1 << uint(bit)

				comparable := make([]uint64, gBlk)
				for b := 0; b < gBlk; b++ {
					comparable[b] = gReach[v][b] | gRev[v][b]
				}
				comparable[v>>6] |= 1 << (uint(v) & 63)
				for b := 0; b < gBlk; b++ {
					if R[b]&^comparable[b] != 0 {
						fmt.Println(-1)
						return
					}
				}
			}
		}

		// Select closest ancestor.
		candidates := 0
		chosen := -1
		for blkIdx := 0; blkIdx < gBlk; blkIdx++ {
			val := R[blkIdx]
			if val == 0 {
				continue
			}
			base := blkIdx << 6
			for val != 0 {
				bit := bits.TrailingZeros64(val)
				v := base + bit
				val &^= 1 << uint(bit)

				hasLower := false
				for b := 0; b < gBlk; b++ {
					if R[b]&gRev[v][b] != 0 {
						hasLower = true
						break
					}
				}
				if hasLower {
					continue
				}
				candidates++
				chosen = v
			}
		}
		if candidates != 1 {
			fmt.Println(-1)
			return
		}
		parent[u] = chosen
	}

	// Build actual next-hop mapping for cities.
	next := make([]int, n)
	for i := range next {
		next[i] = -1
	}

	// Representatives for groups: pick first city.
	rep := make([]int, gCount)
	for i := 0; i < gCount; i++ {
		rep[i] = groupNodes[i][0]
	}

	// Cycle groups.
	for i := 0; i < gCount; i++ {
		if !groupCyc[i] {
			continue
		}
		nodes := groupNodes[i]
		if len(nodes) < 2 {
			fmt.Println(-1)
			return
		}
		for j := 0; j < len(nodes); j++ {
			next[nodes[j]] = nodes[(j+1)%len(nodes)]
		}
	}

	// Non-cycle edges.
	for i := 0; i < gCount; i++ {
		if groupCyc[i] {
			continue
		}
		if parent[i] != -1 {
			next[rep[i]] = rep[parent[i]]
		}
	}

	// Verify reachability for all passengers and collect precedence constraints.
	used := make([]bool, n)
	for i := 0; i < n; i++ {
		if next[i] != -1 {
			used[i] = true
		}
	}

	precAdj := make([][]int, n)
	precIndeg := make([]int, n)

	for _, e := range edges {
		a, b := e[0], e[1]
		cur := a
		seq := make([]int, 0)
		for steps := 0; steps <= n; steps++ {
			if cur == b {
				break
			}
			if next[cur] == -1 {
				fmt.Println(-1)
				return
			}
			seq = append(seq, cur)
			cur = next[cur]
		}
		if cur != b {
			fmt.Println(-1)
			return
		}
		for i := 0; i+1 < len(seq); i++ {
			u := seq[i]
			v := seq[i+1]
			precAdj[u] = append(precAdj[u], v)
		}
	}

	// Deduplicate precedence edges and compute indegrees only for used nodes.
	for i := 0; i < n; i++ {
		if len(precAdj[i]) == 0 {
			continue
		}
		seen := make(map[int]struct{}, len(precAdj[i]))
		uniq := precAdj[i][:0]
		for _, v := range precAdj[i] {
			if _, ok := seen[v]; ok {
				continue
			}
			seen[v] = struct{}{}
			uniq = append(uniq, v)
		}
		precAdj[i] = uniq
	}
	for u := 0; u < n; u++ {
		if !used[u] {
			continue
		}
		for _, v := range precAdj[u] {
			if used[v] {
				precIndeg[v]++
			}
		}
	}

	// Topological order of launches respecting precedence.
	launchOrder := make([]int, 0)
	lq := make([]int, 0)
	for i := 0; i < n; i++ {
		if used[i] && precIndeg[i] == 0 {
			lq = append(lq, i)
		}
	}
	for len(lq) > 0 {
		v := lq[0]
		lq = lq[1:]
		launchOrder = append(launchOrder, v)
		for _, to := range precAdj[v] {
			if !used[to] {
				continue
			}
			precIndeg[to]--
			if precIndeg[to] == 0 {
				lq = append(lq, to)
			}
		}
	}

	totalEdges := 0
	for i := 0; i < n; i++ {
		if used[i] {
			totalEdges++
		}
	}

	if len(launchOrder) != totalEdges {
		fmt.Println(-1)
		return
	}

	out := bufio.NewWriter(os.Stdout)
	fmt.Fprintln(out, totalEdges)
	for _, v := range launchOrder {
		fmt.Fprintf(out, "%d %d\n", v+1, next[v]+1)
	}
	out.Flush()
}
