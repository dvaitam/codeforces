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

type DSU struct {
	parent []int
	size   []int64
}

func NewDSU(n int) *DSU {
	p := make([]int, n)
	sz := make([]int64, n)
	for i := 0; i < n; i++ {
		p[i] = i
		sz[i] = 1
	}
	return &DSU{parent: p, size: sz}
}

func (d *DSU) Find(x int) int {
	if d.parent[x] != x {
		d.parent[x] = d.Find(d.parent[x])
	}
	return d.parent[x]
}

func (d *DSU) Union(x, y int, weight int) int64 {
	rx := d.Find(x)
	ry := d.Find(y)
	if rx == ry {
		return 0
	}
	if d.size[rx] < d.size[ry] {
		rx, ry = ry, rx
	}
	contrib := int64(weight) * d.size[rx] * d.size[ry]
	d.parent[ry] = rx
	d.size[rx] += d.size[ry]
	return contrib
}

func runBinary(bin, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func solveD(n int, a []int, edges [][2]int) string {
	g := make([][]int, n)
	for _, e := range edges {
		u, v := e[0]-1, e[1]-1
		g[u] = append(g[u], v)
		g[v] = append(g[v], u)
	}
	ord := make([]int, n)
	for i := 0; i < n; i++ {
		ord[i] = i
	}
	sort.Slice(ord, func(i, j int) bool { return a[ord[i]] > a[ord[j]] })
	dsu := NewDSU(n)
	active := make([]bool, n)
	var sum int64
	for _, u := range ord {
		active[u] = true
		for _, v := range g[u] {
			if active[v] {
				sum += dsu.Union(u, v, a[u])
			}
		}
	}
	denom := float64(n) * float64(n-1)
	avg := 2.0 * float64(sum) / denom
	return fmt.Sprintf("%.6f", avg)
}

func generateCase(rng *rand.Rand) (string, string) {
	n := rng.Intn(8) + 2
	maxEdges := n * (n - 1) / 2
	m := rng.Intn(maxEdges-n+1) + n - 1 // ensure connected by at least tree size
	// create tree first for connectivity
	edges := make([][2]int, 0, m)
	for i := 2; i <= n; i++ {
		p := rng.Intn(i-1) + 1
		edges = append(edges, [2]int{p, i})
	}
	used := make(map[[2]int]bool)
	for _, e := range edges {
		a, b := e[0], e[1]
		if a > b {
			a, b = b, a
		}
		used[[2]int{a, b}] = true
	}
	for len(edges) < m {
		a := rng.Intn(n) + 1
		b := rng.Intn(n) + 1
		if a == b {
			continue
		}
		aa, bb := a, b
		if aa > bb {
			aa, bb = bb, aa
		}
		if used[[2]int{aa, bb}] {
			continue
		}
		used[[2]int{aa, bb}] = true
		edges = append(edges, [2]int{a, b})
	}
	vals := make([]int, n)
	for i := range vals {
		vals[i] = rng.Intn(100)
	}
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d %d\n", n, len(edges))
	for i, v := range vals {
		if i > 0 {
			sb.WriteByte(' ')
		}
		fmt.Fprintf(&sb, "%d", v)
	}
	sb.WriteByte('\n')
	for _, e := range edges {
		fmt.Fprintf(&sb, "%d %d\n", e[0], e[1])
	}
	input := sb.String()
	exp := solveD(n, vals, edges)
	return input, exp
}

func manualCase() (string, string) {
	n := 3
	vals := []int{10, 20, 30}
	edges := [][2]int{{1, 2}, {2, 3}}
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d %d\n", n, len(edges))
	fmt.Fprintf(&sb, "10 20 30\n")
	for _, e := range edges {
		fmt.Fprintf(&sb, "%d %d\n", e[0], e[1])
	}
	input := sb.String()
	exp := solveD(n, vals, edges)
	return input, exp
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	var cases [][2]string
	in1, ex1 := manualCase()
	cases = append(cases, [2]string{in1, ex1})
	for i := 0; i < 100; i++ {
		in, exp := generateCase(rng)
		cases = append(cases, [2]string{in, exp})
	}
	for idx, tc := range cases {
		out, err := runBinary(bin, tc[0])
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", idx+1, err, tc[0])
			os.Exit(1)
		}
		if strings.TrimSpace(out) != tc[1] {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\ninput:\n%s", idx+1, tc[1], out, tc[0])
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(cases))
}
