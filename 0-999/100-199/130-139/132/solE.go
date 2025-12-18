package main

import (
	"fmt"
	"math/bits"
)

const INF = int64(1e18)

type Edge struct {
	to   int
	rev  int
	cap  int
	flow int
	cost int64
}

type MCMF struct {
	n          int
	graph      [][]Edge
	dist       []int64
	parentEdge []int
	parentNode []int
}

func NewMCMF(n int) *MCMF {
	return &MCMF{
		n:     n,
		graph: make([][]Edge, n),
	}
}

func (mc *MCMF) AddEdge(u, v, cap int, cost int64) {
	mc.graph[u] = append(mc.graph[u], Edge{to: v, rev: len(mc.graph[v]), cap: cap, cost: cost})
	mc.graph[v] = append(mc.graph[v], Edge{to: u, rev: len(mc.graph[u]) - 1, cap: 0, cost: -cost})
}

func (mc *MCMF) SPFA(s, t int) bool {
	mc.dist = make([]int64, mc.n)
	mc.parentNode = make([]int, mc.n)
	mc.parentEdge = make([]int, mc.n)
	for i := range mc.dist {
		mc.dist[i] = INF
	}
	inQueue := make([]bool, mc.n)
	queue := []int{s}
	mc.dist[s] = 0
	inQueue[s] = true

	qIdx := 0
	for qIdx < len(queue) {
		u := queue[qIdx]
		qIdx++
		inQueue[u] = false

		for i := range mc.graph[u] {
			e := &mc.graph[u][i]
			if e.cap-e.flow > 0 && mc.dist[e.to] > mc.dist[u]+e.cost {
				mc.dist[e.to] = mc.dist[u] + e.cost
				mc.parentNode[e.to] = u
				mc.parentEdge[e.to] = i
				if !inQueue[e.to] {
					queue = append(queue, e.to)
					inQueue[e.to] = true
				}
			}
		}
	}
	return mc.dist[t] != INF
}

func (mc *MCMF) Solve(s, t, minFlow int) (int, int64) {
	totalFlow := 0
	var totalCost int64 = 0

	for mc.SPFA(s, t) {
		if totalFlow >= minFlow && mc.dist[t] >= 0 {
			break
		}

		flow := int(1e9)
		curr := t
		for curr != s {
			p := mc.parentNode[curr]
			idx := mc.parentEdge[curr]
			if rem := mc.graph[p][idx].cap - mc.graph[p][idx].flow; rem < flow {
				flow = rem
			}
			curr = p
		}

		totalFlow += flow
		totalCost += int64(flow) * mc.dist[t]
		curr = t
		for curr != s {
			p := mc.parentNode[curr]
			idx := mc.parentEdge[curr]
			mc.graph[p][idx].flow += flow
			rev := mc.graph[p][idx].rev
			mc.graph[curr][rev].flow -= flow
			curr = p
		}
	}
	return totalFlow, totalCost
}

func main() {
	var n, m int
	if _, err := fmt.Scan(&n, &m); err != nil {
		return
	}
	a := make([]int, n)
	var baseCost int64
	for i := 0; i < n; i++ {
		fmt.Scan(&a[i])
		baseCost += int64(bits.OnesCount(uint(a[i])))
	}

	S := 0
	T := 2*n + 1
	mc := NewMCMF(2*n + 2)

	for i := 0; i < n; i++ {
		mc.AddEdge(S, i+1, 1, 0)
		mc.AddEdge(n+1+i, T, 1, 0)
		for j := i + 1; j < n; j++ {
			var cost int64 = 0
			if a[i] == a[j] {
				cost = -int64(bits.OnesCount(uint(a[j])))
			}
			mc.AddEdge(i+1, n+1+j, 1, cost)
		}
	}

	minFlow := 0
	if n > m {
		minFlow = n - m
	}

	_, reduction := mc.Solve(S, T, minFlow)
	finalPenalty := baseCost + reduction

	match := make([]int, n)
	for i := range match {
		match[i] = -1
	}
	hasIncoming := make([]bool, n)

	for u := 1; u <= n; u++ {
		for _, e := range mc.graph[u] {
			if e.to > n && e.to <= 2*n && e.flow == 1 && e.cap == 1 {
				v := e.to - (n + 1)
				match[u-1] = v
			hasIncoming[v] = true
			}
		}
	}

	chainID := make([]int, n)
	currentChain := 0
	for i := 0; i < n; i++ {
		if !hasIncoming[i] {
			curr := i
			for curr != -1 {
				chainID[curr] = currentChain
				curr = match[curr]
			}
			currentChain++
		}
	}

	var program []string
	varVals := make(map[int]int)

	for i := 0; i < n; i++ {
		cid := chainID[i]
		vName := string('a' + byte(cid))
		val := a[i]

		currVal, ok := varVals[cid]
		if !ok || currVal != val {
			program = append(program, fmt.Sprintf("%s=%d", vName, val))
			varVals[cid] = val
		}
		program = append(program, fmt.Sprintf("print(%s)", vName))
	}

	fmt.Printf("%d %d\n", len(program), finalPenalty)
	for _, line := range program {
		fmt.Println(line)
	}
}