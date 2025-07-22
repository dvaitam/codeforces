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

type edge struct{ u, v int }

type dsu struct{ p []int }

func newDSU(n int) *dsu {
	p := make([]int, n)
	for i := 0; i < n; i++ {
		p[i] = i
	}
	return &dsu{p}
}
func (d *dsu) find(x int) int {
	for d.p[x] != x {
		x = d.p[x]
	}
	return x
}
func (d *dsu) union(a, b int) {
	ra, rb := d.find(a), d.find(b)
	if ra != rb {
		d.p[rb] = ra
	}
}

func components(n int, edges []edge) int {
	d := newDSU(n)
	for _, e := range edges {
		d.union(e.u, e.v)
	}
	seen := make(map[int]bool)
	for i := 0; i < n; i++ {
		seen[d.find(i)] = true
	}
	return len(seen)
}

func generateCase(rng *rand.Rand) (string, string) {
	n := rng.Intn(5) + 2 // 2..6
	maxEdges := n * (n - 1) / 2
	m := rng.Intn(maxEdges) + 1
	edges := make([]edge, 0, m)
	exist := make(map[[2]int]bool)
	for len(edges) < m {
		u := rng.Intn(n)
		v := rng.Intn(n)
		if u == v {
			continue
		}
		a, b := u, v
		if a > b {
			a, b = b, a
		}
		if exist[[2]int{a, b}] {
			continue
		}
		exist[[2]int{a, b}] = true
		edges = append(edges, edge{u, v})
	}
	q := rng.Intn(5) + 1
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d %d\n", n, m)
	for _, e := range edges {
		fmt.Fprintf(&sb, "%d %d\n", e.u+1, e.v+1)
	}
	fmt.Fprintf(&sb, "%d\n", q)
	var out strings.Builder
	for i := 0; i < q; i++ {
		l := rng.Intn(m) + 1
		r := rng.Intn(m-l+1) + l
		fmt.Fprintf(&sb, "%d %d\n", l, r)
		var remaining []edge
		for idx, e := range edges {
			pos := idx + 1
			if pos < l || pos > r {
				remaining = append(remaining, e)
			}
		}
		comp := components(n, remaining)
		fmt.Fprintf(&out, "%d\n", comp)
	}
	return sb.String(), strings.TrimSpace(out.String())
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
	got := strings.TrimSpace(out.String())
	if got != expected {
		return fmt.Errorf("expected:\n%s\n\ngot:\n%s", expected, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierD.go /path/to/binary")
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
