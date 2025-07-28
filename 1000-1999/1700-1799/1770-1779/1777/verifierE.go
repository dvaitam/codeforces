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

type dsu struct {
	p  []int
	sz []int
}

func newDSU(n int) *dsu {
	d := &dsu{p: make([]int, n), sz: make([]int, n)}
	for i := 0; i < n; i++ {
		d.p[i] = i
		d.sz[i] = 1
	}
	return d
}

func (d *dsu) find(x int) int {
	if d.p[x] != x {
		d.p[x] = d.find(d.p[x])
	}
	return d.p[x]
}

func (d *dsu) union(a, b int) {
	a = d.find(a)
	b = d.find(b)
	if a == b {
		return
	}
	if d.sz[a] < d.sz[b] {
		a, b = b, a
	}
	d.p[b] = a
	d.sz[a] += d.sz[b]
}

type Edge struct {
	u, v int
	w    int
}

func can(edges []Edge, n int, x int) bool {
	d := newDSU(n)
	for _, e := range edges {
		if e.w <= x {
			d.union(e.u, e.v)
		}
	}
	compID := make(map[int]int)
	id := 0
	for i := 0; i < n; i++ {
		r := d.find(i)
		if _, ok := compID[r]; !ok {
			compID[r] = id
			id++
		}
	}
	indeg := make([]int, id)
	for _, e := range edges {
		if e.w > x {
			a := compID[d.find(e.u)]
			b := compID[d.find(e.v)]
			if a != b {
				indeg[b]++
			}
		}
	}
	cnt := 0
	for i := 0; i < id; i++ {
		if indeg[i] == 0 {
			cnt++
			if cnt > 1 {
				return false
			}
		}
	}
	return cnt == 1
}

func solveCase(n, m int, edges []Edge) int {
	d := newDSU(n)
	for _, e := range edges {
		d.union(e.u, e.v)
	}
	root := d.find(0)
	for i := 1; i < n; i++ {
		if d.find(i) != root {
			return -1
		}
	}
	lo, hi := 0, int(1e9)
	for lo < hi {
		mid := (lo + hi) / 2
		if can(edges, n, mid) {
			hi = mid
		} else {
			lo = mid + 1
		}
	}
	return lo
}

func generateCase(rng *rand.Rand) (string, string) {
	n := rng.Intn(5) + 1
	m := rng.Intn(5) + n - 1
	edges := make([]Edge, m)
	// start with tree
	for i := 0; i < n-1; i++ {
		edges[i] = Edge{u: i, v: i + 1, w: rng.Intn(20)}
	}
	for i := n - 1; i < m; i++ {
		u := rng.Intn(n)
		v := rng.Intn(n)
		for u == v {
			v = rng.Intn(n)
		}
		edges[i] = Edge{u: u, v: v, w: rng.Intn(20)}
	}
	ans := solveCase(n, m, edges)
	var sb strings.Builder
	sb.WriteString("1\n")
	sb.WriteString(fmt.Sprintf("%d %d\n", n, m))
	for _, e := range edges {
		sb.WriteString(fmt.Sprintf("%d %d %d\n", e.u+1, e.v+1, e.w))
	}
	return sb.String(), fmt.Sprintf("%d\n", ans)
}

func runCase(exe, input, expected string) error {
	cmd := exec.Command(exe)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	outStr := strings.TrimSpace(out.String())
	exp := strings.TrimSpace(expected)
	if outStr != exp {
		return fmt.Errorf("expected %s got %s", exp, outStr)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	exe := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		in, exp := generateCase(rng)
		if err := runCase(exe, in, exp); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
