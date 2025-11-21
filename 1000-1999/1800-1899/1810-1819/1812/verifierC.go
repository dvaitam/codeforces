package main

import (
	"bytes"
	"fmt"
	"io"
	"math/big"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
)

var piDigits = []int{3, 1, 4, 1, 5, 9, 2, 6, 5, 3, 5, 8, 9, 7, 9, 3, 2, 3, 8, 4, 6, 2, 6, 4, 3, 3, 8, 3, 2, 7, 9, 5}

func runProgram(path string, input []byte) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(path, ".go") {
		cmd = exec.Command("go", "run", path)
	} else {
		cmd = exec.Command(path)
	}
	cmd.Stdin = bytes.NewReader(input)
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, stderr.String())
	}
	return out.String(), nil
}

func parseInput(data []byte) ([][]int64, error) {
	reader := bytes.NewReader(data)
	var t int
	if _, err := fmt.Fscan(reader, &t); err != nil {
		return nil, fmt.Errorf("failed to read t: %v", err)
	}
	if t > len(piDigits) {
		return nil, fmt.Errorf("too many test cases: %d > %d", t, len(piDigits))
	}
	tests := make([][]int64, t)
	for i := 0; i < t; i++ {
		n := piDigits[i]
		arr := make([]int64, n)
		for j := 0; j < n; j++ {
			if _, err := fmt.Fscan(reader, &arr[j]); err != nil {
				return nil, fmt.Errorf("failed to read a_%d for test %d: %v", j+1, i+1, err)
			}
		}
		tests[i] = arr
	}
	return tests, nil
}

func computeProducts(tests [][]int64) []string {
	ans := make([]string, len(tests))
	for i, arr := range tests {
		prod := big.NewInt(1)
		for _, v := range arr {
			prod.Mul(prod, big.NewInt(v))
		}
		ans[i] = prod.String()
	}
	return ans
}

func parseOutput(out string, t int) ([]string, error) {
	lines := strings.Fields(out)
	if len(lines) < t {
		return nil, fmt.Errorf("expected %d outputs, got %d", t, len(lines))
	}
	return lines[:t], nil
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	target := os.Args[len(os.Args)-1]
	if target == "--" {
		fmt.Println("usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}

	inputData, err := io.ReadAll(os.Stdin)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to read input: %v\n", err)
		os.Exit(1)
	}
	tests, err := parseInput(inputData)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
	expected := computeProducts(tests)

	_, src, _, _ := runtime.Caller(0)
	baseDir := filepath.Dir(src)
	refPath := filepath.Join(baseDir, "1812C.go")

	refOut, err := runProgram(refPath, inputData)
	if err != nil {
		fmt.Fprintf(os.Stderr, "reference runtime error: %v\n", err)
		os.Exit(1)
	}
	refAns, err := parseOutput(refOut, len(tests))
	if err != nil {
		fmt.Fprintf(os.Stderr, "reference output parse error: %v\n", err)
		os.Exit(1)
	}
	for i, val := range refAns {
		if val != expected[i] {
			fmt.Fprintf(os.Stderr, "reference mismatch on test %d: expected %s got %s\n", i+1, expected[i], val)
			os.Exit(1)
		}
	}

	targetOut, err := runProgram(target, inputData)
	if err != nil {
		fmt.Fprintf(os.Stderr, "target runtime error: %v\n", err)
		os.Exit(1)
	}
	ans, err := parseOutput(targetOut, len(tests))
	if err != nil {
		fmt.Fprintf(os.Stderr, "target output parse error: %v\n", err)
		os.Exit(1)
	}
	for i, val := range ans {
		if val != expected[i] {
			fmt.Fprintf(os.Stderr, "wrong answer on test %d: expected %s got %s\n", i+1, expected[i], val)
			os.Exit(1)
		}
	}
	fmt.Println("all tests passed")
}
