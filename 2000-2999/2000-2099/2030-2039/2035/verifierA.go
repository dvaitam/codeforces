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

const refSource = "./2035A.go"

type testCase struct {
	input string
}

type testInstance struct {
	n int64
	m int64
	r int64
	c int64
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierA.go /path/to/binary")
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
	tmp, err := os.CreateTemp("", "2035A-ref-*")
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
	rng := rand.New(rand.NewSource(20352035))
	var tests []testCase

	tests = append(tests, sampleTest())
	tests = append(tests, makeTest([]testInstance{
		{n: 1, m: 1, r: 1, c: 1},
		{n: 2, m: 3, r: 1, c: 2},
		{n: 5, m: 5, r: 5, c: 5},
	}))

	for i := 0; i < 30; i++ {
		tests = append(tests, randomCase(rng, rng.Intn(5)+1))
	}

	tests = append(tests, edgeCase())

	return tests
}

func sampleTest() testCase {
	return makeTest([]testInstance{
		{2, 3, 1, 2},
		{2, 2, 2, 1},
		{1, 1, 1, 1},
		{1_000_000, 1_000_000, 1, 1},
	})
}

func makeTest(instances []testInstance) testCase {
	var b strings.Builder
	fmt.Fprintln(&b, len(instances))
	for _, inst := range instances {
		fmt.Fprintf(&b, "%d %d %d %d\n", inst.n, inst.m, inst.r, inst.c)
	}
	return testCase{input: b.String()}
}

func randomCase(rng *rand.Rand, maxCases int) testCase {
	if maxCases < 1 {
		maxCases = 1
	}
	t := rng.Intn(maxCases) + 1
	inst := make([]testInstance, t)
	for i := 0; i < t; i++ {
		n := rng.Int63n(1_000_000) + 1
		m := rng.Int63n(1_000_000) + 1
		r := rng.Int63n(n) + 1
		c := rng.Int63n(m) + 1
		inst[i] = testInstance{n: n, m: m, r: r, c: c}
	}
	return makeTest(inst)
}

func edgeCase() testCase {
	return makeTest([]testInstance{
		{n: 1_000_000, m: 1_000_000, r: 1_000_000, c: 1_000_000},
		{n: 1_000_000, m: 1_000_000, r: 1, c: 1_000_000},
		{n: 1_000_000, m: 1_000_000, r: 1_000_000, c: 1},
	})
}
