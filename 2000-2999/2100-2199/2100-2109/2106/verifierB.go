package main

import (
	"bytes"
	"context"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

const (
	refSource   = "./2106B.go"
	totalNLimit = 200000
	defaultTime = 15 * time.Second
)

type testCase struct {
	n int
	x int
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierB.go /path/to/candidate")
		os.Exit(1)
	}
	candidate := os.Args[1]

	tests := buildTests()
	input := buildInput(tests)

	refBin, cleanup, err := buildReference()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to build reference: %v\n", err)
		os.Exit(1)
	}
	defer cleanup()

	refOut, err := runProgram(refBin, input)
	if err != nil {
		fmt.Fprintf(os.Stderr, "reference runtime error: %v\n", err)
		os.Exit(1)
	}
	refPerms, err := parsePerms(refOut, tests)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to parse reference output: %v\noutput:\n%s\n", err, refOut)
		os.Exit(1)
	}

	expectedCounts := make([]int, len(tests))
	for i, tc := range tests {
		if err := validatePermutation(refPerms[i], tc.n); err != nil {
			fmt.Fprintf(os.Stderr, "reference produced invalid permutation on test %d: %v\n", i+1, err)
			os.Exit(1)
		}
		expectedCounts[i] = countColor(refPerms[i], tc.x)
	}

	candOut, err := runProgram(candidate, input)
	if err != nil {
		fmt.Fprintf(os.Stderr, "candidate runtime error: %v\n", err)
		os.Exit(1)
	}
	candPerms, err := parsePerms(candOut, tests)
	if err != nil {
		fmt.Fprintf(os.Stderr, "invalid candidate output: %v\noutput:\n%s\n", err, candOut)
		os.Exit(1)
	}

	for i, tc := range tests {
		perm := candPerms[i]
		if err := validatePermutation(perm, tc.n); err != nil {
			fmt.Fprintf(os.Stderr, "test %d: invalid permutation: %v\ninput:\n%s", i+1, err, singleInput(tc))
			os.Exit(1)
		}
		got := countColor(perm, tc.x)
		if got != expectedCounts[i] {
			fmt.Fprintf(os.Stderr, "test %d: wrong number of color %d cells, expected %d got %d\ninput:\n%s", i+1, tc.x, expectedCounts[i], got, singleInput(tc))
			os.Exit(1)
		}
	}

	fmt.Printf("All %d tests passed.\n", len(tests))
}

func buildReference() (string, func(), error) {
	dir, err := os.MkdirTemp("", "verifier-2106B-*")
	if err != nil {
		return "", nil, err
	}
	bin := filepath.Join(dir, "ref2106B")
	cmd := exec.Command("go", "build", "-o", bin, refSource)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		os.RemoveAll(dir)
		return "", nil, fmt.Errorf("go build failed: %v\n%s", err, out.String())
	}
	return bin, func() { os.RemoveAll(dir) }, nil
}

func buildTests() []testCase {
	tests := deterministicTests()
	total := 0
	for _, tc := range tests {
		total += tc.n
	}

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for total < totalNLimit {
		n := rng.Intn(4000) + 1
		if total+n > totalNLimit {
			n = totalNLimit - total
		}
		if n <= 0 {
			break
		}
		x := rng.Intn(n + 1)
		tests = append(tests, testCase{n: n, x: x})
		total += n
	}
	return tests
}

func deterministicTests() []testCase {
	return []testCase{
		{n: 1, x: 0},
		{n: 2, x: 0},
		{n: 2, x: 2},
		{n: 4, x: 2},
		{n: 5, x: 0},
		{n: 5, x: 4},
		{n: 6, x: 3},
		{n: 8, x: 7},
		{n: 10, x: 5},
	}
}

func buildInput(tests []testCase) string {
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d\n", len(tests))
	for _, tc := range tests {
		fmt.Fprintf(&sb, "%d %d\n", tc.n, tc.x)
	}
	return sb.String()
}

func singleInput(tc testCase) string {
	return fmt.Sprintf("1\n%d %d\n", tc.n, tc.x)
}

func runProgram(path, input string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), defaultTime)
	defer cancel()
	cmd := commandFor(ctx, path)
	cmd.Stdin = strings.NewReader(input)
	var outBuf, errBuf bytes.Buffer
	cmd.Stdout = &outBuf
	cmd.Stderr = &errBuf
	if err := cmd.Run(); err != nil {
		if ctx.Err() == context.DeadlineExceeded {
			return "", fmt.Errorf("timeout")
		}
		return "", fmt.Errorf("%v\n%s", err, errBuf.String())
	}
	return outBuf.String(), nil
}

func commandFor(ctx context.Context, path string) *exec.Cmd {
	switch filepath.Ext(path) {
	case ".go":
		return exec.CommandContext(ctx, "go", "run", path)
	case ".py":
		return exec.CommandContext(ctx, "python3", path)
	default:
		return exec.CommandContext(ctx, path)
	}
}

func parsePerms(out string, tests []testCase) ([][]int, error) {
	fields := strings.Fields(out)
	idx := 0
	res := make([][]int, len(tests))
	for i, tc := range tests {
		if idx+tc.n > len(fields) {
			return nil, fmt.Errorf("not enough numbers for test %d", i+1)
		}
		cur := make([]int, tc.n)
		for j := 0; j < tc.n; j++ {
			v, err := strconv.Atoi(fields[idx+j])
			if err != nil {
				return nil, fmt.Errorf("invalid integer %q", fields[idx+j])
			}
			cur[j] = v
		}
		res[i] = cur
		idx += tc.n
	}
	if idx != len(fields) {
		return nil, fmt.Errorf("extra tokens in output")
	}
	return res, nil
}

func validatePermutation(p []int, n int) error {
	if len(p) != n {
		return fmt.Errorf("expected %d elements, got %d", n, len(p))
	}
	seen := make([]bool, n)
	for i, v := range p {
		if v < 0 || v >= n {
			return fmt.Errorf("value %d out of range at position %d", v, i+1)
		}
		if seen[v] {
			return fmt.Errorf("duplicate value %d", v)
		}
		seen[v] = true
	}
	return nil
}

func countColor(p []int, target int) int {
	n := len(p)
	seen := make([]bool, n+1)
	mex := 0
	cnt := 0
	for _, v := range p {
		seen[v] = true
		for seen[mex] {
			mex++
		}
		if mex == target {
			cnt++
		}
	}
	return cnt
}
