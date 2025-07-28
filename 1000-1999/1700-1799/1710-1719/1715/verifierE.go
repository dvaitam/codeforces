package main

import (
	"bytes"
	"container/heap"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

type EdgeE struct {
	to int
	w  int64
}

type item struct {
	node int
	dist int64
}

type priorityQueue []*item

func (pq priorityQueue) Len() int            { return len(pq) }
func (pq priorityQueue) Less(i, j int) bool  { return pq[i].dist < pq[j].dist }
func (pq priorityQueue) Swap(i, j int)       { pq[i], pq[j] = pq[j], pq[i] }
func (pq *priorityQueue) Push(x interface{}) { *pq = append(*pq, x.(*item)) }
func (pq *priorityQueue) Pop() interface{} {
	old := *pq
	n := len(old)
	it := old[n-1]
	*pq = old[:n-1]
	return it
}

type Line struct {
	m, b int64
}

func (ln Line) value(x int64) int64 { return ln.m*x + ln.b }

type Node struct {
	ln          Line
	left, right *Node
}

func insert(node *Node, l, r int64, ln Line) *Node {
	if node == nil {
		return &Node{ln: ln}
	}
	mid := (l + r) >> 1
	leftBetter := ln.value(l) < node.ln.value(l)
	midBetter := ln.value(mid) < node.ln.value(mid)
	if midBetter {
		node.ln, ln = ln, node.ln
	}
	if l == r {
		return node
	}
	if leftBetter != midBetter {
		node.left = insert(node.left, l, mid, ln)
	} else {
		node.right = insert(node.right, mid+1, r, ln)
	}
	return node
}

func query(node *Node, l, r, x int64) int64 {
	if node == nil {
		return 1<<63 - 1
	}
	res := node.ln.value(x)
	if l == r {
		return res
	}
	mid := (l + r) >> 1
	if x <= mid {
		val := query(node.left, l, mid, x)
		if val < res {
			res = val
		}
	} else {
		val := query(node.right, mid+1, r, x)
		if val < res {
			res = val
		}
	}
	return res
}

type LiChao struct {
	root        *Node
	left, right int64
}

func NewLiChao(l, r int64) *LiChao     { return &LiChao{left: l, right: r} }
func (lc *LiChao) Insert(ln Line)      { lc.root = insert(lc.root, lc.left, lc.right, ln) }
func (lc *LiChao) Query(x int64) int64 { return query(lc.root, lc.left, lc.right, x) }

func dijkstra(start []int64, g [][]EdgeE) []int64 {
	n := len(g) - 1
	dist := make([]int64, n+1)
	copy(dist, start)
	pq := &priorityQueue{}
	heap.Init(pq)
	for i := 1; i <= n; i++ {
		if dist[i] < (1<<63 - 1) {
			heap.Push(pq, &item{node: i, dist: dist[i]})
		}
	}
	for pq.Len() > 0 {
		it := heap.Pop(pq).(*item)
		if it.dist != dist[it.node] {
			continue
		}
		u := it.node
		d := it.dist
		for _, e := range g[u] {
			nd := d + e.w
			if nd < dist[e.to] {
				dist[e.to] = nd
				heap.Push(pq, &item{node: e.to, dist: nd})
			}
		}
	}
	return dist
}

func solveE(n, m, k int, edges [][3]int) []int64 {
	g := make([][]EdgeE, n+1)
	for _, e := range edges {
		u := e[0]
		v := e[1]
		w := int64(e[2])
		g[u] = append(g[u], EdgeE{v, w})
		g[v] = append(g[v], EdgeE{u, w})
	}
	const inf int64 = 1<<63 - 1
	dist := make([]int64, n+1)
	for i := 1; i <= n; i++ {
		dist[i] = inf
	}
	dist[1] = 0
	dist = dijkstra(dist, g)
	for iter := 0; iter < k; iter++ {
		tree := NewLiChao(1, int64(n))
		for i := 1; i <= n; i++ {
			if dist[i] < inf {
				x := int64(i)
				tree.Insert(Line{m: -2 * x, b: dist[i] + x*x})
			}
		}
		newDist := make([]int64, n+1)
		for i := 1; i <= n; i++ {
			x := int64(i)
			newDist[i] = x*x + tree.Query(x)
		}
		dist = dijkstra(newDist, g)
	}
	return dist[1:]
}

type testCaseE struct {
	n, m, k int
	edges   [][3]int
}

func genCaseE(rng *rand.Rand) testCaseE {
	n := rng.Intn(5) + 2
	m := rng.Intn(6) + 1
	k := rng.Intn(3) + 1
	edges := make([][3]int, m)
	for i := 0; i < m; i++ {
		u := rng.Intn(n) + 1
		v := rng.Intn(n) + 1
		for v == u {
			v = rng.Intn(n) + 1
		}
		w := rng.Intn(10) + 1
		edges[i] = [3]int{u, v, w}
	}
	return testCaseE{n: n, m: m, k: k, edges: edges}
}

func runCaseE(bin string, tc testCaseE) error {
	var input strings.Builder
	input.WriteString(fmt.Sprintf("%d %d %d\n", tc.n, tc.m, tc.k))
	for _, e := range tc.edges {
		input.WriteString(fmt.Sprintf("%d %d %d\n", e[0], e[1], e[2]))
	}
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input.String())
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	fields := strings.Fields(strings.TrimSpace(out.String()))
	if len(fields) != tc.n {
		return fmt.Errorf("expected %d numbers got %d", tc.n, len(fields))
	}
	expect := solveE(tc.n, tc.m, tc.k, tc.edges)
	for i, f := range fields {
		val, err := strconv.ParseInt(f, 10, 64)
		if err != nil {
			return fmt.Errorf("invalid int %q", f)
		}
		if val != expect[i] {
			return fmt.Errorf("mismatch expected %v got %v", expect, fields)
		}
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for t := 0; t < 100; t++ {
		tc := genCaseE(rng)
		if err := runCaseE(bin, tc); err != nil {
			fmt.Fprintf(os.Stderr, "test %d failed: %v\ninput:\n", t+1, err)
			var inp strings.Builder
			inp.WriteString(fmt.Sprintf("%d %d %d\n", tc.n, tc.m, tc.k))
			for _, e := range tc.edges {
				inp.WriteString(fmt.Sprintf("%d %d %d\n", e[0], e[1], e[2]))
			}
			fmt.Fprint(os.Stderr, inp.String())
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
