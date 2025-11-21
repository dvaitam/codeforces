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
	refSource    = "2000-2999/2000-2099/2090-2099/2094/2094G.go"
	totalQLimit  = 200000
	defaultTime  = 20 * time.Second
	maxAppendVal = 1_000_000
)

type op struct {
	t int
	k int
}

type testCase struct {
	q   int
	ops []op
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierG.go /path/to/candidate")
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
	refAns, err := parseOutputs(refOut, tests)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to parse reference output: %v\noutput:\n%s\n", err, refOut)
		os.Exit(1)
	}

	candOut, err := runProgram(candidate, input)
	if err != nil {
		fmt.Fprintf(os.Stderr, "candidate runtime error: %v\n", err)
		os.Exit(1)
	}
	candAns, err := parseOutputs(candOut, tests)
	if err != nil {
		fmt.Fprintf(os.Stderr, "invalid candidate output: %v\noutput:\n%s\n", err, candOut)
		os.Exit(1)
	}

	for i := range tests {
		exp := refAns[i]
		got := candAns[i]
		for j := 0; j < len(exp) && j < len(got); j++ {
			if exp[j] != got[j] {
				fmt.Fprintf(os.Stderr, "test %d operation %d mismatch: expected %d got %d\ninput:\n%s", i+1, j+1, exp[j], got[j], singleInput(tests[i]))
				os.Exit(1)
			}
		}
		if len(exp) != len(got) {
			fmt.Fprintf(os.Stderr, "test %d length mismatch: expected %d outputs got %d\ninput:\n%s", i+1, len(exp), len(got), singleInput(tests[i]))
			os.Exit(1)
		}
	}

	fmt.Printf("All %d tests passed.\n", len(tests))
}

func buildReference() (string, func(), error) {
	dir, err := os.MkdirTemp("", "verifier-2094G-*")
	if err != nil {
		return "", nil, err
	}
	bin := filepath.Join(dir, "ref2094G")
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
	usedQ := 0
	for _, tc := range tests {
		usedQ += tc.q
	}
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for usedQ < totalQLimit {
		remain := totalQLimit - usedQ
		q := rng.Intn(40000) + 1
		if q > remain {
			q = remain
		}
		ops := make([]op, 0, q)
		// first operation must be append
		ops = append(ops, op{t: 3, k: rng.Intn(maxAppendVal) + 1})
		currLen := 1
		for len(ops) < q {
			t := rng.Intn(3) + 1
			if t == 3 || currLen == 0 {
				k := rng.Intn(maxAppendVal) + 1
				ops = append(ops, op{t: 3, k: k})
				currLen++
				continue
			}
			ops = append(ops, op{t: t})
			// length unchanged for t=1,2
		}
		tests = append(tests, testCase{q: q, ops: ops})
		usedQ += q
	}
	return tests
}

func deterministicTests() []testCase {
	return []testCase{
		{
			q: 1,
			ops: []op{
				{t: 3, k: 5},
			},
		},
		{
			q: 3,
			ops: []op{
				{t: 3, k: 1},
				{t: 3, k: 2},
				{t: 1},
			},
		},
		{
			q: 5,
			ops: []op{
				{t: 3, k: 1},
				{t: 3, k: 2},
				{t: 3, k: 3},
				{t: 2},
				{t: 1},
			},
		},
		{
			q: 6,
			ops: []op{
				{t: 3, k: 10},
				{t: 1},
				{t: 2},
				{t: 3, k: 7},
				{t: 1},
				{t: 2},
			},
		},
	}
}

func buildInput(tests []testCase) string {
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d\n", len(tests))
	for _, tc := range tests {
		fmt.Fprintf(&sb, "%d\n", tc.q)
		for _, op := range tc.ops {
			if op.t == 3 {
				fmt.Fprintf(&sb, "3 %d\n", op.k)
			} else {
				fmt.Fprintf(&sb, "%d\n", op.t)
			}
		}
	}
	return sb.String()
}

func singleInput(tc testCase) string {
	var sb strings.Builder
	fmt.Fprintf(&sb, "1\n%d\n", tc.q)
	for _, op := range tc.ops {
		if op.t == 3 {
			fmt.Fprintf(&sb, "3 %d\n", op.k)
		} else {
			fmt.Fprintf(&sb, "%d\n", op.t)
		}
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

func parseOutputs(out string, tests []testCase) ([][]int64, error) {
	fields := strings.Fields(out)
	idx := 0
	res := make([][]int64, len(tests))
	for i, tc := range tests {
		if idx+tc.q > len(fields) {
			return nil, fmt.Errorf("not enough outputs for test %d", i+1)
		}
		cur := make([]int64, tc.q)
		for j := 0; j < tc.q; j++ {
			v, err := strconv.ParseInt(fields[idx+j], 10, 64)
			if err != nil {
				return nil, fmt.Errorf("invalid integer %q", fields[idx+j])
			}
			cur[j] = v
		}
		res[i] = cur
		idx += tc.q
	}
	if idx != len(fields) {
		return nil, fmt.Errorf("extra tokens in output")
	}
	return res, nil
}
