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

// refSource points to the local reference solution to avoid GOPATH resolution.
const refSource = "542F.go"

type testCase struct {
	input string
}

type task struct {
	t int
	q int
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierF.go /path/to/binary")
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
	tmp, err := os.CreateTemp("", "542F-ref-*")
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
	if strings.HasSuffix(path, ".go") {
		return exec.Command("go", "run", path)
	}
	return exec.Command(path)
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
	rng := rand.New(rand.NewSource(542542))
	var tests []testCase

	tests = append(tests, sampleTest())
	tests = append(tests, makeTest(1, 1, []task{{1, 1}}))
	tests = append(tests, makeTest(1, 2, []task{{2, 10}}))

	for i := 0; i < 40; i++ {
		tests = append(tests, randomCase(rng, rng.Intn(5)+1, 50))
	}

	tests = append(tests, limitCase())

	return tests
}

func sampleTest() testCase {
	return makeTest(5, 5, []task{{1, 1}, {1, 1}, {2, 2}, {3, 3}, {4, 4}})
}

func makeTest(n int, T int, items []task) testCase {
	var b strings.Builder
	fmt.Fprintf(&b, "%d %d\n", n, T)
	for _, it := range items {
		fmt.Fprintf(&b, "%d %d\n", it.t, it.q)
	}
	return testCase{input: b.String()}
}

func randomCase(rng *rand.Rand, n int, T int) testCase {
	if n < 1 {
		n = 1
	}
	if T < 1 {
		T = 1
	}
	if n > 1000 {
		n = 1000
	}
	if T > 100 {
		T = 100
	}

	items := make([]task, n)
	for i := 0; i < n; i++ {
		items[i] = task{
			t: rng.Intn(T) + 1,
			q: rng.Intn(1000) + 1,
		}
	}
	return makeTest(n, T, items)
}

func limitCase() testCase {
	n := 1000
	T := 100
	items := make([]task, n)
	for i := 0; i < n; i++ {
		items[i] = task{t: 1, q: 1000}
	}
	return makeTest(n, T, items)
}
