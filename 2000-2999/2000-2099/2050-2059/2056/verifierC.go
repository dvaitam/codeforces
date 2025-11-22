package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"math/big"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

const refSource = "2000-2999/2000-2099/2050-2059/2056/2056C.go"

func main() {
	if len(os.Args) != 2 {
		fail("usage: go run verifierC.go /path/to/candidate")
	}
	candidate := os.Args[1]

	inputData, err := io.ReadAll(os.Stdin)
	if err != nil {
		fail("failed to read input: %v", err)
	}

	tests, err := parseInput(inputData)
	if err != nil {
		fail("failed to parse input: %v", err)
	}

	refBin, err := buildReference()
	if err != nil {
		fail("failed to build reference: %v", err)
	}
	defer os.Remove(refBin)

	if _, err := runProgram(exec.Command(refBin), inputData); err != nil {
		fail("reference execution failed: %v", err)
	}

	candOut, err := runProgram(commandFor(candidate), inputData)
	if err != nil {
		fail("candidate execution failed: %v", err)
	}

	if err := validateOutputs(candOut, tests); err != nil {
		fail("%v", err)
	}

	fmt.Println("OK")
}

type testCase struct {
	n int
}

func parseInput(data []byte) ([]testCase, error) {
	reader := bufio.NewReader(bytes.NewReader(data))
	var t int
	if _, err := fmt.Fscan(reader, &t); err != nil {
		return nil, err
	}
	tests := make([]testCase, t)
	for i := 0; i < t; i++ {
		if _, err := fmt.Fscan(reader, &tests[i].n); err != nil {
			return nil, err
		}
	}
	return tests, nil
}

func validateOutputs(out string, tests []testCase) error {
	reader := bufio.NewReader(strings.NewReader(out))
	for idx, tc := range tests {
		arr := make([]int, tc.n)
		for i := 0; i < tc.n; i++ {
			if _, err := fmt.Fscan(reader, &arr[i]); err != nil {
				if err == io.EOF {
					return fmt.Errorf("test %d: expected %d integers, got %d", idx+1, tc.n, i)
				}
				return fmt.Errorf("test %d: failed to read integer: %v", idx+1, err)
			}
			if arr[i] < 1 || arr[i] > tc.n {
				return fmt.Errorf("test %d: value %d out of range 1..%d", idx+1, arr[i], tc.n)
			}
		}
		if err := checkPalindromic(arr, tc.n); err != nil {
			return fmt.Errorf("test %d: %v", idx+1, err)
		}
	}
	if extra, err := readToken(reader); err == nil {
		return fmt.Errorf("unexpected extra output token %q", extra)
	} else if err != io.EOF {
		return err
	}
	return nil
}

func checkPalindromic(a []int, n int) error {
	if len(a) != n {
		return fmt.Errorf("array length mismatch: expected %d, got %d", n, len(a))
	}
	// Compute longest palindromic subsequence length.
	lpsLen := longestPalSubseqLength(a)
	count := countPalSubseqLength(a, lpsLen)
	if count.Cmp(big.NewInt(int64(n))) <= 0 {
		return fmt.Errorf("g(a)=%s is not greater than n=%d", count.String(), n)
	}
	return nil
}

func longestPalSubseqLength(a []int) int {
	n := len(a)
	dp := make([][]int, n)
	for i := range dp {
		dp[i] = make([]int, n)
		dp[i][i] = 1
	}
	for length := 2; length <= n; length++ {
		for i := 0; i+length-1 < n; i++ {
			j := i + length - 1
			if a[i] == a[j] {
				if length == 2 {
					dp[i][j] = 2
				} else {
					dp[i][j] = dp[i+1][j-1] + 2
				}
			} else {
				if dp[i+1][j] > dp[i][j-1] {
					dp[i][j] = dp[i+1][j]
				} else {
					dp[i][j] = dp[i][j-1]
				}
			}
		}
	}
	return dp[0][n-1]
}

func countPalSubseqLength(a []int, target int) *big.Int {
	n := len(a)
	counts := make([][][]big.Int, n)
	for i := 0; i < n; i++ {
		counts[i] = make([][]big.Int, n)
		for j := i; j < n; j++ {
			counts[i][j] = make([]big.Int, target+1)
		}
	}

	// Base cases len = 1
	for i := 0; i < n; i++ {
		counts[i][i][0].SetInt64(1)
		if target >= 1 {
			counts[i][i][1].SetInt64(1)
		}
	}

	// Fill for lengths >= 2
	for length := 2; length <= n; length++ {
		for i := 0; i+length-1 < n; i++ {
			j := i + length - 1
			for l := 0; l <= target; l++ {
				var val big.Int
				val.Add(getCount(counts, i+1, j, l), getCount(counts, i, j-1, l))
				val.Sub(&val, getCount(counts, i+1, j-1, l))
				if a[i] == a[j] && l >= 2 {
					val.Add(&val, getCount(counts, i+1, j-1, l-2))
				}
				counts[i][j][l].Set(&val)
			}
		}
	}

	return &counts[0][n-1][target]
}

var (
	bigZero = big.NewInt(0)
	bigOne  = big.NewInt(1)
)

func getCount(counts [][][]big.Int, i, j, l int) *big.Int {
	if l < 0 {
		return bigZero
	}
	if i > j {
		if l == 0 {
			return bigOne
		}
		return bigZero
	}
	return &counts[i][j][l]
}

func buildReference() (string, error) {
	tmp, err := os.CreateTemp("", "2056C-ref-*")
	if err != nil {
		return "", err
	}
	tmp.Close()

	cmd := exec.Command("go", "build", "-o", tmp.Name(), filepath.Clean(refSource))
	var combined bytes.Buffer
	cmd.Stdout = &combined
	cmd.Stderr = &combined
	if err := cmd.Run(); err != nil {
		os.Remove(tmp.Name())
		return "", fmt.Errorf("%v\n%s", err, combined.String())
	}
	return tmp.Name(), nil
}

func runProgram(cmd *exec.Cmd, input []byte) (string, error) {
	cmd.Stdin = bytes.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return out.String(), err
}

func commandFor(path string) *exec.Cmd {
	switch strings.ToLower(filepath.Ext(path)) {
	case ".go":
		return exec.Command("go", "run", path)
	case ".py":
		return exec.Command("python3", path)
	default:
		return exec.Command(path)
	}
}

func readToken(r *bufio.Reader) (string, error) {
	var b strings.Builder
	for {
		ch, err := r.ReadByte()
		if err != nil {
			return "", err
		}
		if ch > ' ' {
			b.WriteByte(ch)
			break
		}
	}
	for {
		ch, err := r.ReadByte()
		if err != nil {
			if err == io.EOF {
				return b.String(), nil
			}
			return "", err
		}
		if ch <= ' ' {
			return b.String(), nil
		}
		b.WriteByte(ch)
	}
}

func fail(format string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, format+"\n", args...)
	os.Exit(1)
}
