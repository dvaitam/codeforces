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
	d := &DSU{p: make([]int, n)}
	for i := range d.p {
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

func expected(n int, e1, e2 [][2]int) string {
	d1 := NewDSU(n)
	d2 := NewDSU(n)
	for _, e := range e1 {
		d1.Union(e[0]-1, e[1]-1)
	}
	for _, e := range e2 {
		d2.Union(e[0]-1, e[1]-1)
	}
	var ans [][2]int
	for i := 0; i < n; i++ {
		for j := i + 1; j < n; j++ {
			if d1.Find(i) != d1.Find(j) && d2.Find(i) != d2.Find(j) {
				ans = append(ans, [2]int{i + 1, j + 1})
				d1.Union(i, j)
				d2.Union(i, j)
			}
		}
	}
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", len(ans)))
	for _, e := range ans {
		sb.WriteString(fmt.Sprintf("%d %d\n", e[0], e[1]))
	}
	return strings.TrimSpace(sb.String())
}

func randomForest(n int, rng *rand.Rand) [][2]int {
	d := NewDSU(n)
	edges := make([][2]int, 0)
	m := rng.Intn(n)
	for len(edges) < m {
		u := rng.Intn(n)
		v := rng.Intn(n)
		if u == v || d.Find(u) == d.Find(v) {
			continue
		}
		d.Union(u, v)
		edges = append(edges, [2]int{u + 1, v + 1})
	}
	return edges
}

func generateCase(rng *rand.Rand) (string, string) {
	n := rng.Intn(8) + 2
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
	return sb.String(), expected(n, e1, e2)
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierD1.go /path/to/binary")
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
