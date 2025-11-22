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

type query struct {
	t int
	v int
}

type testCase struct {
	queries []query
}

func callerFile() (string, bool) {
	_, file, _, ok := runtime.Caller(0)
	return file, ok
}

func buildOracle() (string, func(), error) {
	_, file, ok := callerFile()
	if !ok {
		return "", nil, fmt.Errorf("cannot determine verifier path")
	}
	dir := filepath.Dir(file)
	tmpDir, err := os.MkdirTemp("", "oracle-2025G-")
	if err != nil {
		return "", nil, err
	}
	outPath := filepath.Join(tmpDir, "oracleG")
	cmd := exec.Command("go", "build", "-o", outPath, "2025G.go")
	cmd.Dir = dir
	if out, err := cmd.CombinedOutput(); err != nil {
		os.RemoveAll(tmpDir)
		return "", nil, fmt.Errorf("failed to build oracle: %v\n%s", err, out)
	}
	cleanup := func() {
		os.RemoveAll(tmpDir)
	}
	return outPath, cleanup, nil
}

func runBinary(bin, input string) (string, error) {
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
	return stdout.String(), nil
}

func buildInput(tc testCase) string {
	var sb strings.Builder
	sb.Grow(len(tc.queries) * 24)
	sb.WriteString(strconv.Itoa(len(tc.queries)))
	sb.WriteByte('\n')
	for _, q := range tc.queries {
		sb.WriteString(strconv.Itoa(q.t))
		sb.WriteByte(' ')
		sb.WriteString(strconv.Itoa(q.v))
		sb.WriteByte('\n')
	}
	return sb.String()
}

func parseAnswers(out string, n int) ([]int64, error) {
	fields := strings.Fields(out)
	if len(fields) != n {
		return nil, fmt.Errorf("expected %d answers, got %d", n, len(fields))
	}
	ans := make([]int64, n)
	for i, f := range fields {
		val, err := strconv.ParseInt(f, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("invalid integer %q: %v", f, err)
		}
		ans[i] = val
	}
	return ans, nil
}

func compareAnswers(exp, got []int64) error {
	for i := range exp {
		if exp[i] != got[i] {
			return fmt.Errorf("at position %d expected %d, got %d", i+1, exp[i], got[i])
		}
	}
	return nil
}

func deterministicTests() []testCase {
	return []testCase{
		{
			queries: []query{
				{t: 2, v: 5},
				{t: 1, v: 4},
				{t: 1, v: 10},
			},
		},
		{
			queries: []query{
				{t: 1, v: 1},
				{t: 1, v: 1},
				{t: 2, v: 1},
				{t: 2, v: 2},
				{t: 1, v: 5},
				{t: 2, v: 3},
			},
		},
	}
}

func randomTests() []testCase {
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	tests := make([]testCase, 0, 120)
	for len(tests) < cap(tests) {
		q := rng.Intn(90) + 10 // [10,99]
		if rng.Intn(5) == 0 {
			q = rng.Intn(3000) + 50 // occasional bigger cases
		}
		queries := make([]query, q)
		for i := 0; i < q; i++ {
			t := 1
			if rng.Intn(2) == 0 {
				t = 2
			}
			val := rng.Intn(1_000_000_000) + 1
			if rng.Intn(6) == 0 {
				val = rng.Intn(20) + 1
			}
			queries[i] = query{t: t, v: val}
		}
		tests = append(tests, testCase{queries: queries})
	}
	return tests
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierG.go /path/to/binary")
		os.Exit(1)
	}
	target := os.Args[1]

	oracle, cleanup, err := buildOracle()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
	defer cleanup()

	tests := append(deterministicTests(), randomTests()...)
	for idx, tc := range tests {
		input := buildInput(tc)
		expOut, err := runBinary(oracle, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "oracle failed on test %d: %v\ninput:\n%s", idx+1, err, input)
			os.Exit(1)
		}
		gotOut, err := runBinary(target, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "target runtime error on test %d: %v\ninput:\n%s", idx+1, err, input)
			os.Exit(1)
		}
		expAns, err := parseAnswers(expOut, len(tc.queries))
		if err != nil {
			fmt.Fprintf(os.Stderr, "oracle output invalid on test %d: %v\n%s", idx+1, err, expOut)
			os.Exit(1)
		}
		gotAns, err := parseAnswers(gotOut, len(tc.queries))
		if err != nil {
			fmt.Fprintf(os.Stderr, "target output invalid on test %d: %v\n%s", idx+1, err, gotOut)
			os.Exit(1)
		}
		if err := compareAnswers(expAns, gotAns); err != nil {
			fmt.Fprintf(os.Stderr, "test %d mismatch: %v\ninput:\n%s", idx+1, err, input)
			os.Exit(1)
		}
	}

	fmt.Println("All tests passed.")
}
