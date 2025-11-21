package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

type testCase struct {
	n int
	s int
	r int
}

func runBinary(bin string, input string) (string, error) {
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

func deterministicTests() []testCase {
	return []testCase{
		{n: 2, s: 7, r: 3},
		{n: 3, s: 10, r: 6},
		{n: 5, s: 20, r: 14},
	}
}

func randomTest(rng *rand.Rand) testCase {
	for {
		n := rng.Intn(5) + 2
		values := make([]int, n)
		sum := 0
		maxVal := 0
		for i := 0; i < n; i++ {
			val := rng.Intn(6) + 1
			values[i] = val
			sum += val
			if val > maxVal {
				maxVal = val
			}
		}
		s := sum
		r := s - maxVal
		if r >= 1 && s <= 300 {
			return testCase{n: n, s: s, r: r}
		}
	}
}

func buildInput(tests []testCase) string {
	var sb strings.Builder
	sb.WriteString(strconv.Itoa(len(tests)))
	sb.WriteByte('\n')
	for _, tc := range tests {
		sb.WriteString(fmt.Sprintf("%d %d %d\n", tc.n, tc.s, tc.r))
	}
	return sb.String()
}

func parseLines(out string, expected int) ([]string, error) {
	out = strings.TrimSpace(out)
	lines := strings.Split(out, "\n")
	if len(lines) != expected {
		return nil, fmt.Errorf("expected %d lines, got %d", expected, len(lines))
	}
	return lines, nil
}

func validateLine(line string, tc testCase) error {
	fields := strings.Fields(line)
	if len(fields) != tc.n {
		return fmt.Errorf("expected %d numbers, got %d", tc.n, len(fields))
	}
	sum := 0
	maxVal := 0
	countMax := 0
	for _, f := range fields {
		val, err := strconv.Atoi(f)
		if err != nil {
			return fmt.Errorf("invalid integer %q", f)
		}
		if val < 1 || val > 6 {
			return fmt.Errorf("value %d out of range [1,6]", val)
		}
		sum += val
		if val == maxVal {
			countMax++
		}
		if val > maxVal {
			maxVal = val
			countMax = 1
		}
	}
	if sum != tc.s {
		return fmt.Errorf("sum mismatch: expected %d got %d", tc.s, sum)
	}
	requiredMax := tc.s - tc.r
	if maxVal != requiredMax {
		return fmt.Errorf("max mismatch: expected %d got %d", requiredMax, maxVal)
	}
	if countMax == 0 {
		return fmt.Errorf("no die with required max value %d", requiredMax)
	}
	if sum-maxVal != tc.r {
		return fmt.Errorf("sum after removing max mismatch: expected %d got %d", tc.r, sum-maxVal)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	target := os.Args[1]

	tests := deterministicTests()
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 200; i++ {
		tests = append(tests, randomTest(rng))
	}

	input := buildInput(tests)
	out, err := runBinary(target, input)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
	lines, err := parseLines(out, len(tests))
	if err != nil {
		fmt.Fprintf(os.Stderr, "output format error: %v\n", err)
		os.Exit(1)
	}
	for i, line := range lines {
		if err := validateLine(line, tests[i]); err != nil {
			fmt.Fprintf(os.Stderr, "test case %d failed: %v\ninput: %d %d %d\n", i+1, err, tests[i].n, tests[i].s, tests[i].r)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(tests))
}
