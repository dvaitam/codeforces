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
	refSource   = "2000-2999/2000-2099/2070-2079/2075/2075C.go"
	totalLimitN = 200000
	totalLimitM = 200000
	maxN        = 200000
	maxM        = 200000
	maxAi       = 200000
	defaultTime = 20 * time.Second
)

type testCase struct {
	n int
	m int
	a []int
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierC.go /path/to/candidate")
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
			fmt.Fprintf(os.Stderr, "wrong answer on test %d\nexpected: %d\nfound: %d\ninput:\n%s", i+1, refAns[i], candAns[i], singleInput(tests[i]))
			os.Exit(1)
		}
	}

	fmt.Printf("All %d tests passed.\n", len(tests))
}

func buildReference() (string, func(), error) {
	dir, err := os.MkdirTemp("", "verifier-2075C-*")
	if err != nil {
		return "", nil, err
	}
	bin := filepath.Join(dir, "ref2075C")
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
	totalN, totalM := 0, 0
	for _, tc := range tests {
		totalN += tc.n
		totalM += tc.m
	}

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for totalN < totalLimitN && totalM < totalLimitM {
		n := rng.Intn(maxN) + 1
		m := rng.Intn(maxM) + 1
		if n < 2 {
			n = 2
		}
		if m < 2 {
			m = 2
		}
		if totalN+n > totalLimitN {
			n = totalLimitN - totalN
		}
		if totalM+m > totalLimitM {
			m = totalLimitM - totalM
		}
		if n <= 0 || m <= 0 {
			break
		}
		a := make([]int, m)
		for i := 0; i < m; i++ {
			a[i] = rng.Intn(minInt(maxAi, n)) + 1
		}
		tests = append(tests, testCase{n: n, m: m, a: a})
		totalN += n
		totalM += m
	}
	return tests
}

func deterministicTests() []testCase {
	return []testCase{
		{n: 2, m: 2, a: []int{1, 1}},
		{n: 3, m: 2, a: []int{3, 3}},
		{n: 4, m: 3, a: []int{1, 2, 4}},
		{n: 5, m: 2, a: []int{2, 5}},
		{n: 6, m: 4, a: []int{1, 3, 4, 6}},
		{n: 8, m: 5, a: []int{8, 1, 4, 4, 7}},
		{n: 10, m: 3, a: []int{10, 10, 10}},
		{n: 12, m: 6, a: []int{2, 2, 2, 2, 2, 2}},
	}
}

func buildInput(tests []testCase) string {
	var sb strings.Builder
	sb.Grow(len(tests) * 32)
	fmt.Fprintf(&sb, "%d\n", len(tests))
	for _, tc := range tests {
		fmt.Fprintf(&sb, "%d %d\n", tc.n, tc.m)
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

func singleInput(tc testCase) string {
	var sb strings.Builder
	fmt.Fprintf(&sb, "1\n%d %d\n", tc.n, tc.m)
	for i, v := range tc.a {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(strconv.Itoa(v))
	}
	sb.WriteByte('\n')
	return sb.String()
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

func parseOutputs(out string, t int) ([]int64, error) {
	fields := strings.Fields(out)
	if len(fields) != t {
		return nil, fmt.Errorf("expected %d numbers, got %d", t, len(fields))
	}
	res := make([]int64, t)
	for i, f := range fields {
		v, err := strconv.ParseInt(f, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("invalid integer %q", f)
		}
		res[i] = v
	}
	return res, nil
}

func minInt(a, b int) int {
	if a < b {
		return a
	}
	return b
}
