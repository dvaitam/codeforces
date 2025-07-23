package main

import (
	"bufio"
	"fmt"
	"os"
)

type Edge struct {
	to int
	id int
}

type EdgeData struct {
	u, v int
	w    int
}

var (
	n, m  int
	s, t  int
	g     [][]Edge
	edges []EdgeData
)

func bfs(start int, skip []bool) []bool {
	vis := make([]bool, n+1)
	q := make([]int, 0, n)
	vis[start] = true
	q = append(q, start)
	for len(q) > 0 {
		v := q[0]
		q = q[1:]
		for _, e := range g[v] {
			if skip[e.id] {
				continue
			}
			if !vis[e.to] {
				vis[e.to] = true
				q = append(q, e.to)
			}
		}
	}
	return vis
}

func findBridges(skip []bool) []bool {
	isBridge := make([]bool, m+1)
	visited := make([]bool, n+1)
	tin := make([]int, n+1)
	low := make([]int, n+1)
	timer := 0

	var dfs func(v, peid int)
	dfs = func(v, peid int) {
		visited[v] = true
		tin[v] = timer
		low[v] = timer
		timer++
		for _, e := range g[v] {
			if skip[e.id] || e.id == peid {
				continue
			}
			if visited[e.to] {
				if tin[e.to] < low[v] {
					low[v] = tin[e.to]
				}
			} else {
				dfs(e.to, e.id)
				if low[e.to] < low[v] {
					low[v] = low[e.to]
				}
				if low[e.to] > tin[v] {
					isBridge[e.id] = true
				}
			}
		}
	}

	dfs(s, -1)
	return isBridge
}

func onPath(visS, visT []bool, id int) bool {
	u := edges[id].u
	v := edges[id].v
	return (visS[u] && visT[v]) || (visS[v] && visT[u])
}

func main() {
	in := bufio.NewReader(os.Stdin)
	fmt.Fscan(in, &n, &m)
	fmt.Fscan(in, &s, &t)
	g = make([][]Edge, n+1)
	edges = make([]EdgeData, m+1)
	for i := 1; i <= m; i++ {
		var u, v, w int
		fmt.Fscan(in, &u, &v, &w)
		edges[i] = EdgeData{u, v, w}
		g[u] = append(g[u], Edge{to: v, id: i})
		g[v] = append(g[v], Edge{to: u, id: i})
	}

	skip := make([]bool, m+1)
	// initial connectivity check
	visS0 := bfs(s, skip)
	if !visS0[t] {
		fmt.Println(0)
		fmt.Println(0)
		return
	}
	visT0 := bfs(t, skip)

	// bridges in original graph
	bridges := findBridges(skip)
	bestCost := int64(1<<63 - 1)
	bestEdges := []int{}
	for i := 1; i <= m; i++ {
		if bridges[i] && onPath(visS0, visT0, i) {
			if int64(edges[i].w) < bestCost {
				bestCost = int64(edges[i].w)
				bestEdges = []int{i}
			}
		}
	}

	candidateEdges := make([]int, 0)
	for i := 1; i <= m; i++ {
		if onPath(visS0, visT0, i) && !bridges[i] {
			candidateEdges = append(candidateEdges, i)
		}
	}

	visS := make([]bool, 0)
	visT := make([]bool, 0)
	for _, id1 := range candidateEdges {
		if int64(edges[id1].w) >= bestCost {
			continue
		}
		skip[id1] = true
		visS = bfs(s, skip)
		if !visS[t] {
			if int64(edges[id1].w) < bestCost {
				bestCost = int64(edges[id1].w)
				bestEdges = []int{id1}
			}
			skip[id1] = false
			continue
		}
		visT = bfs(t, skip)
		bridges2 := findBridges(skip)
		minSecond := int64(1<<63 - 1)
		minId2 := -1
		for i := 1; i <= m; i++ {
			if skip[i] || !bridges2[i] {
				continue
			}
			if onPath(visS, visT, i) {
				if int64(edges[i].w) < minSecond {
					minSecond = int64(edges[i].w)
					minId2 = i
				}
			}
		}
		if minId2 != -1 {
			total := int64(edges[id1].w) + minSecond
			if total < bestCost {
				bestCost = total
				if id1 < minId2 {
					bestEdges = []int{id1, minId2}
				} else {
					bestEdges = []int{minId2, id1}
				}
			}
		}
		skip[id1] = false
	}

	if bestCost == int64(1<<63-1) {
		fmt.Println(-1)
		return
	}

	fmt.Println(bestCost)
	fmt.Println(len(bestEdges))
	if len(bestEdges) > 0 {
		for i, id := range bestEdges {
			if i > 0 {
				fmt.Print(" ")
			}
			fmt.Print(id)
		}
		fmt.Println()
	}
}
