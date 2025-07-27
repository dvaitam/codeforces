package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
	"time"
)

const inf = 100000000

type edge struct {
	u int
	v int
	w int
}

// flow network structures

type Network struct {
	n    int
	to   []int
	cost []int
	cap  []int
	dist []int
	par  []int
	g    [][]int
}

func NewNetwork(n int) *Network {
	g := make([][]int, n)
	return &Network{n: n, g: g, to: []int{}, cost: []int{}, cap: []int{}, dist: make([]int, n), par: make([]int, n)}
}

func (net *Network) AddEdge(a, b, cst, cp int) {
	net.g[a] = append(net.g[a], len(net.to))
	net.to = append(net.to, b)
	net.cap = append(net.cap, cp)
	net.cost = append(net.cost, cst)
}

func (net *Network) Add(a, b, cst, cp int) {
	net.AddEdge(a, b, cst, cp)
	net.AddEdge(b, a, -cst, 0)
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func (net *Network) Mincost(total int) {
	n := net.n
	inQ := make([]bool, n)
	for total > 0 {
		for i := 0; i < n; i++ {
			net.dist[i] = inf
			inQ[i] = false
		}
		net.dist[0] = 0
		q := []int{0}
		inQ[0] = true
		for len(q) > 0 {
			v := q[0]
			q = q[1:]
			inQ[v] = false
			for _, it := range net.g[v] {
				if net.cap[it] > 0 {
					u := net.to[it]
					c := net.cost[it]
					if net.dist[v]+c < net.dist[u] {
						net.dist[u] = net.dist[v] + c
						net.par[u] = it
						if !inQ[u] {
							q = append(q, u)
							inQ[u] = true
						}
					}
				}
			}
		}
		if net.dist[1] >= inf {
			break
		}
		mn := inf
		path := []int{}
		t := 1
		for t != 0 {
			it := net.par[t]
			path = append(path, it)
			mn = min(mn, net.cap[it])
			t = net.to[it^1]
		}
		for _, it := range path {
			net.cap[it] -= mn
			net.cap[it^1] += mn
		}
		total -= mn
	}
}

func solveCaseG(n int, edges []edge) []int {
	me := NewNetwork(n + 2)
	b := make([]int, n)
	for _, e := range edges {
		me.Add(e.u+2, e.v+2, -1, inf)
		b[e.u] += e.w
		b[e.v] -= e.w
	}
	total := 0
	for i := 0; i < n; i++ {
		if b[i] > 0 {
			total += b[i]
			me.Add(0, i+2, 0, b[i])
		} else if b[i] < 0 {
			me.Add(i+2, 1, 0, -b[i])
		}
	}
	me.Mincost(total)
	ans := make([]int, n)
	mn := inf
	for i := 0; i < n; i++ {
		ans[i] = me.dist[i+2]
		if ans[i] < mn {
			mn = ans[i]
		}
	}
	for i := 0; i < n; i++ {
		ans[i] -= mn
	}
	return ans
}

type testCaseG struct {
	n     int
	edges []edge
}

func buildInputG(tc testCaseG) string {
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d %d\n", tc.n, len(tc.edges))
	for _, e := range tc.edges {
		fmt.Fprintf(&sb, "%d %d %d\n", e.u+1, e.v+1, e.w)
	}
	return sb.String()
}

func costFor(edges []edge, vals []int) (int64, bool) {
	var cost int64
	for _, e := range edges {
		if vals[e.u] <= vals[e.v] {
			return 0, false
		}
		cost += int64(e.w) * int64(vals[e.u]-vals[e.v])
	}
	return cost, true
}

func runCaseG(bin string, tc testCaseG) error {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(buildInputG(tc))
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	fields := strings.Fields(out.String())
	if len(fields) != tc.n {
		return fmt.Errorf("expected %d numbers", tc.n)
	}
	vals := make([]int, tc.n)
	for i := 0; i < tc.n; i++ {
		if _, err := fmt.Sscan(fields[i], &vals[i]); err != nil {
			return fmt.Errorf("bad output")
		}
		if vals[i] < 0 || vals[i] > 1_000_000_000 {
			return fmt.Errorf("value out of range")
		}
	}
	outCost, ok := costFor(tc.edges, vals)
	if !ok {
		return fmt.Errorf("constraints violated")
	}
	best := solveCaseG(tc.n, tc.edges)
	bestCost, _ := costFor(tc.edges, best)
	if outCost != bestCost {
		return fmt.Errorf("non optimal answer")
	}
	return nil
}

func generateDAG(rng *rand.Rand, n int) []edge {
	edges := []edge{}
	for i := 0; i < n; i++ {
		for j := i + 1; j < n; j++ {
			if rng.Intn(2) == 0 {
				w := rng.Intn(5) + 1
				edges = append(edges, edge{u: i, v: j, w: w})
			}
		}
	}
	return edges
}

func generateCasesG() []testCaseG {
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	cases := make([]testCaseG, 0, 100)
	cases = append(cases, testCaseG{n: 2, edges: []edge{{0, 1, 1}}})
	for len(cases) < 100 {
		n := rng.Intn(5) + 2
		edges := generateDAG(rng, n)
		cases = append(cases, testCaseG{n: n, edges: edges})
	}
	return cases
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierG.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	cases := generateCasesG()
	for i, tc := range cases {
		if err := runCaseG(bin, tc); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", i+1, err)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
