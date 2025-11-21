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

type edge struct {
	u, v int
	w    int64
}

type testCase struct {
	n, m, p int
	s       []int
	edges   []edge
}

func buildOracle() (string, func(), error) {
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		return "", nil, fmt.Errorf("cannot determine verifier path")
	}
	dir := filepath.Dir(file)
	tmpDir, err := os.MkdirTemp("", "oracle-2021E2-")
	if err != nil {
		return "", nil, err
	}
	outPath := filepath.Join(tmpDir, "oracleE2")
	cmd := exec.Command("go", "build", "-o", outPath, "2021E2.go")
	cmd.Dir = dir
	if out, err := cmd.CombinedOutput(); err != nil {
		os.RemoveAll(tmpDir)
		return "", nil, fmt.Errorf("failed to build oracle: %v\n%s", err, out)
	}
	cleanup := func() {
		os.RemoveAll(tmpDir)
	}
	return outPath, cleanup, nil
}

func runBinary(bin, input string) (string, error) {
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
	return stdout.String(), nil
}

func buildInput(tests []testCase) string {
	var sb strings.Builder
	sb.Grow(len(tests) * 512)
	sb.WriteString(strconv.Itoa(len(tests)))
	sb.WriteByte('\n')
	for _, tc := range tests {
		sb.WriteString(fmt.Sprintf("%d %d %d\n", tc.n, tc.m, tc.p))
		for i, val := range tc.s {
			if i > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(strconv.Itoa(val))
		}
		sb.WriteByte('\n')
		for _, e := range tc.edges {
			sb.WriteString(fmt.Sprintf("%d %d %d\n", e.u, e.v, e.w))
		}
	}
	return sb.String()
}

func parseOutput(out string, tests []testCase) ([][]int64, error) {
	fields := strings.Fields(out)
	idx := 0
	res := make([][]int64, len(tests))
	for i, tc := range tests {
		if idx+tc.n > len(fields) {
			return nil, fmt.Errorf("not enough answers for test %d", i+1)
		}
		ans := make([]int64, tc.n)
		for j := 0; j < tc.n; j++ {
			val, err := strconv.ParseInt(fields[idx+j], 10, 64)
			if err != nil {
				return nil, fmt.Errorf("invalid integer %q", fields[idx+j])
			}
			ans[j] = val
		}
		res[i] = ans
		idx += tc.n
	}
	if idx != len(fields) {
		return nil, fmt.Errorf("extra output values")
	}
	return res, nil
}

func deterministicTests() []testCase {
	return []testCase{
		{
			n: 4, m: 3, p: 2,
			s: []int{2, 3},
			edges: []edge{
				{1, 2, 5},
				{2, 3, 7},
				{3, 4, 1},
			},
		},
		{
			n: 5, m: 5, p: 3,
			s: []int{2, 4, 5},
			edges: []edge{
				{1, 2, 3},
				{2, 3, 4},
				{3, 4, 5},
				{4, 5, 6},
				{1, 5, 2},
			},
		},
	}
}

func randomTests() []testCase {
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	tests := make([]testCase, 0, 120)
	for len(tests) < cap(tests) {
		n := rng.Intn(40) + 2
		m := n - 1 + rng.Intn(n)
		if m > n*(n-1)/2 {
			m = n * (n - 1) / 2
		}
		edges := randomConnectedGraph(rng, n, m)
		p := rng.Intn(n) + 1
		used := make(map[int]bool)
		s := make([]int, p)
		for i := 0; i < p; i++ {
			var v int
			for {
				v = rng.Intn(n) + 1
				if !used[v] {
					used[v] = true
					break
				}
			}
			s[i] = v
		}
		tests = append(tests, testCase{n: n, m: len(edges), p: p, s: s, edges: edges})
	}
	return tests
}

func randomConnectedGraph(rng *rand.Rand, n, m int) []edge {
	edges := make([]edge, 0, m)
	for i := 2; i <= n; i++ {
		p := rng.Intn(i-1) + 1
		edges = append(edges, edge{u: i, v: p, w: randWeight(rng)})
	}
	existing := make(map[[2]int]bool)
	for _, e := range edges {
		existing[[2]int{min(e.u, e.v), max(e.u, e.v)}] = true
	}
	for len(edges) < m {
		u := rng.Intn(n) + 1
		v := rng.Intn(n) + 1
		if u == v {
			continue
		}
		key := [2]int{min(u, v), max(u, v)}
		if existing[key] {
			continue
		}
		existing[key] = true
		edges = append(edges, edge{u: key[0], v: key[1], w: randWeight(rng)})
	}
	return edges
}

func randWeight(rng *rand.Rand) int64 {
	return randRangeInt64(rng, 1, 1_000_000_000)
}

func randRangeInt64(rng *rand.Rand, lo, hi int64) int64 {
	return lo + rng.Int63n(hi-lo+1)
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func compareAnswers(expected, actual [][]int64) error {
	if len(expected) != len(actual) {
		return fmt.Errorf("test count mismatch")
	}
	for i := range expected {
		if len(expected[i]) != len(actual[i]) {
			return fmt.Errorf("test %d answer length mismatch", i+1)
		}
		for j := range expected[i] {
			if expected[i][j] != actual[i][j] {
				return fmt.Errorf("test %d k=%d mismatch: expected %d, got %d", i+1, j+1, expected[i][j], actual[i][j])
			}
		}
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierE2.go /path/to/binary")
		os.Exit(1)
	}
	target := os.Args[1]

	oracle, cleanup, err := buildOracle()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
	defer cleanup()

	tests := append(deterministicTests(), randomTests()...)
	input := buildInput(tests)

	expectedOut, err := runBinary(oracle, input)
	if err != nil {
		fmt.Fprintf(os.Stderr, "oracle failed: %v\ninput:\n%s", err, input)
		os.Exit(1)
	}
	actualOut, err := runBinary(target, input)
	if err != nil {
		fmt.Fprintf(os.Stderr, "target runtime error: %v\ninput:\n%s", err, input)
		os.Exit(1)
	}

	expectedAns, err := parseOutput(expectedOut, tests)
	if err != nil {
		fmt.Fprintf(os.Stderr, "oracle output invalid: %v\n%s", err, expectedOut)
		os.Exit(1)
	}
	actualAns, err := parseOutput(actualOut, tests)
	if err != nil {
		fmt.Fprintf(os.Stderr, "target output invalid: %v\n%s", err, actualOut)
		os.Exit(1)
	}

	if err := compareAnswers(expectedAns, actualAns); err != nil {
		fmt.Fprintf(os.Stderr, "%v\ninput:\n%s", err, input)
		os.Exit(1)
	}

	fmt.Println("All tests passed.")
}
