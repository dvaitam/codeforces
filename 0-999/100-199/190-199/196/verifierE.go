package main

import (
	"bufio"
	"container/heap"
	"fmt"
	"os"
	"os/exec"
	"sort"
	"strings"
)

const inf64 int64 = (1<<63 - 1) / 2

type edge struct {
	v int
	w int64
}

type origEdge struct {
	u, v int
	w    int64
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
func (h *hp) Pop() interface{} {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[:n-1]
	return x
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
func (h *msHP) Pop() interface{} {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[:n-1]
	return x
}

type mstEdge struct {
	u, v int
	w    int64
}

type dsu struct {
	p []int
}

func newDSU(n int) *dsu {
	p := make([]int, n+1)
	for i := range p {
		p[i] = i
	}
	return &dsu{p: p}
}

func (d *dsu) find(x int) int {
	if d.p[x] != x {
		d.p[x] = d.find(d.p[x])
	}
	return d.p[x]
}

func (d *dsu) union(a, b int) bool {
	ra, rb := d.find(a), d.find(b)
	if ra == rb {
		return false
	}
	d.p[rb] = ra
	return true
}

func dijkstra(n int, adj [][]edge, src int) []int64 {
	dist := make([]int64, n+1)
	for i := 1; i <= n; i++ {
		dist[i] = inf64
	}
	dist[src] = 0
	pq := &hp{}
	heap.Push(pq, item{u: src, d: 0})
	for pq.Len() > 0 {
		cur := heap.Pop(pq).(item)
		u, d := cur.u, cur.d
		if d != dist[u] {
			continue
		}
		for _, e := range adj[u] {
			nd := d + e.w
			if nd < dist[e.v] {
				dist[e.v] = nd
				heap.Push(pq, item{u: e.v, d: nd})
			}
		}
	}
	return dist
}

func multiSource(n int, adj [][]edge, sources []int) ([]int64, []int) {
	dist := make([]int64, n+1)
	label := make([]int, n+1)
	for i := 1; i <= n; i++ {
		dist[i] = inf64
		label[i] = 0
	}
	pq := &msHP{}
	for _, s := range sources {
		dist[s] = 0
		label[s] = s
		heap.Push(pq, msItem{u: s, d: 0})
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
				heap.Push(pq, msItem{u: e.v, d: nd})
			}
		}
	}
	return dist, label
}

func solveCase(n, m int, edges []origEdge, portals []int) int64 {
	adj := make([][]edge, n+1)
	for _, e := range edges {
		adj[e.u] = append(adj[e.u], edge{v: e.v, w: e.w})
		adj[e.v] = append(adj[e.v], edge{v: e.u, w: e.w})
	}

	dist1 := dijkstra(n, adj, 1)
	dist, label := multiSource(n, adj, portals)

	var bridge []mstEdge
	for _, e := range edges {
		lu, lv := label[e.u], label[e.v]
		if lu != lv {
			w := dist[e.u] + e.w + dist[e.v]
			bridge = append(bridge, mstEdge{u: lu, v: lv, w: w})
		}
	}
	sort.Slice(bridge, func(i, j int) bool { return bridge[i].w < bridge[j].w })

	ds := newDSU(n)
	var tot int64
	cnt := 0
	for _, e := range bridge {
		if ds.union(e.u, e.v) {
			tot += e.w
			cnt++
			if cnt == len(portals)-1 {
				break
			}
		}
	}

	minD1 := inf64
	for _, p := range portals {
		if dist1[p] < minD1 {
			minD1 = dist1[p]
		}
	}
	if minD1 == inf64 {
		minD1 = 0
	}
	return tot + minD1
}

type testCase struct {
	n, m    int
	edges   []origEdge
	portals []int
}

func parseTests() ([]testCase, error) {
	reader := bufio.NewReader(strings.NewReader(testcasesE))
	var t int
	if _, err := fmt.Fscan(reader, &t); err != nil {
		return nil, err
	}
	tests := make([]testCase, t)
	for i := 0; i < t; i++ {
		var n, m int
		if _, err := fmt.Fscan(reader, &n, &m); err != nil {
			return nil, err
		}
		edges := make([]origEdge, m)
		for j := 0; j < m; j++ {
			var u, v int
			var w int64
			if _, err := fmt.Fscan(reader, &u, &v, &w); err != nil {
				return nil, err
			}
			edges[j] = origEdge{u: u, v: v, w: w}
		}
		var k int
		if _, err := fmt.Fscan(reader, &k); err != nil {
			return nil, err
		}
		portals := make([]int, k)
		for j := 0; j < k; j++ {
			if _, err := fmt.Fscan(reader, &portals[j]); err != nil {
				return nil, err
			}
		}
		tests[i] = testCase{n: n, m: m, edges: edges, portals: portals}
	}
	return tests, nil
}

func buildInput(tc testCase) string {
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d %d\n", tc.n, tc.m)
	for _, e := range tc.edges {
		fmt.Fprintf(&sb, "%d %d %d\n", e.u, e.v, e.w)
	}
	fmt.Fprintf(&sb, "%d\n", len(tc.portals))
	for i, p := range tc.portals {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprint(p))
	}
	sb.WriteByte('\n')
	return sb.String()
}

func runCandidate(bin, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	out, err := cmd.CombinedOutput()
	return strings.TrimSpace(string(out)), err
}

const testcasesE = `100
3 3
1 2 13
1 3 9
2 3 5
2
1 3
4 5
2 3 13
2 4 14
1 3 1
1 2 3
3 4 8
3
4 2 1
3 3
1 2 2
1 3 2
2 3 9
3
3 2 1
3 2
1 2 14
1 3 15
2
1 2
5 9
3 5 7
1 3 11
1 2 18
4 5 3
3 4 5
1 4 2
2 5 9
1 5 15
2 3 10
2
4 2
2 1
1 2 2
1
1
6 5
1 4 4
2 6 16
1 3 11
1 2 19
3 5 14
2
2 4
5 7
1 4 14
3 5 4
2 3 15
3 4 5
2 5 16
1 2 8
1 3 9
5
4 1 2 3 5
5 5
4 5 17
1 4 9
3 4 11
1 2 11
2 3 16
2
1 4
6 6
2 6 1
3 4 16
4 5 13
1 2 2
1 5 9
1 3 15
1
4
5 5
1 2 7
2 4 18
4 5 17
1 3 20
1 4 3
1
5
2 1
1 2 20
1
1
4 3
2 3 3
1 2 10
1 4 18
4
4 1 2 3
3 3
1 2 4
1 3 4
2 3 17
2
2 3
2 1
1 2 16
2
1 2
4 4
2 4 9
1 3 4
2 3 16
1 2 5
1
4
6 8
1 3 14
4 5 10
1 6 2
3 5 18
2 3 20
2 4 12
2 5 6
1 2 17
3
6 3 2
5 4
1 2 20
3 5 20
2 4 7
1 3 10
1
1
5 9
3 5 13
1 5 20
2 4 8
1 4 10
1 3 4
2 3 11
4 5 19
2 5 16
1 2 11
2
3 5
3 3
1 2 13
1 3 15
2 3 11
3
3 2 1
5 8
3 5 3
2 3 2
2 5 6
3 4 12
2 4 20
1 2 5
1 5 12
1 3 15
1
5
3 3
1 2 7
1 3 13
2 3 16
1
1
4 6
3 4 6
1 3 5
2 3 6
1 2 13
1 4 9
2 4 2
2
1 2
6 7
3 5 16
3 6 3
1 3 7
2 3 4
2 4 19
1 2 14
1 5 11
5
4 3 1 6 5
5 8
1 5 17
2 4 20
4 5 13
1 2 2
2 3 8
2 5 10
1 3 3
3 5 1
1
2
5 7
3 5 7
3 4 10
1 5 1
2 3 18
2 5 14
1 3 6
1 2 17
5
1 2 4 5 3
2 1
1 2 15
1
1
2 1
1 2 4
2
2 1
2 1
1 2 6
1
2
2 1
1 2 20
2
1 2
3 2
1 2 6
1 3 9
3
3 2 1
3 2
2 3 18
1 2 16
2
2 3
5 4
1 2 6
1 3 12
1 4 10
2 5 12
5
2 3 5 4 1
3 2
1 2 6
1 3 1
1
1
3 3
1 2 18
2 3 2
1 3 15
1
3
5 5
2 5 2
1 5 16
3 4 2
1 2 11
2 3 20
2
5 2
6 7
1 2 1
1 4 19
4 6 11
1 6 5
2 5 20
2 3 14
5 6 12
2
6 3
6 9
2 3 18
1 2 9
1 6 2
5 6 7
3 5 2
1 4 18
1 5 12
2 6 20
2 5 12
3
4 3 2
5 6
1 5 11
1 3 17
2 3 14
1 2 5
1 4 1
3 5 17
2
2 1
2 1
1 2 14
1
1
6 7
1 5 20
5 6 14
1 6 6
1 2 9
3 5 12
3 4 4
1 3 19
2
5 3
6 8
1 2 20
3 6 7
3 4 6
2 3 6
4 5 20
1 5 19
2 4 19
5 6 19
3
3 6 2
4 3
2 4 4
1 3 6
1 2 11
3
1 3 4
2 1
1 2 8
1
2
5 7
3 4 3
2 4 14
3 5 2
1 2 15
1 5 8
2 3 10
1 4 1
3
3 5 4
3 2
1 2 12
1 3 13
3
1 2 3
6 8
1 6 10
3 6 13
4 6 12
1 2 13
2 5 1
2 4 13
1 3 3
3 4 13
6
1 3 4 5 6 2
5 9
1 4 14
4 5 4
3 4 15
1 3 14
2 4 4
3 5 18
2 5 16
1 2 8
1 5 8
3
4 3 1
3 3
2 3 18
1 2 5
1 3 16
2
3 2
4 3
1 3 18
1 2 9
1 4 1
4
3 1 4 2
3 3
2 3 15
1 3 9
1 2 11
1
3
2 1
1 2 4
1
1
5 8
3 5 13
4 5 4
3 4 19
2 5 17
1 5 16
1 2 18
1 3 6
1 4 1
4
2 3 5 1
4 6
1 4 17
2 4 18
1 2 10
2 3 7
1 3 6
3 4 7
1
3
6 8
2 6 16
1 4 13
2 4 8
5 6 11
3 6 2
2 3 1
3 5 20
1 2 17
2
2 5
5 8
1 3 12
1 4 16
1 5 6
1 2 5
2 4 16
4 5 15
3 4 17
3 5 4
2
1 3
3 2
1 2 18
1 3 7
2
3 2
4 6
1 3 5
2 4 1
3 4 2
1 2 11
2 3 7
1 4 11
3
3 2 4
5 6
1 4 20
1 2 19
1 3 14
2 4 19
3 4 4
1 5 8
4
4 5 3 2
6 6
2 3 18
5 6 20
1 3 10
2 4 9
1 5 15
1 2 17
4
3 5 4 6
4 3
2 3 12
1 2 20
3 4 12
2
2 1
6 7
2 4 14
4 6 8
2 3 15
1 2 9
2 6 9
3 4 11
2 5 19
4
2 6 4 1
5 8
3 5 13
2 4 5
1 2 3
1 3 17
3 4 5
1 5 18
2 3 20
1 4 5
5
3 4 5 1 2
3 3
1 3 18
1 2 14
2 3 4
1
2
6 7
3 5 6
3 6 15
2 5 20
1 2 15
1 3 10
1 6 17
3 4 7
4
2 1 3 6
6 7
1 3 18
3 6 3
1 5 17
5 6 2
1 2 2
4 5 3
2 4 9
4
4 6 5 3
5 6
1 4 13
3 4 2
2 3 14
1 3 19
3 5 20
1 2 17
2
5 4
4 3
2 3 8
1 2 14
3 4 19
4
1 4 3 2
5 5
2 3 9
1 3 8
1 2 19
1 4 19
2 5 6
2
3 1
6 5
4 6 15
1 2 13
2 5 3
2 3 10
3 4 13
4
6 5 4 3
4 5
1 2 7
1 4 20
2 4 8
3 4 5
2 3 7
3
2 1 3
3 3
1 2 6
2 3 9
1 3 4
1
2
5 5
1 2 10
1 5 4
1 4 10
1 3 14
3 4 13
3
2 5 4
2 1
1 2 10
2
2 1
6 6
1 2 10
5 6 1
4 5 13
2 5 20
2 4 3
2 3 16
4
2 3 5 6
6 6
1 3 5
1 2 13
2 4 7
1 6 5
3 5 14
2 5 9
6
4 2 1 3 6 5
2 1
1 2 3
1
1
4 6
3 4 19
2 3 6
1 3 11
1 4 19
1 2 5
2 4 16
2
4 3
5 5
1 2 7
1 4 13
3 5 9
1 3 13
3 4 1
1
3
5 7
4 5 14
1 2 6
1 3 14
1 4 12
1 5 15
2 4 12
3 5 10
1
5
4 6
1 2 20
2 4 20
1 4 2
2 3 17
3 4 17
1 3 9
3
2 3 1
6 9
3 4 10
2 6 19
1 2 19
1 6 6
2 5 11
1 3 17
3 5 5
1 4 2
3 6 1
2
2 5
3 3
1 3 1
1 2 20
2 3 6
3
3 1 2
2 1
1 2 11
2
2 1
3 3
1 2 1
2 3 1
1 3 11
2
2 3
5 7
1 5 14
2 5 5
1 4 10
2 3 17
2 4 16
1 2 14
4 5 15
2
4 2
3 3
1 2 12
2 3 9
1 3 15
2
2 3
2 1
1 2 19
1
2
3 2
1 2 5
2 3 6
1
2
3 2
2 3 2
1 2 4
3
1 3 2
2 1
1 2 5
2
2 1
4 3
1 3 19
3 4 16
1 2 11
1
3
2 1
1 2 14
1
1
3 3
1 2 13
1 3 14
2 3 4
3
2 3 1
3 3
2 3 12
1 2 18
1 3 5
1
2
3 3
2 3 11
1 2 3
1 3 5
3
1 2 3
5 8
1 5 20
4 5 17
1 3 11
1 2 19
3 5 1
1 4 18
2 5 19
2 4 6
2
2 1
3 3
1 2 1
2 3 12
1 3 10
1
3
5 4
3 5 18
1 2 8
1 3 20
1 4 1
1
3
2 1
1 2 10
1
1`

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: verifierE /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]

	tests, err := parseTests()
	if err != nil {
		fmt.Fprintln(os.Stderr, "parse error:", err)
		os.Exit(1)
	}

	for i, tc := range tests {
		want := solveCase(tc.n, tc.m, tc.edges, tc.portals)
		input := buildInput(tc)
		got, err := runCandidate(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d runtime error: %v\n%s\n", i+1, err, got)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != fmt.Sprint(want) {
			fmt.Fprintf(os.Stderr, "case %d failed\ninput:\n%s\nexpected: %d\ngot: %s\n", i+1, input, want, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(tests))
}
