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
	arr []int64
}

const maxTotalN = 200000

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	candidate := os.Args[1]

	refBin, err := buildReference()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer os.Remove(refBin)

	tests := buildTests()
	input := serializeInput(tests)
	expected, err := runAndParse(refBin, input, totalLength(tests))
	if err != nil {
		fmt.Fprintf(os.Stderr, "reference failed: %v\n", err)
		os.Exit(1)
	}

	got, err := runAndParse(candidate, input, totalLength(tests))
	if err != nil {
		fmt.Fprintf(os.Stderr, "candidate failed: %v\n", err)
		os.Exit(1)
	}

	for tIdx, tc := range tests {
		for i := 0; i < len(tc.arr); i++ {
			pos := prefixOffset(tests, tIdx, i)
			if expected[pos] != got[pos] {
				fmt.Fprintf(os.Stderr, "Mismatch in test %d prefix %d: expected %d got %d\n", tIdx+1, i+1, expected[pos], got[pos])
				fmt.Fprintf(os.Stderr, "array: %v\n", tc.arr)
				os.Exit(1)
			}
		}
	}

	fmt.Println("All tests passed")
}

func buildReference() (string, error) {
	const refName = "./ref_2035D.bin"
	cmd := exec.Command("go", "build", "-o", refName, "2035D.go")
	if out, err := cmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("failed to build reference: %v\n%s", err, string(out))
	}
	return refName, nil
}

func runAndParse(target, input string, expectedCount int) ([]int64, error) {
	out, err := runProgram(target, input)
	if err != nil {
		return nil, err
	}
	fields := strings.Fields(out)
	if len(fields) != expectedCount {
		return nil, fmt.Errorf("expected %d outputs, got %d (output: %q)", expectedCount, len(fields), out)
	}
	res := make([]int64, expectedCount)
	for i, f := range fields {
		val, err := strconv.ParseInt(f, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("invalid integer %q: %v", f, err)
		}
		res[i] = val
	}
	return res, nil
}

func runProgram(target, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(target, ".go") {
		cmd = exec.Command("go", "run", target)
	} else {
		cmd = exec.Command(target)
	}
	cmd.Stdin = strings.NewReader(input)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\nstdout:\n%s\nstderr:\n%s", err, stdout.String(), stderr.String())
	}
	return stdout.String(), nil
}

func buildTests() []testCase {
	tests := deterministicTests()
	total := totalLength(tests)
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for total < maxTotalN {
		length := rng.Intn(2000) + 1
		if total+length > maxTotalN {
			length = maxTotalN - total
		}
		arr := make([]int64, length)
		for i := 0; i < length; i++ {
			switch rng.Intn(6) {
			case 0:
				arr[i] = 1
			case 1:
				arr[i] = int64(rng.Intn(1_000_000_000) + 1)
			default:
				arr[i] = int64(rng.Intn(1000) + 1)
			}
		}
		tests = append(tests, testCase{arr: arr})
		total += length
	}
	return tests
}

func deterministicTests() []testCase {
	return []testCase{
		{arr: []int64{1}},
		{arr: []int64{1, 2, 3}},
		{arr: []int64{2, 4, 8, 16}},
		{arr: []int64{3, 6, 12, 24}},
		{arr: []int64{5, 7, 9, 11}},
		{arr: []int64{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}},
		{arr: []int64{1_000_000_000, 1_000_000_000, 1_000_000_000}},
	}
}

func serializeInput(tests []testCase) string {
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

func totalLength(tests []testCase) int {
	total := 0
	for _, tc := range tests {
		total += len(tc.arr)
	}
	return total
}

func prefixOffset(tests []testCase, testIdx, prefix int) int {
	offset := 0
	for i := 0; i < testIdx; i++ {
		offset += len(tests[i].arr)
	}
	return offset + prefix
}
