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

type query struct {
	l int
	r int
}

type testCase struct {
	n, k, q int
	a       []int
	queries []query
}

func buildOracle() (string, func(), error) {
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		return "", nil, fmt.Errorf("cannot determine verifier path")
	}
	dir := filepath.Dir(file)
	tmpDir, err := os.MkdirTemp("", "oracle-2009G2-")
	if err != nil {
		return "", nil, err
	}
	outPath := filepath.Join(tmpDir, "oracleG2")
	cmd := exec.Command("go", "build", "-o", outPath, "2009G2.go")
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
		sb.WriteString(fmt.Sprintf("%d %d %d\n", tc.n, tc.k, tc.q))
		for i, v := range tc.a {
			if i > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(strconv.Itoa(v))
		}
		sb.WriteByte('\n')
		for _, qu := range tc.queries {
			sb.WriteString(fmt.Sprintf("%d %d\n", qu.l, qu.r))
		}
	}
	return sb.String()
}

func parseOutput(out string, totalAnswers int) ([]int64, error) {
	fields := strings.Fields(out)
	if len(fields) != totalAnswers {
		return nil, fmt.Errorf("expected %d answers, got %d", totalAnswers, len(fields))
	}
	res := make([]int64, totalAnswers)
	for i, f := range fields {
		val, err := strconv.ParseInt(f, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("invalid integer %q: %v", f, err)
		}
		res[i] = val
	}
	return res, nil
}

func deterministicTests() []testCase {
	return []testCase{
		{
			n: 7, k: 5, q: 3,
			a: []int{1, 2, 3, 2, 1, 2, 3},
			queries: []query{
				{1, 7},
				{2, 7},
				{3, 7},
			},
		},
		{
			n: 8, k: 4, q: 2,
			a: []int{4, 3, 1, 1, 2, 4, 3, 2},
			queries: []query{
				{3, 6},
				{1, 5},
			},
		},
	}
}

func randomTests() []testCase {
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	tests := make([]testCase, 0, 100)
	for len(tests) < cap(tests) {
		n := rng.Intn(40) + 5
		k := rng.Intn(n) + 1
		if k > n {
			k = n
		}
		q := rng.Intn(40) + 1
		a := make([]int, n)
		for i := 0; i < n; i++ {
			a[i] = rng.Intn(n) + 1
		}
		queries := make([]query, q)
		for i := 0; i < q; i++ {
			l := rng.Intn(n) + 1
			r := rng.Intn(n-l+1) + l
			if r < l+k-1 {
				r = l + k - 1
			}
			if r > n {
				r = n
				l = max(1, r-k+1)
			}
			queries[i] = query{l: l, r: r}
		}
		tests = append(tests, testCase{n: n, k: k, q: q, a: a, queries: queries})
	}
	return tests
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func compareAnswers(expected, actual []int64) error {
	if len(expected) != len(actual) {
		return fmt.Errorf("answer count mismatch: expected %d, got %d", len(expected), len(actual))
	}
	for i := range expected {
		if expected[i] != actual[i] {
			return fmt.Errorf("answer %d mismatch: expected %d, got %d", i+1, expected[i], actual[i])
		}
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierG2.go /path/to/binary")
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

	totalAnswers := 0
	for _, tc := range tests {
		totalAnswers += tc.q
	}

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

	expectedAns, err := parseOutput(expectedOut, totalAnswers)
	if err != nil {
		fmt.Fprintf(os.Stderr, "oracle output invalid: %v\n%s", err, expectedOut)
		os.Exit(1)
	}
	actualAns, err := parseOutput(actualOut, totalAnswers)
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
