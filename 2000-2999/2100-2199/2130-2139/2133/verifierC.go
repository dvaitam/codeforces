package main

import (
	"bufio"
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

const (
	refSource2133C = "2000-2999/2100-2199/2130-2139/2133/2133C.go"
)

type testCase struct {
	n   int
	adj [][]int
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	candidate := os.Args[1]

	refBin, err := buildReference(refSource2133C)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to build reference: %v\n", err)
		os.Exit(1)
	}
	defer os.Remove(refBin)

	tests := buildTests()
	input := serializeTests(tests)

	// Sanity-check reference solver.
	if err := runAndValidate(refBin, tests, input); err != nil {
		fmt.Fprintf(os.Stderr, "reference failed: %v\n", err)
		os.Exit(1)
	}

	if err := runAndValidate(candidate, tests, input); err != nil {
		fmt.Fprintf(os.Stderr, "candidate failed: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("All %d tests passed\n", len(tests))
}

func buildReference(source string) (string, error) {
	tmp, err := os.CreateTemp("", "2133C-ref-*")
	if err != nil {
		return "", err
	}
	tmp.Close()

	srcPath, err := resolveSourcePath(source)
	if err != nil {
		os.Remove(tmp.Name())
		return "", err
	}

	cmd := exec.Command("go", "build", "-o", tmp.Name(), srcPath)
	var buf bytes.Buffer
	cmd.Stdout = &buf
	cmd.Stderr = &buf
	if err := cmd.Run(); err != nil {
		os.Remove(tmp.Name())
		return "", fmt.Errorf("%v\n%s", err, buf.String())
	}
	return tmp.Name(), nil
}

func resolveSourcePath(path string) (string, error) {
	if filepath.IsAbs(path) {
		return path, nil
	}
	cwd, err := os.Getwd()
	if err != nil {
		return "", err
	}
	return filepath.Join(cwd, path), nil
}

func runAndValidate(target string, tests []testCase, input string) error {
	outStr, err := runProgram(target, input)
	if err != nil {
		return err
	}
	paths, err := parseOutput(outStr, tests)
	if err != nil {
		return err
	}

	for i, tc := range tests {
		if err := validatePath(tc, paths[i]); err != nil {
			return fmt.Errorf("test %d: %v", i+1, err)
		}
	}
	return nil
}

func runProgram(target, input string) (string, error) {
	cmd := commandFor(target)
	cmd.Stdin = strings.NewReader(input)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\nstdout:\n%s\nstderr:\n%s", err, stdout.String(), stderr.String())
	}
	return stdout.String(), nil
}

func commandFor(path string) *exec.Cmd {
	if strings.HasSuffix(path, ".go") {
		return exec.Command("go", "run", path)
	}
	return exec.Command(path)
}

func parseOutput(out string, tests []testCase) ([][]int, error) {
	reader := bufio.NewReader(strings.NewReader(out))
	ans := make([][]int, len(tests))
	for idx := range tests {
		var k int
		if _, err := fmt.Fscan(reader, &k); err != nil {
			return nil, fmt.Errorf("test %d: missing path length (%v)", idx+1, err)
		}
		if k < 1 || k > tests[idx].n {
			return nil, fmt.Errorf("test %d: invalid path length %d", idx+1, k)
		}
		path := make([]int, k)
		for i := 0; i < k; i++ {
			if _, err := fmt.Fscan(reader, &path[i]); err != nil {
				return nil, fmt.Errorf("test %d: expected %d nodes, got %d (%v)", idx+1, k, i, err)
			}
		}
		ans[idx] = path
	}
	var extra string
	if _, err := fmt.Fscan(reader, &extra); err == nil {
		return nil, fmt.Errorf("extra output detected starting with %q", extra)
	}
	return ans, nil
}

func validatePath(tc testCase, path []int) error {
	n := tc.n
	for i, v := range path {
		if v < 1 || v > n {
			return fmt.Errorf("node %d out of range: %d", i+1, v)
		}
	}
	// Build adjacency set for quick lookup.
	adjSet := make([]map[int]struct{}, n)
	for i := 0; i < n; i++ {
		if len(tc.adj[i]) > 0 {
			adjSet[i] = make(map[int]struct{}, len(tc.adj[i]))
			for _, v := range tc.adj[i] {
				adjSet[i][v] = struct{}{}
			}
		}
	}
	for i := 0; i+1 < len(path); i++ {
		u := path[i] - 1
		v := path[i+1] - 1
		if adjSet[u] == nil {
			return fmt.Errorf("edge %d->%d not present", path[i], path[i+1])
		}
		if _, ok := adjSet[u][v]; !ok {
			return fmt.Errorf("edge %d->%d not present", path[i], path[i+1])
		}
	}

	maxLen := longestPathLength(tc)
	if len(path) != maxLen {
		return fmt.Errorf("path length %d not maximum (expected %d)", len(path), maxLen)
	}
	return nil
}

func longestPathLength(tc testCase) int {
	order := topoOrder(tc)
	dp := make([]int, tc.n)
	for i := len(order) - 1; i >= 0; i-- {
		u := order[i]
		best := 1
		for _, v := range tc.adj[u] {
			if 1+dp[v] > best {
				best = 1 + dp[v]
			}
		}
		dp[u] = best
	}
	res := 1
	for _, v := range dp {
		if v > res {
			res = v
		}
	}
	return res
}

func topoOrder(tc testCase) []int {
	indeg := make([]int, tc.n)
	for u := 0; u < tc.n; u++ {
		for _, v := range tc.adj[u] {
			indeg[v]++
		}
	}
	q := make([]int, 0, tc.n)
	for i := 0; i < tc.n; i++ {
		if indeg[i] == 0 {
			q = append(q, i)
		}
	}
	for head := 0; head < len(q); head++ {
		u := q[head]
		for _, v := range tc.adj[u] {
			indeg[v]--
			if indeg[v] == 0 {
				q = append(q, v)
			}
		}
	}
	return q
}

func serializeTests(tests []testCase) string {
	var sb strings.Builder
	sb.WriteString(strconv.Itoa(len(tests)))
	sb.WriteByte('\n')
	for _, tc := range tests {
		sb.WriteString(strconv.Itoa(tc.n))
		sb.WriteByte('\n')
		for i := 0; i < tc.n; i++ {
			sb.WriteString(strconv.Itoa(len(tc.adj[i])))
			for _, v := range tc.adj[i] {
				sb.WriteByte(' ')
				sb.WriteString(strconv.Itoa(v + 1))
			}
			sb.WriteByte('\n')
		}
	}
	return sb.String()
}

func buildTests() []testCase {
	tests := make([]testCase, 0)
	totalCube := 0
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))

	add := func(tc testCase) {
		tests = append(tests, tc)
		totalCube += tc.n * tc.n * tc.n
	}

	// Deterministic small cases.
	add(lineCase(2))
	add(lineCase(5))
	add(starCase(6))
	add(forkCase())

	// Random mid-size DAGs.
	for totalCube < 50_000 {
		n := rng.Intn(25) + 5
		tc := randomDAG(n, 0.2, rng)
		add(tc)
	}

	// A larger stress case within budget.
	if totalCube < 100_000 {
		add(randomDAG(60, 0.12, rng))
	}

	return tests
}

func lineCase(n int) testCase {
	adj := make([][]int, n)
	for i := 0; i+1 < n; i++ {
		adj[i] = append(adj[i], i+1)
	}
	return testCase{n: n, adj: adj}
}

func starCase(n int) testCase {
	adj := make([][]int, n)
	for i := 1; i < n; i++ {
		adj[0] = append(adj[0], i)
	}
	return testCase{n: n, adj: adj}
}

func forkCase() testCase {
	// 1->2->3, 1->4->5, 2->5
	n := 5
	adj := make([][]int, n)
	adj[0] = []int{1, 3}
	adj[1] = []int{2, 4}
	adj[3] = []int{4}
	return testCase{n: n, adj: adj}
}

func randomDAG(n int, prob float64, rng *rand.Rand) testCase {
	order := rng.Perm(n)
	pos := make([]int, n)
	for i, v := range order {
		pos[v] = i
	}
	adj := make([][]int, n)
	for u := 0; u < n; u++ {
		for v := 0; v < n; v++ {
			if pos[u] < pos[v] && rng.Float64() < prob {
				adj[u] = append(adj[u], v)
			}
		}
	}
	// Ensure at least one edge for interest.
	if countEdges(adj) == 0 && n >= 2 {
		adj[order[0]] = append(adj[order[0]], order[1])
	}
	return testCase{n: n, adj: adj}
}

func countEdges(adj [][]int) int {
	total := 0
	for _, list := range adj {
		total += len(list)
	}
	return total
}
