package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"sort"
	"strings"
	"time"
)

const LOGH = 19

type edge struct{ u, v, c, t int }
type query struct{ v, x int }

type caseH struct {
	n       int
	e       []int
	edges   []edge
	q       int
	queries []query
}

func genCase(rng *rand.Rand) caseH {
	n := rng.Intn(6) + 2
	e := make([]int, n+1)
	for i := 1; i <= n; i++ {
		e[i] = rng.Intn(10) + 1
	}
	edges := make([]edge, 0, n-1)
	for i := 2; i <= n; i++ {
		p := rng.Intn(i-1) + 1
		edges = append(edges, edge{p, i, rng.Intn(10) + 1, rng.Intn(10) + 1})
	}
	qn := rng.Intn(5) + 1
	qs := make([]query, qn)
	for i := 0; i < qn; i++ {
		qs[i] = query{rng.Intn(10) + 1, rng.Intn(n) + 1}
	}
	return caseH{n, e, edges, qn, qs}
}

var (
	g [][]struct {
		to   int
		toll int
	}
	up    [][]int
	mx    [][]int
	depth []int
)

func maxInt(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func buildLCA(n int) {
	up = make([][]int, LOGH)
	mx = make([][]int, LOGH)
	for i := 0; i < LOGH; i++ {
		up[i] = make([]int, n+1)
		mx[i] = make([]int, n+1)
	}
	depth = make([]int, n+1)
	q := []int{1}
	for len(q) > 0 {
		v := q[0]
		q = q[1:]
		for _, e := range g[v] {
			if e.to == up[0][v] {
				continue
			}
			up[0][e.to] = v
			mx[0][e.to] = e.toll
			depth[e.to] = depth[v] + 1
			q = append(q, e.to)
		}
	}
	for k := 1; k < LOGH; k++ {
		for v := 1; v <= n; v++ {
			anc := up[k-1][v]
			up[k][v] = up[k-1][anc]
			if anc != 0 {
				mx[k][v] = maxInt(mx[k-1][v], mx[k-1][anc])
			} else {
				mx[k][v] = mx[k-1][v]
			}
		}
	}
}

func maxEdge(u, v int) int {
	if u == v {
		return 0
	}
	res := 0
	if depth[u] < depth[v] {
		u, v = v, u
	}
	diff := depth[u] - depth[v]
	for k := LOGH - 1; k >= 0; k-- {
		if diff&(1<<k) != 0 {
			if mx[k][u] > res {
				res = mx[k][u]
			}
			u = up[k][u]
		}
	}
	if u == v {
		return res
	}
	for k := LOGH - 1; k >= 0; k-- {
		if up[k][u] != up[k][v] {
			if mx[k][u] > res {
				res = mx[k][u]
			}
			if mx[k][v] > res {
				res = mx[k][v]
			}
			u = up[k][u]
			v = up[k][v]
		}
	}
	if mx[0][u] > res {
		res = mx[0][u]
	}
	if mx[0][v] > res {
		res = mx[0][v]
	}
	return res
}

type DSU struct{ parent, size, emax, a, b []int }

func newDSU(n int, e []int) *DSU {
	parent := make([]int, n+1)
	size := make([]int, n+1)
	emax := make([]int, n+1)
	a := make([]int, n+1)
	b := make([]int, n+1)
	for i := 1; i <= n; i++ {
		parent[i] = i
		size[i] = 1
		emax[i] = e[i]
		a[i] = i
		b[i] = i
	}
	return &DSU{parent, size, emax, a, b}
}

func (d *DSU) find(x int) int {
	if d.parent[x] != x {
		d.parent[x] = d.find(d.parent[x])
	}
	return d.parent[x]
}

func (d *DSU) union(u, v int) {
	ru := d.find(u)
	rv := d.find(v)
	if ru == rv {
		return
	}
	if d.size[ru] < d.size[rv] {
		ru, rv = rv, ru
	}
	d.parent[rv] = ru
	d.size[ru] += d.size[rv]
	if d.emax[ru] < d.emax[rv] {
		d.emax[ru] = d.emax[rv]
		d.a[ru] = d.a[rv]
		d.b[ru] = d.b[rv]
	} else if d.emax[ru] == d.emax[rv] {
		nodes := []int{d.a[ru], d.b[ru], d.a[rv], d.b[rv]}
		bestA := d.a[ru]
		bestB := d.b[ru]
		bestD := maxEdge(bestA, bestB)
		for i := 0; i < len(nodes); i++ {
			for j := i + 1; j < len(nodes); j++ {
				dtmp := maxEdge(nodes[i], nodes[j])
				if dtmp > bestD {
					bestD = dtmp
					bestA = nodes[i]
					bestB = nodes[j]
				}
			}
		}
		d.a[ru] = bestA
		d.b[ru] = bestB
	}
}

func expected(tc caseH) []string {
	g = make([][]struct {
		to   int
		toll int
	}, tc.n+1)
	for _, e := range tc.edges {
		g[e.u] = append(g[e.u], struct {
			to   int
			toll int
		}{e.v, e.t})
		g[e.v] = append(g[e.v], struct {
			to   int
			toll int
		}{e.u, e.t})
	}
	buildLCA(tc.n)
	edgesCopy := append([]edge(nil), tc.edges...)
	qs := make([]struct{ v, x, idx int }, tc.q)
	for i := 0; i < tc.q; i++ {
		qs[i] = struct{ v, x, idx int }{tc.queries[i].v, tc.queries[i].x, i}
	}
	sort.Slice(edgesCopy, func(i, j int) bool { return edgesCopy[i].c > edgesCopy[j].c })
	sort.Slice(qs, func(i, j int) bool { return qs[i].v > qs[j].v })
	d := newDSU(tc.n, tc.e)
	ansE := make([]int, tc.q)
	ansT := make([]int, tc.q)
	ei := 0
	for _, qu := range qs {
		for ei < len(edgesCopy) && edgesCopy[ei].c >= qu.v {
			d.union(edgesCopy[ei].u, edgesCopy[ei].v)
			ei++
		}
		root := d.find(qu.x)
		ansE[qu.idx] = d.emax[root]
		tollA := maxEdge(qu.x, d.a[root])
		tollB := maxEdge(qu.x, d.b[root])
		if tollA > tollB {
			ansT[qu.idx] = tollA
		} else {
			ansT[qu.idx] = tollB
		}
	}
	res := make([]string, tc.q)
	for i := 0; i < tc.q; i++ {
		res[i] = fmt.Sprintf("%d %d", ansE[i], ansT[i])
	}
	return res
}

func runCase(bin string, tc caseH) error {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d\n", tc.n, tc.q))
	for i := 1; i <= tc.n; i++ {
		if i > 1 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprintf("%d", tc.e[i]))
	}
	sb.WriteByte('\n')
	for _, e := range tc.edges {
		sb.WriteString(fmt.Sprintf("%d %d %d %d\n", e.u, e.v, e.c, e.t))
	}
	for _, qu := range tc.queries {
		sb.WriteString(fmt.Sprintf("%d %d\n", qu.v, qu.x))
	}
	input := sb.String()
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	outLines := strings.Split(strings.TrimSpace(out.String()), "\n")
	exp := expected(tc)
	if len(outLines) != len(exp) {
		return fmt.Errorf("expected %d lines got %d", len(exp), len(outLines))
	}
	for i, line := range outLines {
		fields := strings.Fields(line)
		if len(fields) != 2 {
			return fmt.Errorf("line %d should contain two ints", i+1)
		}
		if fields[0]+" "+fields[1] != exp[i] {
			return fmt.Errorf("line %d expected %s got %s", i+1, exp[i], line)
		}
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierH.go /path/to/binary")
		os.Exit(1)
	}
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	bin := os.Args[1]
	for i := 0; i < 100; i++ {
		tc := genCase(rng)
		if err := runCase(bin, tc); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", i+1, err)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
