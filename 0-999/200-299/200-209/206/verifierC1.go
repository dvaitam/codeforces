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
	refSourceC1 = "206C1.go"
	refBinaryC1 = "ref206C1.bin"
	totalTests  = 80
)

type operation struct {
	t int
	v int
	c byte
}

type testCase struct {
	ops []operation
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierC1.go /path/to/binary")
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
		input := formatInput(tc)

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

		refVals, err := parseOutput(refOut, len(tc.ops))
		if err != nil {
			fmt.Printf("failed to parse reference output on test %d: %v\noutput:\n%s\n", i+1, err, refOut)
			printInput(input)
			os.Exit(1)
		}
		candVals, err := parseOutput(candOut, len(tc.ops))
		if err != nil {
			fmt.Printf("failed to parse candidate output on test %d: %v\noutput:\n%s\n", i+1, err, candOut)
			printInput(input)
			os.Exit(1)
		}

		if len(refVals) != len(candVals) {
			fmt.Printf("test %d failed: expected %d outputs, got %d\n", i+1, len(refVals), len(candVals))
			printInput(input)
			fmt.Println("Reference output:")
			fmt.Println(refOut)
			fmt.Println("Candidate output:")
			fmt.Println(candOut)
			os.Exit(1)
		}
		for idx := range refVals {
			if refVals[idx] != candVals[idx] {
				fmt.Printf("test %d failed at line %d: expected %d, got %d\n", i+1, idx+1, refVals[idx], candVals[idx])
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
	cmd := exec.Command("go", "build", "-o", refBinaryC1, refSourceC1)
	if out, err := cmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("%v\n%s", err, string(out))
	}
	return filepath.Join(".", refBinaryC1), nil
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

func formatInput(tc testCase) []byte {
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d\n", len(tc.ops))
	for _, op := range tc.ops {
		sb.WriteString(fmt.Sprintf("%d %d %c\n", op.t, op.v, op.c))
	}
	return []byte(sb.String())
}

func parseOutput(out string, expected int) ([]int64, error) {
	fields := strings.Fields(out)
	if len(fields) != expected {
		return nil, fmt.Errorf("expected %d integers, got %d", expected, len(fields))
	}
	res := make([]int64, expected)
	for i, f := range fields {
		val, err := strconv.ParseInt(f, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("invalid integer %q: %v", f, err)
		}
		res[i] = val
	}
	return res, nil
}

func generateTests() []testCase {
	tests := []testCase{
		{ops: []operation{{t: 1, v: 1, c: 'a'}}},
		{
			ops: []operation{
				{t: 2, v: 1, c: 'a'},
				{t: 1, v: 1, c: 'a'},
				{t: 2, v: 2, c: 'a'},
				{t: 1, v: 2, c: 'a'},
			},
		},
		{
			ops: []operation{
				{t: 2, v: 1, c: 'b'},
				{t: 2, v: 2, c: 'a'},
				{t: 1, v: 1, c: 'a'},
				{t: 1, v: 2, c: 'b'},
				{t: 2, v: 3, c: 'b'},
			},
		},
	}

	rnd := rand.New(rand.NewSource(time.Now().UnixNano()))
	for len(tests) < totalTests-3 {
		n := rnd.Intn(400) + 1
		ops := randomOperations(rnd, n)
		tests = append(tests, testCase{ops: ops})
	}

	tests = append(tests, testCase{ops: randomOperations(rand.New(rand.NewSource(1)), 2000)})
	tests = append(tests, testCase{ops: randomOperations(rand.New(rand.NewSource(2)), 7000)})
	tests = append(tests, testCase{ops: randomOperations(rand.New(rand.NewSource(3)), 100000)})

	return tests
}

func randomOperations(rnd *rand.Rand, n int) []operation {
	ops := make([]operation, 0, n)
	size1, size2 := 1, 1
	for i := 0; i < n; i++ {
		t := rnd.Intn(2) + 1
		if t == 1 {
			v := rnd.Intn(size1) + 1
			c := byte('a' + rnd.Intn(26))
			ops = append(ops, operation{t: 1, v: v, c: c})
			size1++
		} else {
			v := rnd.Intn(size2) + 1
			c := byte('a' + rnd.Intn(26))
			ops = append(ops, operation{t: 2, v: v, c: c})
			size2++
		}
	}
	return ops
}

func printInput(in []byte) {
	fmt.Println("Input used:")
	fmt.Println(string(in))
}
