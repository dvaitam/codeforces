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

const refSource = "2089C1.go"
const mod = 1000000007

type testCase struct {
	n int
	l int
	k int
}

func main() {
	args := os.Args[1:]
	if len(args) == 2 && args[0] == "--" {
		args = args[1:]
	}
	if len(args) != 1 {
		fmt.Println("usage: go run verifierC1.go /path/to/binary")
		os.Exit(1)
	}
	candidate := args[0]

	refBin, err := buildReference()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer os.Remove(refBin)

	tests := buildTests()
	input := buildInput(tests)

	refOut, err := runProgram(refBin, input)
	if err != nil {
		fmt.Fprintf(os.Stderr, "reference failed: %v\n", err)
		os.Exit(1)
	}
	refAns, err := parseOutput(refOut, tests)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to parse reference output: %v\n%s", err, refOut)
		os.Exit(1)
	}

	candOut, err := runProgram(candidate, input)
	if err != nil {
		fmt.Fprintf(os.Stderr, "candidate runtime error: %v\n", err)
		os.Exit(1)
	}
	candAns, err := parseOutput(candOut, tests)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to parse candidate output: %v\n%s", err, candOut)
		os.Exit(1)
	}

	for idx := range tests {
		r := refAns[idx]
		c := candAns[idx]
		if len(r) != len(c) {
			tc := tests[idx]
			fmt.Fprintf(os.Stderr, "test %d output length mismatch (expected %d got %d) for n=%d l=%d k=%d\n", idx+1, len(r), len(c), tc.n, tc.l, tc.k)
			os.Exit(1)
		}
		for i := range r {
			if r[i] != c[i]%mod {
				tc := tests[idx]
				fmt.Fprintf(os.Stderr, "test %d mismatch at position %d: expected %d got %d (n=%d l=%d k=%d)\n", idx+1, i+1, r[i], c[i]%mod, tc.n, tc.l, tc.k)
				fmt.Fprintf(os.Stderr, "candidate outputs: %v\n", c)
				fmt.Fprintf(os.Stderr, "reference outputs: %v\n", r)
				os.Exit(1)
			}
		}
	}
	fmt.Println("All tests passed")
}

func buildReference() (string, error) {
	refPath, err := referencePath()
	if err != nil {
		return "", err
	}
	tmp, err := os.CreateTemp("", "ref_2089C1_*.bin")
	if err != nil {
		return "", fmt.Errorf("failed to create temp file: %v", err)
	}
	tmp.Close()

	cmd := exec.Command("go", "build", "-o", tmp.Name(), refPath)
	if out, err := cmd.CombinedOutput(); err != nil {
		os.Remove(tmp.Name())
		return "", fmt.Errorf("failed to build reference: %v\n%s", err, string(out))
	}
	return tmp.Name(), nil
}

func referencePath() (string, error) {
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		return "", fmt.Errorf("failed to locate verifier path")
	}
	dir := filepath.Dir(file)
	return filepath.Join(dir, refSource), nil
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
	return strings.TrimSpace(stdout.String()), nil
}

func buildTests() []testCase {
	var tests []testCase
	add := func(n, l int) {
		tests = append(tests, testCase{n: n, l: l, k: 0})
	}

	// Sample-inspired cases
	add(3, 1)
	add(3, 2)
	add(2, 5)
	add(9, 10)
	add(4, 1)

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	totalL := 1 + 2 + 5 + 10 + 1
	for len(tests) < 40 {
		n := rng.Intn(100) + 1
		remaining := 5000 - totalL
		if remaining <= 0 {
			break
		}
		l := rng.Intn(remaining) + 1
		add(n, l)
		totalL += l
	}

	if totalL < 5000 {
		add(100, 5000-totalL)
	}

	return tests
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func buildInput(tests []testCase) string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", len(tests)))
	for _, tc := range tests {
		sb.WriteString(fmt.Sprintf("%d %d %d\n", tc.n, tc.l, tc.k))
	}
	return sb.String()
}

func parseOutput(out string, tests []testCase) ([][]int64, error) {
	lines := strings.Split(strings.TrimSpace(out), "\n")
	if len(lines) != len(tests) {
		return nil, fmt.Errorf("expected %d lines, got %d", len(tests), len(lines))
	}
	res := make([][]int64, len(tests))
	for idx, line := range lines {
		fields := strings.Fields(line)
		if len(fields) != tests[idx].n {
			return nil, fmt.Errorf("line %d expected %d numbers, got %d", idx+1, tests[idx].n, len(fields))
		}
		row := make([]int64, len(fields))
		for i, f := range fields {
			val, err := strconv.ParseInt(f, 10, 64)
			if err != nil {
				return nil, fmt.Errorf("invalid integer %q on line %d", f, idx+1)
			}
			val %= mod
			if val < 0 {
				val += mod
			}
			row[i] = val
		}
		res[idx] = row
	}
	return res, nil
}
