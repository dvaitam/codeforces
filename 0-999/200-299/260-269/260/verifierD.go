package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

type edge struct {
	u, v, w int
}

type dsu struct {
	parent []int
	rank   []int
}

func newDSU(n int) *dsu {
	d := &dsu{parent: make([]int, n+1), rank: make([]int, n+1)}
	for i := 1; i <= n; i++ {
		d.parent[i] = i
	}
	return d
}

func (d *dsu) find(x int) int {
	if d.parent[x] != x {
		d.parent[x] = d.find(d.parent[x])
	}
	return d.parent[x]
}

func (d *dsu) union(x, y int) {
	x = d.find(x)
	y = d.find(y)
	if x == y {
		return
	}
	if d.rank[x] < d.rank[y] {
		x, y = y, x
	}
	d.parent[y] = x
	if d.rank[x] == d.rank[y] {
		d.rank[x]++
	}
}

func generateCase(rng *rand.Rand) (int, []int, []int) {
	n := rng.Intn(10) + 2
	edges := make([]edge, 0, n-1)
	for i := 2; i <= n; i++ {
		p := rng.Intn(i-1) + 1
		w := rng.Intn(20)
		edges = append(edges, edge{u: i, v: p, w: w})
	}
	color := make([]int, n+1)
	for i := 2; i <= n; i++ {
		color[i] = 1 - color[edges[i-2].v]
	}
	sum := make([]int, n+1)
	for _, e := range edges {
		sum[e.u] += e.w
		sum[e.v] += e.w
	}
	input := make([]int, 0, 2*n)
	for i := 1; i <= n; i++ {
		input = append(input, color[i], sum[i])
	}
	return n, color[1:], sum[1:]
}

func parseOutput(out string) ([]edge, error) {
	fields := strings.Fields(strings.TrimSpace(out))
	if len(fields)%3 != 0 {
		return nil, fmt.Errorf("expected triples")
	}
	m := len(fields) / 3
	res := make([]edge, m)
	for i := 0; i < m; i++ {
		u, err1 := strconv.Atoi(fields[3*i])
		v, err2 := strconv.Atoi(fields[3*i+1])
		w, err3 := strconv.Atoi(fields[3*i+2])
		if err1 != nil || err2 != nil || err3 != nil {
			return nil, fmt.Errorf("bad numbers")
		}
		res[i] = edge{u: u, v: v, w: w}
	}
	return res, nil
}

func verify(n int, color, sum []int, edges []edge) error {
	if len(edges) != n-1 {
		return fmt.Errorf("expected %d edges got %d", n-1, len(edges))
	}
	d := newDSU(n)
	cur := make([]int, n+1)
	for _, e := range edges {
		if e.u < 1 || e.u > n || e.v < 1 || e.v > n || e.u == e.v || e.w < 0 {
			return fmt.Errorf("invalid edge")
		}
		if color[e.u-1] == color[e.v-1] {
			return fmt.Errorf("edge with same colors")
		}
		cur[e.u] += e.w
		cur[e.v] += e.w
		d.union(e.u, e.v)
	}
	root := d.find(1)
	for i := 2; i <= n; i++ {
		if d.find(i) != root {
			return fmt.Errorf("graph not connected")
		}
	}
	for i := 1; i <= n; i++ {
		if cur[i] != sum[i-1] {
			return fmt.Errorf("sum mismatch at %d", i)
		}
	}
	return nil
}

func runCase(bin string, n int, color, sum []int) error {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", n))
	for i := 0; i < n; i++ {
		sb.WriteString(fmt.Sprintf("%d %d\n", color[i], sum[i]))
	}

	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(sb.String())
	var out bytes.Buffer
	cmd.Stdout = &out
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\nstderr: %s", err, stderr.String())
	}
	edges, err := parseOutput(out.String())
	if err != nil {
		return err
	}
	return verify(n, color, sum, edges)
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		n, color, sum := generateCase(rng)
		if err := runCase(bin, n, color, sum); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", i+1, err)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
