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

const refSource = "./207B2.go"

type testCase struct {
	input string
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierB2.go /path/to/binary")
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
	tmp, err := os.CreateTemp("", "207B2-ref-*")
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

func equalTokens(a, b string) bool {
	ta := strings.Fields(a)
	tb := strings.Fields(b)
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
	rng := rand.New(rand.NewSource(207020702070))
	var tests []testCase

	tests = append(tests, sliceCase([]int{2, 1, 1}))
	tests = append(tests, sliceCase([]int{2, 2, 2, 2, 2}))
	tests = append(tests, sliceCase([]int{1, 1}))
	tests = append(tests, sliceCase([]int{250000, 1, 1, 1, 1}))
	tests = append(tests, sliceCase([]int{1, 250000, 1, 250000}))

	for i := 0; i < 10; i++ {
		n := rng.Intn(9) + 2 // 2..10
		tests = append(tests, randomCase(rng, n))
	}

	for i := 0; i < 10; i++ {
		n := rng.Intn(400) + 50 // 50..449
		tests = append(tests, randomCase(rng, n))
	}

	for i := 0; i < 5; i++ {
		n := rng.Intn(4000) + 1000 // 1000..4999
		tests = append(tests, randomCase(rng, n))
	}

	tests = append(tests, staircaseCase(10000))
	tests = append(tests, alternatingCase(20000))
	tests = append(tests, constantCase(250000, 250000))
	tests = append(tests, randomCase(rng, 200000))
	tests = append(tests, blockCase(250000, 120000))

	return tests
}

func sliceCase(a []int) testCase {
	var b strings.Builder
	fmt.Fprintln(&b, len(a))
	for i, v := range a {
		if i > 0 {
			b.WriteByte(' ')
		}
		fmt.Fprint(&b, v)
	}
	b.WriteByte('\n')
	return testCase{input: b.String()}
}

func randomCase(rng *rand.Rand, n int) testCase {
	var b strings.Builder
	b.Grow(n*7 + 20)
	fmt.Fprintln(&b, n)
	for i := 0; i < n; i++ {
		if i > 0 {
			b.WriteByte(' ')
		}
		fmt.Fprint(&b, rng.Intn(250000)+1)
	}
	b.WriteByte('\n')
	return testCase{input: b.String()}
}

func constantCase(n, value int) testCase {
	if value < 1 {
		value = 1
	} else if value > 250000 {
		value = 250000
	}
	var b strings.Builder
	b.Grow(n*7 + 20)
	fmt.Fprintln(&b, n)
	for i := 0; i < n; i++ {
		if i > 0 {
			b.WriteByte(' ')
		}
		fmt.Fprint(&b, value)
	}
	b.WriteByte('\n')
	return testCase{input: b.String()}
}

func staircaseCase(n int) testCase {
	if n < 2 {
		n = 2
	}
	var b strings.Builder
	b.Grow(n*7 + 20)
	fmt.Fprintln(&b, n)
	for i := 0; i < n; i++ {
		if i > 0 {
			b.WriteByte(' ')
		}
		val := i + 1
		if val > 250000 {
			val = 250000
		}
		fmt.Fprint(&b, val)
	}
	b.WriteByte('\n')
	return testCase{input: b.String()}
}

func alternatingCase(n int) testCase {
	if n < 2 {
		n = 2
	}
	var b strings.Builder
	b.Grow(n*7 + 20)
	fmt.Fprintln(&b, n)
	for i := 0; i < n; i++ {
		if i > 0 {
			b.WriteByte(' ')
		}
		if i%2 == 0 {
			fmt.Fprint(&b, 1)
		} else {
			fmt.Fprint(&b, 250000)
		}
	}
	b.WriteByte('\n')
	return testCase{input: b.String()}
}

func blockCase(n, pivot int) testCase {
	if n < 2 {
		n = 2
	}
	if pivot < 0 {
		pivot = 0
	} else if pivot > n {
		pivot = n
	}
	var b strings.Builder
	b.Grow(n*7 + 20)
	fmt.Fprintln(&b, n)
	for i := 0; i < n; i++ {
		if i > 0 {
			b.WriteByte(' ')
		}
		if i < pivot {
			fmt.Fprint(&b, 1)
		} else {
			fmt.Fprint(&b, 250000)
		}
	}
	b.WriteByte('\n')
	return testCase{input: b.String()}
}
