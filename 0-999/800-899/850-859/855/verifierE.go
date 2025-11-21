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

const (
	refSource        = "855E.go"
	tempOraclePrefix = "oracle-855E-"
	randomTestsCount = 80
)

type query struct {
	base int
	l    int64
	r    int64
}

type testCase struct {
	name    string
	queries []query
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	candidate := os.Args[1]

	oraclePath, cleanup, err := buildReference()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to build reference: %v\n", err)
		os.Exit(1)
	}
	defer cleanup()

	tests := deterministicTests()
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	tests = append(tests, randomTests(randomTestsCount, rng)...)

	for idx, tc := range tests {
		input := formatInput(tc)

		oracleOut, err := runProgram(oraclePath, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "oracle runtime error on test %d (%s): %v\n", idx+1, tc.name, err)
			fmt.Println("Input:")
			fmt.Print(input)
			os.Exit(1)
		}
		expected, err := parseOutput(oracleOut, len(tc.queries))
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to parse oracle output on test %d (%s): %v\noutput:\n%s", idx+1, tc.name, err, oracleOut)
			os.Exit(1)
		}

		candOut, err := runProgram(candidate, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "candidate runtime error on test %d (%s): %v\n", idx+1, tc.name, err)
			fmt.Println("Input:")
			fmt.Print(input)
			os.Exit(1)
		}
		got, err := parseOutput(candOut, len(tc.queries))
		if err != nil {
			fmt.Fprintf(os.Stderr, "candidate output parse error on test %d (%s): %v\noutput:\n%s", idx+1, tc.name, err, candOut)
			fmt.Println("Input:")
			fmt.Print(input)
			os.Exit(1)
		}

		for i := range expected {
			if expected[i] != got[i] {
				fmt.Fprintf(os.Stderr, "test %d (%s) failed at query %d: expected %d got %d\n", idx+1, tc.name, i+1, expected[i], got[i])
				fmt.Println("Input:")
				fmt.Print(input)
				fmt.Println("Candidate output:")
				fmt.Print(candOut)
				os.Exit(1)
			}
		}
	}

	fmt.Printf("All %d tests passed\n", len(tests))
}

func buildReference() (string, func(), error) {
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		return "", nil, fmt.Errorf("failed to determine verifier path")
	}
	dir := filepath.Dir(file)
	tmpDir, err := os.MkdirTemp("", tempOraclePrefix)
	if err != nil {
		return "", nil, err
	}
	outPath := filepath.Join(tmpDir, "oracleE")
	cmd := exec.Command("go", "build", "-o", outPath, refSource)
	cmd.Dir = dir
	if out, err := cmd.CombinedOutput(); err != nil {
		os.RemoveAll(tmpDir)
		return "", nil, fmt.Errorf("build oracle failed: %v\n%s", err, out)
	}
	cleanup := func() {
		os.RemoveAll(tmpDir)
	}
	return outPath, cleanup, nil
}

func runProgram(bin, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, stderr.String())
	}
	return strings.TrimSpace(stdout.String()), nil
}

func formatInput(tc testCase) string {
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d\n", len(tc.queries))
	for _, q := range tc.queries {
		fmt.Fprintf(&sb, "%d %d %d\n", q.base, q.l, q.r)
	}
	return sb.String()
}

func parseOutput(out string, q int) ([]int64, error) {
	lines := strings.Fields(out)
	if len(lines) != q {
		return nil, fmt.Errorf("expected %d answers, got %d", q, len(lines))
	}
	res := make([]int64, q)
	for i, line := range lines {
		val, err := strconv.ParseInt(line, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("invalid integer %q", line)
		}
		res[i] = val
	}
	return res, nil
}

func deterministicTests() []testCase {
	return []testCase{
		{
			name: "single_small",
			queries: []query{
				{base: 2, l: 1, r: 10},
			},
		},
		{
			name: "multi_base",
			queries: []query{
				{base: 2, l: 4, r: 9},
				{base: 3, l: 1, r: 50},
				{base: 10, l: 100, r: 1000},
			},
		},
		{
			name: "edge_large_range",
			queries: []query{
				{base: 10, l: 1, r: 1_000_000_000_000_000_000},
			},
		},
		{
			name: "dense_queries",
			queries: []query{
				{base: 2, l: 1, r: 64},
				{base: 5, l: 1, r: 5000},
				{base: 7, l: 100, r: 10000},
				{base: 9, l: 12345, r: 54321},
				{base: 10, l: 99999, r: 1000000},
			},
		},
	}
}

func randomTests(count int, rng *rand.Rand) []testCase {
	tests := make([]testCase, 0, count)
	for i := 0; i < count; i++ {
		queryCount := rng.Intn(20) + 1
		queries := make([]query, queryCount)
		for j := 0; j < queryCount; j++ {
			base := rng.Intn(9) + 2
			l := rng.Int63n(1_000_000_000) + 1
			r := l + rng.Int63n(1_000_000_000)
			queries[j] = query{base: base, l: l, r: r}
		}
		tests = append(tests, testCase{
			name:    fmt.Sprintf("random_%d", i+1),
			queries: queries,
		})
	}
	return tests
}
