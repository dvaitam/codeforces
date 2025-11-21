package main

import (
	"bytes"
	"fmt"
	"math"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

type testCase struct {
	arr []int64
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
	return strings.TrimSpace(stdout.String()), nil
}

func deterministicTests() []testCase {
	return []testCase{
		{arr: []int64{5}},
		{arr: []int64{2, 1}},
		{arr: []int64{-5, -4, -3, -2, -1}},
		{arr: []int64{1, 2, -10, 3}},
		{arr: []int64{5, 4, 3, 2, 1}},
		{arr: []int64{0, -1, 2, -3, 4, -5}},
	}
}

func randomTest(rng *rand.Rand) testCase {
	n := rng.Intn(20) + 1
	arr := make([]int64, n)
	for i := 0; i < n; i++ {
		arr[i] = int64(rng.Intn(41) - 20)
	}
	return testCase{arr: arr}
}

func buildInput(tests []testCase) string {
	var sb strings.Builder
	sb.WriteString(strconv.Itoa(len(tests)))
	sb.WriteByte('\n')
	for _, tc := range tests {
		sb.WriteString(strconv.Itoa(len(tc.arr)))
		sb.WriteByte('\n')
		for i, v := range tc.arr {
			if i > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(strconv.FormatInt(v, 10))
		}
		sb.WriteByte('\n')
	}
	return sb.String()
}

func parseOperations(out string, tests []testCase) ([][][2]int, error) {
	tokens := strings.Fields(out)
	pos := 0
	ops := make([][][2]int, len(tests))
	for idx, tc := range tests {
		if pos >= len(tokens) {
			return nil, fmt.Errorf("test %d: missing operation count", idx+1)
		}
		k, err := strconv.Atoi(tokens[pos])
		if err != nil {
			return nil, fmt.Errorf("test %d: invalid operation count %q", idx+1, tokens[pos])
		}
		pos++
		if k < 0 || k > 31 {
			return nil, fmt.Errorf("test %d: invalid k=%d (should be between 0 and 31)", idx+1, k)
		}
		cur := make([][2]int, k)
		for i := 0; i < k; i++ {
			if pos+1 >= len(tokens) {
				return nil, fmt.Errorf("test %d: insufficient tokens for operation %d", idx+1, i+1)
			}
			x, err := strconv.Atoi(tokens[pos])
			if err != nil {
				return nil, fmt.Errorf("test %d: invalid i value %q", idx+1, tokens[pos])
			}
			pos++
			y, err := strconv.Atoi(tokens[pos])
			if err != nil {
				return nil, fmt.Errorf("test %d: invalid j value %q", idx+1, tokens[pos])
			}
			pos++
			if x < 1 || x > len(tc.arr) || y < 1 || y > len(tc.arr) {
				return nil, fmt.Errorf("test %d: operation indices out of range (%d,%d)", idx+1, x, y)
			}
			cur[i] = [2]int{x - 1, y - 1}
		}
		ops[idx] = cur
	}
	if pos != len(tokens) {
		return nil, fmt.Errorf("extra tokens at the end of output")
	}
	return ops, nil
}

func applyOperations(tc testCase, ops [][2]int) ([]int64, error) {
	arr := append([]int64(nil), tc.arr...)
	for _, op := range ops {
		i, j := op[0], op[1]
		arr[i] += arr[j]
		if arr[i] > math.MaxInt64/2 || arr[i] < math.MinInt64/2 {
			return nil, fmt.Errorf("overflow detected")
		}
	}
	return arr, nil
}

func isNonDecreasing(arr []int64) bool {
	for i := 0; i+1 < len(arr); i++ {
		if arr[i] > arr[i+1] {
			return false
		}
	}
	return true
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierA2.go /path/to/binary")
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

	allOps, err := parseOperations(out, tests)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}

	for idx, tc := range tests {
		finalArr, err := applyOperations(tc, allOps[idx])
		if err != nil {
			fmt.Fprintf(os.Stderr, "test %d: %v\n", idx+1, err)
			os.Exit(1)
		}
		if !isNonDecreasing(finalArr) {
			fmt.Fprintf(os.Stderr, "test %d: final array is not non-decreasing\n", idx+1)
			os.Exit(1)
		}
	}

	fmt.Printf("All %d tests passed\n", len(tests))
}
