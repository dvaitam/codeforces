package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

const (
	refSource        = "1221E.go"
	tempOraclePrefix = "oracle-1221E-"
	randomTestCount  = 50
	maxStringLen     = 200
)

type query struct {
	a int
	b int
	s string
}

type testCase struct {
	queries []query
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	candidate := os.Args[1]

	oracle, cleanup, err := buildOracle()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to build oracle: %v\n", err)
		os.Exit(1)
	}
	defer cleanup()

	tests := deterministicTests()
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	tests = append(tests, randomTests(rng, randomTestCount)...)

	for idx, tc := range tests {
		input := formatInput(tc)
		expOut, err := runProgram(oracle, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "oracle runtime error on test %d: %v\n", idx+1, err)
			fmt.Println("Input:")
			fmt.Print(input)
			os.Exit(1)
		}
		gotOut, err := runProgram(candidate, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "candidate runtime error on test %d: %v\n", idx+1, err)
			fmt.Println("Input:")
			fmt.Print(input)
			os.Exit(1)
		}
		exp := normalizeOutputs(expOut, len(tc.queries))
		got := normalizeOutputs(gotOut, len(tc.queries))
		if len(exp) != len(tc.queries) {
			fmt.Fprintf(os.Stderr, "oracle produced %d answers, expected %d\n", len(exp), len(tc.queries))
			os.Exit(1)
		}
		if len(got) != len(tc.queries) {
			fmt.Fprintf(os.Stderr, "candidate produced %d answers, expected %d\n", len(got), len(tc.queries))
			fmt.Println("Input:")
			fmt.Print(input)
			fmt.Println("Candidate output:")
			fmt.Print(gotOut)
			os.Exit(1)
		}
		for i := range exp {
			if exp[i] != got[i] {
				fmt.Fprintf(os.Stderr, "test %d query %d failed: expected %s got %s\n", idx+1, i+1, exp[i], got[i])
				fmt.Println("Input:")
				fmt.Print(input)
				fmt.Println("Candidate output:")
				fmt.Print(gotOut)
				os.Exit(1)
			}
		}
	}

	fmt.Printf("All %d tests passed\n", len(tests))
}

func buildOracle() (string, func(), error) {
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

func runProgram(path, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(path, ".go") {
		cmd = exec.Command("go", "run", path)
	} else {
		cmd = exec.Command(path)
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

func formatInput(tc testCase) string {
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d\n", len(tc.queries))
	for _, q := range tc.queries {
		fmt.Fprintf(&sb, "%d %d\n", q.a, q.b)
		sb.WriteString(q.s)
		sb.WriteByte('\n')
	}
	return sb.String()
}

func deterministicTests() []testCase {
	return []testCase{
		{queries: []query{
			{a: 3, b: 2, s: "..X.."},
		}},
		{queries: []query{
			{a: 5, b: 1, s: "....."},
		}},
		{queries: []query{
			{a: 4, b: 2, s: "XXXX"},
		}},
		{queries: []query{
			{a: 6, b: 2, s: ".X...XX..X."},
			{a: 7, b: 3, s: "...XX..XX...."},
			{a: 10, b: 1, s: "............X"},
		}},
	}
}

func randomTests(rng *rand.Rand, count int) []testCase {
	tests := make([]testCase, 0, count)
	for t := 0; t < count; t++ {
		qCount := rng.Intn(30) + 1
		tc := testCase{queries: make([]query, qCount)}
		for i := 0; i < qCount; i++ {
			b := rng.Intn(50) + 1
			a := b + rng.Intn(50) + 1
			length := rng.Intn(maxStringLen) + 1
			var sb strings.Builder
			for j := 0; j < length; j++ {
				if rng.Intn(2) == 0 {
					sb.WriteByte('.')
				} else {
					sb.WriteByte('X')
				}
			}
			tc.queries[i] = query{a: a, b: b, s: sb.String()}
		}
		tests = append(tests, tc)
	}
	return tests
}

func normalizeOutputs(out string, expected int) []string {
	fields := strings.Fields(out)
	result := make([]string, 0, len(fields))
	for _, f := range fields {
		res := strings.ToUpper(f)
		if res == "YES" || res == "NO" {
			result = append(result, res)
		}
	}
	if len(result) > expected {
		result = result[:expected]
	}
	return result
}
