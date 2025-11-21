package main

import (
	"bufio"
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
	refSource        = "1431I.go"
	tempOraclePrefix = "oracle-1431I-"
	randomTests      = 40
)

type testCase struct {
	n       int
	m       int
	q       int
	mat     []string
	queries []string
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierI.go /path/to/binary")
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
	tests = append(tests, randomTestsCases(rng, randomTests)...)

	for idx, tc := range tests {
		input := formatInput(tc)
		expOut, err := runProgram(oracle, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "oracle runtime error on test %d: %v\n", idx+1, err)
			fmt.Println("Input:")
			fmt.Print(input)
			os.Exit(1)
		}

		candOut, err := runProgram(candidate, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "candidate runtime error on test %d: %v\n", idx+1, err)
			fmt.Println("Input:")
			fmt.Print(input)
			os.Exit(1)
		}

		expVals := parseAnswers(expOut, tc.q)
		gotVals := parseAnswers(candOut, tc.q)
		if expVals == nil || gotVals == nil {
			fmt.Fprintf(os.Stderr, "failed to parse outputs on test %d\n", idx+1)
			fmt.Println("Input:")
			fmt.Print(input)
			fmt.Println("Expected:")
			fmt.Print(expOut)
			fmt.Println("Got:")
			fmt.Print(candOut)
			os.Exit(1)
		}
		for i := 0; i < tc.q; i++ {
			if expVals[i] != gotVals[i] {
				fmt.Fprintf(os.Stderr, "test %d query %d failed: expected %d got %d\n", idx+1, i+1, expVals[i], gotVals[i])
				fmt.Println("Input:")
				fmt.Print(input)
				fmt.Println("Expected:")
				fmt.Print(expOut)
				fmt.Println("Got:")
				fmt.Print(candOut)
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
	outPath := filepath.Join(tmpDir, "oracleI")
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

func parseAnswers(out string, q int) []int {
	fields := strings.Fields(out)
	if len(fields) < q {
		return nil
	}
	ans := make([]int, q)
	for i := 0; i < q; i++ {
		v, err := strconv.Atoi(fields[i])
		if err != nil {
			return nil
		}
		ans[i] = v
	}
	return ans
}

func formatInput(tc testCase) string {
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d %d %d\n", tc.n, tc.m, tc.q)
	for _, row := range tc.mat {
		sb.WriteString(row)
		sb.WriteByte('\n')
	}
	for _, qu := range tc.queries {
		sb.WriteString(qu)
		sb.WriteByte('\n')
	}
	return sb.String()
}

func deterministicTests() []testCase {
	return []testCase{
		{
			n:       2,
			m:       2,
			q:       2,
			mat:     []string{"aa", "bb"},
			queries: []string{"aa", "bb"},
		},
		{
			n:       3,
			m:       3,
			q:       2,
			mat:     []string{"abc", "def", "ghi"},
			queries: []string{"abc", "ghi"},
		},
		{
			n:       2,
			m:       3,
			q:       2,
			mat:     []string{"abc", "bca"},
			queries: []string{"aaa", "bbb"},
		},
	}
}

func randomTestsCases(rng *rand.Rand, count int) []testCase {
	tests := make([]testCase, 0, count)
	for t := 0; t < count; t++ {
		n := rng.Intn(5) + 2
		m := rng.Intn(5) + 2
		mat := make([]string, n)
		for i := 0; i < n; i++ {
			var sb strings.Builder
			for j := 0; j < m; j++ {
				sb.WriteByte(byte('a' + rng.Intn(3)))
			}
			mat[i] = sb.String()
		}
		q := rng.Intn(5) + 1
		queries := make([]string, q)
		for i := 0; i < q; i++ {
			var sb strings.Builder
			for j := 0; j < m; j++ {
				sb.WriteByte(byte('a' + rng.Intn(3)))
			}
			queries[i] = sb.String()
		}
		tests = append(tests, testCase{n: n, m: m, q: q, mat: mat, queries: queries})
	}
	return tests
}
