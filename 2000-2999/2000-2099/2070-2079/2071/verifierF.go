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
	k int
	a []int
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierF.go /path/to/binary")
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
	tests := buildTestSuite(rng)
	input := serializeTests(tests)

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
			fmt.Fprintf(os.Stderr, "Mismatch on test %d: expected %d got %d\nInput:\n%s", i+1, expected[i], got[i], input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}

func buildReference() (string, error) {
	const refBin = "./ref_2071F.bin"
	cmd := exec.Command("go", "build", "-o", refBin, "2071F.go")
	if out, err := cmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("failed to build reference: %v\n%s", err, string(out))
	}
	return refBin, nil
}

func runAndParse(target, input string, tests int) ([]int64, error) {
	out, err := runProgram(target, input)
	if err != nil {
		return nil, err
	}
	return parseAnswers(out, tests)
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

func parseAnswers(out string, tests int) ([]int64, error) {
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

func buildTestSuite(rng *rand.Rand) []testCase {
	tests := deterministicTests()
	totalN := 0
	for _, tc := range tests {
		totalN += tc.n
	}

	for len(tests) < 200 && totalN < 200000 {
		n := rng.Intn(1000) + 1
		if totalN+n > 200000 {
			n = 200000 - totalN
		}
		k := rng.Intn(n)
		if rng.Intn(5) == 0 {
			k = n - 1
		}
		tc := randomTest(rng, n, k)
		tests = append(tests, tc)
		totalN += n
	}
	return tests
}

func deterministicTests() []testCase {
	return []testCase{
		{n: 1, k: 0, a: []int{1}},
		{n: 2, k: 1, a: []int{1, 1}},
		{n: 5, k: 0, a: []int{2, 1, 4, 5, 2}},
		{n: 5, k: 3, a: []int{2, 1, 4, 5, 2}},
		{n: 6, k: 1, a: []int{1, 2, 3, 4, 5, 11}},
		{n: 11, k: 6, a: []int{6, 3, 8, 5, 8, 3, 2, 1, 2, 7, 11}},
		{n: 14, k: 3, a: []int{3, 2, 3, 5, 5, 2, 6, 7, 4, 8, 10, 1, 8, 9}},
		{n: 2, k: 0, a: []int{1, 1000000000}},
		{n: 3, k: 1, a: []int{1000000000, 1, 1}},
	}
}

func randomTest(rng *rand.Rand, n, k int) testCase {
	a := make([]int, n)
	maxVal := 1_000_000_000
	for i := 0; i < n; i++ {
		switch rng.Intn(6) {
		case 0:
			a[i] = 1
		case 1:
			a[i] = maxVal
		default:
			a[i] = rng.Intn(maxVal) + 1
		}
	}
	return testCase{n: n, k: k, a: a}
}

func serializeTests(tests []testCase) string {
	var sb strings.Builder
	sb.WriteString(strconv.Itoa(len(tests)))
	sb.WriteByte('\n')
	for idx, tc := range tests {
		if idx > 0 {
			sb.WriteByte('\n')
		}
		sb.WriteString(fmt.Sprintf("%d %d\n", tc.n, tc.k))
		for i, v := range tc.a {
			if i > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(strconv.Itoa(v))
		}
		sb.WriteByte('\n')
	}
	return sb.String()
}
