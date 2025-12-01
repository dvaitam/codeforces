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
)

const refSource = "./2165A.go"

type testCase struct {
	arr []int64
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}

	refBin, err := buildReference()
	if err != nil {
		fmt.Fprintln(os.Stderr, "failed to build reference solution:", err)
		os.Exit(1)
	}
	defer os.Remove(refBin)

	tests := generateTests()
	input := buildInput(tests)

	expectedOut, err := runProgram(refBin, input)
	if err != nil {
		fmt.Fprintf(os.Stderr, "reference execution failed: %v\noutput:\n%s\n", err, expectedOut)
		os.Exit(1)
	}
	candOut, err := runCandidate(os.Args[1], input)
	if err != nil {
		fmt.Fprintf(os.Stderr, "candidate execution failed: %v\noutput:\n%s\n", err, candOut)
		os.Exit(1)
	}

	expected, err := parseOutputs(expectedOut, len(tests))
	if err != nil {
		fmt.Fprintf(os.Stderr, "could not parse reference output: %v\n", err)
		os.Exit(1)
	}

	got, err := parseOutputs(candOut, len(tests))
	if err != nil {
		fmt.Fprintf(os.Stderr, "could not parse candidate output: %v\n", err)
		os.Exit(1)
	}

	for i := range tests {
		if expected[i] != got[i] {
			fmt.Fprintf(os.Stderr, "wrong answer on test %d: expected %d, got %d\n", i+1, expected[i], got[i])
			os.Exit(1)
		}
	}

	fmt.Printf("All %d tests passed.\n", len(tests))
}

func buildReference() (string, error) {
	tmp, err := os.CreateTemp("", "2165A-ref-*")
	if err != nil {
		return "", err
	}
	tmp.Close()

	cmd := exec.Command("go", "build", "-o", tmp.Name(), filepath.Clean(refSource))
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		os.Remove(tmp.Name())
		return "", fmt.Errorf("%v\n%s", err, out.String())
	}
	return tmp.Name(), nil
}

func runCandidate(path, input string) (string, error) {
	cmd := commandFor(path)
	return runWithInput(cmd, input)
}

func runProgram(path, input string) (string, error) {
	cmd := exec.Command(path)
	return runWithInput(cmd, input)
}

func commandFor(path string) *exec.Cmd {
	switch filepath.Ext(path) {
	case ".go":
		return exec.Command("go", "run", path)
	case ".py":
		return exec.Command("python3", path)
	default:
		return exec.Command(path)
	}
}

func runWithInput(cmd *exec.Cmd, input string) (string, error) {
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return out.String(), err
}

func parseOutputs(output string, t int) ([]int64, error) {
	tokens := strings.Fields(output)
	if len(tokens) < t {
		return nil, fmt.Errorf("expected %d numbers but found %d", t, len(tokens))
	}
	if len(tokens) > t {
		return nil, fmt.Errorf("extra output detected starting with %q", tokens[t])
	}
	ans := make([]int64, t)
	for i := 0; i < t; i++ {
		val, err := strconv.ParseInt(tokens[i], 10, 64)
		if err != nil {
			return nil, fmt.Errorf("token %q is not an integer", tokens[i])
		}
		ans[i] = val
	}
	return ans, nil
}

func generateTests() []testCase {
	const limit = 200000
	rng := rand.New(rand.NewSource(21652065))
	var tests []testCase
	total := 0

	add := func(arr []int64) {
		if len(arr) < 2 {
			return
		}
		if total+len(arr) > limit {
			return
		}
		snapshot := append([]int64(nil), arr...)
		tests = append(tests, testCase{arr: snapshot})
		total += len(arr)
	}

	add([]int64{1, 1, 3, 2})
	add([]int64{20, 27})
	add([]int64{1, 1, 4, 5, 1, 4, 1})

	add([]int64{0, 0})
	add([]int64{5, 5})
	add([]int64{0, 1, 0})
	add([]int64{7, 3, 7, 3})
	add([]int64{1000000000, 0, 1000000000})

	for i := 0; i < 60; i++ {
		n := 2 + rng.Intn(10)
		add(randomCase(rng, n))
		if total >= limit {
			return tests
		}
	}

	for i := 0; i < 40; i++ {
		n := 12 + rng.Intn(40)
		add(randomCase(rng, n))
		if total >= limit {
			return tests
		}
	}

	for i := 0; i < 20; i++ {
		n := 50 + rng.Intn(150)
		add(randomCase(rng, n))
		if total >= limit {
			return tests
		}
	}

	for i := 0; i < 10; i++ {
		n := 200 + rng.Intn(800)
		add(randomCase(rng, n))
		if total >= limit {
			return tests
		}
	}

	remaining := limit - total
	if remaining >= 2 {
		add(randomCase(rng, remaining))
	}

	return tests
}

func randomCase(rng *rand.Rand, n int) []int64 {
	arr := make([]int64, n)
	for i := range arr {
		arr[i] = rng.Int63n(1_000_000_001)
	}
	return arr
}

func buildInput(tests []testCase) string {
	var b strings.Builder
	fmt.Fprintf(&b, "%d\n", len(tests))
	for _, tc := range tests {
		fmt.Fprintf(&b, "%d\n", len(tc.arr))
		writeArray(&b, tc.arr)
	}
	return b.String()
}

func writeArray(b *strings.Builder, arr []int64) {
	for i, v := range arr {
		if i > 0 {
			b.WriteByte(' ')
		}
		fmt.Fprintf(b, "%d", v)
	}
	b.WriteByte('\n')
}
