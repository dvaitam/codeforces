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
	refSource   = "./2101B.go"
	totalNLimit = 200000
	defaultTime = 20 * time.Second
)

type testCase struct {
	n int
	a []int
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
		if len(exp) != len(got) {
			fmt.Fprintf(os.Stderr, "test %d: expected %d numbers, got %d\ninput:\n%s", i+1, len(exp), len(got), singleInput(tests[i]))
			os.Exit(1)
		}
		for j := range exp {
			if exp[j] != got[j] {
				fmt.Fprintf(os.Stderr, "test %d position %d mismatch: expected %d got %d\ninput:\n%s", i+1, j+1, exp[j], got[j], singleInput(tests[i]))
				os.Exit(1)
			}
		}
	}

	fmt.Printf("All %d tests passed.\n", len(tests))
}

func buildReference() (string, func(), error) {
	dir, err := os.MkdirTemp("", "verifier-2101B-*")
	if err != nil {
		return "", nil, err
	}
	bin := filepath.Join(dir, "ref2101B")
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
		n := rng.Intn(5000) + 4
		if total+n > totalNLimit {
			n = totalNLimit - total
		}
		if n < 4 {
			break
		}
		p := randPermutation(rng, n)
		tests = append(tests, testCase{n: n, a: p})
		total += n
	}
	return tests
}

func deterministicTests() []testCase {
	return []testCase{
		{n: 4, a: []int{1, 2, 3, 4}},
		{n: 4, a: []int{4, 3, 2, 1}},
		{n: 5, a: []int{2, 1, 4, 3, 5}},
		{n: 6, a: []int{6, 5, 4, 3, 2, 1}},
		{n: 7, a: []int{3, 1, 4, 2, 7, 6, 5}},
		{n: 8, a: []int{8, 1, 2, 3, 4, 5, 6, 7}},
	}
}

func randPermutation(rng *rand.Rand, n int) []int {
	p := make([]int, n)
	for i := 0; i < n; i++ {
		p[i] = i + 1
	}
	rng.Shuffle(n, func(i, j int) { p[i], p[j] = p[j], p[i] })
	return p
}

func buildInput(tests []testCase) string {
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d\n", len(tests))
	for _, tc := range tests {
		fmt.Fprintf(&sb, "%d\n", tc.n)
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
	fmt.Fprintf(&sb, "1\n%d\n", tc.n)
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

func parseOutputs(out string, tests []testCase) ([][]int, error) {
	fields := strings.Fields(out)
	idx := 0
	res := make([][]int, len(tests))
	for i, tc := range tests {
		if idx+tc.n > len(fields) {
			return nil, fmt.Errorf("not enough outputs for test %d", i+1)
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
