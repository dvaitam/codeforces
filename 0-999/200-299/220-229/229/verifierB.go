package main

import (
	"bytes"
	"container/heap"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"sort"
	"strconv"
	"strings"
	"time"
)

type edgeB struct {
	a, b int
	w    int64
}

type testCaseB struct {
	n     int
	edges []edgeB
	forb  [][]int64
	ans   int64
}

type pqItem struct {
	node int
	dist int64
}

type pq []pqItem

func (p pq) Len() int            { return len(p) }
func (p pq) Less(i, j int) bool  { return p[i].dist < p[j].dist }
func (p pq) Swap(i, j int)       { p[i], p[j] = p[j], p[i] }
func (p *pq) Push(x interface{}) { *p = append(*p, x.(pqItem)) }
func (p *pq) Pop() interface{} {
	old := *p
	n := len(old)
	it := old[n-1]
	*p = old[:n-1]
	return it
}

type interval struct{ l, r int64 }

func solveB(tc testCaseB) int64 {
	n := tc.n
	adj := make([][]edgeB, n+1)
	for _, e := range tc.edges {
		adj[e.a] = append(adj[e.a], e)
		adj[e.b] = append(adj[e.b], edgeB{e.b, e.a, e.w})
	}
	forb := make([][]interval, n+1)
	for i := 1; i <= n; i++ {
		t := tc.forb[i]
		if len(t) == 0 {
			continue
		}
		ivs := make([]interval, 0, len(t))
		start := t[0]
		prev := t[0]
		for j := 1; j < len(t); j++ {
			if t[j] == prev+1 {
				prev = t[j]
			} else {
				ivs = append(ivs, interval{start, prev})
				start = t[j]
				prev = t[j]
			}
		}
		ivs = append(ivs, interval{start, prev})
		forb[i] = ivs
	}
	const INF int64 = 1 << 62
	dist := make([]int64, n+1)
	for i := range dist {
		dist[i] = INF
	}
	dist[1] = 0
	pq := &pq{}
	heap.Push(pq, pqItem{1, 0})
	for pq.Len() > 0 {
		it := heap.Pop(pq).(pqItem)
		u := it.node
		d := it.dist
		if d != dist[u] {
			continue
		}
		if u == n {
			break
		}
		dep := d
		ivs := forb[u]
		if len(ivs) > 0 {
			idx := sort.Search(len(ivs), func(i int) bool { return ivs[i].l > d }) - 1
			if idx >= 0 && ivs[idx].l <= d && d <= ivs[idx].r {
				dep = ivs[idx].r + 1
			}
		}
		for _, e := range adj[u] {
			v := e.b
			nd := dep + e.w
			if nd < dist[v] {
				dist[v] = nd
				heap.Push(pq, pqItem{v, nd})
			}
		}
	}
	if dist[n] == INF {
		return -1
	}
	return dist[n]
}

func genCaseB(rng *rand.Rand) testCaseB {
	n := rng.Intn(4) + 2
	maxEdges := n * (n - 1) / 2
	m := rng.Intn(maxEdges + 1)
	edges := make([]edgeB, 0, m)
	exist := make(map[[2]int]bool)
	for len(edges) < m {
		a := rng.Intn(n) + 1
		b := rng.Intn(n) + 1
		if a == b {
			continue
		}
		if a > b {
			a, b = b, a
		}
		key := [2]int{a, b}
		if exist[key] {
			continue
		}
		exist[key] = true
		w := int64(rng.Intn(10) + 1)
		edges = append(edges, edgeB{a, b, w})
	}
	forb := make([][]int64, n+1)
	for i := 1; i <= n; i++ {
		k := rng.Intn(3)
		t := make([]int64, k)
		cur := int64(0)
		for j := 0; j < k; j++ {
			cur += int64(rng.Intn(4) + 1)
			t[j] = cur
		}
		forb[i] = t
	}
	tc := testCaseB{n: n, edges: edges, forb: forb}
	tc.ans = solveB(tc)
	return tc
}

func runCaseB(bin string, tc testCaseB) error {
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d %d\n", tc.n, len(tc.edges))
	for _, e := range tc.edges {
		fmt.Fprintf(&sb, "%d %d %d\n", e.a, e.b, e.w)
	}
	for i := 1; i <= tc.n; i++ {
		fmt.Fprintf(&sb, "%d", len(tc.forb[i]))
		for _, v := range tc.forb[i] {
			fmt.Fprintf(&sb, " %d", v)
		}
		sb.WriteByte('\n')
	}
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(sb.String())
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	fields := strings.Fields(strings.TrimSpace(out.String()))
	if len(fields) != 1 {
		return fmt.Errorf("expected 1 number got %d", len(fields))
	}
	val, err := strconv.ParseInt(fields[0], 10, 64)
	if err != nil {
		return fmt.Errorf("invalid integer: %v", err)
	}
	if val != tc.ans {
		return fmt.Errorf("expected %d got %d", tc.ans, val)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		tc := genCaseB(rng)
		if err := runCaseB(bin, tc); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", i+1, err)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
