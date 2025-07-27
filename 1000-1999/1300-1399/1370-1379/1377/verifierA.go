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

type graph struct {
	n     int
	edges []edge
}

func generateGraphs() []graph {
	r := rand.New(rand.NewSource(1377))
	tests := make([]graph, 100)
	for i := 0; i < 100; i++ {
		n := r.Intn(8) + 2 // 2..9 nodes
		m := r.Intn(n*(n-1)/2) + 1
		seen := make(map[[2]int]bool)
		edges := make([]edge, 0, m)
		for len(edges) < m {
			a := r.Intn(n)
			b := r.Intn(n)
			if a == b {
				continue
			}
			if a > b {
				a, b = b, a
			}
			key := [2]int{a, b}
			if !seen[key] {
				seen[key] = true
				edges = append(edges, edge{a, b})
			}
		}
		// ensure every node has at least one incident edge
		deg := make([]int, n)
		for _, e := range edges {
			deg[e.u]++
			deg[e.v]++
		}
		for v := 0; v < n; v++ {
			if deg[v] == 0 {
				u := v
				for u == v {
					u = r.Intn(n)
				}
				a, b := v, u
				if a > b {
					a, b = b, a
				}
				k := [2]int{a, b}
				if !seen[k] {
					seen[k] = true
					edges = append(edges, edge{a, b})
				}
			}
		}
		tests[i] = graph{n, edges}
	}
	return tests
}

func runCase(bin string, g graph) error {
	var input strings.Builder
	for _, e := range g.edges {
		fmt.Fprintf(&input, "%d %d\n", e.u, e.v)
	}
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input.String())
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	outStr := strings.TrimSpace(out.String())
	if outStr == "" {
		return fmt.Errorf("empty output")
	}
	lines := strings.Split(outStr, "\n")
	seen := make([]bool, g.n)
	for _, line := range lines {
		fields := strings.Fields(line)
		if len(fields) == 0 {
			return fmt.Errorf("empty line in output")
		}
		for _, f := range fields {
			var x int
			if _, err := fmt.Sscan(f, &x); err != nil {
				return fmt.Errorf("invalid number %q", f)
			}
			if x < 0 || x >= g.n {
				return fmt.Errorf("node %d out of range", x)
			}
			if seen[x] {
				return fmt.Errorf("node %d appears multiple times", x)
			}
			seen[x] = true
		}
	}
	for i, ok := range seen {
		if !ok {
			return fmt.Errorf("node %d missing", i)
		}
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rand.Seed(time.Now().UnixNano())
	cases := generateGraphs()
	for i, g := range cases {
		if err := runCase(bin, g); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", i+1, err)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
