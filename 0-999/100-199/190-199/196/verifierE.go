package main

import (
	"bufio"
	"bytes"
	"container/heap"
	"fmt"
	"os"
	"os/exec"
	"sort"
)

const INF64 int64 = (1<<63 - 1) / 2

type edge struct {
	v int
	w int64
}

type item struct {
	u int
	d int64
}

type hp []item

func (h hp) Len() int            { return len(h) }
func (h hp) Less(i, j int) bool  { return h[i].d < h[j].d }
func (h hp) Swap(i, j int)       { h[i], h[j] = h[j], h[i] }
func (h *hp) Push(x interface{}) { *h = append(*h, x.(item)) }
func (h *hp) Pop() interface{}   { old := *h; n := len(old); x := old[n-1]; *h = old[:n-1]; return x }

func dijkstra(n int, adj [][]edge, src int) []int64 {
	dist := make([]int64, n+1)
	for i := 1; i <= n; i++ {
		dist[i] = INF64
	}
	dist[src] = 0
	pq := &hp{}
	heap.Push(pq, item{src, 0})
	for pq.Len() > 0 {
		cur := heap.Pop(pq).(item)
		u := cur.u
		d := cur.d
		if d != dist[u] {
			continue
		}
		for _, e := range adj[u] {
			nd := d + e.w
			if nd < dist[e.v] {
				dist[e.v] = nd
				heap.Push(pq, item{e.v, nd})
			}
		}
	}
	return dist
}

type msItem struct {
	u int
	d int64
}

type msHP []msItem

func (h msHP) Len() int            { return len(h) }
func (h msHP) Less(i, j int) bool  { return h[i].d < h[j].d }
func (h msHP) Swap(i, j int)       { h[i], h[j] = h[j], h[i] }
func (h *msHP) Push(x interface{}) { *h = append(*h, x.(msItem)) }
func (h *msHP) Pop() interface{}   { old := *h; n := len(old); x := old[n-1]; *h = old[:n-1]; return x }

func multiSource(n int, adj [][]edge, srcs []int) ([]int64, []int) {
	dist := make([]int64, n+1)
	label := make([]int, n+1)
	for i := 1; i <= n; i++ {
		dist[i] = INF64
		label[i] = 0
	}
	pq := &msHP{}
	for _, s := range srcs {
		dist[s] = 0
		label[s] = s
		heap.Push(pq, msItem{s, 0})
	}
	for pq.Len() > 0 {
		cur := heap.Pop(pq).(msItem)
		u, d := cur.u, cur.d
		if d != dist[u] {
			continue
		}
		for _, e := range adj[u] {
			nd := d + e.w
			if nd < dist[e.v] {
				dist[e.v] = nd
				label[e.v] = label[u]
				heap.Push(pq, msItem{e.v, nd})
			}
		}
	}
	return dist, label
}

type mstEdge struct {
	u, v int
	w    int64
}

func solveCase(n, m int, edges []mstEdge, portals []int) int64 {
	adj := make([][]edge, n+1)
	orig := make([]mstEdge, 0, len(edges))
	for _, e := range edges {
		adj[e.u] = append(adj[e.u], edge{e.v, e.w})
		adj[e.v] = append(adj[e.v], edge{e.u, e.w})
		orig = append(orig, e)
	}
	dist1 := dijkstra(n, adj, 1)
	dist, label := multiSource(n, adj, portals)
	var bridge []mstEdge
	for _, e := range orig {
		lu, lv := label[e.u], label[e.v]
		if lu != lv {
			w := dist[e.u] + e.w + dist[e.v]
			bridge = append(bridge, mstEdge{lu, lv, w})
		}
	}
	sort.Slice(bridge, func(i, j int) bool { return bridge[i].w < bridge[j].w })
	parent := make([]int, n+1)
	for i := 1; i <= n; i++ {
		parent[i] = i
	}
	var find func(int) int
	find = func(x int) int {
		if parent[x] != x {
			parent[x] = find(parent[x])
		}
		return parent[x]
	}
	union := func(a, b int) bool {
		ra, rb := find(a), find(b)
		if ra == rb {
			return false
		}
		parent[rb] = ra
		return true
	}
	var tot int64
	cnt := 0
	for _, e := range bridge {
		if union(e.u, e.v) {
			tot += e.w
			cnt++
			if cnt == len(portals)-1 {
				break
			}
		}
	}
	minD1 := INF64
	for _, p := range portals {
		if dist1[p] < minD1 {
			minD1 = dist1[p]
		}
	}
	if minD1 == INF64 {
		minD1 = 0
	}
	return tot + minD1
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	data, err := os.ReadFile("testcasesE.txt")
	if err != nil {
		fmt.Println("could not read testcasesE.txt:", err)
		os.Exit(1)
	}
	scan := bufio.NewScanner(bytes.NewReader(data))
	scan.Split(bufio.ScanWords)
	if !scan.Scan() {
		fmt.Println("invalid test file")
		os.Exit(1)
	}
	var t int
	fmt.Sscanf(scan.Text(), "%d", &t)
	casesN := make([]int, t)
	casesM := make([]int, t)
	allEdges := make([][]mstEdge, t)
	ports := make([][]int, t)
	for i := 0; i < t; i++ {
		if !scan.Scan() {
			fmt.Println("bad test file")
			os.Exit(1)
		}
		fmt.Sscanf(scan.Text(), "%d", &casesN[i])
		scan.Scan()
		fmt.Sscanf(scan.Text(), "%d", &casesM[i])
		es := make([]mstEdge, casesM[i])
		for j := 0; j < casesM[i]; j++ {
			scan.Scan()
			var u int
			fmt.Sscanf(scan.Text(), "%d", &u)
			scan.Scan()
			var v int
			fmt.Sscanf(scan.Text(), "%d", &v)
			scan.Scan()
			var w int64
			fmt.Sscanf(scan.Text(), "%d", &w)
			es[j] = mstEdge{u, v, w}
		}
		allEdges[i] = es
		scan.Scan()
		var k int
		fmt.Sscanf(scan.Text(), "%d", &k)
		ps := make([]int, k)
		for j := 0; j < k; j++ {
			scan.Scan()
			fmt.Sscanf(scan.Text(), "%d", &ps[j])
		}
		ports[i] = ps
	}
	cmd := exec.Command(os.Args[1])
	cmd.Stdin = bytes.NewReader(data)
	out, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println("execution failed:", err)
		os.Exit(1)
	}
	outScan := bufio.NewScanner(bytes.NewReader(out))
	outScan.Split(bufio.ScanWords)
	for i := 0; i < t; i++ {
		if !outScan.Scan() {
			fmt.Printf("missing output for test %d\n", i+1)
			os.Exit(1)
		}
		var got int64
		fmt.Sscanf(outScan.Text(), "%d", &got)
		expect := solveCase(casesN[i], casesM[i], allEdges[i], ports[i])
		if got != expect {
			fmt.Printf("test %d failed: expected %d got %d\n", i+1, expect, got)
			os.Exit(1)
		}
	}
	if outScan.Scan() {
		fmt.Println("extra output detected")
		os.Exit(1)
	}
	fmt.Println("All tests passed!")
}
