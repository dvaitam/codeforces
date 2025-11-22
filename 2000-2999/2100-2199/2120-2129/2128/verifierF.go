package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"time"
)

const refSource = "2128F.go"

type edge struct {
	u, v int
	l, r int
}

type testCase struct {
	name     string
	n, m, k  int
	edges    []edge
	expected string // optional expected? unused, rely on reference
}

func main() {
	if len(os.Args) != 2 && !(len(os.Args) == 3 && os.Args[1] == "--") {
		fmt.Fprintln(os.Stderr, "usage: go run verifierF.go /path/to/candidate")
		os.Exit(1)
	}
	candPath := os.Args[len(os.Args)-1]

	refBin, cleanupRef, err := buildReference()
	if err != nil {
		fmt.Fprintln(os.Stderr, "failed to build reference solution:", err)
		os.Exit(1)
	}
	defer cleanupRef()

	candBin, cleanupCand, err := prepareCandidate(candPath)
	if err != nil {
		fmt.Fprintln(os.Stderr, "failed to prepare candidate:", err)
		os.Exit(1)
	}
	defer cleanupCand()

	tests := buildTests()
	input := buildInput(tests)

	refOut, err := runProgram(refBin, input)
	if err != nil {
		fmt.Fprintf(os.Stderr, "reference runtime error: %v\n", err)
		os.Exit(1)
	}
	refAns, err := parseOutputs(refOut, len(tests))
	if err != nil {
		fmt.Fprintf(os.Stderr, "reference produced invalid output: %v\n%s", err, refOut)
		os.Exit(1)
	}

	candOut, err := runProgram(candBin, input)
	if err != nil {
		fmt.Fprintf(os.Stderr, "candidate runtime error: %v\ninput preview:\n%s\n", err, previewInput(input))
		os.Exit(1)
	}
	candAns, err := parseOutputs(candOut, len(tests))
	if err != nil {
		fmt.Fprintf(os.Stderr, "candidate produced invalid output: %v\noutput:\n%s\ninput preview:\n%s\n", err, candOut, previewInput(input))
		os.Exit(1)
	}

	for i := range tests {
		if refAns[i] != candAns[i] {
			fmt.Fprintf(os.Stderr, "wrong answer on test %d (%s): expected %s got %s\n", i+1, tests[i].name, refAns[i], candAns[i])
			fmt.Fprintln(os.Stderr, previewInput(buildInput([]testCase{tests[i]})))
			os.Exit(1)
		}
	}

	fmt.Printf("All %d tests passed.\n", len(tests))
}

func buildReference() (string, func(), error) {
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		return "", nil, fmt.Errorf("cannot determine verifier directory")
	}
	dir := filepath.Dir(file)

	tmpDir, err := os.MkdirTemp("", "ref2128F-")
	if err != nil {
		return "", nil, err
	}
	binPath := filepath.Join(tmpDir, "ref")

	cmd := exec.Command("go", "build", "-o", binPath, refSource)
	cmd.Dir = dir
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		os.RemoveAll(tmpDir)
		return "", nil, fmt.Errorf("reference build failed: %v\n%s", err, stderr.String())
	}

	cleanup := func() { _ = os.RemoveAll(tmpDir) }
	return binPath, cleanup, nil
}

func prepareCandidate(path string) (string, func(), error) {
	abs, err := filepath.Abs(path)
	if err != nil {
		return "", nil, err
	}
	if filepath.Ext(abs) != ".go" {
		return abs, func() {}, nil
	}

	tmpDir, err := os.MkdirTemp("", "cand2128F-")
	if err != nil {
		return "", nil, err
	}
	binPath := filepath.Join(tmpDir, "candidate")

	cmd := exec.Command("go", "build", "-o", binPath, abs)
	cmd.Dir = filepath.Dir(abs)
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		os.RemoveAll(tmpDir)
		return "", nil, fmt.Errorf("candidate build failed: %v\n%s", err, stderr.String())
	}

	cleanup := func() { _ = os.RemoveAll(tmpDir) }
	return binPath, cleanup, nil
}

func runProgram(bin, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return out.String(), err
}

func parseOutputs(out string, expected int) ([]string, error) {
	tokens := strings.Fields(out)
	if len(tokens) != expected {
		return nil, fmt.Errorf("expected %d tokens, got %d", expected, len(tokens))
	}
	res := make([]string, expected)
	for i, tok := range tokens {
		up := strings.ToUpper(tok)
		if up == "YES" || up == "Y" {
			res[i] = "YES"
		} else if up == "NO" || up == "N" {
			res[i] = "NO"
		} else {
			return nil, fmt.Errorf("invalid answer %q", tok)
		}
	}
	return res, nil
}

func buildInput(tests []testCase) string {
	var sb strings.Builder
	sb.WriteString(strconv.Itoa(len(tests)))
	sb.WriteByte('\n')
	for _, tc := range tests {
		sb.WriteString(fmt.Sprintf("%d %d %d\n", tc.n, tc.m, tc.k))
		for _, e := range tc.edges {
			sb.WriteString(fmt.Sprintf("%d %d %d %d\n", e.u, e.v, e.l, e.r))
		}
	}
	return sb.String()
}

func buildTests() []testCase {
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	tests := []testCase{}
	tests = append(tests, tinyTriangles()...)
	tests = append(tests, pathVsVia(), parallelEdges(), tightBounds())
	tests = append(tests, randomTests(rng, 8, 6, 20, 0.3, "rand_small")...)
	tests = append(tests, randomTests(rng, 6, 50, 200, 0.6, "rand_mid")...)
	tests = append(tests, randomTests(rng, 4, 800, 1500, 0.9, "rand_big")...)
	tests = append(tests, stressTests(rng)...)
	return tests
}

func tinyTriangles() []testCase {
	// Minimal graphs where k is in the middle of alternative paths.
	return []testCase{
		{
			name: "simple_yes",
			n:    4, m: 4, k: 2,
			edges: []edge{
				{1, 2, 1, 5},
				{2, 3, 1, 5},
				{3, 4, 1, 5},
				{1, 4, 10, 10},
			},
		},
		{
			name: "simple_no",
			n:    4, m: 4, k: 2,
			edges: []edge{
				{1, 2, 5, 5},
				{2, 3, 5, 5},
				{3, 4, 5, 5},
				{1, 4, 15, 15},
			},
		},
		{
			name: "square_alt",
			n:    4, m: 5, k: 3,
			edges: []edge{
				{1, 2, 1, 3},
				{2, 4, 1, 3},
				{1, 3, 2, 2},
				{3, 4, 2, 2},
				{2, 3, 4, 10},
			},
		},
	}
}

func pathVsVia() testCase {
	// Long path plus shortcut through k.
	n := 7
	edges := []edge{
		{1, 2, 5, 5},
		{2, 3, 5, 5},
		{3, 4, 5, 5},
		{4, 5, 5, 5},
		{5, 6, 5, 5},
		{6, 7, 5, 5},
		{2, 6, 1, 20},
	}
	return testCase{name: "path_vs_via", n: n, m: len(edges), k: 4, edges: edges}
}

func parallelEdges() testCase {
	edges := []edge{
		{1, 2, 1, 1},
		{1, 3, 100, 100},
		{3, 2, 1, 2}, // ensure different path with strict bounds
		{2, 4, 1, 1},
	}
	return testCase{name: "detour_costs", n: 4, m: len(edges), k: 2, edges: edges}
}

func tightBounds() testCase {
	n := 6
	edges := []edge{
		{1, 2, 10, 10},
		{2, 3, 10, 10},
		{3, 6, 10, 10},
		{1, 4, 5, 50},
		{4, 5, 5, 50},
		{5, 6, 5, 50},
		{2, 5, 1, 1},
	}
	return testCase{name: "tight_bounds", n: n, m: len(edges), k: 3, edges: edges}
}

func randomTests(rng *rand.Rand, count, minN, maxN int, extraDensity float64, tag string) []testCase {
	res := make([]testCase, 0, count)
	for i := 0; i < count; i++ {
		n := rng.Intn(maxN-minN+1) + minN
		k := rng.Intn(n-2) + 2 // 2..n-1
		treeEdges := randomTree(rng, n)
		edgeSet := make(map[[2]int]bool)
		edges := make([]edge, 0, n-1)
		for _, p := range treeEdges {
			edgeSet[sortedPair(p[0], p[1])] = true
			edges = append(edges, randomEdgeWeight(rng, p[0], p[1]))
		}
		// add extra edges for density
		maxExtra := int(float64(n) * extraDensity)
		for j := 0; j < maxExtra; j++ {
			u := rng.Intn(n) + 1
			v := rng.Intn(n) + 1
			if u == v {
				continue
			}
			key := sortedPair(u, v)
			if edgeSet[key] {
				continue
			}
			edgeSet[key] = true
			edges = append(edges, randomEdgeWeight(rng, u, v))
		}
		res = append(res, testCase{
			name:  fmt.Sprintf("%s_%d_n%d", tag, i+1, n),
			n:     n,
			m:     len(edges),
			k:     k,
			edges: edges,
		})
	}
	return res
}

func stressTests(rng *rand.Rand) []testCase {
	res := []testCase{}
	// Large sparse
	res = append(res, largeCase(rng, 70000, 80000, "large_sparse"))
	// Large denser
	res = append(res, largeCase(rng, 90000, 110000, "large_dense"))
	return res
}

func largeCase(rng *rand.Rand, n, m int, name string) testCase {
	tree := randomTree(rng, n)
	edges := make([]edge, 0, m)
	edgeSet := make(map[[2]int]bool)
	for _, p := range tree {
		edgeSet[sortedPair(p[0], p[1])] = true
		edges = append(edges, randomEdgeWeight(rng, p[0], p[1]))
	}
	// add random extra edges until m reached
	for len(edges) < m {
		u := rng.Intn(n) + 1
		v := rng.Intn(n) + 1
		if u == v {
			continue
		}
		key := sortedPair(u, v)
		if edgeSet[key] {
			continue
		}
		edgeSet[key] = true
		edges = append(edges, randomEdgeWeight(rng, u, v))
	}
	k := rng.Intn(n-2) + 2
	return testCase{name: name, n: n, m: len(edges), k: k, edges: edges}
}

func randomTree(rng *rand.Rand, n int) [][2]int {
	parents := make([][2]int, 0, n-1)
	for v := 2; v <= n; v++ {
		p := rng.Intn(v-1) + 1
		parents = append(parents, [2]int{p, v})
	}
	return parents
}

func randomEdgeWeight(rng *rand.Rand, u, v int) edge {
	l := rng.Intn(1_000_000_000) + 1
	r := l + rng.Intn(1_000_000_000-l+1)
	return edge{u: u, v: v, l: l, r: r}
}

func sortedPair(a, b int) [2]int {
	if a > b {
		a, b = b, a
	}
	return [2]int{a, b}
}

func previewInput(input string) string {
	lines := strings.Split(input, "\n")
	if len(lines) > 12 {
		return strings.Join(lines[:12], "\n") + "\n..."
	}
	return input
}
