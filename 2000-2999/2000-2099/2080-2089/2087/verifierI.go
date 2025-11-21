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

type edge struct {
	u, v int
}

type testCase struct {
	input  string
	n, m   int
	maxDeg int
	edges  []edge
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierI.go /path/to/binary")
		os.Exit(1)
	}
	candidate := os.Args[1]

	tests := generateTests()
	for i, tc := range tests {
		out, err := runCandidate(candidate, tc.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "runtime error on test %d: %v\ninput:\n%s", i+1, err, tc.input)
			os.Exit(1)
		}
		if err := validate(tc, out); err != nil {
			fmt.Fprintf(os.Stderr, "wrong answer on test %d: %v\ninput:\n%soutput:\n%s\n", i+1, err, tc.input, out)
			os.Exit(1)
		}
	}

	fmt.Printf("All %d tests passed.\n", len(tests))
}

func runCandidate(path, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(path, ".go") {
		cmd = exec.Command("go", "run", path)
	} else {
		cmd = exec.Command(path)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return out.String(), fmt.Errorf("%v\n%s", err, stderr.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func validate(tc testCase, output string) error {
	fields := strings.Fields(output)
	ptr := 0
	nextInt := func() (int, error) {
		if ptr >= len(fields) {
			return 0, fmt.Errorf("not enough tokens")
		}
		var x int
		_, err := fmt.Sscan(fields[ptr], &x)
		ptr++
		return x, err
	}

	k, err := nextInt()
	if err != nil {
		return fmt.Errorf("failed to read k: %v", err)
	}
	if k < 0 {
		return fmt.Errorf("k must be non-negative")
	}

	added := make([]edge, 0, k)
	for i := 0; i < k; i++ {
		u, err := nextInt()
		if err != nil {
			return fmt.Errorf("failed to read added edge %d u: %v", i+1, err)
		}
		v, err := nextInt()
		if err != nil {
			return fmt.Errorf("failed to read added edge %d v: %v", i+1, err)
		}
		if u < 1 || u > tc.n || v < 1 || v > tc.n {
			return fmt.Errorf("added edge %d has vertex out of range", i+1)
		}
		if u == v {
			return fmt.Errorf("added edge %d is a self-loop", i+1)
		}
		added = append(added, edge{u, v})
	}

	c, err := nextInt()
	if err != nil {
		return fmt.Errorf("failed to read c: %v", err)
	}
	if c <= 0 {
		return fmt.Errorf("c must be positive")
	}

	assign := make([]int, tc.m+k)
	for i := 0; i < len(assign); i++ {
		val, err := nextInt()
		if err != nil {
			return fmt.Errorf("failed to read assignment %d: %v", i+1, err)
		}
		assign[i] = val
	}
	if ptr != len(fields) {
		return fmt.Errorf("unexpected extra tokens at the end")
	}

	// Minimality checks.
	if c != tc.maxDeg {
		return fmt.Errorf("c=%d but minimal required is %d", c, tc.maxDeg)
	}
	kMin := tc.n*c - tc.m
	if k != kMin {
		return fmt.Errorf("k=%d but minimal required is %d", k, kMin)
	}

	// Assignment validity.
	for i, v := range assign {
		if v < 1 || v > c {
			return fmt.Errorf("assignment %d has invalid cycle id %d", i+1, v)
		}
	}

	allEdges := make([]edge, 0, tc.m+k)
	allEdges = append(allEdges, tc.edges...)
	allEdges = append(allEdges, added...)

	perColor := make([]int, c+1)
	for _, col := range assign {
		perColor[col]++
	}
	for col := 1; col <= c; col++ {
		if perColor[col] != tc.n {
			return fmt.Errorf("cycle %d has %d edges, expected %d", col, perColor[col], tc.n)
		}
	}

	for col := 1; col <= c; col++ {
		next := make([]int, tc.n+1)
		indeg := make([]int, tc.n+1)
		outdeg := make([]int, tc.n+1)
		for idx, e := range allEdges {
			if assign[idx] != col {
				continue
			}
			outdeg[e.u]++
			indeg[e.v]++
			if next[e.u] != 0 {
				return fmt.Errorf("vertex %d has multiple outgoing edges in cycle %d", e.u, col)
			}
			next[e.u] = e.v
		}
		for v := 1; v <= tc.n; v++ {
			if indeg[v] != 1 || outdeg[v] != 1 {
				return fmt.Errorf("vertex %d has indeg=%d outdeg=%d in cycle %d", v, indeg[v], outdeg[v], col)
			}
		}
		visited := make([]bool, tc.n+1)
		cur := 1
		for steps := 0; steps < tc.n; steps++ {
			if visited[cur] {
				return fmt.Errorf("cycle %d repeats a vertex before visiting all nodes", col)
			}
			visited[cur] = true
			cur = next[cur]
		}
		if cur != 1 {
			return fmt.Errorf("cycle %d does not return to the start after %d steps", col, tc.n)
		}
		for v := 1; v <= tc.n; v++ {
			if !visited[v] {
				return fmt.Errorf("cycle %d misses vertex %d", col, v)
			}
		}
	}

	return nil
}

func generateTests() []testCase {
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	var tests []testCase

	// Tiny graphs.
	tests = append(tests, dagPath(2))
	tests = append(tests, dagStar(4))

	// A couple of deterministic scenarios.
	tests = append(tests, dagComplete(4))
	tests = append(tests, dagComplete(6))

	for i := 0; i < 40; i++ {
		n := rng.Intn(8) + 2 // 2..9
		tests = append(tests, randomDag(rng, n, 0.35))
	}

	// Larger stress cases.
	tests = append(tests, randomDag(rng, 25, 0.25))
	tests = append(tests, randomDag(rng, 40, 0.35))
	tests = append(tests, randomDag(rng, 70, 0.2))
	tests = append(tests, randomDag(rng, 100, 0.15))

	return tests
}

func dagPath(n int) testCase {
	var edges []edge
	for i := 1; i < n; i++ {
		edges = append(edges, edge{i, i + 1})
	}
	return buildCase(n, edges)
}

func dagStar(n int) testCase {
	var edges []edge
	for i := 2; i <= n; i++ {
		edges = append(edges, edge{1, i})
	}
	return buildCase(n, edges)
}

func dagComplete(n int) testCase {
	var edges []edge
	for i := 1; i <= n; i++ {
		for j := i + 1; j <= n; j++ {
			edges = append(edges, edge{i, j})
		}
	}
	return buildCase(n, edges)
}

func randomDag(rng *rand.Rand, n int, prob float64) testCase {
	order := rng.Perm(n)
	pos := make([]int, n)
	for idx, val := range order {
		pos[val] = idx
	}
	var edges []edge
	for i := 0; i < n; i++ {
		for j := 0; j < n; j++ {
			if pos[i] < pos[j] && rng.Float64() < prob {
				edges = append(edges, edge{i + 1, j + 1})
			}
		}
	}
	if len(edges) == 0 {
		// Ensure at least one edge.
		u := order[0] + 1
		v := order[1] + 1
		edges = append(edges, edge{u, v})
	}
	return buildCase(n, edges)
}

func buildCase(n int, edges []edge) testCase {
	unique := make(map[[2]int]struct{})
	var filt []edge
	for _, e := range edges {
		key := [2]int{e.u, e.v}
		if e.u == e.v {
			continue
		}
		if _, ok := unique[key]; ok {
			continue
		}
		unique[key] = struct{}{}
		filt = append(filt, e)
	}
	edges = filt
	m := len(edges)

	indeg := make([]int, n+1)
	outdeg := make([]int, n+1)
	for _, e := range edges {
		outdeg[e.u]++
		indeg[e.v]++
	}
	maxDeg := 0
	for i := 1; i <= n; i++ {
		if indeg[i] > maxDeg {
			maxDeg = indeg[i]
		}
		if outdeg[i] > maxDeg {
			maxDeg = outdeg[i]
		}
	}
	if maxDeg == 0 {
		maxDeg = 1
	}

	var b strings.Builder
	fmt.Fprintf(&b, "%d %d\n", n, m)
	for _, e := range edges {
		fmt.Fprintf(&b, "%d %d\n", e.u, e.v)
	}

	return testCase{
		input:  b.String(),
		n:      n,
		m:      m,
		maxDeg: maxDeg,
		edges:  edges,
	}
}
