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
	refSource        = "1381A1.go"
	tempOraclePrefix = "oracle-1381A1-"
	randomTests      = 60
)

type testCase struct {
	n int
	a string
	b string
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierA1.go /path/to/binary")
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
		_, err := runProgram(oracle, input)
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

		ops, err := parseCandidate(candOut, tc.n)
		if err != nil {
			fmt.Fprintf(os.Stderr, "candidate output parse error on test %d: %v\n", idx+1, err)
			fmt.Println("Input:")
			fmt.Print(input)
			fmt.Println("Output:")
			fmt.Print(candOut)
			os.Exit(1)
		}
		if err := verifyOps(tc, ops); err != nil {
			fmt.Fprintf(os.Stderr, "test %d failed: %v\n", idx+1, err)
			fmt.Println("Input:")
			fmt.Print(input)
			fmt.Println("Output:")
			fmt.Print(candOut)
			os.Exit(1)
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
	outPath := filepath.Join(tmpDir, "oracleA1")
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
	return stdout.String(), nil
}

func parseCandidate(out string, n int) ([]int, error) {
	scanner := bufio.NewScanner(strings.NewReader(out))
	scanner.Split(bufio.ScanWords)
	if !scanner.Scan() {
		return nil, fmt.Errorf("empty output")
	}
	k, err := strconv.Atoi(scanner.Text())
	if err != nil {
		return nil, fmt.Errorf("failed to parse operation count: %v", err)
	}
	if k < 0 || k > 3*n {
		return nil, fmt.Errorf("operation count %d exceeds limit 3n", k)
	}
	ops := make([]int, 0, k)
	for i := 0; i < k; i++ {
		if !scanner.Scan() {
			return nil, fmt.Errorf("expected %d operations, got %d", k, len(ops))
		}
		val, err := strconv.Atoi(scanner.Text())
		if err != nil {
			return nil, fmt.Errorf("invalid operation %q", scanner.Text())
		}
		if val < 1 || val > n {
			return nil, fmt.Errorf("operation %d out of range [1,%d]", val, n)
		}
		ops = append(ops, val)
	}
	return ops, nil
}

func verifyOps(tc testCase, ops []int) error {
	a := []byte(tc.a)
	for _, p := range ops {
		for i := 0; i < p; i++ {
			if a[i] == '0' {
				a[i] = '1'
			} else {
				a[i] = '0'
			}
		}
		for i, j := 0, p-1; i < j; i, j = i+1, j-1 {
			a[i], a[j] = a[j], a[i]
		}
	}
	if string(a) != tc.b {
		return fmt.Errorf("operations did not transform a to b")
	}
	return nil
}

func formatInput(tc testCase) string {
	return fmt.Sprintf("%d\n%d\n%s\n%s\n", 1, tc.n, tc.a, tc.b)
}

func deterministicTests() []testCase {
	return []testCase{
		{n: 1, a: "0", b: "1"},
		{n: 1, a: "1", b: "1"},
		{n: 2, a: "01", b: "10"},
		{n: 3, a: "000", b: "111"},
		{n: 4, a: "0101", b: "0011"},
		{n: 5, a: "11111", b: "00000"},
	}
}

func randomTestsCases(rng *rand.Rand, count int) []testCase {
	tests := make([]testCase, 0, count)
	totalLen := 0
	for len(tests) < count {
		n := rng.Intn(20) + 1
		if totalLen+n > 1000 {
			break
		}
		totalLen += n
		var sb strings.Builder
		for i := 0; i < n; i++ {
			if rng.Intn(2) == 0 {
				sb.WriteByte('0')
			} else {
				sb.WriteByte('1')
			}
		}
		a := sb.String()
		sb.Reset()
		for i := 0; i < n; i++ {
			if rng.Intn(2) == 0 {
				sb.WriteByte('0')
			} else {
				sb.WriteByte('1')
			}
		}
		b := sb.String()
		tests = append(tests, testCase{n: n, a: a, b: b})
	}
	return tests
}
