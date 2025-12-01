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

const refSource = "./2065G.go"
const maxN = 200000

type testCase struct {
	input string
}

type caseSpec struct {
	n   int
	arr []int
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
	tmp, err := os.CreateTemp("", "2065G-ref-*")
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
	rng := rand.New(rand.NewSource(20652065))

	tests = append(tests, sampleTest())

	tests = append(tests, buildInput([]caseSpec{
		{n: 2, arr: []int{2, 2}},
		{n: 3, arr: []int{2, 3, 2}},
		{n: 5, arr: []int{2, 4, 6, 3, 5}},
	}))

	tests = append(tests, buildInput([]caseSpec{
		{n: 6, arr: []int{6, 6, 6, 6, 6, 6}},
		{n: 6, arr: []int{2, 3, 5, 7, 11, 13}},
	}))

	tests = append(tests, buildInput([]caseSpec{
		{n: 8, arr: []int{2, 2, 3, 3, 4, 4, 5, 5}},
	}))

	for i := 0; i < 15; i++ {
		tests = append(tests, randomBatch(rng, 5, 2000))
	}

	for i := 0; i < 8; i++ {
		tests = append(tests, randomBatch(rng, 8, 20000))
	}

	tests = append(tests, randomBatch(rng, 6, maxN))

	return tests
}

func sampleTest() testCase {
	return buildInput([]caseSpec{
		{n: 4, arr: []int{2, 3, 4, 6}},
		{n: 6, arr: []int{2, 3, 4, 5, 6, 9}},
		{n: 9, arr: []int{2, 4, 5, 7, 8, 9, 3, 5, 5}},
	})
}

func buildInput(cases []caseSpec) testCase {
	var b strings.Builder
	fmt.Fprintf(&b, "%d\n", len(cases))
	for _, cs := range cases {
		fmt.Fprintf(&b, "%d\n", cs.n)
		for i, v := range cs.arr {
			if i > 0 {
				b.WriteByte(' ')
			}
			fmt.Fprintf(&b, "%d", v)
		}
		b.WriteByte('\n')
	}
	return testCase{input: b.String()}
}

func randomBatch(rng *rand.Rand, maxCases, maxLen int) testCase {
	t := rng.Intn(maxCases) + 1
	var specs []caseSpec
	remaining := maxN
	for i := 0; i < t; i++ {
		minRemaining := t - i - 1
		maxNForCase := remaining - minRemaining
		if maxNForCase < 2 {
			maxNForCase = 2
		}
		if maxNForCase > maxLen {
			maxNForCase = maxLen
		}
		n := rng.Intn(maxNForCase-1) + 2
		remaining -= n
		arr := make([]int, n)
		for j := 0; j < n; j++ {
			arr[j] = rng.Intn(n-1) + 2
		}
		specs = append(specs, caseSpec{n: n, arr: arr})
	}
	return buildInput(specs)
}
