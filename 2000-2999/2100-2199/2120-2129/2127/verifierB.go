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

const refSource = "./2127B.go"
const maxN = 200000

type testInput struct {
	raw string
	t   int
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierB.go /path/to/binary")
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

	for i, ti := range tests {
		expect, err := runProgram(refBin, ti.raw)
		if err != nil {
			fmt.Fprintf(os.Stderr, "reference failed on test %d: %v\ninput:\n%s", i+1, err, ti.raw)
			os.Exit(1)
		}
		got, err := runCandidate(candidate, ti.raw)
		if err != nil {
			fmt.Fprintf(os.Stderr, "runtime error on test %d: %v\ninput:\n%soutput:\n%s\n", i+1, err, ti.raw, got)
			os.Exit(1)
		}

		if !validOutput(ti.t, got) {
			fmt.Fprintf(os.Stderr, "invalid output format on test %d\ninput:\n%s\noutput:\n%s\n", i+1, ti.raw, got)
			os.Exit(1)
		}

		if !equalTokens(expect, got) {
			fmt.Fprintf(os.Stderr, "wrong answer on test %d\ninput:\n%s\nexpected:\n%s\ngot:\n%s\n", i+1, ti.raw, expect, got)
			os.Exit(1)
		}
	}

	fmt.Printf("All %d tests passed.\n", len(tests))
}

func buildReference() (string, error) {
	tmp, err := os.CreateTemp("", "2127B-ref-*")
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

func validOutput(t int, output string) bool {
	tokens := strings.Fields(output)
	if len(tokens) != t {
		return false
	}
	for _, tok := range tokens {
		if _, err := atoi(tok); err != nil {
			return false
		}
	}
	return true
}

func atoi(s string) (int, error) {
	var x int
	_, err := fmt.Sscan(s, &x)
	return x, err
}

func generateTests() []testInput {
	var tests []testInput
	rng := rand.New(rand.NewSource(21272127))

	tests = append(tests, buildInput([]caseSpec{
		{n: 3, x: 1, s: "..#"},
		{n: 4, x: 2, s: "...."},
		{n: 5, x: 3, s: "##..#"},
		{n: 6, x: 4, s: "#...#."},
	}))

	tests = append(tests, buildInput([]caseSpec{
		{n: 2, x: 1, s: ".."},
		{n: 2, x: 2, s: ".."},
		{n: 3, x: 2, s: ".#."},
	}))

	for i := 0; i < 12; i++ {
		tests = append(tests, randomBatch(rng, 20, 2000))
	}
	tests = append(tests, randomBatch(rng, 40, maxN))

	return tests
}

type caseSpec struct {
	n int
	x int
	s string
}

func buildInput(cases []caseSpec) testInput {
	var b strings.Builder
	fmt.Fprintf(&b, "%d\n", len(cases))
	for _, cs := range cases {
		fmt.Fprintf(&b, "%d %d\n", cs.n, cs.x)
		fmt.Fprintf(&b, "%s\n", cs.s)
	}
	return testInput{raw: b.String(), t: len(cases)}
}

func randomBatch(rng *rand.Rand, maxCases int, maxNPerCase int) testInput {
	t := rng.Intn(maxCases) + 1
	remaining := maxN
	var cases []caseSpec
	for i := 0; i < t; i++ {
		minRemaining := t - i - 1
		maxAllowed := remaining - minRemaining
		if maxAllowed < 2 {
			break
		}
		if maxAllowed > maxNPerCase {
			maxAllowed = maxNPerCase
		}
		n := rng.Intn(maxAllowed-1) + 2
		remaining -= n
		x := rng.Intn(n) + 1
		s := randomGrid(rng, n, x)
		cases = append(cases, caseSpec{n: n, x: x, s: s})
	}
	if len(cases) == 0 {
		cases = append(cases, caseSpec{n: 2, x: 1, s: ".."})
	}
	return buildInput(cases)
}

func randomGrid(rng *rand.Rand, n, x int) string {
	b := make([]byte, n)
	emptyCount := 0
	for i := 0; i < n; i++ {
		if rng.Intn(3) == 0 {
			b[i] = '#'
		} else {
			b[i] = '.'
			emptyCount++
		}
	}
	if b[x-1] == '#' {
		b[x-1] = '.'
		emptyCount++
	}
	if emptyCount < 2 {
		// ensure at least two empty cells
		for i := 0; i < n && emptyCount < 2; i++ {
			if b[i] == '#' {
				b[i] = '.'
				emptyCount++
			}
		}
	}
	return string(b)
}
