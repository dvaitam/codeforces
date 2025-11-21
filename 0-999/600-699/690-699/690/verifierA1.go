package main

import (
	"bytes"
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
	refSourceA1 = "690A1.go"
	refBinaryA1 = "ref690A1.bin"
	totalTests  = 120
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierA1.go /path/to/binary")
		os.Exit(1)
	}
	candidate := os.Args[1]

	refPath, err := buildReference()
	if err != nil {
		fmt.Println("failed to build reference:", err)
		os.Exit(1)
	}
	defer os.Remove(refPath)

	tests := generateTests()
	for idx, n := range tests {
		input := fmt.Sprintf("%d\n", n)

		refOut, err := runProgram(refPath, input)
		if err != nil {
			fmt.Printf("reference runtime error on test %d: %v\n", idx+1, err)
			printInput(input)
			os.Exit(1)
		}
		candOut, err := runProgram(candidate, input)
		if err != nil {
			fmt.Printf("candidate runtime error on test %d: %v\n", idx+1, err)
			printInput(input)
			os.Exit(1)
		}

		refVal, err := parseOutput(refOut)
		if err != nil {
			fmt.Printf("failed to parse reference output on test %d: %v\noutput:\n%s\n", idx+1, err, refOut)
			printInput(input)
			os.Exit(1)
		}
		candVal, err := parseOutput(candOut)
		if err != nil {
			fmt.Printf("failed to parse candidate output on test %d: %v\noutput:\n%s\n", idx+1, err, candOut)
			printInput(input)
			os.Exit(1)
		}

		if refVal != candVal {
			fmt.Printf("test %d failed: expected %d, got %d\n", idx+1, refVal, candVal)
			printInput(input)
			fmt.Println("Reference output:")
			fmt.Println(refOut)
			fmt.Println("Candidate output:")
			fmt.Println(candOut)
			os.Exit(1)
		}
	}

	fmt.Printf("All %d tests passed\n", len(tests))
}

func buildReference() (string, error) {
	cmd := exec.Command("go", "build", "-o", refBinaryA1, refSourceA1)
	if out, err := cmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("%v\n%s", err, string(out))
	}
	return filepath.Join(".", refBinaryA1), nil
}

func runProgram(path string, input string) (string, error) {
	cmd := exec.Command(path)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	return strings.TrimSpace(out.String()), err
}

func parseOutput(out string) (int64, error) {
	if out == "" {
		return 0, fmt.Errorf("empty output")
	}
	val, err := strconv.ParseInt(strings.Fields(out)[0], 10, 64)
	if err != nil {
		return 0, fmt.Errorf("invalid integer %q: %v", out, err)
	}
	return val, nil
}

func generateTests() []int64 {
	tests := []int64{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	rnd := rand.New(rand.NewSource(time.Now().UnixNano()))
	for len(tests) < totalTests-3 {
		tests = append(tests, randInt64(rnd, 1, 1_000_000_000))
	}
	tests = append(tests, 999_999_999, 1_000_000_000, 500_000_001)
	return tests
}

func randInt64(rnd *rand.Rand, lo, hi int64) int64 {
	if lo == hi {
		return lo
	}
	return lo + rnd.Int63n(hi-lo+1)
}

func printInput(in string) {
	fmt.Println("Input used:")
	fmt.Print(in)
}
