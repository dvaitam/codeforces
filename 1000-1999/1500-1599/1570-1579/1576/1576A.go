package main

import (
	"bufio"
	"fmt"
	"os"
)

type Edge struct {
	id    int
	group int
	u     int
	v     int
	dist  int
	cap   int
}

func main() {
	in := bufio.NewReader(os.Stdin)

	var nodeCount, edgeCount, constrCount, flowCount int
	if _, err := fmt.Fscan(in, &nodeCount, &edgeCount, &constrCount, &flowCount); err != nil {
		return
	}

	edges := make([]Edge, edgeCount)
	for i := 0; i < edgeCount; i++ {
		fmt.Fscan(in, &edges[i].id, &edges[i].group, &edges[i].u, &edges[i].v, &edges[i].dist, &edges[i].cap)
	}

	// skip constrained pairs
	for i := 0; i < constrCount; i++ {
		var a, b, c int
		fmt.Fscan(in, &a, &b, &c)
	}

	type Flow struct {
		id   int
		s    int
		t    int
		rate int
	}
	flows := make([]Flow, flowCount)
	for i := 0; i < flowCount; i++ {
		fmt.Fscan(in, &flows[i].id, &flows[i].s, &flows[i].t, &flows[i].rate)
	}

	// Build adjacency from source to dest edges
	// Map pair (u,v) -> list of edge indices
	adj := make(map[[2]int][]int)
	for i, e := range edges {
		adj[[2]int{e.u, e.v}] = append(adj[[2]int{e.u, e.v}], i)
		adj[[2]int{e.v, e.u}] = append(adj[[2]int{e.v, e.u}], i)
	}

	const SFL = 200
	const GFL = 100

	nodeFlowCount := make([]int, nodeCount)
	groupFlowCount := make(map[int]int)
	edgeCapLeft := make([]int, edgeCount)
	for i, e := range edges {
		edgeCapLeft[i] = e.cap
	}

	type Path struct {
		flowID int
		edgeID int
	}
	var results []Path

	for _, f := range flows {
		if nodeFlowCount[f.s] >= SFL || nodeFlowCount[f.t] >= SFL {
			continue
		}
		candidates := adj[[2]int{f.s, f.t}]
		ok := false
		for _, ei := range candidates {
			e := edges[ei]
			if edgeCapLeft[ei] < f.rate {
				continue
			}
			if groupFlowCount[e.group] >= GFL {
				continue
			}
			// Accept this edge
			edgeCapLeft[ei] -= f.rate
			nodeFlowCount[f.s]++
			nodeFlowCount[f.t]++
			groupFlowCount[e.group]++
			results = append(results, Path{flowID: f.id, edgeID: e.id})
			ok = true
			break
		}
		if !ok {
			// skip this flow
		}
	}

	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()
	fmt.Fprintln(out, len(results))
	for _, p := range results {
		fmt.Fprintf(out, "%d %d\n", p.flowID, p.edgeID)
	}
}
