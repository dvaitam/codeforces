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
	refSource   = "2000-2999/2100-2199/2120-2129/2126/2126D.go"
	totalNLimit = 100000
	defaultTime = 20 * time.Second
	maxVal      = 1_000_000_000
)

type testCase struct {
	n    int
	k    int
	l    []int
	r    []int
	real []int
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierD.go /path/to/candidate")
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

	for i := range tests {
		if candAns[i] != refAns[i] {
			fmt.Fprintf(os.Stderr, "wrong answer on test %d: expected %d got %d\ninput:\n%s", i+1, refAns[i], candAns[i], singleInput(tests[i]))
			os.Exit(1)
		}
	}

	fmt.Printf("All %d tests passed.\n", len(tests))
}

func buildReference() (string, func(), error) {
	dir, err := os.MkdirTemp("", "verifier-2126D-*")
	if err != nil {
		return "", nil, err
	}
	bin := filepath.Join(dir, "ref2126D")
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
		k := rng.Intn(maxVal + 1)
		l := make([]int, n)
		r := make([]int, n)
		real := make([]int, n)
		for i := 0; i < n; i++ {
			x := rng.Intn(maxVal + 1)
			y := rng.Intn(maxVal + 1)
			if x > y {
				x, y = y, x
			}
			z := x + rng.Intn(y-x+1)
			l[i], r[i], real[i] = x, y, z
		}
		tests = append(tests, testCase{n: n, k: k, l: l, r: r, real: real})
		total += n
	}
	return tests
}

func deterministicTests() []testCase {
	return []testCase{
		{
			n: 3, k: 1,
			l:    []int{1, 1, 1},
			r:    []int{1, 2, 10},
			real: []int{1, 2, 10},
		},
		{
			n: 1, k: 0,
			l:    []int{0},
			r:    []int{0},
			real: []int{0},
		},
		{
			n: 2, k: 5,
			l:    []int{0, 6},
			r:    []int{5, 10},
			real: []int{5, 6},
		},
		{
			n: 4, k: 3,
			l:    []int{1, 0, 2, 3},
			r:    []int{5, 3, 6, 10},
			real: []int{2, 3, 4, 5},
		},
		{
			n: 5, k: 4,
			l:    []int{0, 2, 4, 6, 8},
			r:    []int{1, 3, 5, 7, 9},
			real: []int{1, 2, 3, 4, 9},
		},
	}
}

func buildInput(tests []testCase) string {
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d\n", len(tests))
	for _, tc := range tests {
		fmt.Fprintf(&sb, "%d %d\n", tc.n, tc.k)
		for i := 0; i < tc.n; i++ {
			fmt.Fprintf(&sb, "%d %d %d\n", tc.l[i], tc.r[i], tc.real[i])
		}
	}
	return sb.String()
}

func singleInput(tc testCase) string {
	var sb strings.Builder
	fmt.Fprintf(&sb, "1\n%d %d\n", tc.n, tc.k)
	for i := 0; i < tc.n; i++ {
		fmt.Fprintf(&sb, "%d %d %d\n", tc.l[i], tc.r[i], tc.real[i])
	}
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
