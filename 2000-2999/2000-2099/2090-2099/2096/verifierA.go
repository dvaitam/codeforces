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

const refSource = "2000-2999/2000-2099/2090-2099/2096/2096A.go"

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
		got, err := runCandidate(candidate, tc.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "runtime error on test %d: %v\ninput:\n%soutput:\n%s\n", i+1, err, tc.input, got)
			os.Exit(1)
		}

		if ok, reason := checkValidity(tc.input, got); !ok {
			fmt.Fprintf(os.Stderr, "invalid output on test %d: %s\ninput:\n%s\noutput:\n%s\n", i+1, reason, tc.input, got)
			os.Exit(1)
		}
	}

	fmt.Printf("All %d tests passed.\n", len(tests))
}

func buildReference() (string, error) {
	tmp, err := os.CreateTemp("", "2096A-ref-*")
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

func checkValidity(rawInput, output string) (bool, string) {
	inTokens := strings.Fields(rawInput)
	if len(inTokens) == 0 {
		return false, "empty input tokens"
	}
	t, err := atoi(inTokens[0])
	if err != nil || t <= 0 {
		return false, "bad test count"
	}
	idx := 1
	outTokens := strings.Fields(output)
	pos := 0

	for caseIdx := 0; caseIdx < t; caseIdx++ {
		if idx >= len(inTokens) {
			return false, "missing n"
		}
		n, err := atoi(inTokens[idx])
		if err != nil || n < 2 {
			return false, "invalid n"
		}
		idx++
		if idx >= len(inTokens) {
			return false, "missing s"
		}
		s := inTokens[idx]
		idx++
		if len(s) != n-1 {
			return false, "bad s length"
		}

		if pos+n > len(outTokens) {
			return false, "not enough output numbers"
		}
		perm := make([]int, n)
		seen := make([]bool, n+1)
		for i := 0; i < n; i++ {
			val, err := atoi(outTokens[pos+i])
			if err != nil || val < 1 || val > n || seen[val] {
				return false, "output not a permutation"
			}
			seen[val] = true
			perm[i] = val
		}
		pos += n

		for i := 0; i < n-1; i++ {
			switch s[i] {
			case '<':
				for j := 0; j <= i; j++ {
					if perm[i+1] >= perm[j] {
						return false, "constraint < violated"
					}
				}
			case '>':
				for j := 0; j <= i; j++ {
					if perm[i+1] <= perm[j] {
						return false, "constraint > violated"
					}
				}
			default:
				return false, "invalid character in s"
			}
		}
	}

	if pos != len(outTokens) {
		return false, "extra output tokens"
	}
	if idx != len(inTokens) {
		return false, "unused input tokens"
	}
	return true, ""
}

func atoi(s string) (int, error) {
	var x int
	_, err := fmt.Sscan(s, &x)
	return x, err
}

func generateTests() []testCase {
	var tests []testCase
	rng := rand.New(rand.NewSource(20962096))

	tests = append(tests, buildInput([]caseSpec{
		{n: 2, s: "<"},
		{n: 2, s: ">"},
		{n: 3, s: "<<"},
		{n: 3, s: ">>"},
		{n: 3, s: "<>"},
		{n: 3, s: "><"},
	}))

	tests = append(tests, buildInput([]caseSpec{
		{n: 5, s: "<<<>"}, // mostly decreasing then jump
		{n: 5, s: ">><<"},
		{n: 6, s: "<<<<<"},
		{n: 6, s: ">>>>>"},
	}))

	for i := 0; i < 10; i++ {
		tests = append(tests, randomBatch(rng, 5, 20))
	}

	tests = append(tests, randomBatch(rng, 10, 100))

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
	for i := 0; i < t; i++ {
		n := rng.Intn(maxN-1) + 2
		sb := make([]byte, n-1)
		for j := range sb {
			if rng.Intn(2) == 0 {
				sb[j] = '<'
			} else {
				sb[j] = '>'
			}
		}
		specs = append(specs, caseSpec{n: n, s: string(sb)})
	}
	return buildInput(specs)
}
