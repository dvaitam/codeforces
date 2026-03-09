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


// countWays enumerates all spanning trees of the graph and counts those with exactly k leaves.
func countWays(n int, edges [][2]int, k int) int64 {
	// deduplicate edges
	type Edge [2]int
	seen := make(map[Edge]bool)
	var deduped [][2]int
	for _, e := range edges {
		a, b := e[0], e[1]
		if a > b {
			a, b = b, a
		}
		key := Edge{a, b}
		if !seen[key] {
			seen[key] = true
			deduped = append(deduped, [2]int{a, b})
		}
	}
	edges = deduped

	m := len(edges)
	need := n - 1
	var total int64
	sel := make([]int, 0, need)

	var rec func(start int)
	rec = func(start int) {
		if len(sel) == need {
			// check acyclicity via union-find; n-1 acyclic edges => spanning tree
			parent := make([]int, n)
			for i := range parent {
				parent[i] = i
			}
			var find func(int) int
			find = func(x int) int {
				if parent[x] != x {
					parent[x] = find(parent[x])
				}
				return parent[x]
			}
			deg := make([]int, n)
			for _, idx := range sel {
				u, v := edges[idx][0], edges[idx][1]
				pu, pv := find(u), find(v)
				if pu == pv {
					return // cycle => not a tree
				}
				parent[pu] = pv
				deg[u]++
				deg[v]++
			}
			leaves := 0
			for i := 0; i < n; i++ {
				if deg[i] == 1 {
					leaves++
				}
			}
			if leaves == k {
				total++
			}
			return
		}
		if m-start < need-len(sel) {
			return
		}
		for i := start; i < m; i++ {
			sel = append(sel, i)
			rec(i + 1)
			sel = sel[:len(sel)-1]
		}
	}
	rec(0)
	return total
}

func generateGraph(rng *rand.Rand) (string, int64) {
	n := rng.Intn(5) + 3
	// build connected graph with no duplicate edges
	type Edge [2]int
	edgeSet := make(map[Edge]bool)
	var edges [][2]int
	addEdge := func(a, b int) {
		if a > b {
			a, b = b, a
		}
		key := Edge{a, b}
		if !edgeSet[key] {
			edgeSet[key] = true
			edges = append(edges, [2]int{a, b})
		}
	}
	for i := 1; i < n; i++ {
		p := rng.Intn(i)
		addEdge(p, i)
	}
	for i := 0; i < n; i++ {
		for j := i + 1; j < n; j++ {
			if rng.Float64() < 0.3 {
				addEdge(i, j)
			}
		}
	}
	m := len(edges)
	k := rng.Intn(n-2) + 2
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d %d\n", n, m, k))
	for _, e := range edges {
		sb.WriteString(fmt.Sprintf("%d %d\n", e[0]+1, e[1]+1))
	}
	exp := countWays(n, edges, k)
	return sb.String(), exp
}

func runCase(exe, input string, expected int64) error {
	cmd := exec.Command(exe)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	outStr := strings.TrimSpace(out.String())
	got, err := strconv.ParseInt(outStr, 10, 64)
	if err != nil {
		return fmt.Errorf("bad output: %v", err)
	}
	if got != expected {
		return fmt.Errorf("expected %d got %d", expected, got)
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
		in, exp := generateGraph(rng)
		if err := runCase(exe, in, exp); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
