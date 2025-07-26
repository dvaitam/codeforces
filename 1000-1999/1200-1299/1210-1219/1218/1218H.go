package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

type NodeEvent struct {
	depth   int
	rootPos int
}

type Query struct {
	m     int64
	yPos  int
	index int
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n int
	if _, err := fmt.Fscan(in, &n); err != nil {
		return
	}
	to := make([]int, n+1)
	rev := make([][]int, n+1)
	indeg := make([]int, n+1)
	for i := 1; i <= n; i++ {
		fmt.Fscan(in, &to[i])
		rev[to[i]] = append(rev[to[i]], i)
		indeg[to[i]]++
	}

	// find cycle nodes using Kahn's algorithm
	q := make([]int, 0)
	for i := 1; i <= n; i++ {
		if indeg[i] == 0 {
			q = append(q, i)
		}
	}
	for front := 0; front < len(q); front++ {
		v := q[front]
		u := to[v]
		indeg[u]--
		if indeg[u] == 0 {
			q = append(q, u)
		}
	}
	onCycle := make([]bool, n+1)
	for i := 1; i <= n; i++ {
		if indeg[i] > 0 {
			onCycle[i] = true
		}
	}

	// extract cycles and assign positions
	cycleID := make([]int, n+1)
	pos := make([]int, n+1)
	var cycleLen []int
	var cycles [][]int
	id := 0
	visited := make([]bool, n+1)
	for i := 1; i <= n; i++ {
		if onCycle[i] && !visited[i] {
			id++
			nodes := []int{}
			cur := i
			for {
				visited[cur] = true
				cycleID[cur] = id
				pos[cur] = len(nodes)
				nodes = append(nodes, cur)
				cur = to[cur]
				if cur == i {
					break
				}
			}
			cycles = append(cycles, nodes)
			cycleLen = append(cycleLen, len(nodes))
		}
	}
	compCount := id

	// prepare structures per component
	eventsByComp := make([][]NodeEvent, compCount+1)

	depth := make([]int, n+1)
	root := make([]int, n+1)
	comp := make([]int, n+1)
	tin := make([]int, n+1)
	tout := make([]int, n+1)
	timer := 0

	depthTin := map[int][]int{}

	type frame struct{ v, idx int }
	for cid := 1; cid <= compCount; cid++ {
		nodes := cycles[cid-1]
		L := len(nodes)
		for j := 0; j < L; j++ {
			start := nodes[j]
			root[start] = start
			comp[start] = cid
			depth[start] = 0
			st := []frame{{start, 0}}
			for len(st) > 0 {
				top := &st[len(st)-1]
				v := top.v
				if top.idx == 0 {
					tin[v] = timer
					timer++
					eventsByComp[cid] = append(eventsByComp[cid], NodeEvent{depth: depth[v], rootPos: j})
					depthTin[depth[v]] = append(depthTin[depth[v]], tin[v])
				}
				if top.idx < len(rev[v]) {
					u := rev[v][top.idx]
					top.idx++
					if onCycle[u] {
						continue
					}
					root[u] = start
					comp[u] = cid
					depth[u] = depth[v] + 1
					st = append(st, frame{u, 0})
				} else {
					tout[v] = timer
					st = st[:len(st)-1]
				}
			}
		}
	}

	maxDepth := 0
	for d := range depthTin {
		if d > maxDepth {
			maxDepth = d
		}
	}
	for d := range depthTin {
		sort.Ints(depthTin[d])
	}

	var Q int
	fmt.Fscan(in, &Q)
	ans := make([]int, Q)
	queriesByComp := make([][]Query, compCount+1)

	for idx := 0; idx < Q; idx++ {
		var m int64
		var y int
		fmt.Fscan(in, &m, &y)
		if onCycle[y] {
			cid := comp[y]
			queriesByComp[cid] = append(queriesByComp[cid], Query{m: m, yPos: pos[y], index: idx})
		} else {
			d := int64(depth[y]) + m
			if d <= int64(maxDepth) {
				D := int(d)
				vec := depthTin[D]
				l := sort.SearchInts(vec, tin[y])
				r := sort.SearchInts(vec, tout[y])
				ans[idx] = r - l
			} else {
				ans[idx] = 0
			}
		}
	}

	// process cycle queries per component
	for cid := 1; cid <= compCount; cid++ {
		qs := queriesByComp[cid]
		if len(qs) == 0 {
			continue
		}
		ev := eventsByComp[cid]
		sort.Slice(ev, func(i, j int) bool { return ev[i].depth < ev[j].depth })
		sort.Slice(qs, func(i, j int) bool { return qs[i].m < qs[j].m })
		L := cycleLen[cid-1]
		rotSum := make([]int, L)
		idxEv := 0
		for _, qv := range qs {
			for idxEv < len(ev) && int64(ev[idxEv].depth) <= qv.m {
				e := ev[idxEv]
				k := (e.depth%L - e.rootPos + L) % L
				rotSum[k]++
				idxEv++
			}
			km := int((qv.m - int64(qv.yPos)) % int64(L))
			if km < 0 {
				km += L
			}
			ans[qv.index] = rotSum[km]
		}
	}

	for i := 0; i < Q; i++ {
		fmt.Fprintln(out, ans[i])
	}
}
