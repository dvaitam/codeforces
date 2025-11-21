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
	refSource        = "958B1.go"
	tempOraclePrefix = "oracle-958B1-"
	randomTestsCount = 120
	maxRandomN       = 1000
)

type edge struct {
	u int
	v int
}

type testCase struct {
	name  string
	n     int
	edges []edge
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierB1.go /path/to/binary")
		os.Exit(1)
	}
	candidate := os.Args[1]

	oraclePath, cleanup, err := buildReference()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to build reference: %v\n", err)
		os.Exit(1)
	}
	defer cleanup()

	tests := deterministicTests()
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	tests = append(tests, randomTests(randomTestsCount, rng)...)
	tests = append(tests, largeTests()...)

	for idx, tc := range tests {
		input := formatInput(tc)

		expOut, err := runProgram(oraclePath, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "oracle runtime error on test %d (%s): %v\n", idx+1, tc.name, err)
			fmt.Println("Input:")
			fmt.Print(input)
			os.Exit(1)
		}
		expected, err := parseAnswer(expOut)
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to parse oracle output on test %d (%s): %v\noutput:\n%s", idx+1, tc.name, err, expOut)
			os.Exit(1)
		}

		gotOut, err := runProgram(candidate, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "candidate runtime error on test %d (%s): %v\n", idx+1, tc.name, err)
			fmt.Println("Input:")
			fmt.Print(input)
			os.Exit(1)
		}
		got, err := parseAnswer(gotOut)
		if err != nil {
			fmt.Fprintf(os.Stderr, "candidate output parse error on test %d (%s): %v\noutput:\n%s", idx+1, tc.name, err, gotOut)
			fmt.Println("Input:")
			fmt.Print(input)
			os.Exit(1)
		}

		if expected != got {
			fmt.Fprintf(os.Stderr, "test %d (%s) failed: expected %d got %d\n", idx+1, tc.name, expected, got)
			fmt.Println("Input:")
			fmt.Print(input)
			fmt.Println("Candidate output:")
			fmt.Print(gotOut)
			os.Exit(1)
		}
	}

	fmt.Printf("All %d tests passed\n", len(tests))
}

func buildReference() (string, func(), error) {
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		return "", nil, fmt.Errorf("failed to determine verifier path")
	}
	dir := filepath.Dir(file)
	tmpDir, err := os.MkdirTemp("", tempOraclePrefix)
	if err != nil {
		return "", nil, err
	}
	outPath := filepath.Join(tmpDir, "oracleB1")
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

func runProgram(bin, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
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

func parseAnswer(out string) (int, error) {
	out = strings.TrimSpace(out)
	if out == "" {
		return 0, fmt.Errorf("empty output")
	}
	val, err := strconv.Atoi(out)
	if err != nil {
		return 0, fmt.Errorf("invalid integer %q", out)
	}
	return val, nil
}

func formatInput(tc testCase) string {
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d\n", tc.n)
	for _, e := range tc.edges {
		fmt.Fprintf(&sb, "%d %d\n", e.u, e.v)
	}
	return sb.String()
}

func deterministicTests() []testCase {
	return []testCase{
		{name: "two_nodes", n: 2, edges: []edge{{1, 2}}},
		{name: "line_5", n: 5, edges: []edge{{1, 2}, {2, 3}, {3, 4}, {4, 5}}},
		{name: "star_5", n: 5, edges: []edge{{1, 2}, {1, 3}, {1, 4}, {1, 5}}},
		{name: "balanced_tree", n: 7, edges: []edge{{1, 2}, {1, 3}, {2, 4}, {2, 5}, {3, 6}, {3, 7}}},
	}
}

func randomTests(count int, rng *rand.Rand) []testCase {
	tests := make([]testCase, 0, count)
	for i := 0; i < count; i++ {
		n := rng.Intn(maxRandomN-2) + 2
		edges := randomTree(n, rng)
		tests = append(tests, testCase{
			name:  fmt.Sprintf("random_%d", i+1),
			n:     n,
			edges: edges,
		})
	}
	return tests
}

func largeTests() []testCase {
	n := maxRandomN
	lineEdges := make([]edge, 0, n-1)
	for i := 2; i <= n; i++ {
		lineEdges = append(lineEdges, edge{i - 1, i})
	}
	starEdges := make([]edge, 0, n-1)
	for i := 2; i <= n; i++ {
		starEdges = append(starEdges, edge{1, i})
	}
	return []testCase{
		{name: "large_line", n: n, edges: lineEdges},
		{name: "large_star", n: n, edges: starEdges},
	}
}

func randomTree(n int, rng *rand.Rand) []edge {
	edges := make([]edge, 0, n-1)
	for v := 2; v <= n; v++ {
		u := rng.Intn(v-1) + 1
		edges = append(edges, edge{u: u, v: v})
	}
	rng.Shuffle(len(edges), func(i, j int) {
		edges[i], edges[j] = edges[j], edges[i]
	})
	return edges
}
