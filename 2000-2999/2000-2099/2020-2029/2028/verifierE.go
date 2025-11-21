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
}

type testCase struct {
	n     int
	edges []edge
}

const mod int64 = 998244353

func buildOracle() (string, func(), error) {
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		return "", nil, fmt.Errorf("cannot determine verifier path")
	}
	dir := filepath.Dir(file)
	tmpDir, err := os.MkdirTemp("", "oracle-2028E-")
	if err != nil {
		return "", nil, err
	}
	outPath := filepath.Join(tmpDir, "oracleE")
	cmd := exec.Command("go", "build", "-o", outPath, "2028E.go")
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
	sb.Grow(len(tests) * 256)
	sb.WriteString(strconv.Itoa(len(tests)))
	sb.WriteByte('\n')
	for _, tc := range tests {
		sb.WriteString(strconv.Itoa(tc.n))
		sb.WriteByte('\n')
		for _, e := range tc.edges {
			sb.WriteString(fmt.Sprintf("%d %d\n", e.u, e.v))
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
			return nil, fmt.Errorf("not enough outputs for test %d", i+1)
		}
		ans := make([]int64, tc.n)
		for j := 0; j < tc.n; j++ {
			val, err := strconv.ParseInt(fields[idx+j], 10, 64)
			if err != nil {
				return nil, fmt.Errorf("invalid integer %q", fields[idx+j])
			}
			val %= mod
			if val < 0 {
				val += mod
			}
			ans[j] = val
		}
		res[i] = ans
		idx += tc.n
	}
	if idx != len(fields) {
		return nil, fmt.Errorf("extra values in output")
	}
	return res, nil
}

func deterministicTests() []testCase {
	return []testCase{
		{
			n: 2,
			edges: []edge{
				{1, 2},
			},
		},
		{
			n: 5,
			edges: []edge{
				{1, 2},
				{2, 3},
				{2, 4},
				{4, 5},
			},
		},
		{
			n: 9,
			edges: []edge{
				{1, 2}, {1, 3}, {3, 4}, {4, 5}, {5, 6}, {3, 7}, {7, 8}, {8, 9},
			},
		},
	}
}

func randomTree(rng *rand.Rand, n int) []edge {
	edges := make([]edge, 0, n-1)
	for i := 2; i <= n; i++ {
		p := rng.Intn(i-1) + 1
		edges = append(edges, edge{u: p, v: i})
	}
	return edges
}

func randomTests() []testCase {
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	tests := make([]testCase, 0, 120)
	totalN := 0
	for len(tests) < cap(tests) && totalN < 180000 {
		n := rng.Intn(2000) + 2
		if totalN+n > 200000 {
			break
		}
		edges := randomTree(rng, n)
		tests = append(tests, testCase{n: n, edges: edges})
		totalN += n
	}
	return tests
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
				return fmt.Errorf("test %d vertex %d mismatch: expected %d, got %d", i+1, j+1, expected[i][j], actual[i][j])
			}
		}
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierE.go /path/to/binary")
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
