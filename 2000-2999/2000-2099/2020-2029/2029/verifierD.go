package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
)

const refSource = "2000-2999/2000-2099/2020-2029/2029/2029D.go"

type testCase struct {
	n     int
	m     int
	edges [][2]int
}

type operation struct {
	a, b, c int
}

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierD.go /path/to/candidate")
		os.Exit(1)
	}
	candidate := os.Args[len(os.Args)-1]

	refBin, err := buildReference()
	if err != nil {
		fmt.Fprintln(os.Stderr, "failed to build reference:", err)
		os.Exit(1)
	}
	defer os.Remove(refBin)

	tests := generateTests()
	input := buildInput(tests)

	refOut, err := runProgram(refBin, input)
	if err != nil {
		fmt.Fprintf(os.Stderr, "reference failed: %v\noutput:\n%s\n", err, refOut)
		os.Exit(1)
	}
	candOut, err := runCandidate(candidate, input)
	if err != nil {
		fmt.Fprintf(os.Stderr, "candidate failed: %v\noutput:\n%s\n", err, candOut)
		os.Exit(1)
	}

	refOps, err := parseOutputs(refOut, len(tests))
	if err != nil {
		fmt.Fprintf(os.Stderr, "could not parse reference output: %v\n", err)
		os.Exit(1)
	}
	candOps, err := parseOutputs(candOut, len(tests))
	if err != nil {
		fmt.Fprintf(os.Stderr, "could not parse candidate output: %v\n", err)
		os.Exit(1)
	}

	for i, tc := range tests {
		if len(candOps[i]) > len(refOps[i]) {
			fmt.Fprintf(os.Stderr, "test %d: used %d operations, expected at most %d\n", i+1, len(candOps[i]), len(refOps[i]))
			os.Exit(1)
		}
		if err := validate(tc, candOps[i]); err != nil {
			fmt.Fprintf(os.Stderr, "test %d: %v\n", i+1, err)
			os.Exit(1)
		}
	}

	fmt.Printf("All %d tests passed.\n", len(tests))
}

func buildReference() (string, error) {
	tmp, err := os.CreateTemp("", "2029D-ref-*")
	if err != nil {
		return "", err
	}
	tmp.Close()
	cmd := exec.Command("go", "build", "-o", tmp.Name(), filepath.Clean(refSource))
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		os.Remove(tmp.Name())
		return "", fmt.Errorf("%v\n%s", err, out.String())
	}
	return tmp.Name(), nil
}

func runCandidate(path, input string) (string, error) {
	cmd := commandFor(path)
	return runWithInput(cmd, input)
}

func runProgram(path, input string) (string, error) {
	cmd := exec.Command(path)
	return runWithInput(cmd, input)
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

func runWithInput(cmd *exec.Cmd, input string) (string, error) {
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return out.String(), err
}

func parseOutputs(output string, t int) ([][]operation, error) {
	lines := strings.Split(strings.TrimSpace(output), "\n")
	idx := 0
	ans := make([][]operation, t)
	for i := 0; i < t; i++ {
		for idx < len(lines) && strings.TrimSpace(lines[idx]) == "" {
			idx++
		}
		if idx >= len(lines) {
			return nil, fmt.Errorf("missing output for test %d", i+1)
		}
		k, err := strconv.Atoi(strings.TrimSpace(lines[idx]))
		if err != nil {
			return nil, fmt.Errorf("test %d: invalid operation count", i+1)
		}
		if k < 0 {
			return nil, fmt.Errorf("test %d: negative operation count", i+1)
		}
		idx++
		ops := make([]operation, 0, k)
		for j := 0; j < k; j++ {
			for idx < len(lines) && strings.TrimSpace(lines[idx]) == "" {
				idx++
			}
			if idx >= len(lines) {
				return nil, fmt.Errorf("test %d: missing operation line %d", i+1, j+1)
			}
			fields := strings.Fields(lines[idx])
			if len(fields) != 3 {
				return nil, fmt.Errorf("test %d: operation %d must have 3 integers", i+1, j+1)
			}
			a, err1 := strconv.Atoi(fields[0])
			b, err2 := strconv.Atoi(fields[1])
			c, err3 := strconv.Atoi(fields[2])
			if err1 != nil || err2 != nil || err3 != nil {
				return nil, fmt.Errorf("test %d: invalid integers in operation %d", i+1, j+1)
			}
			ops = append(ops, operation{a, b, c})
			idx++
		}
		ans[i] = ops
	}
	return ans, nil
}

func validate(tc testCase, ops []operation) error {
	limit := 2 * max(tc.n, tc.m)
	if len(ops) > limit {
		return fmt.Errorf("used %d operations, limit is %d", len(ops), limit)
	}
	graph := make([][]bool, tc.n+1)
	for i := range graph {
		graph[i] = make([]bool, tc.n+1)
	}
	for _, e := range tc.edges {
		u, v := e[0], e[1]
		graph[u][v] = !graph[u][v]
		graph[v][u] = !graph[v][u]
	}
	for idx, op := range ops {
		a, b, c := op.a, op.b, op.c
		if !checkTriple(tc.n, a, b, c) {
			return fmt.Errorf("operation %d has invalid vertices %d %d %d", idx+1, a, b, c)
		}
		toggle := func(u, v int) {
			graph[u][v] = !graph[u][v]
			graph[v][u] = !graph[v][u]
		}
		toggle(a, b)
		toggle(b, c)
		toggle(a, c)
	}
	return checkCool(graph, tc.n)
}

func checkTriple(n, a, b, c int) bool {
	if a < 1 || a > n || b < 1 || b > n || c < 1 || c > n {
		return false
	}
	if a == b || b == c || a == c {
		return false
	}
	return true
}

func checkCool(graph [][]bool, n int) error {
	edgeCount := 0
	for i := 1; i <= n; i++ {
		for j := i + 1; j <= n; j++ {
			if graph[i][j] {
				edgeCount++
			}
		}
	}
	if edgeCount == 0 {
		return nil
	}
	seen := make([]bool, n+1)
	var queue []int
	queue = append(queue, 1)
	seen[1] = true
	for len(queue) > 0 {
		u := queue[0]
		queue = queue[1:]
		for v := 1; v <= n; v++ {
			if graph[u][v] && !seen[v] {
				seen[v] = true
				queue = append(queue, v)
			}
		}
	}
	for i := 1; i <= n; i++ {
		if !seen[i] {
			return fmt.Errorf("graph is disconnected")
		}
	}
	if edgeCount != n-1 {
		return fmt.Errorf("graph has %d edges, expected %d", edgeCount, n-1)
	}
	return nil
}

func generateTests() []testCase {
	rng := rand.New(rand.NewSource(2029))
	var tests []testCase
	totalN := 0
	totalM := 0

	add := func(tc testCase) {
		if totalN+tc.n > 100000 || totalM+tc.m > 200000 {
			return
		}
		tests = append(tests, tc)
		totalN += tc.n
		totalM += tc.m
	}

	add(makeTree(3))
	add(makeTree(5))
	add(completeGraph(4))
	add(completeGraph(6))
	add(randomGraph(rng, 8, 10))

	for totalN < 100000 && totalM < 200000 {
		n := rng.Intn(300) + 3
		m := rng.Intn(min(2*n, 200)) + n - 1
		add(randomGraph(rng, n, m))
		if len(tests) > 200 {
			break
		}
	}

	return tests
}

func makeTree(n int) testCase {
	edges := make([][2]int, n-1)
	for i := 2; i <= n; i++ {
		edges[i-2] = [2]int{i, i - 1}
	}
	return testCase{n: n, m: len(edges), edges: edges}
}

func completeGraph(n int) testCase {
	var edges [][2]int
	for i := 1; i <= n; i++ {
		for j := i + 1; j <= n; j++ {
			edges = append(edges, [2]int{i, j})
		}
	}
	return testCase{n: n, m: len(edges), edges: edges}
}

func randomGraph(rng *rand.Rand, n, m int) testCase {
	maxEdges := n * (n - 1) / 2
	if m > maxEdges {
		m = maxEdges
	}
	edges := make(map[[2]int]struct{})
	for len(edges) < m {
		u := rng.Intn(n) + 1
		v := rng.Intn(n) + 1
		if u == v {
			continue
		}
		if u > v {
			u, v = v, u
		}
		key := [2]int{u, v}
		if _, exists := edges[key]; !exists {
			edges[key] = struct{}{}
		}
	}
	list := make([][2]int, 0, len(edges))
	for e := range edges {
		list = append(list, e)
	}
	return testCase{n: n, m: len(list), edges: list}
}

func buildInput(tests []testCase) string {
	var b strings.Builder
	fmt.Fprintf(&b, "%d\n", len(tests))
	for _, tc := range tests {
		fmt.Fprintf(&b, "%d %d\n", tc.n, tc.m)
		for _, e := range tc.edges {
			fmt.Fprintf(&b, "%d %d\n", e[0], e[1])
		}
	}
	return b.String()
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
