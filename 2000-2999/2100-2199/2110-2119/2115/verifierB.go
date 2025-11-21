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
	x int
	y int
	z int
}

type testCase struct {
	n       int
	q       int
	a       []int64
	queries []query
}

func buildOracle() (string, func(), error) {
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		return "", nil, fmt.Errorf("cannot determine verifier path")
	}
	dir := filepath.Dir(file)
	tmpDir, err := os.MkdirTemp("", "oracle-2115B-")
	if err != nil {
		return "", nil, err
	}
	outPath := filepath.Join(tmpDir, "oracleB")
	cmd := exec.Command("go", "build", "-o", outPath, "2115B.go")
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
	totalN := 0
	for _, tc := range tests {
		totalN += tc.n
	}
	sb.Grow(totalN*16 + len(tests)*32)
	sb.WriteString(strconv.Itoa(len(tests)))
	sb.WriteByte('\n')
	for _, tc := range tests {
		sb.WriteString(fmt.Sprintf("%d %d\n", tc.n, tc.q))
		for i, v := range tc.a {
			if i > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(strconv.FormatInt(v, 10))
		}
		sb.WriteByte('\n')
		for _, qu := range tc.queries {
			sb.WriteString(fmt.Sprintf("%d %d %d\n", qu.x, qu.y, qu.z))
		}
	}
	return sb.String()
}

func parseOutput(out string, tests []testCase) ([][]int64, error) {
	tokens := strings.Fields(out)
	idx := 0
	res := make([][]int64, len(tests))
	for i, tc := range tests {
		if idx+tc.n > len(tokens) {
			return nil, fmt.Errorf("test %d: not enough numbers in output", i+1)
		}
		arr := make([]int64, tc.n)
		for j := 0; j < tc.n; j++ {
			v, err := strconv.ParseInt(tokens[idx+j], 10, 64)
			if err != nil {
				return nil, fmt.Errorf("test %d: invalid integer %q", i+1, tokens[idx+j])
			}
			arr[j] = v
		}
		idx += tc.n
		res[i] = arr
	}
	if idx != len(tokens) {
		return nil, fmt.Errorf("extra tokens detected after parsing output")
	}
	return res, nil
}

func compareAnswers(expected, actual [][]int64) error {
	if len(expected) != len(actual) {
		return fmt.Errorf("answer test count mismatch")
	}
	for i := range expected {
		if len(expected[i]) != len(actual[i]) {
			return fmt.Errorf("test %d length mismatch", i+1)
		}
		for j := range expected[i] {
			if expected[i][j] != actual[i][j] {
				return fmt.Errorf("test %d position %d mismatch: expected %d, got %d", i+1, j+1, expected[i][j], actual[i][j])
			}
		}
	}
	return nil
}

func deterministicTests() []testCase {
	return []testCase{
		{n: 1, q: 0, a: []int64{5}},
		{n: 3, q: 2, a: []int64{1, 2, 3}, queries: []query{{1, 1, 1}, {1, 2, 2}}},
		{n: 5, q: 4, a: []int64{10, 20, 30, 40, 50}, queries: []query{{1, 5, 3}, {2, 4, 4}, {3, 3, 1}, {1, 1, 1}}},
	}
}

func randomTests(totalN int) []testCase {
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	tests := make([]testCase, 0, 128)
	used := 0
	for used < totalN {
		remain := totalN - used
		n := rng.Intn(min(2000, remain)) + 1
		q := rng.Intn(200) // keep small for input size
		a := make([]int64, n)
		for i := 0; i < n; i++ {
			a[i] = rng.Int63n(1_000_000) - 500_000
		}
		qs := make([]query, q)
		for i := 0; i < q; i++ {
			l := rng.Intn(n) + 1
			r := rng.Intn(n) + 1
			if l > r {
				l, r = r, l
			}
			z := rng.Intn(n) + 1
			qs[i] = query{x: l, y: r, z: z}
		}
		tests = append(tests, testCase{n: n, q: q, a: a, queries: qs})
		used += n
	}
	return tests
}

func totalN(tests []testCase) int {
	sum := 0
	for _, tc := range tests {
		sum += tc.n
	}
	return sum
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	target := os.Args[1]

	oracle, cleanup, err := buildOracle()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
	defer cleanup()

	tests := deterministicTests()
	const nLimit = 100_000
	used := totalN(tests)
	if used < nLimit {
		tests = append(tests, randomTests(nLimit-used)...)
	}

	input := buildInput(tests)

	expOut, err := runBinary(oracle, input)
	if err != nil {
		fmt.Fprintf(os.Stderr, "oracle failed: %v\ninput:\n%s", err, input)
		os.Exit(1)
	}
	expectedAns, err := parseOutput(expOut, tests)
	if err != nil {
		fmt.Fprintf(os.Stderr, "oracle output invalid: %v\n%s", err, expOut)
		os.Exit(1)
	}

	actOut, err := runBinary(target, input)
	if err != nil {
		fmt.Fprintf(os.Stderr, "target runtime error: %v\ninput:\n%s", err, input)
		os.Exit(1)
	}
	actualAns, err := parseOutput(actOut, tests)
	if err != nil {
		fmt.Fprintf(os.Stderr, "target output invalid: %v\n%s", err, actOut)
		os.Exit(1)
	}

	if err := compareAnswers(expectedAns, actualAns); err != nil {
		fmt.Fprintf(os.Stderr, "%v\ninput:\n%s", err, input)
		os.Exit(1)
	}

	fmt.Println("All tests passed.")
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
