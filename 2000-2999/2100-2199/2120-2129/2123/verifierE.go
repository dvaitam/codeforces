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

const refSource = "./2123E.go"
const maxN = 200000

type testCase struct {
	input string
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierE.go /path/to/binary")
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
	tmp, err := os.CreateTemp("", "2123E-ref-*")
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
	rng := rand.New(rand.NewSource(21232123))

	tests = append(tests, sampleTest())

	tests = append(tests, buildInput([]testSpec{
		{n: 1, a: []int{0}},
		{n: 2, a: []int{0, 0}},
		{n: 3, a: []int{1, 2, 3}},
	}))

	tests = append(tests, buildInput([]testSpec{
		{n: 5, a: []int{0, 1, 0, 1, 2}},
		{n: 6, a: []int{3, 2, 0, 4, 5, 1}},
	}))

	for i := 0; i < 10; i++ {
		tests = append(tests, randomBatch(rng, 20, 2000))
	}
	tests = append(tests, randomBatch(rng, 50, maxN))

	return tests
}

type testSpec struct {
	n int
	a []int
}

func sampleTest() testCase {
	return testCase{
		input: "5\n" +
			"5\n1 0 0 1 2\n" +
			"6\n3 2 0 4 5 1\n" +
			"6\n1 2 0 1 3 2\n" +
			"5\n0 3 4 1 5\n" +
			"5\n0 0 0 0 0\n",
	}
}

func buildInput(specs []testSpec) testCase {
	var b strings.Builder
	fmt.Fprintf(&b, "%d\n", len(specs))
	for _, ts := range specs {
		fmt.Fprintf(&b, "%d\n", ts.n)
		for i, v := range ts.a {
			if i > 0 {
				b.WriteByte(' ')
			}
			fmt.Fprintf(&b, "%d", v)
		}
		b.WriteByte('\n')
	}
	return testCase{input: b.String()}
}

func randomBatch(rng *rand.Rand, maxCases, maxNPerCase int) testCase {
	t := rng.Intn(maxCases) + 1
	remaining := maxN
	var specs []testSpec
	for i := 0; i < t; i++ {
		minRemaining := t - i - 1
		maxAllowed := remaining - minRemaining
		if maxAllowed < 1 {
			break
		}
		if maxAllowed > maxNPerCase {
			maxAllowed = maxNPerCase
		}
		n := rng.Intn(maxAllowed) + 1
		remaining -= n
		arr := make([]int, n)
		for j := 0; j < n; j++ {
			arr[j] = rng.Intn(n + 1)
		}
		specs = append(specs, testSpec{n: n, a: arr})
	}
	if len(specs) == 0 {
		specs = append(specs, testSpec{n: 1, a: []int{0}})
	}
	return buildInput(specs)
}
