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
	b []int64
}

const maxTotalN = 200000

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierA.go /path/to/binary")
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

	for i := range tests {
		if expected[i] != got[i] {
			fmt.Fprintf(os.Stderr, "test %d mismatch: expected %d got %d\nn=%d b=%v\n",
				i+1, expected[i], got[i], tests[i].n, tests[i].b)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}

func buildReference() (string, error) {
	const refName = "./ref_1320A.bin"
	cmd := exec.Command("go", "build", "-o", refName, "1320A.go")
	if out, err := cmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("failed to build reference: %v\n%s", err, string(out))
	}
	return refName, nil
}

func runAndParse(target, input string, tests int) ([]int64, error) {
	out, err := runProgram(target, input)
	if err != nil {
		return nil, err
	}
	fields := strings.Fields(out)
	if len(fields) != tests {
		return nil, fmt.Errorf("expected %d answers, got %d (output: %q)", tests, len(fields), out)
	}
	ans := make([]int64, tests)
	for i, f := range fields {
		val, err := strconv.ParseInt(f, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("invalid integer %q: %v", f, err)
		}
		ans[i] = val
	}
	return ans, nil
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
		remain := maxTotalN - total
		n := rng.Intn(min(5000, remain)) + 1
		b := make([]int64, n)
		for i := 0; i < n; i++ {
			switch rng.Intn(5) {
			case 0:
				b[i] = int64(rng.Intn(10) + 1)
			case 1:
				b[i] = int64(400000 - rng.Intn(100))
			default:
				b[i] = int64(rng.Intn(400000) + 1)
			}
		}
		tests = append(tests, testCase{n: n, b: b})
		total += n
	}
	return tests
}

func deterministicTests() []testCase {
	return []testCase{
		{n: 1, b: []int64{400000}},
		{n: 6, b: []int64{10, 7, 1, 9, 10, 15}},
		{n: 7, b: []int64{8, 9, 26, 11, 12, 29, 14}},
		{n: 5, b: []int64{1, 2, 3, 4, 5}},
	}
}

func serializeInput(tests []testCase) string {
	var sb strings.Builder
	sb.WriteString(strconv.Itoa(len(tests)))
	sb.WriteByte('\n')
	for _, tc := range tests {
		sb.WriteString(strconv.Itoa(tc.n))
		sb.WriteByte('\n')
		for i, v := range tc.b {
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
		total += tc.n
	}
	return total
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
