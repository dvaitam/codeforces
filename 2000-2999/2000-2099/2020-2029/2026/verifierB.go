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
	a []int64
}

const maxNTotal = 2000

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierB.go /path/to/binary")
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

	expected, err := runAndParse(refBin, input, len(tests))
	if err != nil {
		fmt.Fprintf(os.Stderr, "reference failed: %v\n", err)
		os.Exit(1)
	}

	got, err := runAndParse(candidate, input, len(tests))
	if err != nil {
		fmt.Fprintf(os.Stderr, "candidate failed: %v\n", err)
		os.Exit(1)
	}

	for idx := range tests {
		if expected[idx] != got[idx] {
			fmt.Fprintf(os.Stderr, "test %d mismatch: expected %d got %d\nsequence: %v\n", idx+1, expected[idx], got[idx], tests[idx].a)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}

func buildReference() (string, error) {
	const refName = "./ref_2026B.bin"
	cmd := exec.Command("go", "build", "-o", refName, "2026B.go")
	if out, err := cmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("failed to build reference: %v\n%s", err, string(out))
	}
	return refName, nil
}

func runAndParse(target, input string, cases int) ([]int64, error) {
	out, err := runProgram(target, input)
	if err != nil {
		return nil, err
	}
	fields := strings.Fields(out)
	if len(fields) != cases {
		return nil, fmt.Errorf("expected %d answers, got %d (output: %q)", cases, len(fields), out)
	}
	res := make([]int64, cases)
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
	for total < maxNTotal {
		remaining := maxNTotal - total
		n := rng.Intn(remaining) + 1
		arr := make([]int64, n)
		cur := rand.Int63n(10)
		arr[0] = int64(cur + 1)
		for i := 1; i < n; i++ {
			cur += rand.Int63n(10) + 1
			arr[i] = cur
		}
		tests = append(tests, testCase{a: arr})
		total += n
	}
	return tests
}

func deterministicTests() []testCase {
	return []testCase{
		{a: []int64{1}},
		{a: []int64{1, 2}},
		{a: []int64{1, 3}},
		{a: []int64{1, 2, 3}},
		{a: []int64{2, 4, 9}},
		{a: []int64{1, 5, 8, 10, 13}},
	}
}

func totalLength(tests []testCase) int {
	total := 0
	for _, tc := range tests {
		total += len(tc.a)
	}
	return total
}

func serializeInput(tests []testCase) string {
	var sb strings.Builder
	sb.WriteString(strconv.Itoa(len(tests)))
	sb.WriteByte('\n')
	for _, tc := range tests {
		sb.WriteString(strconv.Itoa(len(tc.a)))
		sb.WriteByte('\n')
		for i, v := range tc.a {
			if i > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(strconv.FormatInt(v, 10))
		}
		sb.WriteByte('\n')
	}
	return sb.String()
}
