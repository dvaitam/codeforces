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
	s string
}

const maxTotalN = 400000

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierC.go /path/to/binary")
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
			fmt.Fprintf(os.Stderr, "Mismatch on test %d: expected %d got %d\nn=%d s=%s\n", i+1, expected[i], got[i], len(tests[i].s), tests[i].s)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}

func buildReference() (string, error) {
	const refName = "./ref_2026C.bin"
	cmd := exec.Command("go", "build", "-o", refName, "2026C.go")
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

	for total < maxTotalN {
		remaining := maxTotalN - total
		n := rng.Intn(min(10000, remaining)) + 1
		if n <= 0 {
			break
		}
		sb := strings.Builder{}
		for i := 0; i < n-1; i++ {
			if rng.Intn(4) == 0 {
				sb.WriteByte('1')
			} else {
				sb.WriteByte('0')
			}
		}
		sb.WriteByte('1') // ensure s_n = '1'
		tests = append(tests, testCase{s: sb.String()})
		total += n
	}
	return tests
}

func deterministicTests() []testCase {
	return []testCase{
		{s: "1"},
		{s: "11"},
		{s: "10"},
		{s: "101"},
		{s: "111111"},
		{s: "1010101"},
	}
}

func totalLength(tests []testCase) int {
	total := 0
	for _, tc := range tests {
		total += len(tc.s)
	}
	return total
}

func serializeInput(tests []testCase) string {
	var sb strings.Builder
	sb.WriteString(strconv.Itoa(len(tests)))
	sb.WriteByte('\n')
	for _, tc := range tests {
		sb.WriteString(strconv.Itoa(len(tc.s)))
		sb.WriteByte('\n')
		sb.WriteString(tc.s)
		sb.WriteByte('\n')
	}
	return sb.String()
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
