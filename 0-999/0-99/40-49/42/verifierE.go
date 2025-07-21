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

type edge struct{ u, v, w int }

type DSU struct{ p, r []int }

func newDSU(n int) *DSU {
	p := make([]int, n)
	r := make([]int, n)
	for i := 0; i < n; i++ {
		p[i] = i
	}
	return &DSU{p, r}
}

func (d *DSU) find(x int) int {
	if d.p[x] != x {
		d.p[x] = d.find(d.p[x])
	}
	return d.p[x]
}

func (d *DSU) union(x, y int) bool {
	rx, ry := d.find(x), d.find(y)
	if rx == ry {
		return false
	}
	if d.r[rx] < d.r[ry] {
		d.p[rx] = ry
	} else if d.r[ry] < d.r[rx] {
		d.p[ry] = rx
	} else {
		d.p[ry] = rx
		d.r[rx]++
	}
	return true
}

func mstCost(n int, edges []edge) int64 {
	sort.Slice(edges, func(i, j int) bool { return edges[i].w < edges[j].w })
	d := newDSU(n)
	var total int64
	cnt := 0
	for _, e := range edges {
		if d.union(e.u, e.v) {
			total += int64(e.w)
			cnt++
			if cnt == n-1 {
				break
			}
		}
	}
	if cnt != n-1 {
		return -1
	}
	return total
}

func generateCase(rng *rand.Rand) (string, int64, int) {
	n := rng.Intn(4) + 2
	maxEdges := n * (n - 1) / 2
	m := rng.Intn(maxEdges) + 1
	edges := make([]edge, m)
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n%d\n", n, m))
	for i := 0; i < m; i++ {
		u := rng.Intn(n)
		v := rng.Intn(n)
		for v == u {
			v = rng.Intn(n)
		}
		w := rng.Intn(20) + 1
		edges[i] = edge{u, v, w}
		sb.WriteString(fmt.Sprintf("%d %d %d\n", u+1, v+1, w))
	}
	q := rng.Intn(5) + 1
	sb.WriteString(fmt.Sprintf("%d\n", q))
	for i := 0; i < q; i++ {
		a := rng.Intn(n)
		b := rng.Intn(n)
		for b == a {
			b = rng.Intn(n)
		}
		sb.WriteString(fmt.Sprintf("%d %d\n", a+1, b+1))
	}
	cost := mstCost(n, edges)
	return sb.String(), cost, q
}

func runCase(bin, input string, cost int64, q int) error {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	lines := strings.Split(strings.TrimSpace(out.String()), "\n")
	if len(lines) != q {
		return fmt.Errorf("expected %d lines, got %d", q, len(lines))
	}
	for i := 0; i < q; i++ {
		var val int64
		if _, err := fmt.Sscan(strings.TrimSpace(lines[i]), &val); err != nil {
			return fmt.Errorf("bad output line %d", i+1)
		}
		if val != cost {
			return fmt.Errorf("line %d: expected %d got %d", i+1, cost, val)
		}
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	bin := os.Args[1]
	for i := 0; i < 100; i++ {
		in, cost, q := generateCase(rng)
		if err := runCase(bin, in, cost, q); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:%s", i+1, err, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
