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

const refSource = "./2117E.go"
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
	tmp, err := os.CreateTemp("", "2117E-ref-*")
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
	rng := rand.New(rand.NewSource(21172117))

	tests = append(tests, buildInput([][]int{
		{2, 1, 1, 2, 2},
	}))

	tests = append(tests, buildInput([][]int{
		{3, 1, 2, 3, 3, 2, 1},
		{4, 1, 1, 1, 1, 2, 2, 2, 2},
	}))

	tests = append(tests, buildInput([][]int{
		{5, 1, 2, 3, 4, 5, 5, 4, 3, 2, 1},
	}))

	for i := 0; i < 8; i++ {
		tests = append(tests, randomBatch(rng, 50, 2000))
	}

	tests = append(tests, randomBatch(rng, 100, maxN))

	return tests
}

// For each test slice: first element n, next n are a, next n are b.
func buildInput(cases [][]int) testCase {
	var b strings.Builder
	fmt.Fprintf(&b, "%d\n", len(cases))
	for _, cs := range cases {
		if len(cs) < 1 {
			continue
		}
		n := cs[0]
		if 1+n+n != len(cs) {
			continue
		}
		fmt.Fprintf(&b, "%d\n", n)
		for i := 1; i <= n; i++ {
			if i > 1 {
				b.WriteByte(' ')
			}
			fmt.Fprintf(&b, "%d", cs[i])
		}
		b.WriteByte('\n')
		for i := 1 + n; i < 1+2*n; i++ {
			if i > 1+n {
				b.WriteByte(' ')
			}
			fmt.Fprintf(&b, "%d", cs[i])
		}
		b.WriteByte('\n')
	}
	return testCase{input: b.String()}
}

func randomBatch(rng *rand.Rand, maxCases int, maxNPerCase int) testCase {
	t := rng.Intn(maxCases) + 1
	var cases [][]int
	totalN := 0
	for i := 0; i < t; i++ {
		remaining := maxN - totalN - (t - i - 1)
		if remaining < 2 {
			break
		}
		nCap := maxNPerCase
		if remaining < nCap {
			nCap = remaining
		}
		n := rng.Intn(nCap-1) + 2
		totalN += n
		arr := make([]int, 1+2*n)
		arr[0] = n
		for j := 0; j < n; j++ {
			arr[1+j] = rng.Intn(n) + 1
		}
		for j := 0; j < n; j++ {
			arr[1+n+j] = rng.Intn(n) + 1
		}
		cases = append(cases, arr)
	}
	return buildInput(cases)
}
