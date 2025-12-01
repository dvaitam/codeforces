package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

const refSource = "./2087G.go"

type testCase struct {
	input string
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierG.go /path/to/binary")
		os.Exit(1)
	}

	refBin, err := buildReference()
	if err != nil {
		fmt.Fprintln(os.Stderr, "failed to build reference:", err)
		os.Exit(1)
	}
	defer os.Remove(refBin)

	candidate := os.Args[1]
	tests := generateTests()

	for i, tc := range tests {
		expect, err := runProgram(refBin, tc.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "reference failed on test %d: %v\ninput:\n%s", i+1, err, tc.input)
			os.Exit(1)
		}

		got, err := runCandidate(candidate, tc.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "runtime error on test %d: %v\ninput:\n%soutput:\n%s\n", i+1, err, tc.input, got)
			os.Exit(1)
		}

		if !equalTokens(expect, got) {
			fmt.Fprintf(os.Stderr, "wrong answer on test %d\ninput:\n%s\nexpected:\n%s\ngot:\n%s\n", i+1, tc.input, expect, got)
			os.Exit(1)
		}
	}

	fmt.Printf("All %d tests passed.\n", len(tests))
}

func buildReference() (string, error) {
	tmp, err := os.CreateTemp("", "2087G-ref-*")
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

func commandFor(path string) *exec.Cmd {
	if strings.HasSuffix(path, ".go") {
		return exec.Command("go", "run", path)
	}
	return exec.Command(path)
}

func runCandidate(path, input string) (string, error) {
	cmd := commandFor(path)
	return runWithInput(cmd, input)
}

func runProgram(path, input string) (string, error) {
	cmd := exec.Command(path)
	return runWithInput(cmd, input)
}

func runWithInput(cmd *exec.Cmd, input string) (string, error) {
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return strings.TrimSpace(out.String()), err
}

func equalTokens(expected, got string) bool {
	ta := strings.Fields(expected)
	tb := strings.Fields(got)
	if len(ta) != len(tb) {
		return false
	}
	for i := range ta {
		if ta[i] != tb[i] {
			return false
		}
	}
	return true
}

func generateTests() []testCase {
	rng := rand.New(rand.NewSource(20872087))
	var tests []testCase

	// Simple and hand-crafted cases.
	tests = append(tests, makeTest([]int64{0}))
	tests = append(tests, makeTest([]int64{5}))
	tests = append(tests, makeTest([]int64{0, 0, 0, 0, 0}))
	tests = append(tests, makeTest([]int64{5, 4, 3, 2, 1, 0}))
	tests = append(tests, makeTest([]int64{5, 9, 0, 3, 2, 0, 2, 2, 9, 9}))
	tests = append(tests, makeTest([]int64{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}))

	// Small random arrays to check combinational choices.
	for i := 0; i < 15; i++ {
		n := rng.Intn(30) + 1
		tests = append(tests, makeTest(randomArray(rng, n, 100)))
	}

	// Medium random arrays with wider values.
	for i := 0; i < 5; i++ {
		n := rng.Intn(800) + 200
		tests = append(tests, makeTest(randomArray(rng, n, 1_000_000)))
	}

	// Structured mid-sized cases.
	tests = append(tests, makeTest(increasingArray(1000, 3, 1_000_000)))
	tests = append(tests, makeTest(alternatingArray(1000, 1_000_000)))
	tests = append(tests, makeTest(blockArray(rng, 1200)))

	// Large edge-focused cases.
	tests = append(tests, makeTest(fillArray(100000, 0)))
	tests = append(tests, makeTest(randomArray(rng, 200000, 1_000_000)))

	return tests
}

func makeTest(a []int64) testCase {
	return testCase{input: buildInput(a)}
}

func buildInput(a []int64) string {
	var b strings.Builder
	fmt.Fprintf(&b, "%d\n", len(a))
	for i, v := range a {
		if i > 0 {
			b.WriteByte(' ')
		}
		fmt.Fprintf(&b, "%d", v)
	}
	b.WriteByte('\n')
	return b.String()
}

func randomArray(rng *rand.Rand, n int, maxVal int64) []int64 {
	if maxVal <= 0 {
		maxVal = 1
	}
	a := make([]int64, n)
	for i := 0; i < n; i++ {
		a[i] = rng.Int63n(maxVal + 1)
	}
	return a
}

func increasingArray(n int, step int64, capVal int64) []int64 {
	a := make([]int64, n)
	cur := int64(0)
	for i := 0; i < n; i++ {
		a[i] = cur
		if cur+step <= capVal {
			cur += step
		}
	}
	return a
}

func alternatingArray(n int, high int64) []int64 {
	a := make([]int64, n)
	for i := 0; i < n; i++ {
		if i%2 == 0 {
			a[i] = high
		} else {
			a[i] = 0
		}
	}
	return a
}

func blockArray(rng *rand.Rand, n int) []int64 {
	a := make([]int64, n)
	cur := int64(rng.Intn(50))
	i := 0
	for i < n {
		block := rng.Intn(40) + 1
		val := cur + int64(rng.Intn(7)) - 3
		if val < 0 {
			val = 0
		}
		if val > 1_000_000 {
			val = 1_000_000
		}
		for j := 0; j < block && i < n; j++ {
			a[i] = val
			i++
		}
		cur = val
	}
	return a
}

func fillArray(n int, val int64) []int64 {
	a := make([]int64, n)
	for i := range a {
		a[i] = val
	}
	return a
}
