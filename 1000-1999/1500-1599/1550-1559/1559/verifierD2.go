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

type DSU struct{ p []int }

func NewDSU(n int) *DSU {
	d := &DSU{p: make([]int, n+1)}
	for i := 1; i <= n; i++ {
		d.p[i] = i
	}
	return d
}

func (d *DSU) Find(x int) int {
	if d.p[x] != x {
		d.p[x] = d.Find(d.p[x])
	}
	return d.p[x]
}

func (d *DSU) Union(a, b int) {
	ra, rb := d.Find(a), d.Find(b)
	if ra != rb {
		d.p[ra] = rb
	}
}

func run(bin, input string) (string, error) {
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

func solve(n int, e1, e2 [][2]int) string {
	d1 := NewDSU(n)
	d2 := NewDSU(n)
	for _, e := range e1 {
		d1.Union(e[0], e[1])
	}
	for _, e := range e2 {
		d2.Union(e[0], e[1])
	}
	type pair struct{ x, y int }
	var res []pair
	for i := 2; i <= n; i++ {
		if d1.Find(i) != d1.Find(1) && d2.Find(i) != d2.Find(1) {
			res = append(res, pair{i, 1})
			d1.Union(i, 1)
			d2.Union(i, 1)
		}
	}
	var s1, s2 []int
	r1 := d1.Find(1)
	r2 := d2.Find(1)
	for i := 2; i <= n; i++ {
		if d1.Find(i) == i && d1.Find(i) != r1 {
			s1 = append(s1, i)
		}
		if d2.Find(i) == i && d2.Find(i) != r2 {
			s2 = append(s2, i)
		}
	}
	t := len(s1)
	if len(s2) < t {
		t = len(s2)
	}
	for i := 0; i < t; i++ {
		res = append(res, pair{s1[i], s2[i]})
	}
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", len(res)))
	for _, p := range res {
		sb.WriteString(fmt.Sprintf("%d %d\n", p.x, p.y))
	}
	return strings.TrimSpace(sb.String())
}

func randomForest(n int, rng *rand.Rand) [][2]int {
	d := NewDSU(n)
	edges := make([][2]int, 0)
	m := rng.Intn(n)
	for len(edges) < m {
		u := rng.Intn(n) + 1
		v := rng.Intn(n) + 1
		if u == v || d.Find(u) == d.Find(v) {
			continue
		}
		d.Union(u, v)
		edges = append(edges, [2]int{u, v})
	}
	return edges
}

func generateCase(rng *rand.Rand) (string, string) {
	n := rng.Intn(20) + 2
	e1 := randomForest(n, rng)
	e2 := randomForest(n, rng)
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d %d\n", n, len(e1), len(e2)))
	for _, e := range e1 {
		sb.WriteString(fmt.Sprintf("%d %d\n", e[0], e[1]))
	}
	for _, e := range e2 {
		sb.WriteString(fmt.Sprintf("%d %d\n", e[0], e[1]))
	}
	return sb.String(), solve(n, e1, e2)
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierD2.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		input, want := generateCase(rng)
		got, err := run(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d: %v\ninput:\n%s", i+1, err, input)
			os.Exit(1)
		}
		if got != want {
			fmt.Fprintf(os.Stderr, "case %d failed: expected\n%s\ngot\n%s\ninput:\n%s", i+1, want, got, input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
