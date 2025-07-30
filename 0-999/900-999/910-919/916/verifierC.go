package main

import (
	"bufio"
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

// isPrime checks primality for small values used in the tests.
func isPrime(n int64) bool {
	if n < 2 {
		return false
	}
	for i := int64(2); i*i <= n; i++ {
		if n%i == 0 {
			return false
		}
	}
	return true
}

func genCase(rng *rand.Rand) string {
	n := rng.Intn(8) + 2
	maxE := n * (n - 1) / 2
	m := rng.Intn(maxE-(n-1)+1) + (n - 1)
	return fmt.Sprintf("%d %d\n", n, m)
}

func run(bin, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	return strings.TrimSpace(out.String()), nil
}

type edge struct {
	u, v int
	w    int64
}

type adjEdge struct {
	to int
	w  int64
}

type item struct {
	v int
	d int64
}

type pq []item

func (p pq) Len() int            { return len(p) }
func (p pq) Less(i, j int) bool  { return p[i].d < p[j].d }
func (p pq) Swap(i, j int)       { p[i], p[j] = p[j], p[i] }
func (p *pq) Push(x interface{}) { *p = append(*p, x.(item)) }
func (p *pq) Pop() interface{} {
	old := *p
	v := old[len(old)-1]
	*p = old[:len(old)-1]
	return v
}

type dsu struct {
	p []int
	r []int
}

func newDSU(n int) *dsu {
	d := &dsu{p: make([]int, n), r: make([]int, n)}
	for i := range d.p {
		d.p[i] = i
	}
	return d
}

func (d *dsu) find(x int) int {
	if d.p[x] != x {
		d.p[x] = d.find(d.p[x])
	}
	return d.p[x]
}

func (d *dsu) union(a, b int) bool {
	a = d.find(a)
	b = d.find(b)
	if a == b {
		return false
	}
	if d.r[a] < d.r[b] {
		a, b = b, a
	}
	d.p[b] = a
	if d.r[a] == d.r[b] {
		d.r[a]++
	}
	return true
}

func parseAndCheck(n, m int, out string) error {
	scanner := bufio.NewScanner(strings.NewReader(out))
	if !scanner.Scan() {
		return fmt.Errorf("no output")
	}
	fields := strings.Fields(scanner.Text())
	if len(fields) != 2 {
		return fmt.Errorf("first line should contain 2 numbers")
	}
	sp, err := strconv.ParseInt(fields[0], 10, 64)
	if err != nil {
		return fmt.Errorf("bad shortest path: %v", err)
	}
	mst, err := strconv.ParseInt(fields[1], 10, 64)
	if err != nil {
		return fmt.Errorf("bad MST weight: %v", err)
	}
	edges := make([]edge, 0, m)
	seen := make(map[[2]int]bool)
	for i := 0; i < m; i++ {
		if !scanner.Scan() {
			return fmt.Errorf("expected %d edges, got %d", m, i)
		}
		f := strings.Fields(scanner.Text())
		if len(f) != 3 {
			return fmt.Errorf("edge %d should contain 3 numbers", i+1)
		}
		u, err1 := strconv.Atoi(f[0])
		v, err2 := strconv.Atoi(f[1])
		w, err3 := strconv.ParseInt(f[2], 10, 64)
		if err1 != nil || err2 != nil || err3 != nil {
			return fmt.Errorf("invalid edge on line %d", i+2)
		}
		if u < 1 || u > n || v < 1 || v > n {
			return fmt.Errorf("edge %d contains invalid vertex", i+1)
		}
		if u == v {
			return fmt.Errorf("edge %d is a loop", i+1)
		}
		if w < 1 || w > 1_000_000_000 {
			return fmt.Errorf("edge %d has invalid weight", i+1)
		}
		a, b := u, v
		if a > b {
			a, b = b, a
		}
		if seen[[2]int{a, b}] {
			return fmt.Errorf("duplicate edge %d-%d", a, b)
		}
		seen[[2]int{a, b}] = true
		edges = append(edges, edge{u, v, w})
	}
	if scanner.Scan() {
		extra := strings.TrimSpace(scanner.Text())
		if extra != "" {
			return fmt.Errorf("extra output after edges")
		}
	}

	// compute MST weight using Kruskal
	ecopy := append([]edge(nil), edges...)
	sort.Slice(ecopy, func(i, j int) bool { return ecopy[i].w < ecopy[j].w })
	d := newDSU(n + 1)
	var mstSum int64
	cnt := 0
	for _, e := range ecopy {
		if d.union(e.u, e.v) {
			mstSum += e.w
			cnt++
			if cnt == n-1 {
				break
			}
		}
	}
	if cnt != n-1 {
		return fmt.Errorf("graph is not connected")
	}
	if mstSum != mst {
		return fmt.Errorf("MST weight mismatch: got %d expected %d", mstSum, mst)
	}
	if !isPrime(mst) {
		return fmt.Errorf("MST weight %d is not prime", mst)
	}

	// compute shortest path from 1 to n using Dijkstra
	adj := make([][]adjEdge, n+1)
	for _, e := range edges {
		adj[e.u] = append(adj[e.u], adjEdge{e.v, e.w})
		adj[e.v] = append(adj[e.v], adjEdge{e.u, e.w})
	}
	const inf = int64(1<<62 - 1)
	dist := make([]int64, n+1)
	for i := range dist {
		dist[i] = inf
	}
	dist[1] = 0
	pq := &pq{{1, 0}}
	heap.Init(pq)
	for pq.Len() > 0 {
		it := heap.Pop(pq).(item)
		if it.d != dist[it.v] {
			continue
		}
		if it.v == n {
			break
		}
		for _, e := range adj[it.v] {
			nd := it.d + e.w
			if nd < dist[e.to] {
				dist[e.to] = nd
				heap.Push(pq, item{e.to, nd})
			}
		}
	}
	if dist[n] == inf {
		return fmt.Errorf("no path from 1 to %d", n)
	}
	if dist[n] != sp {
		return fmt.Errorf("shortest path mismatch: got %d expected %d", sp, dist[n])
	}
	if !isPrime(sp) {
		return fmt.Errorf("shortest path %d is not prime", sp)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		in := genCase(rng)
		out, err := run(bin, in)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}
		var n, m int
		fmt.Sscan(in, &n, &m)
		if err := parseAndCheck(n, m, out); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s\noutput:\n%s\n", i+1, err, in, out)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
