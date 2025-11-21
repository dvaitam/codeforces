package main

import (
	"bufio"
	"bytes"
	"fmt"
	"math"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"sort"
	"strconv"
	"strings"
)

type testCase struct {
	name  string
	input string
}

var verifierDir string

func init() {
	if _, file, _, ok := runtime.Caller(0); ok {
		verifierDir = filepath.Dir(file)
	} else {
		verifierDir = "."
	}
}

func buildReference() (string, error) {
	outPath := filepath.Join(verifierDir, "ref1916D.bin")
	cmd := exec.Command("go", "build", "-o", outPath, "1916D.go")
	cmd.Dir = verifierDir
	if out, err := cmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("failed to build reference: %v\n%s", err, out)
	}
	return outPath, nil
}

func runProgram(target, input string) (string, error) {
	if !filepath.IsAbs(target) {
		if abs, err := filepath.Abs(target); err == nil {
			target = abs
		}
	}
	var cmd *exec.Cmd
	if strings.HasSuffix(target, ".go") {
		cmd = exec.Command("go", "run", target)
	} else {
		cmd = exec.Command(target)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return out.String(), fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	return out.String(), nil
}

func parseOutput(output string, tests []int) ([][]string, error) {
	reader := bufio.NewScanner(strings.NewReader(output))
	reader.Buffer(make([]byte, 0, 1024), 1<<20)
	results := make([][]string, len(tests))
	for i, n := range tests {
		results[i] = make([]string, 0, n)
		for len(results[i]) < n {
			if !reader.Scan() {
				if err := reader.Err(); err != nil {
					return nil, fmt.Errorf("failed to read output: %v", err)
				}
				return nil, fmt.Errorf("expected %d numbers for test %d, got %d", n, i+1, len(results[i]))
			}
			line := strings.TrimSpace(reader.Text())
			if line == "" {
				continue
			}
			results[i] = append(results[i], line)
		}
	}
	if reader.Scan() {
		return nil, fmt.Errorf("extra output: %s", reader.Text())
	}
	return results, nil
}

func isPerfectSquare(s string) bool {
	if len(s) == 0 {
		return false
	}
	if s[0] == '0' {
		return false
	}
	value, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return false
	}
	root := int64(math.Round(math.Sqrt(float64(value))))
	return root*root == value
}

func checkTestCase(n int, numbers []string) error {
	if len(numbers) != n {
		return fmt.Errorf("expected %d numbers, got %d", n, len(numbers))
	}
	digitCount := make(map[string]int)
	for idx, num := range numbers {
		if len(num) != n {
			return fmt.Errorf("number %d has length %d, expected %d", idx+1, len(num), n)
		}
		if !isPerfectSquare(num) {
			return fmt.Errorf("number %d (%s) is not a perfect square", idx+1, num)
		}
		// build digit multiset
		bytesDigits := []byte(num)
		sort.Slice(bytesDigits, func(i, j int) bool { return bytesDigits[i] < bytesDigits[j] })
		digitPattern := string(bytesDigits)
		if idx == 0 {
			digitCount[digitPattern] = 1
		} else {
			digitCount[digitPattern]++
		}
		if idx > 0 && digitCount[digitPattern] != idx+1 {
			return fmt.Errorf("multiset mismatch")
		}
	}
	return nil
}

func verifyCase(candidate, reference string, tc testCase) error {
	// parse test case
	reader := strings.NewReader(tc.input)
	var t int
	fmt.Fscan(reader, &t)
	tests := make([]int, t)
	for i := 0; i < t; i++ {
		fmt.Fscan(reader, &tests[i])
	}

	refOut, err := runProgram(reference, tc.input)
	if err != nil {
		return fmt.Errorf("reference error: %v", err)
	}
	refNumbers, err := parseOutput(refOut, tests)
	if err != nil {
		return fmt.Errorf("reference produced invalid output: %v", err)
	}
	for i, n := range tests {
		if err := checkTestCase(n, refNumbers[i]); err != nil {
			return fmt.Errorf("reference invalid for test %d: %v", i+1, err)
		}
	}

	candOut, err := runProgram(candidate, tc.input)
	if err != nil {
		return fmt.Errorf("candidate error: %v", err)
	}
	candNumbers, err := parseOutput(candOut, tests)
	if err != nil {
		return fmt.Errorf("invalid candidate output: %v", err)
	}
	for i, n := range tests {
		if err := checkTestCase(n, candNumbers[i]); err != nil {
			return fmt.Errorf("test %d invalid: %v\ncandidate output:\n%s", i+1, err, candOut)
		}
	}
	return nil
}

func manualTests() []testCase {
	return []testCase{
		{name: "small_cases", input: "3\n1\n3\n5\n"},
		{name: "single_large", input: "1\n9\n"},
	}
}

func main() {
	args := os.Args[1:]
	if len(args) > 0 && args[0] == "--" {
		args = args[1:]
	}
	if len(args) != 1 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	candidate := args[0]

	ref, err := buildReference()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer os.Remove(ref)

	tests := manualTests()
	for i, tc := range tests {
		if err := verifyCase(candidate, ref, tc); err != nil {
			fmt.Fprintf(os.Stderr, "case %d (%s) failed: %v\ninput:\n%s", i+1, tc.name, err, tc.input)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed.\n", len(tests))
}
