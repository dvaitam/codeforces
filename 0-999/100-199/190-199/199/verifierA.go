package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
)

const (
	refSource        = "199A.go"
	tempOraclePrefix = "oracle-199A-"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	candidate := os.Args[1]

	oraclePath, cleanup, err := buildReference()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to build reference: %v\n", err)
		os.Exit(1)
	}
	defer cleanup()

	tests := fibonacciNumbers()
	fibSet := make(map[int64]struct{}, len(tests))
	for _, v := range tests {
		fibSet[int64(v)] = struct{}{}
	}

	for idx, n := range tests {
		if err := verifyCase(idx+1, n, candidate, oraclePath, fibSet); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
	}

	fmt.Printf("All %d tests passed\n", len(tests))
}

func verifyCase(testID, n int, candidate, oracle string, fibSet map[int64]struct{}) error {
	input := fmt.Sprintf("%d\n", n)

	if _, err := runProgram(oracle, input); err != nil {
		return fmt.Errorf("oracle runtime error on test %d (n=%d): %v", testID, n, err)
	}

	out, err := runProgram(candidate, input)
	if err != nil {
		return fmt.Errorf("candidate runtime error on test %d (n=%d): %v", testID, n, err)
	}

	values, err := parseTriple(out)
	if err != nil {
		return fmt.Errorf("candidate output parse error on test %d (n=%d): %v\noutput:\n%s", testID, n, err, out)
	}

	var sum int64
	for i, v := range values {
		if v < 0 {
			return fmt.Errorf("test %d failed: value %d is negative", testID, v)
		}
		if _, ok := fibSet[v]; !ok {
			return fmt.Errorf("test %d failed: value %d at position %d is not a Fibonacci number", testID, v, i+1)
		}
		sum += v
	}

	if sum != int64(n) {
		return fmt.Errorf("test %d failed: expected sum %d, got %d", testID, n, sum)
	}

	return nil
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

	outPath := filepath.Join(tmpDir, "oracleA")
	cmd := exec.Command("go", "build", "-o", outPath, refSource)
	cmd.Dir = dir
	if out, err := cmd.CombinedOutput(); err != nil {
		os.RemoveAll(tmpDir)
		return "", nil, fmt.Errorf("build oracle failed: %v\n%s", err, string(out))
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

func parseTriple(out string) ([3]int64, error) {
	tokens := strings.Fields(out)
	if len(tokens) != 3 {
		return [3]int64{}, fmt.Errorf("expected 3 integers, got %d tokens", len(tokens))
	}
	var res [3]int64
	for i := 0; i < 3; i++ {
		val, err := strconv.ParseInt(tokens[i], 10, 64)
		if err != nil {
			return [3]int64{}, fmt.Errorf("token %q is not an integer", tokens[i])
		}
		res[i] = val
	}
	return res, nil
}

func fibonacciNumbers() []int {
	fibs := []int{0, 1}
	for {
		next := fibs[len(fibs)-1] + fibs[len(fibs)-2]
		if next >= 1000000000 {
			break
		}
		fibs = append(fibs, next)
	}
	return fibs
}
