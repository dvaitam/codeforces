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

const refSource = "2060E.go"

type testCase struct {
	n    int
	m1   int
	m2   int
	fEds [][2]int
	gEds [][2]int
}

func main() {
	args := os.Args[1:]
	if len(args) == 2 && args[0] == "--" {
		args = args[1:]
	}
	if len(args) != 1 {
		fmt.Println("usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	candidate := args[0]

	refBin, err := buildReference()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer os.Remove(refBin)

	tests := buildTests()
	input := buildInput(tests)

	refOut, err := runProgram(refBin, input)
	if err != nil {
		fmt.Fprintf(os.Stderr, "reference failed: %v\n", err)
		os.Exit(1)
	}
	refAns, err := parseOutput(refOut, len(tests))
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to parse reference output: %v\n%s", err, refOut)
		os.Exit(1)
	}

	candOut, err := runProgram(candidate, input)
	if err != nil {
		fmt.Fprintf(os.Stderr, "candidate runtime error: %v\n", err)
		os.Exit(1)
	}
	candAns, err := parseOutput(candOut, len(tests))
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to parse candidate output: %v\n%s", err, candOut)
		os.Exit(1)
	}

	for i := range tests {
		if refAns[i] != candAns[i] {
			tc := tests[i]
			fmt.Fprintf(os.Stderr, "test %d mismatch: expected %d got %d\n", i+1, refAns[i], candAns[i])
			fmt.Fprintf(os.Stderr, "n=%d m1=%d m2=%d\n", tc.n, tc.m1, tc.m2)
			fmt.Fprintf(os.Stderr, "F edges: %s\n", formatEdges(tc.fEds))
			fmt.Fprintf(os.Stderr, "G edges: %s\n", formatEdges(tc.gEds))
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}

func buildReference() (string, error) {
	refPath, err := referencePath()
	if err != nil {
		return "", err
	}
	tmp, err := os.CreateTemp("", "ref_2060E_*.bin")
	if err != nil {
		return "", fmt.Errorf("failed to create temp file: %v", err)
	}
	tmp.Close()

	cmd := exec.Command("go", "build", "-o", tmp.Name(), refPath)
	if out, err := cmd.CombinedOutput(); err != nil {
		os.Remove(tmp.Name())
		return "", fmt.Errorf("failed to build reference: %v\n%s", err, string(out))
	}
	return tmp.Name(), nil
}

func referencePath() (string, error) {
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		return "", fmt.Errorf("failed to locate verifier path")
	}
	dir := filepath.Dir(file)
	return filepath.Join(dir, refSource), nil
}

func runProgram(target, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(target, ".go") {
		cmd = exec.Command("go", "run", target)
	} else {
		cmd = exec.Command(target)
	}
	cmd.Stdin = strings.NewReader(input)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\nstdout:\n%s\nstderr:\n%s", err, stdout.String(), stderr.String())
	}
	return strings.TrimSpace(stdout.String()), nil
}

func buildTests() []testCase {
	var tests []testCase
	add := func(tc testCase) {
		tests = append(tests, tc)
	}

	// Basic edge cases
	add(testCase{n: 1, m1: 0, m2: 0})
	add(makeLineGraphs(3, true))
	add(makeLineGraphs(3, false))
	add(testCase{n: 4, m1: 0, m2: 3, gEds: [][2]int{{1, 2}, {2, 3}, {3, 4}}})
	add(testCase{n: 4, m1: 6, m2: 0, fEds: completeEdges(4)})
	add(testCase{n: 5, m1: 2, m2: 2, fEds: [][2]int{{1, 2}, {3, 4}}, gEds: [][2]int{{2, 3}, {4, 5}}})

	// Larger structured cases
	add(testCase{n: 6, m1: 3, m2: 3, fEds: [][2]int{{1, 2}, {3, 4}, {5, 6}}, gEds: [][2]int{{1, 3}, {2, 5}, {4, 6}}})
	cycle7 := cycleEdges(7)
	star7 := starEdges(7, 4)
	add(testCase{n: 7, m1: len(cycle7), m2: len(star7), fEds: cycle7, gEds: star7})
	path10 := pathEdges(10)
	cycle10 := cycleEdges(10)
	add(testCase{n: 10, m1: len(cycle10), m2: len(path10), fEds: cycle10, gEds: path10})

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))

	for len(tests) < 150 {
		n := rng.Intn(30) + 2 // up to 31 nodes small random
		m1 := rng.Intn(n*(n-1)/2 + 1)
		m2 := rng.Intn(n*(n-1)/2 + 1)
		fEds := randomGraphEdges(rng, n, m1)
		gEds := randomGraphEdges(rng, n, m2)
		add(testCase{n: n, m1: len(fEds), m2: len(gEds), fEds: fEds, gEds: gEds})
	}

	// Stress with many nodes but sparse
	add(testCase{n: 200000, m1: 0, m2: 0})
	chain := chainGraph(200000)
	add(testCase{n: 200000, m1: len(chain), m2: len(chain), fEds: chain, gEds: chain})

	return tests
}

func chainGraph(n int) [][2]int {
	eds := make([][2]int, 0, n-1)
	for i := 1; i < n; i++ {
		eds = append(eds, [2]int{i, i + 1})
	}
	return eds
}

func makeLineGraphs(n int, same bool) testCase {
	path := pathEdges(n)
	if same {
		return testCase{n: n, m1: len(path), m2: len(path), fEds: path, gEds: path}
	}
	return testCase{n: n, m1: 0, m2: len(path), fEds: nil, gEds: path}
}

func pathEdges(n int) [][2]int {
	edges := make([][2]int, 0, n-1)
	for i := 1; i < n; i++ {
		edges = append(edges, [2]int{i, i + 1})
	}
	return edges
}

func cycleEdges(n int) [][2]int {
	edges := pathEdges(n)
	edges = append(edges, [2]int{n, 1})
	return edges
}

func starEdges(n, center int) [][2]int {
	edges := make([][2]int, 0, n-1)
	for i := 1; i <= n; i++ {
		if i == center {
			continue
		}
		edges = append(edges, [2]int{center, i})
	}
	return edges
}

func completeEdges(n int) [][2]int {
	var edges [][2]int
	for i := 1; i <= n; i++ {
		for j := i + 1; j <= n; j++ {
			edges = append(edges, [2]int{i, j})
		}
	}
	return edges
}

func randomGraphEdges(rng *rand.Rand, n, m int) [][2]int {
	maxEdges := n * (n - 1) / 2
	if m > maxEdges {
		m = maxEdges
	}
	seen := make(map[int]struct{}, m)
	var edges [][2]int
	for len(edges) < m {
		u := rng.Intn(n) + 1
		v := rng.Intn(n) + 1
		if u == v {
			continue
		}
		if u > v {
			u, v = v, u
		}
		key := u*n + v
		if _, ok := seen[key]; ok {
			continue
		}
		seen[key] = struct{}{}
		edges = append(edges, [2]int{u, v})
	}
	return edges
}

func buildInput(tests []testCase) string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", len(tests)))
	for _, tc := range tests {
		sb.WriteString(fmt.Sprintf("%d %d %d\n", tc.n, tc.m1, tc.m2))
		for _, e := range tc.fEds {
			sb.WriteString(fmt.Sprintf("%d %d\n", e[0], e[1]))
		}
		for _, e := range tc.gEds {
			sb.WriteString(fmt.Sprintf("%d %d\n", e[0], e[1]))
		}
	}
	return sb.String()
}

func parseOutput(out string, expected int) ([]int, error) {
	fields := strings.Fields(out)
	if len(fields) != expected {
		return nil, fmt.Errorf("expected %d outputs, got %d", expected, len(fields))
	}
	ans := make([]int, expected)
	for i, s := range fields {
		val, err := strconv.Atoi(s)
		if err != nil {
			return nil, fmt.Errorf("invalid integer %q", s)
		}
		ans[i] = val
	}
	return ans, nil
}

func formatEdges(edges [][2]int) string {
	if len(edges) == 0 {
		return "(none)"
	}
	parts := make([]string, len(edges))
	for i, e := range edges {
		parts[i] = fmt.Sprintf("(%d,%d)", e[0], e[1])
	}
	return strings.Join(parts, " ")
}
