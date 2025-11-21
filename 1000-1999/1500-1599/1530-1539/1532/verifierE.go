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
	n   int
	arr []int
}

type parsedAnswer struct {
	k       int
	indices []int
}

func buildOracle() (string, func(), error) {
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		return "", nil, fmt.Errorf("cannot determine verifier path")
	}
	dir := filepath.Dir(file)
	tmpDir, err := os.MkdirTemp("", "oracle-1532E-")
	if err != nil {
		return "", nil, err
	}
	outPath := filepath.Join(tmpDir, "oracleE")
	cmd := exec.Command("go", "build", "-o", outPath, "1532E.go")
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

func buildInput(tc testCase) string {
	var sb strings.Builder
	sb.Grow(tc.n*6 + 16)
	sb.WriteString(strconv.Itoa(tc.n))
	sb.WriteByte('\n')
	for i, v := range tc.arr {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(strconv.Itoa(v))
	}
	sb.WriteByte('\n')
	return sb.String()
}

func parseOutput(output string, n int) (parsedAnswer, error) {
	output = strings.TrimSpace(output)
	if len(output) == 0 {
		return parsedAnswer{}, fmt.Errorf("empty output")
	}
	tokens := strings.Fields(output)
	if len(tokens) == 0 {
		return parsedAnswer{}, fmt.Errorf("no tokens in output")
	}
	k, err := strconv.Atoi(tokens[0])
	if err != nil {
		return parsedAnswer{}, fmt.Errorf("invalid k: %v", err)
	}
	if k < 0 || k > n {
		return parsedAnswer{}, fmt.Errorf("k out of range: %d", k)
	}
	if len(tokens)-1 != k {
		return parsedAnswer{}, fmt.Errorf("expected %d indices, got %d", k, len(tokens)-1)
	}
	indices := make([]int, k)
	seen := make(map[int]bool, k)
	for i := 0; i < k; i++ {
		val, err := strconv.Atoi(tokens[i+1])
		if err != nil {
			return parsedAnswer{}, fmt.Errorf("invalid index %q: %v", tokens[i+1], err)
		}
		if val < 1 || val > n {
			return parsedAnswer{}, fmt.Errorf("index %d out of range", val)
		}
		if seen[val] {
			return parsedAnswer{}, fmt.Errorf("duplicate index %d", val)
		}
		seen[val] = true
		indices[i] = val
	}
	return parsedAnswer{k: k, indices: indices}, nil
}

func deterministicTests() []testCase {
	return []testCase{
		{n: 5, arr: []int{2, 5, 1, 2, 2}},
		{n: 4, arr: []int{8, 3, 5, 2}},
		{n: 5, arr: []int{2, 1, 2, 4, 3}},
		{n: 2, arr: []int{1, 1}},
		{n: 3, arr: []int{1, 2, 3}},
		{n: 6, arr: []int{1, 2, 3, 6, 7, 1}},
		{n: 3, arr: []int{10, 1, 9}},
	}
}

func randomTests() []testCase {
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	tests := make([]testCase, 0, 150)
	for len(tests) < cap(tests) {
		n := rng.Intn(60) + 2
		if rng.Intn(5) == 0 {
			n = rng.Intn(200) + 2
		}
		arr := make([]int, n)
		for i := 0; i < n; i++ {
			val := rng.Intn(1_000_000) + 1
			if rng.Intn(6) == 0 {
				val = rng.Intn(50) + 1
			}
			arr[i] = val
		}
		tests = append(tests, testCase{n: n, arr: arr})
	}
	return tests
}

func toSet(indices []int) map[int]struct{} {
	m := make(map[int]struct{}, len(indices))
	for _, v := range indices {
		m[v] = struct{}{}
	}
	return m
}

func compareAnswers(expected, actual parsedAnswer) error {
	if expected.k != actual.k {
		return fmt.Errorf("k mismatch: expected %d, got %d", expected.k, actual.k)
	}
	expSet := toSet(expected.indices)
	actSet := toSet(actual.indices)
	if len(expSet) != len(actSet) {
		return fmt.Errorf("distinct index count mismatch: expected %d, got %d", len(expSet), len(actSet))
	}
	for idx := range expSet {
		if _, ok := actSet[idx]; !ok {
			return fmt.Errorf("missing index %d", idx)
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
	for idx, tc := range tests {
		input := buildInput(tc)
		expectedOut, err := runBinary(oracle, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "oracle failed on test %d: %v\ninput:\n%s", idx+1, err, input)
			os.Exit(1)
		}
		actualOut, err := runBinary(target, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "target runtime error on test %d: %v\ninput:\n%s", idx+1, err, input)
			os.Exit(1)
		}
		expectedAns, err := parseOutput(expectedOut, tc.n)
		if err != nil {
			fmt.Fprintf(os.Stderr, "oracle output invalid on test %d: %v\noutput:\n%s", idx+1, err, expectedOut)
			os.Exit(1)
		}
		actualAns, err := parseOutput(actualOut, tc.n)
		if err != nil {
			fmt.Fprintf(os.Stderr, "target output invalid on test %d: %v\noutput:\n%s", idx+1, err, actualOut)
			os.Exit(1)
		}
		if err := compareAnswers(expectedAns, actualAns); err != nil {
			fmt.Fprintf(os.Stderr, "test %d mismatch: %v\ninput:\n%s", idx+1, err, input)
			os.Exit(1)
		}
	}

	fmt.Println("All tests passed.")
}
