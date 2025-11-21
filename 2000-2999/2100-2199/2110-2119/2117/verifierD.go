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
	refSource   = "2000-2999/2100-2199/2110-2119/2117/2117D.go"
	totalNLimit = 200000
	defaultTime = 20 * time.Second
	maxA        = 1_000_000_000
)

type testCase struct {
	n int
	a []int
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
		if refAns[i] != candAns[i] {
			fmt.Fprintf(os.Stderr, "wrong answer on test %d: expected %s got %s\ninput:\n%s", i+1, refAns[i], candAns[i], singleInput(tests[i]))
			os.Exit(1)
		}
	}

	fmt.Printf("All %d tests passed.\n", len(tests))
}

func buildReference() (string, func(), error) {
	dir, err := os.MkdirTemp("", "verifier-2117D-*")
	if err != nil {
		return "", nil, err
	}
	bin := filepath.Join(dir, "ref2117D")
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
		n := rng.Intn(5000) + 2
		if total+n > totalNLimit {
			n = totalNLimit - total
		}
		if n < 2 {
			break
		}
		a := make([]int, n)
		for i := 0; i < n; i++ {
			a[i] = rng.Intn(maxA) + 1
		}
		tests = append(tests, testCase{n: n, a: a})
		total += n
	}
	return tests
}

func deterministicTests() []testCase {
	return []testCase{
		{n: 2, a: []int{1, 1}},
		{n: 2, a: []int{1, 2}},
		{n: 3, a: []int{1, 2, 3}},
		{n: 4, a: []int{4, 4, 4, 4}},
		{n: 5, a: []int{10, 9, 8, 7, 6}},
		{n: 6, a: []int{1, 100, 1, 100, 1, 100}},
		{n: 7, a: []int{7, 6, 5, 4, 3, 2, 1}},
	}
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
			sb.WriteString(fmtInt(v))
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
		sb.WriteString(fmtInt(v))
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

func parseOutputs(out string, t int) ([]string, error) {
	fields := strings.Fields(out)
	if len(fields) != t {
		return nil, fmt.Errorf("expected %d tokens, got %d", t, len(fields))
	}
	res := make([]string, t)
	for i, f := range fields {
		res[i] = normalizeAnswer(f)
	}
	return res, nil
}

func normalizeAnswer(s string) string {
	s = strings.ToUpper(s)
	if s == "YES" || s == "Y" {
		return "YES"
	}
	return "NO"
}

func fmtInt(v int) string {
	return strconv.FormatInt(int64(v), 10)
}
