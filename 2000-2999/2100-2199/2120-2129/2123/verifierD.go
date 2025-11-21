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
	refSource   = "2000-2999/2100-2199/2120-2129/2123/2123D.go"
	totalNLimit = 200000
	defaultTime = 15 * time.Second
)

type testCase struct {
	n int
	k int
	s string
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
	dir, err := os.MkdirTemp("", "verifier-2123D-*")
	if err != nil {
		return "", nil, err
	}
	bin := filepath.Join(dir, "ref2123D")
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
		k := rng.Intn(n-1) + 1
		var sb strings.Builder
		for i := 0; i < n; i++ {
			if rng.Intn(2) == 0 {
				sb.WriteByte('0')
			} else {
				sb.WriteByte('1')
			}
		}
		tests = append(tests, testCase{n: n, k: k, s: sb.String()})
		total += n
	}
	return tests
}

func deterministicTests() []testCase {
	return []testCase{
		{n: 2, k: 1, s: "00"},
		{n: 2, k: 1, s: "11"},
		{n: 3, k: 1, s: "101"},
		{n: 4, k: 2, s: "1111"},
		{n: 5, k: 2, s: "10101"},
		{n: 6, k: 5, s: "110110"},
		{n: 7, k: 3, s: "0010011"},
	}
}

func buildInput(tests []testCase) string {
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d\n", len(tests))
	for _, tc := range tests {
		fmt.Fprintf(&sb, "%d %d\n", tc.n, tc.k)
		sb.WriteString(tc.s)
		sb.WriteByte('\n')
	}
	return sb.String()
}

func singleInput(tc testCase) string {
	return fmt.Sprintf("1\n%d %d\n%s\n", tc.n, tc.k, tc.s)
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
		f = strings.ToLower(f)
		if f == "alice" {
			res[i] = "alice"
		} else {
			res[i] = "bob"
		}
	}
	return res, nil
}
