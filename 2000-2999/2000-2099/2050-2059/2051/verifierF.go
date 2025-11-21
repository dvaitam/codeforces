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
	n   int64
	m   int64
	ops []int64
}

func buildOracle() (string, func(), error) {
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		return "", nil, fmt.Errorf("cannot determine verifier path")
	}
	dir := filepath.Dir(file)
	tmpDir, err := os.MkdirTemp("", "oracle-2051F-")
	if err != nil {
		return "", nil, err
	}
	outPath := filepath.Join(tmpDir, "oracleF")
	cmd := exec.Command("go", "build", "-o", outPath, "2051F.go")
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
	totalOps := 0
	for _, tc := range tests {
		totalOps += len(tc.ops)
	}

	var sb strings.Builder
	sb.Grow(totalOps*12 + len(tests)*64)
	sb.WriteString(strconv.Itoa(len(tests)))
	sb.WriteByte('\n')
	for _, tc := range tests {
		sb.WriteString(fmt.Sprintf("%d %d %d\n", tc.n, tc.m, len(tc.ops)))
		for i, v := range tc.ops {
			if i > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(strconv.FormatInt(v, 10))
		}
		sb.WriteByte('\n')
	}
	return sb.String()
}

func parseOutput(out string, expected int) ([]int64, error) {
	fields := strings.Fields(out)
	if len(fields) != expected {
		return nil, fmt.Errorf("expected %d answers, got %d", expected, len(fields))
	}
	res := make([]int64, expected)
	for i, f := range fields {
		v, err := strconv.ParseInt(f, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("invalid integer %q at position %d", f, i+1)
		}
		res[i] = v
	}
	return res, nil
}

func deterministicTests() []testCase {
	return []testCase{
		{n: 6, m: 5, ops: []int64{1, 2, 3}},                                        // sample-style small deck
		{n: 2, m: 1, ops: []int64{2, 1}},                                           // joker starts on top edge
		{n: 10, m: 10, ops: []int64{10, 9, 1, 10, 5}},                              // joker at the bottom, mixed pulls
		{n: 5, m: 3, ops: []int64{3, 3, 3}},                                        // repeatedly pulling joker position
		{n: 8, m: 4, ops: []int64{1, 8, 2, 7, 3, 6, 4, 5, 4}},                      // covers both ends and middle
		{n: 1_000_000_000, m: 500_000_000, ops: []int64{1, 1_000_000_000, 2, 999}}, // stress large n
	}
}

func randomTests(maxTotalOps int) []testCase {
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	tests := make([]testCase, 0, 200)
	total := 0
	for total < maxTotalOps {
		remaining := maxTotalOps - total
		q := rng.Intn(400) + 1
		if q > remaining {
			q = remaining
		}
		n := rng.Int63n(1_000_000_000-1) + 2
		m := rng.Int63n(n) + 1
		ops := make([]int64, q)
		for i := 0; i < q; i++ {
			if rng.Intn(4) == 0 && n < 50 {
				ops[i] = int64(rng.Intn(int(n)) + 1)
			} else {
				ops[i] = rng.Int63n(n) + 1
			}
		}
		tests = append(tests, testCase{n: n, m: m, ops: ops})
		total += q
	}
	return tests
}

func heavyCase(maxTotalOps int) []testCase {
	if maxTotalOps <= 0 {
		return nil
	}
	n := int64(1_000_000_000)
	q := maxTotalOps
	ops := make([]int64, q)
	for i := 0; i < q; i++ {
		if i%2 == 0 {
			ops[i] = 1
		} else {
			ops[i] = n
		}
	}
	return []testCase{{n: n, m: n / 2, ops: ops}}
}

func compareAnswers(tests []testCase, expected, actual []int64) error {
	if len(expected) != len(actual) {
		return fmt.Errorf("answer count mismatch: expected %d, got %d", len(expected), len(actual))
	}
	idx := 0
	for ti, tc := range tests {
		for oi := range tc.ops {
			if expected[idx] != actual[idx] {
				return fmt.Errorf("test %d operation %d mismatch: expected %d, got %d", ti+1, oi+1, expected[idx], actual[idx])
			}
			idx++
		}
	}
	return nil
}

func totalOps(tests []testCase) int {
	total := 0
	for _, tc := range tests {
		total += len(tc.ops)
	}
	return total
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierF.go /path/to/binary")
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
	currentTotal := totalOps(tests)
	const limit = 200_000
	randomBudget := limit - currentTotal - 10_000
	if randomBudget < 0 {
		randomBudget = 0
	}
	tests = append(tests, randomTests(randomBudget)...)
	currentTotal = totalOps(tests)
	if currentTotal < limit {
		tests = append(tests, heavyCase(limit-currentTotal)...)
	}

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

	expectedAns, err := parseOutput(expectedOut, totalOps(tests))
	if err != nil {
		fmt.Fprintf(os.Stderr, "oracle output invalid: %v\n%s", err, expectedOut)
		os.Exit(1)
	}
	actualAns, err := parseOutput(actualOut, totalOps(tests))
	if err != nil {
		fmt.Fprintf(os.Stderr, "target output invalid: %v\n%s", err, actualOut)
		os.Exit(1)
	}

	if err := compareAnswers(tests, expectedAns, actualAns); err != nil {
		fmt.Fprintf(os.Stderr, "%v\ninput:\n%s", err, input)
		os.Exit(1)
	}

	fmt.Println("All tests passed.")
}
