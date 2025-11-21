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
	k int64
	s string
}

const maxTotalN = 200000

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

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	tests := buildTests(rng)
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

	for i := range expected {
		if expected[i] != got[i] {
			fmt.Fprintf(os.Stderr, "test %d mismatch: expected %d got %d\nn=%d k=%d s=%s\n", i+1, expected[i], got[i], tests[i].n, tests[i].k, tests[i].s)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}

func buildReference() (string, error) {
	const refName = "./ref_2042C.bin"
	cmd := exec.Command("go", "build", "-o", refName, "2042C.go")
	if out, err := cmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("failed to build reference: %v\n%s", err, string(out))
	}
	return refName, nil
}

func runAndParse(target, input string, count int) ([]int64, error) {
	out, err := runProgram(target, input)
	if err != nil {
		return nil, err
	}
	fields := strings.Fields(out)
	if len(fields) != count {
		return nil, fmt.Errorf("expected %d answers, got %d (output: %q)", count, len(fields), out)
	}
	ans := make([]int64, count)
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

func buildTests(rng *rand.Rand) []testCase {
	tests := deterministicTests()
	total := 0
	for _, tc := range tests {
		total += tc.n
	}
	for len(tests) < 500 && total < maxTotalN {
		remaining := maxTotalN - total
		n := rng.Intn(min(1000, remaining-1)) + 2
		if n > remaining {
			n = remaining
		}
		s := randomString(rng, n)
		k := randK(rng, n)
		tests = append(tests, testCase{n: n, k: k, s: s})
		total += n
	}
	return tests
}

func deterministicTests() []testCase {
	return []testCase{
		{n: 4, k: 1, s: "1001"},
		{n: 4, k: 1_000_000_000, s: "1010"},
		{n: 4, k: 10, s: "1100"},
		{n: 6, k: 3, s: "001111"},
		{n: 6, k: 3, s: "111010"},
		{n: 5, k: 1, s: "01010"},
		{n: 5, k: 1, s: "10101"},
		{n: 10, k: 20, s: "1111111111"},
		{n: 10, k: 20, s: "0000000000"},
	}
}

func randomString(rng *rand.Rand, n int) string {
	b := make([]byte, n)
	mode := rng.Intn(5)
	for i := 0; i < n; i++ {
		switch mode {
		case 0:
			b[i] = '0'
		case 1:
			b[i] = '1'
		default:
			if rng.Intn(2) == 0 {
				b[i] = '0'
			} else {
				b[i] = '1'
			}
		}
	}
	return string(b)
}

func randK(rng *rand.Rand, n int) int64 {
	if rng.Intn(5) == 0 {
		return int64(rng.Intn(1_000_000_000) + 1)
	}
	maxScore := int64(n * (n - 1) / 2)
	if maxScore <= 0 {
		return 1
	}
	return int64(rng.Int63n(maxScore) + 1)
}

func serializeInput(tests []testCase) string {
	var sb strings.Builder
	sb.WriteString(strconv.Itoa(len(tests)))
	sb.WriteByte('\n')
	for _, tc := range tests {
		sb.WriteString(fmt.Sprintf("%d %d\n%s\n", tc.n, tc.k, tc.s))
	}
	return sb.String()
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
