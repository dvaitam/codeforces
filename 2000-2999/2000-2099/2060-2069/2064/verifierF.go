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
	refSource = "2000-2999/2000-2099/2060-2069/2064/2064F.go"
	totalNCap = 200000
)

type testCase struct {
	n int
	k int
	a []int
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierF.go /path/to/candidate")
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
	refAns, err := parseOutputs(refOut, len(tests))
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to parse reference output: %v\noutput:\n%s\n", err, refOut)
		os.Exit(1)
	}

	candOut, err := runProgram(candidate, input)
	if err != nil {
		fmt.Fprintf(os.Stderr, "candidate runtime error: %v\n", err)
		os.Exit(1)
	}
	candAns, err := parseOutputs(candOut, len(tests))
	if err != nil {
		fmt.Fprintf(os.Stderr, "invalid candidate output: %v\noutput:\n%s\n", err, candOut)
		os.Exit(1)
	}

	for i := range refAns {
		if candAns[i] != refAns[i] {
			fmt.Fprintf(os.Stderr, "wrong answer on test %d\nexpected: %d\nfound: %d\n", i+1, refAns[i], candAns[i])
			os.Exit(1)
		}
	}

	fmt.Printf("All %d tests passed.\n", len(tests))
}

func buildReference() (string, func(), error) {
	dir, err := os.MkdirTemp("", "verifier-2064F-*")
	if err != nil {
		return "", nil, err
	}
	bin := filepath.Join(dir, "ref2064F")
	cmd := exec.Command("go", "build", "-o", bin, refSource)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		os.RemoveAll(dir)
		return "", nil, fmt.Errorf("go build failed: %v\n%s", err, out.String())
	}
	cleanup := func() { os.RemoveAll(dir) }
	return bin, cleanup, nil
}

func buildTests() []testCase {
	tests := deterministicTests()
	total := 0
	for _, tc := range tests {
		total += tc.n
	}

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for total < totalNCap {
		// Bias towards small and medium sizes, but occasionally pick a bigger one.
		n := rng.Intn(4000) + 2
		if rng.Intn(5) == 0 {
			n = rng.Intn(30000) + 2000
		}
		if n > totalNCap-total {
			n = totalNCap - total
		}
		if n < 2 {
			break
		}
		k := rng.Intn(n-1) + n + 1 // n < k < 2n
		a := make([]int, n)
		for i := 0; i < n; i++ {
			a[i] = rng.Intn(n) + 1
		}
		tests = append(tests, testCase{n: n, k: k, a: a})
		total += n
	}
	return tests
}

func deterministicTests() []testCase {
	return []testCase{
		{n: 2, k: 3, a: []int{1, 2}},
		{n: 3, k: 4, a: []int{2, 2, 2}},
		{n: 4, k: 7, a: []int{1, 3, 4, 2}},
		{n: 5, k: 7, a: []int{1, 2, 3, 4, 5}},
		{n: 5, k: 9, a: []int{5, 5, 5, 5, 5}},
		{n: 6, k: 10, a: []int{1, 6, 2, 6, 3, 6}},
		{n: 8, k: 11, a: []int{4, 4, 4, 7, 7, 7, 5, 6}},
	}
}

func buildInput(tests []testCase) string {
	var sb strings.Builder
	sb.Grow(len(tests) * 32)
	fmt.Fprintf(&sb, "%d\n", len(tests))
	for _, tc := range tests {
		fmt.Fprintf(&sb, "%d %d\n", tc.n, tc.k)
		for i, v := range tc.a {
			if i > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(strconv.Itoa(v))
		}
		sb.WriteByte('\n')
	}
	return sb.String()
}

func runProgram(path, input string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
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

func parseOutputs(out string, t int) ([]int64, error) {
	fields := strings.Fields(out)
	if len(fields) != t {
		return nil, fmt.Errorf("expected %d numbers, got %d", t, len(fields))
	}
	ans := make([]int64, t)
	for i, f := range fields {
		v, err := strconv.ParseInt(f, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("invalid integer %q", f)
		}
		ans[i] = v
	}
	return ans, nil
}
