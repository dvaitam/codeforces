package main

import (
	"bufio"
	"fmt"
	"os"
)

type Edge struct {
	u int
	v int
	x int64
}

type Adj struct {
	to int
	id int
}

type TreeAdj struct {
	to    int
	chord int
}

type Operation struct {
	val int64
	typ int // 0 => base tree, 1 => replace tree edge
	rem int
	add int
}

var (
	n, m             int
	mod              int64
	edges            []Edge
	g                [][]Adj
	isBridge         []bool
	timer            int
	tin              []int
	low              []int
	inTree           []bool
	treeEdges        []int
	treeIdx          []int
	parent           []int
	parentEdge       []int
	depth            []int
	compID           []int
	compSum          []int64
	compW            []int64
	operations       []Operation
	assignedSum      []int64
	sumChords        int64
	diff             []int64
	adjTree          [][]TreeAdj
	treeEdgeIsBridge []bool
	modInvCache      map[int64]int64
)

func norm(x int64) int64 {
	x %= mod
	if x < 0 {
		x += mod
	}
	return x
}

func addOperation(typ int, rem, add int, val int64) {
	val = norm(val)
	if val == 0 {
		return
	}
	operations = append(operations, Operation{val: val, typ: typ, rem: rem, add: add})
}

func dfsBridge(v, pe int) {
	timer++
	tin[v] = timer
	low[v] = timer
	for _, e := range g[v] {
		if e.id == pe {
			continue
		}
		to := e.to
		if tin[to] == 0 {
			dfsBridge(to, e.id)
			if low[to] > tin[v] {
				isBridge[e.id] = true
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

func buildTree() {
	visited := make([]bool, n+1)
	stack := []int{1}
	visited[1] = true
	parent[1] = 0
	parentEdge[1] = 0
	depth[1] = 0
	for len(stack) > 0 {
		v := stack[len(stack)-1]
		stack = stack[:len(stack)-1]
		for _, e := range g[v] {
			to := e.to
			if visited[to] {
				continue
			}
			visited[to] = true
			parent[to] = v
			parentEdge[to] = e.id
			depth[to] = depth[v] + 1
			inTree[e.id] = true
			treeEdges = append(treeEdges, e.id)
			stack = append(stack, to)
		}
	}
}

func getPathEdges(u, v int) []int {
	pathU := make([]int, 0)
	pathV := make([]int, 0)
	uu := u
	vv := v
	for depth[uu] > depth[vv] {
		pathU = append(pathU, parentEdge[uu])
		uu = parent[uu]
	}
	for depth[vv] > depth[uu] {
		pathV = append(pathV, parentEdge[vv])
		vv = parent[vv]
	}
	for uu != vv {
		pathU = append(pathU, parentEdge[uu])
		pathV = append(pathV, parentEdge[vv])
		uu = parent[uu]
		vv = parent[vv]
	}
	res := make([]int, 0, len(pathU)+len(pathV))
	res = append(res, pathU...)
	for i := len(pathV) - 1; i >= 0; i-- {
		res = append(res, pathV[i])
	}
	return res
}

func dfsAssignComponents() int {
	compID = make([]int, m+1)
	for i := 1; i <= m; i++ {
		compID[i] = -1
	}
	compCnt := 0
	visited := make([]bool, n+1)
	for v := 1; v <= n; v++ {
		if visited[v] {
			continue
		}
		stack := []int{v}
		visited[v] = true
		hasEdge := false
		for len(stack) > 0 {
			u := stack[len(stack)-1]
			stack = stack[:len(stack)-1]
			for _, e := range g[u] {
				if isBridge[e.id] {
					continue
				}
				if compID[e.id] == -1 {
					compID[e.id] = compCnt
					hasEdge = true
				}
				if !visited[e.to] {
					visited[e.to] = true
					stack = append(stack, e.to)
				}
			}
		}
		if hasEdge {
			compCnt++
		}
	}
	for i := 1; i <= m; i++ {
		if compID[i] == -1 {
			compID[i] = compCnt
			compCnt++
		}
	}
	return compCnt
}

func modPow(a, e int64) int64 {
	result := int64(1)
	base := norm(a)
	exp := e
	for exp > 0 {
		if exp&1 == 1 {
			result = (result * base) % mod
		}
		base = (base * base) % mod
		exp >>= 1
	}
	return result
}

func modInv(a int64) int64 {
	if val, ok := modInvCache[a]; ok {
		return val
	}
	res := modPow(a, mod-2)
	modInvCache[a] = res
	return res
}

func applyTransfer(fromIdx, toIdx int, chordID int, delta int64) {
	delta = norm(delta)
	if delta == 0 {
		return
	}
	fromEdgeID := treeEdges[fromIdx]
	toEdgeID := treeEdges[toIdx]
	addOperation(1, fromEdgeID, chordID, delta)
	addOperation(1, toEdgeID, chordID, -delta)
	diff[fromIdx] = norm(diff[fromIdx] - delta)
	diff[toIdx] = norm(diff[toIdx] + delta)
}

func dfsBalance(cur, parentIdx int, visited []bool) int64 {
	visited[cur] = true
	total := diff[cur]
	for _, adj := range adjTree[cur] {
		if adj.to == parentIdx {
			continue
		}
		if visited[adj.to] {
			continue
		}
		child := dfsBalance(adj.to, cur, visited)
		if child != 0 {
			applyTransfer(adj.to, cur, adj.chord, child)
			total = norm(total + child)
		}
	}
	diff[cur] = total
	return total
}

func main() {
	in := bufio.NewReader(os.Stdin)
	if _, err := fmt.Fscan(in, &n, &m, &mod); err != nil {
		return
	}
	edges = make([]Edge, m+1)
	g = make([][]Adj, n+1)
	for i := 1; i <= m; i++ {
		var u, v int
		var x int64
		fmt.Fscan(in, &u, &v, &x)
		edges[i] = Edge{u: u, v: v, x: norm(x)}
		g[u] = append(g[u], Adj{to: v, id: i})
		g[v] = append(g[v], Adj{to: u, id: i})
	}

	isBridge = make([]bool, m+1)
	tin = make([]int, n+1)
	low = make([]int, n+1)
	timer = 0
	dfsBridge(1, -1)

	inTree = make([]bool, m+1)
	treeEdges = make([]int, 0, n-1)
	parent = make([]int, n+1)
	parentEdge = make([]int, n+1)
	depth = make([]int, n+1)
	buildTree()
	if len(treeEdges) != n-1 {
		fmt.Println(-1)
		return
	}

	treeIdx = make([]int, m+1)
	for i := 1; i <= m; i++ {
		treeIdx[i] = -1
	}
	for idx, eid := range treeEdges {
		treeIdx[eid] = idx
	}

	compCnt := dfsAssignComponents()
	compSum = make([]int64, compCnt)
	compW = make([]int64, compCnt)
	for i := 1; i <= m; i++ {
		comp := compID[i]
		compSum[comp] = norm(compSum[comp] + edges[i].x)
		if inTree[i] {
			compW[comp]++
		}
	}

	modInvCache = make(map[int64]int64)
	lambda := int64(0)
	setLambda := false

	for comp := 0; comp < compCnt; comp++ {
		w := compW[comp] % mod
		s := compSum[comp]
		if w == 0 {
			if s != 0 {
				fmt.Println(-1)
				return
			}
			continue
		}
		candidate := (s * modInv(w)) % mod
		if !setLambda {
			lambda = candidate
			setLambda = true
		} else if lambda != candidate {
			fmt.Println(-1)
			return
		}
	}
	if !setLambda {
		lambda = 0
	}

	assignedSum = make([]int64, len(treeEdges))
	sumChords = 0
	operations = make([]Operation, 0)
	adjTree = make([][]TreeAdj, len(treeEdges))
	for i := range adjTree {
		adjTree[i] = make([]TreeAdj, 0)
	}

	for edgeID := 1; edgeID <= m; edgeID++ {
		if inTree[edgeID] {
			continue
		}
		path := getPathEdges(edges[edgeID].u, edges[edgeID].v)
		if len(path) == 0 {
			fmt.Println(-1)
			return
		}
		assignEdge := path[0]
		assignIdx := treeIdx[assignEdge]
		assignedSum[assignIdx] = norm(assignedSum[assignIdx] + edges[edgeID].x)
		addOperation(1, assignEdge, edgeID, edges[edgeID].x)
		sumChords = norm(sumChords + edges[edgeID].x)
		for i := 0; i+1 < len(path); i++ {
			a := treeIdx[path[i]]
			b := treeIdx[path[i+1]]
			if a == -1 || b == -1 {
				continue
			}
			adjTree[a] = append(adjTree[a], TreeAdj{to: b, chord: edgeID})
			adjTree[b] = append(adjTree[b], TreeAdj{to: a, chord: edgeID})
		}
	}

	addOperation(0, 0, 0, -sumChords)

	diff = make([]int64, len(treeEdges))
	treeEdgeIsBridge = make([]bool, len(treeEdges))
	for idx, eid := range treeEdges {
		treeEdgeIsBridge[idx] = isBridge[eid]
		cur := norm(-assignedSum[idx])
		target := edges[eid].x
		desired := norm(target - lambda)
		diff[idx] = norm(desired - cur)
		if treeEdgeIsBridge[idx] && diff[idx] != 0 {
			fmt.Println(-1)
			return
		}
	}

	visited := make([]bool, len(treeEdges))
	for idx := range treeEdges {
		if treeEdgeIsBridge[idx] || visited[idx] {
			continue
		}
		if dfsBalance(idx, -1, visited) != 0 {
			fmt.Println(-1)
			return
		}
	}

	addOperation(0, 0, 0, lambda)

	if len(operations) > 2*m {
		fmt.Println(-1)
		return
	}

	writer := bufio.NewWriter(os.Stdout)
	fmt.Fprintln(writer, len(operations))
	for _, op := range operations {
		fmt.Fprint(writer, op.val)
		if op.typ == 0 {
			for _, eid := range treeEdges {
				fmt.Fprintf(writer, " %d", eid)
			}
		} else {
			for _, eid := range treeEdges {
				if eid == op.rem {
					continue
				}
				fmt.Fprintf(writer, " %d", eid)
			}
			fmt.Fprintf(writer, " %d", op.add)
		}
		fmt.Fprintln(writer)
	}
	writer.Flush()
}
