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

const (
	refSource     = "2120F.go"
	totalNLimit   = 250
	maxGraphs     = 10
	randomSeedDet = 2120
)

type graph struct {
	edges [][2]int
}

type testCase struct {
	n      int
	k      int
	graphs []graph
}

func main() {
	args := os.Args[1:]
	if len(args) == 2 && args[0] == "--" {
		args = args[1:]
	}
	if len(args) != 1 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierF.go /path/to/candidate")
		os.Exit(1)
	}
	candidate := args[0]

	refBin, err := buildReference()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to build reference: %v\n", err)
		os.Exit(1)
	}
	defer os.Remove(refBin)

	tests := buildTests()
	input := buildInput(tests)

	refOut, err := runProgram(refBin, input)
	if err != nil {
		fmt.Fprintf(os.Stderr, "reference runtime error: %v\n%s\n", err, refOut)
		os.Exit(1)
	}
	expect, err := parseOutputs(refOut, len(tests))
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to parse reference output: %v\n%s\n", err, refOut)
		os.Exit(1)
	}

	candOut, err := runProgram(candidate, input)
	if err != nil {
		fmt.Fprintf(os.Stderr, "candidate runtime error: %v\n%s\n", err, candOut)
		os.Exit(1)
	}
	got, err := parseOutputs(candOut, len(tests))
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to parse candidate output: %v\n%s\n", err, candOut)
		os.Exit(1)
	}

	for i := range tests {
		if expect[i] != got[i] {
			tc := tests[i]
			fmt.Fprintf(os.Stderr, "test %d mismatch: expected %s got %s\n", i+1, expect[i], got[i])
			fmt.Fprintf(os.Stderr, "n=%d k=%d\n", tc.n, tc.k)
			for gIdx, g := range tc.graphs {
				fmt.Fprintf(os.Stderr, " graph %d edges=%d %s\n", gIdx+1, len(g.edges), summarizeEdges(g.edges, tc.n))
			}
			os.Exit(1)
		}
	}

	fmt.Printf("All %d tests passed.\n", len(tests))
}

func buildReference() (string, error) {
	dir, err := verifierDir()
	if err != nil {
		return "", err
	}
	tmp, err := os.CreateTemp("", "2120F-ref-*")
	if err != nil {
		return "", err
	}
	tmp.Close()

	cmd := exec.Command("go", "build", "-o", tmp.Name(), filepath.Join(dir, refSource))
	if out, err := cmd.CombinedOutput(); err != nil {
		os.Remove(tmp.Name())
		return "", fmt.Errorf("%v\n%s", err, string(out))
	}
	return tmp.Name(), nil
}

func verifierDir() (string, error) {
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		return "", fmt.Errorf("cannot determine verifier directory")
	}
	return filepath.Dir(file), nil
}

func runProgram(target, input string) (string, error) {
	cmd := commandFor(target)
	cmd.Stdin = strings.NewReader(input)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()
	if err != nil {
		return stdout.String() + stderr.String(), err
	}
	return stdout.String(), nil
}

func commandFor(path string) *exec.Cmd {
	switch filepath.Ext(path) {
	case ".go":
		return exec.Command("go", "run", path)
	case ".py":
		return exec.Command("python3", path)
	default:
		return exec.Command(path)
	}
}

func parseOutputs(out string, expected int) ([]string, error) {
	tokens := strings.Fields(out)
	if len(tokens) < expected {
		return nil, fmt.Errorf("expected %d tokens, got %d", expected, len(tokens))
	}
	if len(tokens) > expected {
		return nil, fmt.Errorf("extra output starting at token %q", tokens[expected])
	}
	res := make([]string, expected)
	for i := 0; i < expected; i++ {
		tok := strings.ToUpper(tokens[i])
		if tok != "YES" && tok != "NO" {
			return nil, fmt.Errorf("token %q is not YES/NO", tokens[i])
		}
		res[i] = tok
	}
	return res, nil
}

func buildInput(tests []testCase) string {
	var sb strings.Builder
	sb.WriteString(strconv.Itoa(len(tests)))
	sb.WriteByte('\n')
	for _, tc := range tests {
		sb.WriteString(fmt.Sprintf("%d %d\n", tc.n, tc.k))
		for _, g := range tc.graphs {
			sb.WriteString(strconv.Itoa(len(g.edges)))
			sb.WriteByte('\n')
			for _, e := range g.edges {
				sb.WriteString(fmt.Sprintf("%d %d\n", e[0]+1, e[1]+1))
			}
		}
	}
	return sb.String()
}

func buildTests() []testCase {
	tests := make([]testCase, 0)
	totalN := 0
	add := func(tc testCase) {
		if totalN+tc.n > totalNLimit {
			return
		}
		tests = append(tests, tc)
		totalN += tc.n
	}

	for _, tc := range deterministicTests() {
		add(tc)
	}

	// Deterministic random biggish case to vary edges.
	if totalN < totalNLimit-60 {
		add(randomCase(rand.New(rand.NewSource(randomSeedDet)), 60, 4))
	}

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for totalN < totalNLimit {
		n := 5 + rng.Intn(40)
		if n+totalN > totalNLimit {
			n = totalNLimit - totalN
		}
		k := 1 + rng.Intn(maxGraphs)
		add(randomCase(rng, n, k))
	}

	return tests
}

func deterministicTests() []testCase {
	// twin vertices -> No
	noTwinFree := testCase{
		n: 2,
		k: 1,
		graphs: []graph{
			{edges: [][2]int{}},
		},
	}
	// triangle single graph -> Yes
	allConnected := testCase{
		n: 3,
		k: 1,
		graphs: []graph{
			{edges: [][2]int{{0, 1}, {1, 2}, {0, 2}}},
		},
	}
	// mix of twin-free and non twin-free graphs, result No.
	mixed := testCase{
		n: 4,
		k: 2,
		graphs: []graph{
			{edges: [][2]int{{0, 1}, {1, 2}}},
			{edges: [][2]int{}}, // twins exist
		},
	}
	return []testCase{noTwinFree, allConnected, mixed}
}

func randomCase(rng *rand.Rand, n int, k int) testCase {
	graphs := make([]graph, k)
	for i := 0; i < k; i++ {
		edges := randomEdges(rng, n, rng.Float64()*0.4)
		// Occasionally force twins by duplicating adjacency of two vertices.
		if rng.Intn(5) == 0 && n >= 3 {
			u := rng.Intn(n)
			v := (u + 1 + rng.Intn(n-1)) % n
			edges = makeTwins(n, edges, u, v, rng)
		}
		graphs[i] = graph{edges: edges}
	}
	return testCase{n: n, k: k, graphs: graphs}
}

func randomEdges(rng *rand.Rand, n int, p float64) [][2]int {
	edges := make([][2]int, 0)
	for u := 0; u < n; u++ {
		for v := u + 1; v < n; v++ {
			if rng.Float64() < p {
				edges = append(edges, [2]int{u, v})
			}
		}
	}
	return edges
}

func makeTwins(n int, edges [][2]int, u, v int, rng *rand.Rand) [][2]int {
	// Reset adjacency of u and v to match a new random mask.
	adj := make([][]bool, n)
	for i := range adj {
		adj[i] = make([]bool, n)
	}
	for _, e := range edges {
		a, b := e[0], e[1]
		adj[a][b] = true
		adj[b][a] = true
	}
	mask := make([]bool, n)
	for i := 0; i < n; i++ {
		if i == u || i == v {
			continue
		}
		mask[i] = rng.Intn(2) == 0
		adj[u][i] = mask[i]
		adj[i][u] = mask[i]
		adj[v][i] = mask[i]
		adj[i][v] = mask[i]
	}
	adj[u][v] = rng.Intn(2) == 0
	adj[v][u] = adj[u][v]

	newEdges := make([][2]int, 0)
	for i := 0; i < n; i++ {
		for j := i + 1; j < n; j++ {
			if adj[i][j] {
				newEdges = append(newEdges, [2]int{i, j})
			}
		}
	}
	return newEdges
}

func summarizeEdges(edges [][2]int, n int) string {
	const limit = 8
	if len(edges) == 0 {
		return "[]"
	}
	if len(edges) <= limit {
		return fmt.Sprint(edges)
	}
	head := fmt.Sprint(edges[:limit])
	return head[:len(head)-1] + fmt.Sprintf(" ... total=%d]", len(edges))
}
