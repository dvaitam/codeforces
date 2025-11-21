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
	refSource        = "1662E.go"
	tempOraclePrefix = "oracle-1662E-"
	randomTests      = 60
	maxN             = 8
)

type testCase struct {
	n int
	p []int
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
		gotOut, err := runProgram(candidate, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "candidate runtime error on test %d: %v\n", idx+1, err)
			fmt.Println("Input:")
			fmt.Print(input)
			os.Exit(1)
		}
		exp := parseAnswer(expOut)
		got := parseAnswer(gotOut)
		if exp == nil || got == nil || len(exp) != len(got) {
			fmt.Fprintf(os.Stderr, "output parse error on test %d\n", idx+1)
			fmt.Println("Input:")
			fmt.Print(input)
			fmt.Println("Expected:")
			fmt.Print(expOut)
			fmt.Println("Got:")
			fmt.Print(gotOut)
			os.Exit(1)
		}
		for i := range exp {
			if exp[i] != got[i] {
				fmt.Fprintf(os.Stderr, "test %d case %d failed: expected %d got %d\n", idx+1, i+1, exp[i], got[i])
				fmt.Println("Input:")
				fmt.Print(input)
				fmt.Println("Expected:")
				fmt.Print(expOut)
				fmt.Println("Got:")
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
	fmt.Fprintf(&sb, "1\n%d\n", tc.n)
	for i, v := range tc.p {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(strconv.Itoa(v))
	}
	sb.WriteByte('\n')
	return sb.String()
}

func parseAnswer(out string) []int {
	fields := strings.Fields(out)
	ans := make([]int, 0, len(fields))
	for _, f := range fields {
		v, err := strconv.Atoi(f)
		if err != nil {
			return nil
		}
		ans = append(ans, v)
	}
	return ans
}

func deterministicTests() []testCase {
	return []testCase{
		{n: 3, p: []int{1, 2, 3}},
		{n: 3, p: []int{3, 2, 1}},
		{n: 4, p: []int{2, 3, 4, 1}},
		{n: 5, p: []int{5, 4, 3, 2, 1}},
	}
}

func randomTestsCases(rng *rand.Rand, count int) []testCase {
	tests := make([]testCase, 0, count)
	for len(tests) < count {
		n := rng.Intn(maxN-2) + 3
		p := rand.Perm(n)
		for i := range p {
			p[i]++
		}
		tests = append(tests, testCase{n: n, p: p})
	}
	return tests
}
