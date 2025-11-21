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
	a []int64
}

func buildOracle() (string, func(), error) {
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		return "", nil, fmt.Errorf("cannot determine verifier path")
	}
	dir := filepath.Dir(file)
	tmpDir, err := os.MkdirTemp("", "oracle-2102B-")
	if err != nil {
		return "", nil, err
	}
	outPath := filepath.Join(tmpDir, "oracleB")
	cmd := exec.Command("go", "build", "-o", outPath, "2102B.go")
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
	sb.Grow(totalN*12 + len(tests)*24)
	sb.WriteString(strconv.Itoa(len(tests)))
	sb.WriteByte('\n')
	for _, tc := range tests {
		sb.WriteString(strconv.Itoa(tc.n))
		sb.WriteByte('\n')
		for i, v := range tc.a {
			if i > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(strconv.FormatInt(v, 10))
		}
		sb.WriteByte('\n')
	}
	return sb.String()
}

func parseOutput(out string, expected int) ([]string, error) {
	fields := strings.Fields(out)
	if len(fields) != expected {
		return nil, fmt.Errorf("expected %d answers, got %d", expected, len(fields))
	}
	res := make([]string, expected)
	for i, f := range fields {
		lf := strings.ToLower(f)
		if lf != "yes" && lf != "no" {
			return nil, fmt.Errorf("invalid answer %q at position %d", f, i+1)
		}
		res[i] = lf
	}
	return res, nil
}

func compareAnswers(expected, actual []string) error {
	if len(expected) != len(actual) {
		return fmt.Errorf("answer count mismatch")
	}
	for i := range expected {
		if expected[i] != actual[i] {
			return fmt.Errorf("answer %d mismatch: expected %s, got %s", i+1, expected[i], actual[i])
		}
	}
	return nil
}

func deterministicTests() []testCase {
	return []testCase{
		{n: 3, a: []int64{2, 3, 1}},                    // sample yes
		{n: 5, a: []int64{1, 2, 3, 4, 5}},              // yes (1 smallest)
		{n: 4, a: []int64{4, 2, 0, -5}},                // yes after flip a1 -> -4
		{n: 4, a: []int64{5, 100, -2, -3}},             // no
		{n: 1, a: []int64{10}},                         // median always a1
		{n: 2, a: []int64{-2, 5}},                      // need a1 median? ceil(2/2)=1, smallest => yes
		{n: 2, a: []int64{10, -1}},                     // no
		{n: 7, a: []int64{9, 8, 7, 6, 5, 4, 3}},        // a1 largest -> no
		{n: 7, a: []int64{-9, -1, -2, -3, -4, -5, -6}}, // a1 largest abs -> no
	}
}

func randomUniqueArray(n int, rng *rand.Rand) []int64 {
	perm := rng.Perm(2_000_000)
	arr := make([]int64, n)
	for i := 0; i < n; i++ {
		val := int64(perm[i] + 1)
		if rng.Intn(2) == 0 {
			val = -val
		}
		arr[i] = val
	}
	return arr
}

func randomTests(totalN int) []testCase {
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	tests := make([]testCase, 0, 128)
	used := 0
	for used < totalN {
		remain := totalN - used
		n := rng.Intn(min(5000, remain)) + 1
		arr := randomUniqueArray(n, rng)
		tests = append(tests, testCase{n: n, a: arr})
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

	actOut, err := runBinary(target, input)
	if err != nil {
		fmt.Fprintf(os.Stderr, "target runtime error: %v\ninput:\n%s", err, input)
		os.Exit(1)
	}

	expectedAns, err := parseOutput(expOut, len(tests))
	if err != nil {
		fmt.Fprintf(os.Stderr, "oracle output invalid: %v\n%s", err, expOut)
		os.Exit(1)
	}
	actualAns, err := parseOutput(actOut, len(tests))
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
