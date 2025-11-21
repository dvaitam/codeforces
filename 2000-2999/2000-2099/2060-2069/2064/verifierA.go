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

const refSource = "2000-2999/2000-2099/2060-2069/2064/2064A.go"

type testCase struct {
	input string
}

type caseSpec struct {
	n int
	s string
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
	tmp, err := os.CreateTemp("", "2064A-ref-*")
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
	var tests []testCase
	rng := rand.New(rand.NewSource(20642064))

	tests = append(tests, buildInput([]caseSpec{
		{n: 1, s: "0"},
		{n: 1, s: "1"},
		{n: 5, s: "00110"},
		{n: 4, s: "1111"},
	}))

	tests = append(tests, buildInput([]caseSpec{
		{n: 6, s: "000000"},
		{n: 6, s: "111000"},
		{n: 6, s: "000111"},
		{n: 7, s: "1010101"},
	}))

	tests = append(tests, buildInput([]caseSpec{
		{n: 10, s: "0101010101"},
		{n: 10, s: "1111100000"},
		{n: 10, s: "0000011111"},
	}))

	for i := 0; i < 20; i++ {
		tests = append(tests, randomBatch(rng, 5, 40))
	}

	for i := 0; i < 10; i++ {
		tests = append(tests, randomBatch(rng, 10, 200))
	}

	tests = append(tests, randomBatch(rng, 15, 1000))

	return tests
}

func buildInput(cases []caseSpec) testCase {
	var b strings.Builder
	fmt.Fprintf(&b, "%d\n", len(cases))
	for _, cs := range cases {
		fmt.Fprintf(&b, "%d\n%s\n", cs.n, cs.s)
	}
	return testCase{input: b.String()}
}

func randomBatch(rng *rand.Rand, maxCases, maxN int) testCase {
	t := rng.Intn(maxCases) + 1
	var specs []caseSpec
	remaining := 1000
	for i := 0; i < t; i++ {
		minRemaining := t - i - 1
		maxLen := remaining - minRemaining
		if maxLen < 1 {
			maxLen = 1
		}
		if maxLen > maxN {
			maxLen = maxN
		}
		n := rng.Intn(maxLen) + 1
		remaining -= n
		specs = append(specs, caseSpec{n: n, s: randomBinary(rng, n)})
	}
	return buildInput(specs)
}

func randomBinary(rng *rand.Rand, n int) string {
	b := make([]byte, n)
	for i := 0; i < n; i++ {
		if rng.Intn(2) == 1 {
			b[i] = '1'
		} else {
			b[i] = '0'
		}
	}
	return string(b)
}
