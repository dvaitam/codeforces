package main

import (
	"bytes"
	"fmt"
	"math/bits"
	"math/rand"
	"os"
	"os/exec"
	"strings"
	"time"
)

type Edge struct{ u, v int }

func brute(n int, edges []Edge) []int {
	best := make([]int, n+1)
	adj := make([][]int, n)
	for _, e := range edges {
		adj[e.u] = append(adj[e.u], e.v)
		adj[e.v] = append(adj[e.v], e.u)
	}
	totalMask := 1 << n
	for maskA := 0; maskA < totalMask; maskA++ {
		countA := bits.OnesCount(uint(maskA))
		for maskB := 0; maskB < totalMask; maskB++ {
			if maskA&maskB != 0 {
				continue
			}
			ok := true
			for _, e := range edges {
				uInA := maskA>>e.u&1 == 1
				vInA := maskA>>e.v&1 == 1
				uInB := maskB>>e.u&1 == 1
				vInB := maskB>>e.v&1 == 1
				if (uInA && vInB) || (vInA && uInB) {
					ok = false
					break
				}
			}
			if !ok {
				continue
			}
			countB := bits.OnesCount(uint(maskB))
			if countB > best[countA] {
				best[countA] = countB
			}
		}
	}
	return best
}

func runCase(exe string, n int, edges []Edge) error {
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d %d\n", n, len(edges))
	for _, e := range edges {
		fmt.Fprintf(&sb, "%d %d\n", e.u+1, e.v+1)
	}
	input := sb.String()
	cmd := exec.Command(exe)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	fields := strings.Fields(out.String())
	expect := brute(n, edges)
	if len(fields) != n+1 {
		return fmt.Errorf("expected %d numbers got %d", n+1, len(fields))
	}
	for i := 0; i <= n; i++ {
		var v int
		if _, err := fmt.Sscan(fields[i], &v); err != nil {
			return fmt.Errorf("parse error: %v", err)
		}
		if v != expect[i] {
			return fmt.Errorf("index %d expected %d got %d", i, expect[i], v)
		}
	}
	return nil
}

func genTree(rng *rand.Rand, n int) []Edge {
	edges := make([]Edge, 0, n-1)
	for i := 1; i < n; i++ {
		p := rng.Intn(i)
		edges = append(edges, Edge{p, i})
	}
	return edges
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierG.go /path/to/binary")
		os.Exit(1)
	}
	exe := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	cases := []struct {
		n     int
		edges []Edge
	}{
		{1, []Edge{}},
		{2, []Edge{{0, 1}}},
	}
	for i := 0; i < 100; i++ {
		n := rng.Intn(6) + 1
		edges := genTree(rng, n)
		cases = append(cases, struct {
			n     int
			edges []Edge
		}{n, edges})
	}
	for idx, c := range cases {
		if err := runCase(exe, c.n, c.edges); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", idx+1, err)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
