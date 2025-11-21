package main

import (
	"bufio"
	"container/heap"
	"os"
	"sort"
)

const (
	LOG = 20
	INF = int64(1 << 60)
)

type edge struct {
	to int
	w  int
}

type state struct {
	turn int64
	idx  int
	dist int64
	node int
}

type priorityQueue []state

func (pq priorityQueue) Len() int { return len(pq) }

func (pq priorityQueue) Less(i, j int) bool {
	if pq[i].turn != pq[j].turn {
		return pq[i].turn < pq[j].turn
	}
	if pq[i].idx != pq[j].idx {
		return pq[i].idx < pq[j].idx
	}
	return pq[i].dist < pq[j].dist
}

func (pq priorityQueue) Swap(i, j int) { pq[i], pq[j] = pq[j], pq[i] }

func (pq *priorityQueue) Push(x interface{}) {
	*pq = append(*pq, x.(state))
}

func (pq *priorityQueue) Pop() interface{} {
	old := *pq
	n := len(old)
	x := old[n-1]
	*pq = old[:n-1]
	return x
}

type fastScanner struct {
	r *bufio.Reader
}

func newFastScanner() *fastScanner {
	return &fastScanner{r: bufio.NewReader(os.Stdin)}
}

func (fs *fastScanner) nextInt() int {
	sign := 1
	val := 0
	c, err := fs.r.ReadByte()
	for (c < '0' || c > '9') && c != '-' {
		if err != nil {
			return 0
		}
		c, err = fs.r.ReadByte()
	}
	if c == '-' {
		sign = -1
		c, err = fs.r.ReadByte()
	}
	for c >= '0' && c <= '9' {
		val = val*10 + int(c-'0')
		c, err = fs.r.ReadByte()
		if err != nil {
			break
		}
	}
	return sign * val
}

var (
	tree       [][]int
	virtualAdj [][]edge
	tin        []int
	tout       []int
	depth      []int
	up         [][]int
	timer      int
	inVirtual  []bool
	vtNodes    []int
	bestTurn   []int64
	bestVirus  []int
)

func main() {
	fs := newFastScanner()
	n := fs.nextInt()

	tree = make([][]int, n+1)
	for i := 0; i < n-1; i++ {
		u := fs.nextInt()
		v := fs.nextInt()
		tree[u] = append(tree[u], v)
		tree[v] = append(tree[v], u)
	}

	tin = make([]int, n+1)
	tout = make([]int, n+1)
	depth = make([]int, n+1)
	up = make([][]int, LOG)
	for i := 0; i < LOG; i++ {
		up[i] = make([]int, n+1)
	}
	iterativeDFS(1, n)

	virtualAdj = make([][]edge, n+1)
	inVirtual = make([]bool, n+1)
	bestTurn = make([]int64, n+1)
	bestVirus = make([]int, n+1)

	q := fs.nextInt()
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	for ; q > 0; q-- {
		vtNodes = vtNodes[:0]
		k := fs.nextInt()
		m := fs.nextInt()

		virusNodes := make([]int, k+1)
		speeds := make([]int64, k+1)
		nodes := make([]int, 0, k+m)

		for i := 1; i <= k; i++ {
			v := fs.nextInt()
			s := fs.nextInt()
			virusNodes[i] = v
			speeds[i] = int64(s)
			nodes = append(nodes, v)
		}

		targets := make([]int, m)
		for i := 0; i < m; i++ {
			u := fs.nextInt()
			targets[i] = u
			nodes = append(nodes, u)
		}

		buildVirtualTree(nodes)

		for _, v := range vtNodes {
			bestTurn[v] = INF
			bestVirus[v] = 0
		}

		pq := priorityQueue{}
		heap.Init(&pq)
		for i := 1; i <= k; i++ {
			heap.Push(&pq, state{turn: 0, idx: i, dist: 0, node: virusNodes[i]})
		}

		for pq.Len() > 0 {
			cur := heap.Pop(&pq).(state)
			if cur.turn > bestTurn[cur.node] {
				continue
			}
			if cur.turn == bestTurn[cur.node] && bestVirus[cur.node] != 0 && cur.idx >= bestVirus[cur.node] {
				continue
			}
			bestTurn[cur.node] = cur.turn
			bestVirus[cur.node] = cur.idx
			for _, e := range virtualAdj[cur.node] {
				newDist := cur.dist + int64(e.w)
				s := speeds[cur.idx]
				newTurn := (newDist + s - 1) / s
				if newTurn > bestTurn[e.to] {
					continue
				}
				if newTurn == bestTurn[e.to] && bestVirus[e.to] != 0 && cur.idx >= bestVirus[e.to] {
					continue
				}
				heap.Push(&pq, state{turn: newTurn, idx: cur.idx, dist: newDist, node: e.to})
			}
		}

		for i, city := range targets {
			if i > 0 {
				out.WriteByte(' ')
			}
			out.WriteString(intToString(bestVirus[city]))
		}
		out.WriteByte('\n')

		for _, v := range vtNodes {
			inVirtual[v] = false
		}
	}
}

func iterativeDFS(root, n int) {
	type frame struct {
		node    int
		parent  int
		visited bool
	}
	stack := []frame{{node: root, parent: root, visited: false}}
	timer = 0
	for len(stack) > 0 {
		cur := stack[len(stack)-1]
		stack = stack[:len(stack)-1]
		if !cur.visited {
			tin[cur.node] = timer
			timer++
			up[0][cur.node] = cur.parent
			for i := 1; i < LOG; i++ {
				up[i][cur.node] = up[i-1][up[i-1][cur.node]]
			}
			stack = append(stack, frame{node: cur.node, parent: cur.parent, visited: true})
			for i := len(tree[cur.node]) - 1; i >= 0; i-- {
				to := tree[cur.node][i]
				if to == cur.parent {
					continue
				}
				depth[to] = depth[cur.node] + 1
				stack = append(stack, frame{node: to, parent: cur.node, visited: false})
			}
		} else {
			tout[cur.node] = timer
			timer++
		}
	}
}

func isAncestor(u, v int) bool {
	return tin[u] <= tin[v] && tout[v] <= tout[u]
}

func lca(u, v int) int {
	if isAncestor(u, v) {
		return u
	}
	if isAncestor(v, u) {
		return v
	}
	for i := LOG - 1; i >= 0; i-- {
		ancestor := up[i][u]
		if !isAncestor(ancestor, v) {
			u = ancestor
		}
	}
	return up[0][u]
}

func ensureNode(v int) {
	if !inVirtual[v] {
		inVirtual[v] = true
		vtNodes = append(vtNodes, v)
		virtualAdj[v] = virtualAdj[v][:0]
	}
}

func addVirtualEdge(u, v int) {
	w := depth[v] - depth[u]
	if w < 0 {
		w = -w
	}
	virtualAdj[u] = append(virtualAdj[u], edge{to: v, w: w})
	virtualAdj[v] = append(virtualAdj[v], edge{to: u, w: w})
}

func buildVirtualTree(nodes []int) {
	if len(nodes) == 0 {
		return
	}
	sort.Slice(nodes, func(i, j int) bool { return tin[nodes[i]] < tin[nodes[j]] })
	nodes = uniqueSorted(nodes)
	origLen := len(nodes)
	for i := 0; i < origLen-1; i++ {
		nodes = append(nodes, lca(nodes[i], nodes[i+1]))
	}
	sort.Slice(nodes, func(i, j int) bool { return tin[nodes[i]] < tin[nodes[j]] })
	nodes = uniqueSorted(nodes)

	stack := make([]int, 0, len(nodes))
	ensureNode(nodes[0])
	stack = append(stack, nodes[0])
	for i := 1; i < len(nodes); i++ {
		v := nodes[i]
		ensureNode(v)
		for len(stack) > 0 && !isAncestor(stack[len(stack)-1], v) {
			stack = stack[:len(stack)-1]
		}
		if len(stack) == 0 {
			stack = append(stack, v)
			continue
		}
		parent := stack[len(stack)-1]
		addVirtualEdge(parent, v)
		stack = append(stack, v)
	}
}

func uniqueSorted(arr []int) []int {
	if len(arr) == 0 {
		return arr
	}
	idx := 1
	for i := 1; i < len(arr); i++ {
		if arr[i] != arr[i-1] {
			arr[idx] = arr[i]
			idx++
		}
	}
	return arr[:idx]
}

func intToString(x int) string {
	if x == 0 {
		return "0"
	}
	buf := [20]byte{}
	pos := len(buf)
	val := x
	if val < 0 {
		val = -val
	}
	for val > 0 {
		pos--
		buf[pos] = byte('0' + val%10)
		val /= 10
	}
	if x < 0 {
		pos--
		buf[pos] = '-'
	}
	return string(buf[pos:])
}
