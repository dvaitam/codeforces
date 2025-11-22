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
	refSource = "2039F2.go"
	refBinary = "ref2039F2.bin"
)

type testCase struct {
	name string
	ms   []int
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierF2.go /path/to/binary")
		os.Exit(1)
	}
	candidate := os.Args[1]

	refPath, err := buildReference()
	if err != nil {
		fmt.Println("failed to build reference:", err)
		os.Exit(1)
	}
	defer os.Remove(refPath)

	tests := buildTests()
	for idx, tc := range tests {
		input := formatInput(tc)

		refOut, err := runProgram(refPath, input)
		if err != nil {
			fmt.Printf("reference runtime error on test %d (%s): %v\n", idx+1, tc.name, err)
			printInput(input)
			os.Exit(1)
		}
		refVals, err := parseOutputs(refOut, len(tc.ms))
		if err != nil {
			fmt.Printf("failed to parse reference output on test %d (%s): %v\noutput:\n%s\n", idx+1, tc.name, err, refOut)
			printInput(input)
			os.Exit(1)
		}

		candOut, err := runProgram(candidate, input)
		if err != nil {
			fmt.Printf("candidate runtime error on test %d (%s): %v\n", idx+1, tc.name, err)
			printInput(input)
			os.Exit(1)
		}
		candVals, err := parseOutputs(candOut, len(tc.ms))
		if err != nil {
			fmt.Printf("failed to parse candidate output on test %d (%s): %v\noutput:\n%s\n", idx+1, tc.name, err, candOut)
			printInput(input)
			os.Exit(1)
		}

		for i := range refVals {
			if candVals[i] != refVals[i] {
				fmt.Printf("test %d (%s) failed on case %d: expected %d, got %d\n", idx+1, tc.name, i+1, refVals[i], candVals[i])
				printInput(input)
				fmt.Println("Candidate output:")
				fmt.Println(candOut)
				os.Exit(1)
			}
		}
	}

	fmt.Printf("All %d tests passed\n", len(tests))
}

func buildReference() (string, error) {
	cmd := exec.Command("go", "build", "-o", refBinary, refSource)
	if out, err := cmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("%v\n%s", err, string(out))
	}
	return filepath.Join(".", refBinary), nil
}

func runProgram(path string, input []byte) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(path, ".go") {
		cmd = exec.Command("go", "run", path)
	} else {
		cmd = exec.Command(path)
	}
	cmd.Stdin = bytes.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return out.String(), err
	}
	return out.String(), nil
}

func parseOutputs(out string, want int) ([]int64, error) {
	fields := strings.Fields(out)
	if len(fields) != want {
		return nil, fmt.Errorf("expected %d integers, got %d", want, len(fields))
	}
	res := make([]int64, want)
	for i, f := range fields {
		v, err := strconv.ParseInt(f, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("invalid integer %q: %v", f, err)
		}
		res[i] = v
	}
	return res, nil
}

func formatInput(tc testCase) []byte {
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d\n", len(tc.ms))
	for _, v := range tc.ms {
		fmt.Fprintf(&sb, "%d\n", v)
	}
	return []byte(sb.String())
}

func buildTests() []testCase {
	tests := []testCase{
		{name: "sample", ms: []int{2, 5, 9}},
		{name: "single-1", ms: []int{1}},
		{name: "single-2", ms: []int{2}},
		{name: "increasing-small", ms: []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}},
		{name: "pair-extremes", ms: []int{1, 1_000_000}},
		{name: "big-single", ms: []int{1_000_000}},
	}

	// Deterministic pseudo-random tests for repeatability.
	rnd := rand.New(rand.NewSource(20240521))
	for i := 0; i < 25; i++ {
		t := rnd.Intn(8) + 1 // between 1 and 8 test cases
		ms := make([]int, t)
		for j := range ms {
			if i%5 == 0 {
				// bias some tests toward smaller values
				ms[j] = rnd.Intn(50) + 1
			} else {
				ms[j] = rnd.Intn(1_000_000) + 1
			}
		}
		tests = append(tests, testCase{
			name: fmt.Sprintf("rand-%d", i+1),
			ms:   ms,
		})
	}

	// Larger batch to ensure handling of many test cases in one input.
	largeBatch := make([]int, 40)
	for i := range largeBatch {
		if i%3 == 0 {
			largeBatch[i] = 1
		} else if i%3 == 1 {
			largeBatch[i] = 1_000_000
		} else {
			largeBatch[i] = rnd.Intn(500_000) + 1
		}
	}
	tests = append(tests, testCase{name: "bulk-40", ms: largeBatch})

	return tests
}

func printInput(in []byte) {
	fmt.Println("Input used:")
	fmt.Print(string(in))
}
