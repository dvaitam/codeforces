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
	n int64
	k int64
}

const maxTotalTests = 5000

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

	for i := range expected {
		if expected[i] != got[i] {
			fmt.Fprintf(os.Stderr, "Mismatch on test %d: expected %d got %d\n", i+1, expected[i], got[i])
			fmt.Fprintf(os.Stderr, "n=%d k=%d\n", tests[i].n, tests[i].k)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}

func buildReference() (string, error) {
	const refName = "./ref_2053C.bin"
	cmd := exec.Command("go", "build", "-o", refName, "2053C.go")
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
	res := make([]int64, tests)
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
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for len(tests) < maxTotalTests {
		n := randomInt(rng, 1, 2_000_000_000)
		k := randomInt(rng, 1, n)
		tests = append(tests, testCase{n: n, k: k})
	}
	return tests
}

func deterministicTests() []testCase {
	return []testCase{
		{n: 1, k: 1},
		{n: 2, k: 1},
		{n: 2, k: 2},
		{n: 3, k: 1},
		{n: 3, k: 2},
		{n: 3, k: 3},
		{n: 7, k: 2},
		{n: 11, k: 3},
		{n: 55801, k: 689},
		{n: 9, k: 9},
		{n: 1_000_000_000, k: 1},
		{n: 2_000_000_000, k: 1},
		{n: 2_000_000_000, k: 2_000_000_000},
		{n: 1_234_567_890, k: 987_654_321},
	}
}

func randomInt(rng *rand.Rand, lo, hi int64) int64 {
	if lo == hi {
		return lo
	}
	return lo + rng.Int63n(hi-lo+1)
}

func serializeInput(tests []testCase) string {
	var sb strings.Builder
	sb.WriteString(strconv.Itoa(len(tests)))
	sb.WriteByte('\n')
	for _, tc := range tests {
		sb.WriteString(fmt.Sprintf("%d %d\n", tc.n, tc.k))
	}
	return sb.String()
}
