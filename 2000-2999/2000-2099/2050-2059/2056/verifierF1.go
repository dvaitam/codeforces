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

const refSource = "2056F1.go"

type testCase struct {
	input string
}

type caseSpec struct {
	k int
	m int
	n string
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierF1.go /path/to/binary")
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
	tmp, err := os.CreateTemp("", "2056F1-ref-*")
	if err != nil {
		return "", err
	}
	tmp.Close()

	source := filepath.Join(".", refSource)
	cmd := exec.Command("go", "build", "-o", tmp.Name(), source)
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
	rng := rand.New(rand.NewSource(20562056))

	tests = append(tests, buildInput([]caseSpec{
		{k: 1, m: 1, n: "1"},
		{k: 2, m: 3, n: "10"},
		{k: 3, m: 4, n: "101"},
	}))

	tests = append(tests, buildInput([]caseSpec{
		{k: 4, m: 6, n: "1111"},
		{k: 5, m: 5, n: "10001"},
	}))

	tests = append(tests, buildInput([]caseSpec{
		{k: 6, m: 7, n: "110110"},
	}))

	tests = append(tests, buildInput([]caseSpec{
		{k: 8, m: 9, n: "10000001"},
		{k: 7, m: 15, n: "1111111"},
	}))

	tests = append(tests, buildInput([]caseSpec{
		{k: 200, m: 500, n: oneFollowedByZeros(200)},
	}))

	tests = append(tests, buildInput([]caseSpec{
		{k: 200, m: 499, n: strings.Repeat("1", 200)},
	}))

	for i := 0; i < 20; i++ {
		tests = append(tests, randomBatch(rng, 4))
	}

	for i := 0; i < 10; i++ {
		tests = append(tests, randomBatch(rng, 8))
	}

	tests = append(tests, mixedEdgeBatch(rng))

	return tests
}

func buildInput(cases []caseSpec) testCase {
	var b strings.Builder
	fmt.Fprintf(&b, "%d\n", len(cases))
	for _, cs := range cases {
		fmt.Fprintf(&b, "%d %d\n%s\n", cs.k, cs.m, cs.n)
	}
	return testCase{input: b.String()}
}

func randomBatch(rng *rand.Rand, maxCases int) testCase {
	t := rng.Intn(maxCases) + 1
	remaining := 200
	specs := make([]caseSpec, 0, t)
	for i := 0; i < t; i++ {
		minRemaining := t - i - 1
		maxK := remaining - minRemaining
		if maxK <= 0 {
			maxK = 1
		}
		k := rng.Intn(maxK) + 1
		remaining -= k
		m := rng.Intn(500) + 1
		specs = append(specs, caseSpec{k: k, m: m, n: randomBinary(rng, k)})
	}
	return buildInput(specs)
}

func mixedEdgeBatch(rng *rand.Rand) testCase {
	specs := []caseSpec{
		{k: 10, m: 1, n: "1" + strings.Repeat("0", 9)},
		{k: 12, m: 250, n: strings.Repeat("1", 12)},
		{k: 25, m: 500, n: patternedBinary(25)},
		{k: 30, m: 123, n: randomBinary(rng, 30)},
	}
	return buildInput(specs)
}

func randomBinary(rng *rand.Rand, k int) string {
	if k <= 0 {
		return ""
	}
	b := make([]byte, k)
	b[0] = '1'
	for i := 1; i < k; i++ {
		if rng.Intn(2) == 1 {
			b[i] = '1'
		} else {
			b[i] = '0'
		}
	}
	return string(b)
}

func oneFollowedByZeros(k int) string {
	if k <= 0 {
		return ""
	}
	return "1" + strings.Repeat("0", k-1)
}

func patternedBinary(k int) string {
	if k <= 0 {
		return ""
	}
	b := make([]byte, k)
	for i := 0; i < k; i++ {
		if i%3 == 0 {
			b[i] = '1'
		} else {
			b[i] = '0'
		}
	}
	b[0] = '1'
	return string(b)
}
