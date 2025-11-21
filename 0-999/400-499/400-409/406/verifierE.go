package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

const (
	refSourceE = "406E.go"
	refBinaryE = "ref406E.bin"
	totalTests = 60
)

type testCase struct {
	n int64
	m int
	s []int
	f []int64
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	candidate := os.Args[1]

	refPath, err := buildReference()
	if err != nil {
		fmt.Println("failed to build reference:", err)
		os.Exit(1)
	}
	defer os.Remove(refPath)

	tests := generateTests()
	for idx, tc := range tests {
		input := formatInput(tc)

		refOut, err := runProgram(refPath, input)
		if err != nil {
			fmt.Printf("reference runtime error on test %d: %v\n", idx+1, err)
			printInput(input)
			os.Exit(1)
		}
		candOut, err := runProgram(candidate, input)
		if err != nil {
			fmt.Printf("candidate runtime error on test %d: %v\n", idx+1, err)
			printInput(input)
			os.Exit(1)
		}

		refVal, err := parseOutput(refOut)
		if err != nil {
			fmt.Printf("failed to parse reference output on test %d: %v\noutput:\n%s\n", idx+1, err, refOut)
			printInput(input)
			os.Exit(1)
		}
		candVal, err := parseOutput(candOut)
		if err != nil {
			fmt.Printf("failed to parse candidate output on test %d: %v\noutput:\n%s\n", idx+1, err, candOut)
			printInput(input)
			os.Exit(1)
		}

		if refVal != candVal {
			fmt.Printf("test %d failed: expected %d, got %d\n", idx+1, refVal, candVal)
			printInput(input)
			fmt.Println("Reference output:")
			fmt.Println(refOut)
			fmt.Println("Candidate output:")
			fmt.Println(candOut)
			os.Exit(1)
		}
	}

	fmt.Printf("All %d tests passed\n", len(tests))
}

func buildReference() (string, error) {
	cmd := exec.Command("go", "build", "-o", refBinaryE, refSourceE)
	if out, err := cmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("%v\n%s", err, string(out))
	}
	return filepath.Join(".", refBinaryE), nil
}

func runProgram(path string, input []byte) (string, error) {
	cmd := exec.Command(path)
	cmd.Stdin = bytes.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	return out.String(), err
}

func formatInput(tc testCase) []byte {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d\n", tc.n, tc.m))
	for i := 0; i < tc.m; i++ {
		sb.WriteString(fmt.Sprintf("%d %d\n", tc.s[i], tc.f[i]))
	}
	return []byte(sb.String())
}

func parseOutput(out string) (int64, error) {
	fields := strings.Fields(out)
	if len(fields) != 1 {
		return 0, fmt.Errorf("expected 1 integer, got %d", len(fields))
	}
	val, err := strconv.ParseInt(fields[0], 10, 64)
	if err != nil {
		return 0, fmt.Errorf("invalid integer %q: %v", fields[0], err)
	}
	return val, nil
}

func generateTests() []testCase {
	tests := []testCase{
		buildTestCase(1, []int{0, 1, 0}, []int64{1, 1, 1}),
		buildTestCase(5, []int{0, 0, 1}, []int64{3, 1, 4}),
		buildTestCase(10, []int{0, 1, 1, 0}, []int64{1, 10, 5, 7}),
	}

	rnd := rand.New(rand.NewSource(time.Now().UnixNano()))

	addRandom := func(m int, maxN int64) {
		ops := m
		s := make([]int, ops)
		f := make([]int64, ops)
		for i := 0; i < ops; i++ {
			s[i] = rnd.Intn(2)
			f[i] = randInt64(rnd, 1, maxN)
		}
		tests = append(tests, buildTestCase(maxN, s, f))
	}

	for len(tests) < totalTests-3 {
		m := rnd.Intn(500) + 3
		n := randInt64(rnd, 1, 1_000_000_000)
		s := make([]int, m)
		f := make([]int64, m)
		for i := 0; i < m; i++ {
			s[i] = rnd.Intn(2)
			f[i] = randInt64(rnd, 1, n)
		}
		tests = append(tests, buildTestCase(n, s, f))
	}

	addRandom(1000, 1_000_000_000)
	addRandom(100000, 1_000_000_000)
	addRandom(100000, 1_000_000_000)

	return tests
}

func randInt64(rnd *rand.Rand, lo, hi int64) int64 {
	if lo == hi {
		return lo
	}
	return lo + rnd.Int63n(hi-lo+1)
}

func buildTestCase(n int64, s []int, f []int64) testCase {
	return testCase{
		n: n,
		m: len(s),
		s: append([]int(nil), s...),
		f: append([]int64(nil), f...),
	}
}

func printInput(in []byte) {
	fmt.Println("Input used:")
	fmt.Println(string(in))
}
