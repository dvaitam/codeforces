package main

import (
	"bufio"
	"fmt"
	"os"
)

const inf = int(1e9)

type edge struct {
	to    int
	rev   int
	cap   int
	isRev bool
}

type segNode struct {
	l, r        int
	left, right int
}

type flowEdge struct {
	to  int
	cap int
}

type pair struct {
	weapon, ship int
}

type triple struct {
	a, b, c int
}

var (
	graph  [][]edge
	level  []int
	ptr    []int
	seg    []segNode
	segCnt int
)

func addEdge(u, v, cap int) {
	forward := edge{to: v, rev: len(graph[v]), cap: cap}
	graph[u] = append(graph[u], forward)
	backward := edge{to: u, rev: len(graph[u]) - 1, cap: 0, isRev: true}
	graph[v] = append(graph[v], backward)
}

func build(idx, l, r int) {
	seg[idx].l = l
	seg[idx].r = r
	if l == r {
		return
	}
	mid := (l + r) >> 1
	seg[idx].left = segCnt + 1
	segCnt++
	build(seg[idx].left, l, mid)
	seg[idx].right = segCnt + 1
	segCnt++
	build(seg[idx].right, mid+1, r)
	addEdge(idx, seg[idx].left, inf)
	addEdge(idx, seg[idx].right, inf)
}

func findLeaf(idx, pos int) int {
	if seg[idx].l == seg[idx].r {
		return idx
	}
	if pos <= seg[seg[idx].left].r {
		return findLeaf(seg[idx].left, pos)
	}
	return findLeaf(seg[idx].right, pos)
}

func connectInterval(nodeID, idx, L, R int) {
	if idx == 0 {
		return
	}
	if R < seg[idx].l || L > seg[idx].r {
		return
	}
	if L <= seg[idx].l && seg[idx].r <= R {
		addEdge(nodeID, idx, 1)
		return
	}
	connectInterval(nodeID, seg[idx].left, L, R)
	connectInterval(nodeID, seg[idx].right, L, R)
}

func bfs(S, T int) bool {
	for i := range level {
		level[i] = -1
	}
	queue := []int{S}
	level[S] = 0
	for len(queue) > 0 {
		v := queue[0]
		queue = queue[1:]
		for _, e := range graph[v] {
			if e.cap > 0 && level[e.to] == -1 {
				level[e.to] = level[v] + 1
				queue = append(queue, e.to)
			}
		}
	}
	return level[T] != -1
}

func dfs(v, T, pushed int) int {
	if pushed == 0 {
		return 0
	}
	if v == T {
		return pushed
	}
	for ptr[v] < len(graph[v]) {
		i := ptr[v]
		e := &graph[v][i]
		if e.cap > 0 && level[e.to] == level[v]+1 {
			minVal := pushed
			if e.cap < minVal {
				minVal = e.cap
			}
			tr := dfs(e.to, T, minVal)
			if tr > 0 {
				e.cap -= tr
				graph[e.to][e.rev].cap += tr
				return tr
			}
		}
		ptr[v]++
	}
	return 0
}

func dinic(S, T int) int {
	flow := 0
	for bfs(S, T) {
		for i := range ptr {
			ptr[i] = 0
		}
		for {
			pushed := dfs(S, T, inf)
			if pushed == 0 {
				break
			}
			flow += pushed
		}
	}
	return flow
}

func dfsAssignment(u, S, T int, visited []bool, flowGraph [][]*flowEdge, x, y *int) bool {
	if u == T {
		return true
	}
	visited[u] = true
	for _, e := range flowGraph[u] {
		if e.cap == 0 || visited[e.to] {
			continue
		}
		if dfsAssignment(e.to, S, T, visited, flowGraph, x, y) {
			e.cap--
			if e.to == T {
				*y = u
			}
			if u == S {
				*x = e.to
			}
			visited[u] = false
			return true
		}
	}
	visited[u] = false
	return false
}

func nextInt(r *bufio.Reader) int {
	sign, val := 1, 0
	c, _ := r.ReadByte()
	for (c < '0' || c > '9') && c != '-' {
		c, _ = r.ReadByte()
	}
	if c == '-' {
		sign = -1
		c, _ = r.ReadByte()
	}
	for c >= '0' && c <= '9' {
		val = val*10 + int(c-'0')
		c, _ = r.ReadByte()
	}
	return sign * val
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	n := nextInt(reader)
	m := nextInt(reader)

	if m == 0 {
		fmt.Fprintln(writer, 0)
		return
	}

	maxSeg := 4*m + 5
	maxNodes := maxSeg + n + 50
	graph = make([][]edge, maxNodes)
	seg = make([]segNode, maxSeg)
	segCnt = 1
	build(1, 1, m)
	baseWeapon := segCnt
	S := baseWeapon + n + 1
	T := S + 1

	level = make([]int, T+1)
	ptr = make([]int, T+1)

	for pos := 1; pos <= m; pos++ {
		leaf := findLeaf(1, pos)
		addEdge(leaf, T, 1)
	}

	bazooka := make([]triple, n+1)
	isBazooka := make([]bool, n+1)

	for i := 1; i <= n; i++ {
		typ := nextInt(reader)
		nodeID := baseWeapon + i
		switch typ {
		case 0:
			k := nextInt(reader)
			addEdge(S, nodeID, 1)
			for j := 0; j < k; j++ {
				target := nextInt(reader)
				leaf := findLeaf(1, target)
				addEdge(nodeID, leaf, 1)
			}
		case 1:
			l := nextInt(reader)
			r := nextInt(reader)
			addEdge(S, nodeID, 1)
			connectInterval(nodeID, 1, l, r)
		case 2:
			a := nextInt(reader)
			b := nextInt(reader)
			c := nextInt(reader)
			bazooka[i] = triple{a: a, b: b, c: c}
			isBazooka[i] = true
			addEdge(S, nodeID, 2)
			leafA := findLeaf(1, a)
			leafB := findLeaf(1, b)
			leafC := findLeaf(1, c)
			addEdge(nodeID, leafA, 1)
			addEdge(nodeID, leafB, 1)
			addEdge(nodeID, leafC, 1)
		}
	}

	nodeCount := T
	graph = graph[:nodeCount+1]

	maxFlow := dinic(S, T)

	flowGraph := make([][]*flowEdge, nodeCount+1)
	for u := 1; u <= nodeCount; u++ {
		for idx := range graph[u] {
			e := graph[u][idx]
			if e.isRev && e.cap > 0 {
				flowGraph[e.to] = append(flowGraph[e.to], &flowEdge{to: u, cap: e.cap})
			}
		}
	}

	results := make([]pair, 0, maxFlow)
	visited := make([]bool, nodeCount+1)
	for i := 0; i < maxFlow; i++ {
		for j := range visited {
			visited[j] = false
		}
		x, y := -1, -1
		if !dfsAssignment(S, S, T, visited, flowGraph, &x, &y) || x == -1 || y == -1 {
			break
		}
		weaponIdx := x - baseWeapon
		if weaponIdx < 1 || weaponIdx > n {
			continue
		}
		ship := seg[y].l
		results = append(results, pair{weapon: weaponIdx, ship: ship})
	}

	for i := 1; i <= n; i++ {
		if !isBazooka[i] {
			continue
		}
		cnt := 0
		for _, res := range results {
			if res.weapon == i {
				cnt++
			}
		}
		if cnt == 1 {
			targets := bazooka[i]
			for idx := range results {
				ship := results[idx].ship
				if ship == targets.a || ship == targets.b || ship == targets.c {
					if results[idx].weapon != i {
						results[idx].weapon = i
						break
					}
				}
			}
		}
	}

	fmt.Fprintln(writer, len(results))
	for _, res := range results {
		fmt.Fprintf(writer, "%d %d\n", res.weapon, res.ship)
	}
}
