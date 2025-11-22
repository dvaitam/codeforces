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
	n, m int
	A    []int
	X    []int
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierJ.go /path/to/binary")
		os.Exit(1)
	}
	target := os.Args[1]

	ref, cleanup, err := buildReference()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to build reference: %v\n", err)
		os.Exit(1)
	}
	defer cleanup()

	tests := append(deterministicTests(), randomTests()...)
	for i, tc := range tests {
		input := serialize(tc)

		refOut, err := runProgram(ref, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "reference runtime error on test %d: %v\ninput:\n%s", i+1, err, input)
			os.Exit(1)
		}
		refAns, err := parseOutput(refOut)
		if err != nil {
			fmt.Fprintf(os.Stderr, "reference output parse error on test %d: %v\noutput:\n%s", i+1, err, refOut)
			os.Exit(1)
		}

		candOut, err := runProgram(target, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "candidate runtime error on test %d: %v\ninput:\n%s", i+1, err, input)
			os.Exit(1)
		}
		candAns, err := parseOutput(candOut)
		if err != nil {
			fmt.Fprintf(os.Stderr, "candidate output parse error on test %d: %v\noutput:\n%s", i+1, err, candOut)
			os.Exit(1)
		}

		if refAns != candAns {
			fmt.Fprintf(os.Stderr, "test %d mismatch: expected %d got %d\ninput:\n%s", i+1, refAns, candAns, input)
			os.Exit(1)
		}
	}

	fmt.Println("All tests passed")
}

func buildReference() (string, func(), error) {
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		return "", nil, fmt.Errorf("cannot determine current path")
	}
	dir := filepath.Dir(file)
	tmpDir, err := os.MkdirTemp("", "ref-2045J-")
	if err != nil {
		return "", nil, err
	}
	outPath := filepath.Join(tmpDir, "ref2045J")
	cmd := exec.Command("go", "build", "-o", outPath, "2045J.go")
	cmd.Dir = dir
	if out, err := cmd.CombinedOutput(); err != nil {
		os.RemoveAll(tmpDir)
		return "", nil, fmt.Errorf("failed to build reference: %v\n%s", err, string(out))
	}
	cleanup := func() {
		os.RemoveAll(tmpDir)
	}
	return outPath, cleanup, nil
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
	return stdout.String(), nil
}

func serialize(tc testCase) string {
	var sb strings.Builder
	sb.Grow((tc.n + tc.m) * 12)
	sb.WriteString(fmt.Sprintf("%d %d\n", tc.n, tc.m))
	for i, v := range tc.A {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(strconv.Itoa(v))
	}
	sb.WriteByte('\n')
	for i, v := range tc.X {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(strconv.Itoa(v))
	}
	sb.WriteByte('\n')
	return sb.String()
}

func parseOutput(out string) (int64, error) {
	fields := strings.Fields(out)
	if len(fields) == 0 {
		return 0, fmt.Errorf("empty output")
	}
	val, err := strconv.ParseInt(fields[0], 10, 64)
	if err != nil {
		return 0, fmt.Errorf("invalid integer %q: %v", fields[0], err)
	}
	if val < 0 {
		return 0, fmt.Errorf("negative answer %d", val)
	}
	return val, nil
}

func deterministicTests() []testCase {
	return []testCase{
		// Sample 1
		{n: 3, m: 4, A: []int{0, 3, 0}, X: []int{1, 2, 1, 1}},
		// Sample 2
		{n: 5, m: 2, A: []int{0, 7, 13, 22, 24}, X: []int{12, 10}},
		// Sample 3
		{n: 3, m: 3, A: []int{0, 0, 0}, X: []int{1, 2, 3}},
		// Small distinct
		{n: 2, m: 2, A: []int{1, 3}, X: []int{0, 4}},
	}
}

func randomTests() []testCase {
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	tests := make([]testCase, 0, 150)
	for len(tests) < cap(tests) {
		n := rng.Intn(60) + 2
		m := rng.Intn(60) + 2
		A := make([]int, n)
		X := make([]int, m)
		for i := 0; i < n; i++ {
			switch rng.Intn(4) {
			case 0:
				A[i] = rng.Intn(8)
			case 1:
				A[i] = rng.Intn(1 << 10)
			default:
				A[i] = rng.Intn(1 << 20)
			}
		}
		for i := 0; i < m; i++ {
			switch rng.Intn(4) {
			case 0:
				X[i] = rng.Intn(4)
			case 1:
				X[i] = rng.Intn(1 << 10)
			default:
				X[i] = rng.Intn(1 << 20)
			}
		}
		tests = append(tests, testCase{n: n, m: m, A: A, X: X})
	}
	return tests
}
