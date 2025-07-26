package main

import (
	"bufio"
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
	"time"
)

type DSU struct {
	p []int
}

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
	fa := d.Find(a)
	fb := d.Find(b)
	if fa != fb {
		d.p[fa] = fb
	}
}

func generateCase(rng *rand.Rand) (string, int, [][2]int, int) {
	n := rng.Intn(5) + 2
	edges := make([][2]int, 0, n*(n-1)/2)
	used := make(map[[2]int]bool)
	for i := 2; i <= n; i++ {
		edges = append(edges, [2]int{i - 1, i})
		used[[2]int{i - 1, i}] = true
	}
	total := n * (n - 1) / 2
	extras := rng.Intn(total - (n - 1) + 1)
	for len(edges) < n-1+extras {
		u := rng.Intn(n) + 1
		v := rng.Intn(n) + 1
		if u == v {
			continue
		}
		if u > v {
			u, v = v, u
		}
		if used[[2]int{u, v}] {
			continue
		}
		used[[2]int{u, v}] = true
		edges = append(edges, [2]int{u, v})
	}
	deg := make([]int, n+1)
	for _, e := range edges {
		deg[e[0]]++
		deg[e[1]]++
	}
	maxDeg := 0
	for i := 1; i <= n; i++ {
		if deg[i] > maxDeg {
			maxDeg = deg[i]
		}
	}
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d\n", n, len(edges)))
	for _, e := range edges {
		sb.WriteString(fmt.Sprintf("%d %d\n", e[0], e[1]))
	}
	return sb.String(), maxDeg, edges, n
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
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, stderr.String())
	}
	return out.String(), nil
}

func checkOutput(n int, edgesIn [][2]int, expect int, out string) error {
	r := bufio.NewReader(strings.NewReader(out))
	edgesOut := make([][2]int, 0, n-1)
	for i := 0; i < n-1; i++ {
		var u, v int
		_, err := fmt.Fscan(r, &u, &v)
		if err != nil {
			return fmt.Errorf("invalid output")
		}
		edgesOut = append(edgesOut, [2]int{u, v})
	}
	// check no extra integers
	if _, err := fmt.Fscan(r, new(int)); err == nil {
		return fmt.Errorf("extra data")
	}
	allowed := make(map[[2]int]bool)
	for _, e := range edgesIn {
		a, b := e[0], e[1]
		if a > b {
			a, b = b, a
		}
		allowed[[2]int{a, b}] = true
	}
	dsu := NewDSU(n)
	deg := make([]int, n+1)
	for _, e := range edgesOut {
		u, v := e[0], e[1]
		if u < 1 || u > n || v < 1 || v > n {
			return fmt.Errorf("vertex out of range")
		}
		a, b := u, v
		if a > b {
			a, b = b, a
		}
		if !allowed[[2]int{a, b}] {
			return fmt.Errorf("edge not in input")
		}
		if dsu.Find(u) == dsu.Find(v) {
			return fmt.Errorf("not a tree")
		}
		dsu.Union(u, v)
		deg[u]++
		deg[v]++
	}
	root := dsu.Find(1)
	for i := 2; i <= n; i++ {
		if dsu.Find(i) != root {
			return fmt.Errorf("not a tree")
		}
	}
	maxd := 0
	for i := 1; i <= n; i++ {
		if deg[i] > maxd {
			maxd = deg[i]
		}
	}
	if maxd != expect {
		return fmt.Errorf("max degree %d, expected %d", maxd, expect)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 && !(len(os.Args) == 3 && os.Args[1] == "--") {
		fmt.Println("usage: go run verifierF1.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[len(os.Args)-1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		in, expect, edges, n := generateCase(rng)
		out, err := run(bin, in)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}
		if err := checkOutput(n, edges, expect, out); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%soutput:\n%s", i+1, err, in, out)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
