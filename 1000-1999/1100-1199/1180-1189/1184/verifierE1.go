package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"sort"
	"strings"
	"time"
)

const (
	refSource        = "1184E1.go"
	tempOraclePrefix = "oracle-1184E1-"
	randomTestsCount = 40
	maxN             = 2000
	maxM             = 4000
	maxWeight        = 1_000_000_000
)

type edge struct {
	u, v, w int
}

type testCase struct {
	n     int
	m     int
	edges []edge
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierE1.go /path/to/binary")
		os.Exit(1)
	}
	candidate := os.Args[1]

	oracle, cleanup, err := buildOracle()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to build oracle: %v\n", err)
		os.Exit(1)
	}
	defer cleanup()

	tests := deterministicTests()
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	tests = append(tests, randomTests(rng, randomTestsCount)...)

	for idx, tc := range tests {
		input := formatInput(tc)
		exp, err := runProgram(oracle, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "oracle runtime error on test %d: %v\n", idx+1, err)
			fmt.Println("Input:")
			fmt.Print(input)
			os.Exit(1)
		}
		got, err := runProgram(candidate, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "candidate runtime error on test %d: %v\n", idx+1, err)
			fmt.Println("Input:")
			fmt.Print(input)
			os.Exit(1)
		}
		if strings.TrimSpace(exp) != strings.TrimSpace(got) {
			fmt.Fprintf(os.Stderr, "test %d failed: expected %s got %s\n", idx+1, strings.TrimSpace(exp), strings.TrimSpace(got))
			fmt.Println("Input:")
			fmt.Print(input)
			os.Exit(1)
		}
	}

	fmt.Printf("All %d tests passed\n", len(tests))
}

func buildOracle() (string, func(), error) {
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		return "", nil, fmt.Errorf("failed to determine verifier path")
	}
	dir := filepath.Dir(file)
	tmpDir, err := os.MkdirTemp("", tempOraclePrefix)
	if err != nil {
		return "", nil, err
	}
	outPath := filepath.Join(tmpDir, "oracleE1")
	cmd := exec.Command("go", "build", "-o", outPath, refSource)
	cmd.Dir = dir
	if out, err := cmd.CombinedOutput(); err != nil {
		os.RemoveAll(tmpDir)
		return "", nil, fmt.Errorf("build oracle failed: %v\n%s", err, out)
	}
	cleanup := func() {
		os.RemoveAll(tmpDir)
	}
	return outPath, cleanup, nil
}

func runProgram(path, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(path, ".go") {
		cmd = exec.Command("go", "run", path)
	} else {
		cmd = exec.Command(path)
	}
	cmd.Stdin = strings.NewReader(input)
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, stderr.String())
	}
	return strings.TrimSpace(stdout.String()), nil
}

func formatInput(tc testCase) string {
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d %d\n", tc.n, tc.m)
	for _, e := range tc.edges {
		fmt.Fprintf(&sb, "%d %d %d\n", e.u, e.v, e.w)
	}
	return sb.String()
}

func deterministicTests() []testCase {
	return []testCase{
		makeGraph(2, []edge{{1, 2, 0}}),
		makeGraph(3, []edge{{1, 2, 5}, {2, 3, 1}, {3, 1, 10}}),
		makeGraph(4, []edge{
			{1, 2, 3}, {2, 3, 4}, {3, 4, 5}, {4, 1, 6}, {1, 3, 2},
		}),
		makeGraph(5, []edge{
			{1, 2, 1}, {1, 3, 2}, {1, 4, 3}, {1, 5, 4},
			{2, 3, 5}, {3, 4, 6}, {4, 5, 7}, {2, 5, 8},
		}),
	}
}

func makeGraph(n int, edges []edge) testCase {
	return testCase{
		n:     n,
		m:     len(edges),
		edges: append([]edge(nil), edges...),
	}
}

func randomTests(rng *rand.Rand, count int) []testCase {
	tests := make([]testCase, 0, count)
	for i := 0; i < count; i++ {
		n := rng.Intn(maxN-2) + 2
		m := rng.Intn(maxM-n+1) + n
		edges := randomGraph(rng, n, m)
		tests = append(tests, testCase{n: n, m: len(edges), edges: edges})
	}
	return tests
}

func randomGraph(rng *rand.Rand, n, m int) []edge {
	type pair struct{ u, v int }
	used := make(map[pair]struct{})
	edges := make([]edge, 0, m)
	// Ensure connectivity via tree edges
	for v := 2; v <= n; v++ {
		u := rng.Intn(v-1) + 1
		e := edge{u: u, v: v, w: rng.Intn(maxWeight + 1)}
		edges = append(edges, e)
		used[pair{min(u, v), max(u, v)}] = struct{}{}
	}
	// Add extra edges
	for len(edges) < m {
		u := rng.Intn(n) + 1
		v := rng.Intn(n-1) + 1
		if v >= u {
			v++
		}
		a, b := min(u, v), max(u, v)
		p := pair{a, b}
		if _, ok := used[p]; ok {
			continue
		}
		used[p] = struct{}{}
		edges = append(edges, edge{u: a, v: b, w: rng.Intn(maxWeight + 1)})
	}
	sort.Slice(edges, func(i, j int) bool {
		if edges[i].u != edges[j].u {
			return edges[i].u < edges[j].u
		}
		if edges[i].v != edges[j].v {
			return edges[i].v < edges[j].v
		}
		return edges[i].w < edges[j].w
	})
	return edges
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
