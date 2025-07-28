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

type Edge struct {
	u, v int
	w    int
}

type Query struct {
	i, k int
	idx  int
}

type DSU struct {
	parent []int
	size   []int
}

func NewDSU(n int) *DSU {
	d := &DSU{parent: make([]int, n), size: make([]int, n)}
	for i := 0; i < n; i++ {
		d.parent[i] = i
		d.size[i] = 1
	}
	return d
}

func (d *DSU) find(x int) int {
	if d.parent[x] != x {
		d.parent[x] = d.find(d.parent[x])
	}
	return d.parent[x]
}

func (d *DSU) union(x, y int) {
	x = d.find(x)
	y = d.find(y)
	if x == y {
		return
	}
	if d.size[x] < d.size[y] {
		x, y = y, x
	}
	d.parent[y] = x
	d.size[x] += d.size[y]
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func weight(x, y, d int) int {
	diff := abs(x - y)
	return abs(diff - d)
}

func computeExpected(n, q, s, dVal int, a []int, queries []Query) []string {
	edges := make([]Edge, 0, 8*n)
	for i := 0; i < n-1; i++ {
		edges = append(edges, Edge{u: i, v: i + 1, w: weight(a[i], a[i+1], dVal)})
	}
	for i := 0; i < n; i++ {
		x := a[i]
		j := sort.SearchInts(a, x+dVal)
		if j < n {
			edges = append(edges, Edge{u: i, v: j, w: weight(x, a[j], dVal)})
		}
		if j-1 >= 0 {
			edges = append(edges, Edge{u: i, v: j - 1, w: weight(x, a[j-1], dVal)})
		}
		j = sort.SearchInts(a, x-dVal)
		if j < n {
			edges = append(edges, Edge{u: i, v: j, w: weight(x, a[j], dVal)})
		}
		if j-1 >= 0 {
			edges = append(edges, Edge{u: i, v: j - 1, w: weight(x, a[j-1], dVal)})
		}
	}
	sort.Slice(edges, func(i, j int) bool { return edges[i].w < edges[j].w })
	sort.Slice(queries, func(i, j int) bool { return queries[i].k < queries[j].k })
	dsu := NewDSU(n)
	ans := make([]string, q)
	eidx := 0
	for _, qu := range queries {
		for eidx < len(edges) && edges[eidx].w <= qu.k {
			dsu.union(edges[eidx].u, edges[eidx].v)
			eidx++
		}
		if dsu.find(s) == dsu.find(qu.i) {
			ans[qu.idx] = "Yes"
		} else {
			ans[qu.idx] = "No"
		}
	}
	return ans
}

func generateCase(rng *rand.Rand) (string, []string) {
	n := rng.Intn(6) + 2
	q := rng.Intn(5) + 1
	s := rng.Intn(n) + 1
	dVal := rng.Intn(10) + 1
	a := make([]int, n)
	a[0] = rng.Intn(5)
	for i := 1; i < n; i++ {
		a[i] = a[i-1] + rng.Intn(5) + 1
	}
	queries := make([]Query, q)
	for i := 0; i < q; i++ {
		queries[i].i = rng.Intn(n)
		queries[i].k = rng.Intn(15)
		queries[i].idx = i
	}
	exp := computeExpected(n, q, s-1, dVal, a, append([]Query(nil), queries...))
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d %d %d\n", n, q, s, dVal))
	for i, v := range a {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprintf("%d", v))
	}
	sb.WriteByte('\n')
	for i := 0; i < q; i++ {
		sb.WriteString(fmt.Sprintf("%d %d\n", queries[i].i+1, queries[i].k))
	}
	return sb.String(), exp
}

func runCase(bin string, input string, expected []string) error {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	fields := strings.Fields(strings.TrimSpace(out.String()))
	if len(fields) != len(expected) {
		return fmt.Errorf("expected %d lines got %d", len(expected), len(fields))
	}
	for i, f := range fields {
		if f != expected[i] {
			return fmt.Errorf("expected %v got %v", expected, fields)
		}
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierF.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		in, exp := generateCase(rng)
		if err := runCase(bin, in, exp); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
