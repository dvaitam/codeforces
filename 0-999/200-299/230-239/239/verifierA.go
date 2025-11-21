package main

import (
	"bufio"
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"sort"
	"strings"
)

type testCase struct {
	name  string
	input string
	// verification uses input alone; no single expected string due to multiple valid outputs
}

func parseInput(input string) (int64, int64, int64, error) {
	reader := bufio.NewReader(strings.NewReader(input))
	var y, k, n int64
	_, err := fmt.Fscan(reader, &y, &k, &n)
	return y, k, n, err
}

func expectedValues(y, k, n int64) []int64 {
	start := ((y + k) / k) * k
	if start <= y {
		start += k
	}
	if start > n {
		return nil
	}
	var vals []int64
	for m := start; m <= n; m += k {
		vals = append(vals, m-y)
	}
	return vals
}

func validateOutput(y, k, n int64, output string) error {
	output = strings.TrimSpace(output)
	expect := expectedValues(y, k, n)
	if len(expect) == 0 {
		if output != "-1" {
			return fmt.Errorf("expected -1 but got %q", output)
		}
		return nil
	}
	if output == "" {
		return fmt.Errorf("expected positive list but got empty output")
	}
	fields := strings.Fields(output)
	values := make([]int64, len(fields))
	for i, f := range fields {
		var val int64
		if _, err := fmt.Sscan(f, &val); err != nil {
			return fmt.Errorf("failed to parse value %q: %v", f, err)
		}
		values[i] = val
	}
	if !sort.SliceIsSorted(values, func(i, j int) bool { return values[i] < values[j] }) {
		return fmt.Errorf("values are not strictly increasing: %v", values)
	}
	for i := 1; i < len(values); i++ {
		if values[i] == values[i-1] {
			return fmt.Errorf("duplicate value %d", values[i])
		}
	}
	if len(values) != len(expect) {
		return fmt.Errorf("expected %d values but got %d (%v)", len(expect), len(values), values)
	}
	for i, val := range values {
		if val != expect[i] {
			return fmt.Errorf("value mismatch at pos %d: expect %d got %d", i, expect[i], val)
		}
		if val <= 0 {
			return fmt.Errorf("x must be positive, got %d", val)
		}
		if y+val > n {
			return fmt.Errorf("y+x exceeds n: y=%d x=%d n=%d", y, val, n)
		}
		if (y+val)%k != 0 {
			return fmt.Errorf("(y+x) is not divisible by k: y=%d x=%d k=%d", y, val, k)
		}
	}
	return nil
}

func runCandidate(bin, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return out.String(), err
}

func makeCase(name string, y, k, n int64) testCase {
	return testCase{
		name:  name,
		input: fmt.Sprintf("%d %d %d\n", y, k, n),
	}
}

func handcraftedTests() []testCase {
	return []testCase{
		makeCase("min_no_solution", 1, 1, 1),
		makeCase("simple_one", 1, 2, 3),
		makeCase("simple_multiple", 1, 2, 9),
		makeCase("y_equals_n", 5, 2, 5),
		makeCase("k_big", 10, 100, 1000),
		makeCase("exact_upper", 3, 4, 11),
		makeCase("large_values", 999999937, 999999937, 1000000000),
		makeCase("close_bounds", 100, 101, 202),
		makeCase("mid_range", 123456, 7890, 1000000),
		makeCase("equal_k_n", 500000, 500000, 500000000),
	}
}

func randomTests() []testCase {
	rng := rand.New(rand.NewSource(239))
	var tests []testCase
	add := func(prefix string, count int, maxVal int64) {
		for i := 0; i < count; i++ {
			y := int64(rng.Int63n(maxVal) + 1)
			k := int64(rng.Int63n(maxVal) + 1)
			n := y + int64(rng.Int63n(maxVal))
			if n < y {
				n = y
			}
			tests = append(tests, makeCase(fmt.Sprintf("%s_%d", prefix, i+1), y, k, n))
		}
	}
	add("small", 200, 1_000)
	add("medium", 200, 1_000_000)
	add("large", 200, 1_000_000_000)
	return tests
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	tests := append(handcraftedTests(), randomTests()...)
	for idx, tc := range tests {
		y, k, n, err := parseInput(tc.input)
		if err != nil {
			fmt.Printf("failed to parse test %s: %v\n", tc.name, err)
			os.Exit(1)
		}
		output, runErr := runCandidate(bin, tc.input)
		if runErr != nil {
			fmt.Printf("test %d (%s) runtime error: %v\ninput:\n%s", idx+1, tc.name, runErr, tc.input)
			os.Exit(1)
		}
		if err := validateOutput(y, k, n, output); err != nil {
			fmt.Printf("test %d (%s) failed: %v\ninput:\n%soutput:\n%s", idx+1, tc.name, err, tc.input, output)
			os.Exit(1)
		}
	}
	fmt.Println("all tests passed")
}
