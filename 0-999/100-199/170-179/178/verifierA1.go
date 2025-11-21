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
	refSourceA1 = "178A1.go"
	refBinaryA1 = "ref178A1.bin"
	totalTests  = 60
)

type testCase struct {
	n   int
	arr []int
}

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
	for i, tc := range tests {
		input := buildInput(tc)

		refOut, err := runProgram(refPath, input)
		if err != nil {
			fmt.Printf("reference runtime error on test %d: %v\n", i+1, err)
			printInput(input)
			os.Exit(1)
		}

		candOut, err := runProgram(candidate, input)
		if err != nil {
			fmt.Printf("candidate runtime error on test %d: %v\n", i+1, err)
			printInput(input)
			os.Exit(1)
		}

		refVals, err := parseOutput(refOut, tc.n)
		if err != nil {
			fmt.Printf("failed to parse reference output on test %d: %v\noutput:\n%s\n", i+1, err, refOut)
			printInput(input)
			os.Exit(1)
		}
		candVals, err := parseOutput(candOut, tc.n)
		if err != nil {
			fmt.Printf("failed to parse candidate output on test %d: %v\noutput:\n%s\n", i+1, err, candOut)
			printInput(input)
			os.Exit(1)
		}

		if len(refVals) != len(candVals) {
			fmt.Printf("test %d failed: expected %d values, got %d\n", i+1, len(refVals), len(candVals))
			printInput(input)
			fmt.Println("Reference output:")
			fmt.Println(refOut)
			fmt.Println("Candidate output:")
			fmt.Println(candOut)
			os.Exit(1)
		}

		for idx := range refVals {
			if refVals[idx] != candVals[idx] {
				fmt.Printf("test %d failed at position %d: expected %d, got %d\n", i+1, idx+1, refVals[idx], candVals[idx])
				printInput(input)
				fmt.Println("Reference output:")
				fmt.Println(refOut)
				fmt.Println("Candidate output:")
				fmt.Println(candOut)
				os.Exit(1)
			}
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

func runProgram(path string, input []byte) (string, error) {
	cmd := exec.Command(path)
	cmd.Stdin = bytes.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	return out.String(), err
}

func generateTests() []testCase {
	tests := []testCase{
		{n: 1, arr: []int{0}},
		{n: 2, arr: []int{0, 0}},
		{n: 2, arr: []int{7, 0}},
		{n: 5, arr: []int{1, 0, 1, 0, 1}},
		{n: 6, arr: []int{10, 9, 8, 7, 6, 5}},
	}

	rnd := rand.New(rand.NewSource(time.Now().UnixNano()))
	for len(tests) < totalTests-1 {
		n := rnd.Intn(60) + 2
		arr := make([]int, n)
		for i := 0; i < n; i++ {
			arr[i] = rnd.Intn(10001)
		}
		tests = append(tests, testCase{n: n, arr: arr})
	}

	largeN := 100000
	largeArr := make([]int, largeN)
	for i := 0; i < largeN; i++ {
		largeArr[i] = (i * 37) % 10001
	}
	tests = append(tests, testCase{n: largeN, arr: largeArr})
	return tests
}

func buildInput(tc testCase) []byte {
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d\n", tc.n)
	for i, v := range tc.arr {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(strconv.Itoa(v))
	}
	sb.WriteByte('\n')
	return []byte(sb.String())
}

func parseOutput(out string, n int) ([]int, error) {
	fields := strings.Fields(out)
	expected := n - 1
	if len(fields) != expected {
		if expected == 0 && len(fields) == 0 {
			return []int{}, nil
		}
		return nil, fmt.Errorf("expected %d integers, got %d", expected, len(fields))
	}
	res := make([]int, len(fields))
	for i, f := range fields {
		val, err := strconv.Atoi(f)
		if err != nil {
			return nil, fmt.Errorf("invalid integer %q: %v", f, err)
		}
		res[i] = val
	}
	return res, nil
}

func printInput(in []byte) {
	fmt.Println("Input used:")
	fmt.Println(string(in))
}
