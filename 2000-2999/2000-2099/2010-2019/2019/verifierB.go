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

type testCase struct {
	n int
	q int
	x []int64
	k []int64
}

func buildOracle() (string, func(), error) {
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		return "", nil, fmt.Errorf("cannot determine verifier path")
	}
	dir := filepath.Dir(file)
	tmpDir, err := os.MkdirTemp("", "oracle-2019B-")
	if err != nil {
		return "", nil, err
	}
	outPath := filepath.Join(tmpDir, "oracleB")
	cmd := exec.Command("go", "build", "-o", outPath, "2019B.go")
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
		sb.WriteString(fmt.Sprintf("%d %d\n", tc.n, tc.q))
		for i, v := range tc.x {
			if i > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(strconv.FormatInt(v, 10))
		}
		sb.WriteByte('\n')
		for i, val := range tc.k {
			if i > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(strconv.FormatInt(val, 10))
		}
		sb.WriteByte('\n')
	}
	return sb.String()
}

func parseOutput(out string, tests []testCase) ([][]int64, error) {
	fields := strings.Fields(out)
	res := make([][]int64, len(tests))
	idx := 0
	for ti, tc := range tests {
		if idx+tc.q > len(fields) {
			return nil, fmt.Errorf("not enough answers for test %d", ti+1)
		}
		ans := make([]int64, tc.q)
		for j := 0; j < tc.q; j++ {
			val, err := strconv.ParseInt(fields[idx+j], 10, 64)
			if err != nil {
				return nil, fmt.Errorf("invalid integer %q", fields[idx+j])
			}
			ans[j] = val
		}
		res[ti] = ans
		idx += tc.q
	}
	if idx != len(fields) {
		return nil, fmt.Errorf("extra output values")
	}
	return res, nil
}

func deterministicTests() []testCase {
	return []testCase{
		{
			n: 2, q: 2,
			x: []int64{101, 200},
			k: []int64{2, 1},
		},
		{
			n: 6, q: 15,
			x: []int64{1, 2, 3, 5, 6, 7},
			k: []int64{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15},
		},
		{
			n: 5, q: 8,
			x: []int64{2, 4, 6, 8, 10},
			k: []int64{1, 2, 3, 4, 5, 6, 7, 8},
		},
	}
}

func randomTests() []testCase {
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	tests := make([]testCase, 0, 120)
	for len(tests) < cap(tests) {
		n := rng.Intn(60) + 2
		q := rng.Intn(60) + 1
		x := make([]int64, n)
		cur := int64(rng.Intn(5) + 1)
		for i := 0; i < n; i++ {
			incr := int64(rng.Intn(10) + 1)
			cur += incr
			x[i] = cur
			cur = x[i]
		}
		k := make([]int64, q)
		for i := 0; i < q; i++ {
			if rng.Intn(3) == 0 {
				k[i] = int64(rng.Intn(n*n) + 1)
			} else {
				k[i] = randRangeInt64(rng, 1, 1_000_000_000_000)
			}
		}
		tests = append(tests, testCase{n: n, q: q, x: x, k: k})
	}
	return tests
}

func randRangeInt64(rng *rand.Rand, lo, hi int64) int64 {
	return lo + rng.Int63n(hi-lo+1)
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
				return fmt.Errorf("test %d query %d mismatch: expected %d, got %d",
					i+1, j+1, expected[i][j], actual[i][j])
			}
		}
	}
	return nil
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
