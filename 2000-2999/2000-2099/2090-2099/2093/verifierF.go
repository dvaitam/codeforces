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

const refSource = "./2093F.go"

type testCase struct {
	input string
}

type caseData struct {
	a  []string
	bs [][]string
}

var wordCounter int64 = 1

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
	tmp, err := os.CreateTemp("", "2093F-ref-*")
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

func equalTokens(expected, got string) bool {
	ta := strings.Fields(expected)
	tb := strings.Fields(got)
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
	rng := rand.New(rand.NewSource(20932093))
	var tests []testCase

	// Deterministic edge checks.
	tests = append(tests, makeTest([]caseData{
		simpleMatchCase(),
	}))
	tests = append(tests, makeTest([]caseData{
		impossibleSingleNetwork(),
	}))
	tests = append(tests, makeTest([]caseData{
		mixedSmallCase(),
	}))
	tests = append(tests, makeTest([]caseData{
		dualCaseBundle(rng),
	}))

	// Random small possible cases.
	for i := 0; i < 20; i++ {
		n := rng.Intn(9) + 2
		m := rng.Intn(8) + 2
		best := rng.Intn(n + 1)
		if m == 1 {
			best = n
		}
		tests = append(tests, makeTest([]caseData{
			makePossibleCase(rng, n, m, best),
		}))
	}

	// Random small impossible cases.
	for i := 0; i < 10; i++ {
		n := rng.Intn(8) + 2
		m := rng.Intn(7) + 2
		tests = append(tests, makeTest([]caseData{
			makeImpossibleCase(rng, n, m),
		}))
	}

	// Mid-sized mixed bundle.
	tests = append(tests, makeTest([]caseData{
		makePossibleCase(rng, 120, 60, 40),
		makePossibleCase(rng, 180, 90, 90),
		makeImpossibleCase(rng, 150, 40),
	}))

	// Large edge-focused cases within constraints.
	tests = append(tests, makeTest([]caseData{
		makePossibleCase(rng, 400, 400, 250),
	}))
	tests = append(tests, makeTest([]caseData{
		makePossibleCase(rng, 500, 300, 150),
	}))

	return tests
}

func makeTest(cases []caseData) testCase {
	return testCase{input: buildInput(cases)}
}

func buildInput(cases []caseData) string {
	var b strings.Builder
	fmt.Fprintf(&b, "%d\n", len(cases))
	for _, cs := range cases {
		n := len(cs.a)
		m := len(cs.bs)
		fmt.Fprintf(&b, "%d %d\n", n, m)
		for i, v := range cs.a {
			if i > 0 {
				b.WriteByte(' ')
			}
			b.WriteString(v)
		}
		b.WriteByte('\n')
		for i := 0; i < m; i++ {
			for j := 0; j < n; j++ {
				if j > 0 {
					b.WriteByte(' ')
				}
				b.WriteString(cs.bs[i][j])
			}
			b.WriteByte('\n')
		}
	}
	return b.String()
}

func simpleMatchCase() caseData {
	a := []string{"w1"}
	bs := [][]string{{"w1"}}
	return caseData{a: a, bs: bs}
}

func impossibleSingleNetwork() caseData {
	a := []string{"w2", "w3"}
	bs := [][]string{{"x1", "x2"}}
	return caseData{a: a, bs: bs}
}

func mixedSmallCase() caseData {
	a := []string{"aa", "bb", "cc"}
	bs := [][]string{
		{"aa", "xx", "yy"},
		{"zz", "bb", "cc"},
	}
	return caseData{a: a, bs: bs}
}

func dualCaseBundle(rng *rand.Rand) caseData {
	return makePossibleCase(rng, 6, 3, 3)
}

func makePossibleCase(rng *rand.Rand, n, m, bestMatches int) caseData {
	if m < 1 {
		m = 1
	}
	if bestMatches > n {
		bestMatches = n
	}
	if m == 1 {
		bestMatches = n
	}

	a := make([]string, n)
	for i := 0; i < n; i++ {
		a[i] = newWord()
	}

	bs := make([][]string, m)
	for i := 0; i < m; i++ {
		bs[i] = make([]string, n)
		for j := 0; j < n; j++ {
			bs[i][j] = differentWord(a[j])
		}
	}

	covered := make([]bool, n)
	perm := rng.Perm(n)
	bestNet := 0
	for i := 0; i < bestMatches; i++ {
		pos := perm[i]
		bs[bestNet][pos] = a[pos]
		covered[pos] = true
	}

	for j := 0; j < n; j++ {
		if covered[j] {
			continue
		}
		if m == 1 {
			bs[bestNet][j] = a[j]
			covered[j] = true
			continue
		}
		net := 1 + (j % (m - 1))
		bs[net][j] = a[j]
		covered[j] = true
	}

	return caseData{a: a, bs: bs}
}

func makeImpossibleCase(rng *rand.Rand, n, m int) caseData {
	if m < 1 {
		m = 1
	}
	a := make([]string, n)
	for i := 0; i < n; i++ {
		a[i] = newWord()
	}

	bs := make([][]string, m)
	for i := 0; i < m; i++ {
		bs[i] = make([]string, n)
		for j := 0; j < n; j++ {
			bs[i][j] = differentWord(a[j])
		}
	}

	// Ensure some positions are actually matched so candidate cannot shortcut.
	for j := 1; j < n; j++ {
		net := j % m
		bs[net][j] = a[j]
	}

	// Position 0 is unmatched in all networks.
	bs[0][0] = differentWord(a[0])
	return caseData{a: a, bs: bs}
}

func newWord() string {
	w := fmt.Sprintf("w%d", wordCounter)
	wordCounter++
	return w
}

func differentWord(avoid string) string {
	w := newWord()
	if w == avoid {
		w = newWord()
	}
	return w
}
