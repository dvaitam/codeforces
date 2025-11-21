package main

import (
	"bufio"
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

const refSource1999D = "1000-1999/1900-1999/1990-1999/1999/1999D.go"

type testCase struct {
	s string
	t string
}

type testInput struct {
	input string
	cases []testCase
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	candidate := os.Args[1]

	refBin, err := buildReference(refSource1999D)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to build reference: %v\n", err)
		os.Exit(1)
	}
	defer os.Remove(refBin)

	tests := generateTests()
	for idx, test := range tests {
		refOut, err := runProgram(refBin, test.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "reference failed on test %d: %v\ninput:\n%s", idx+1, err, test.input)
			os.Exit(1)
		}
		expected, err := parseOutputs(refOut, len(test.cases), nil)
		if err != nil {
			fmt.Fprintf(os.Stderr, "reference output invalid on test %d: %v\n", idx+1, err)
			os.Exit(1)
		}

		userOut, err := runCandidate(candidate, test.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "candidate failed on test %d: %v\ninput:\n%s", idx+1, err, test.input)
			os.Exit(1)
		}
		actual, err := parseOutputs(userOut, len(test.cases), test.cases)
		if err != nil {
			fmt.Fprintf(os.Stderr, "candidate output invalid on test %d: %v\ninput:\n%soutput:\n%s\n", idx+1, err, test.input, userOut)
			os.Exit(1)
		}

		for caseIdx, tc := range test.cases {
			exp := expected[caseIdx]
			act := actual[caseIdx]
			if !compareCases(exp, act) {
				fmt.Fprintf(os.Stderr, "test %d case %d mismatch:\nexpected:\n%s\ngot:\n%s\ncase input:\n%s",
					idx+1, caseIdx+1, exp.String(), act.String(), formatCase(tc))
				os.Exit(1)
			}
			if exp.ok {
				if err := validateCandidate(tc, act.str); err != nil {
					fmt.Fprintf(os.Stderr, "test %d case %d invalid constructed string: %v\ncase input:\n%s",
						idx+1, caseIdx+1, err, formatCase(tc))
					os.Exit(1)
				}
			}
		}
	}

	fmt.Printf("All %d tests passed\n", len(tests))
}

type verdict struct {
	ok  bool
	str string
}

func (v verdict) String() string {
	if !v.ok {
		return "NO"
	}
	return "YES\n" + v.str
}

func buildReference(source string) (string, error) {
	tmp, err := os.CreateTemp("", "1999D-ref-*")
	if err != nil {
		return "", err
	}
	tmp.Close()

	cmd := exec.Command("go", "build", "-o", tmp.Name(), filepath.Clean(source))
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		os.Remove(tmp.Name())
		return "", fmt.Errorf("%v\n%s", err, out.String())
	}
	return tmp.Name(), nil
}

func runProgram(path, input string) (string, error) {
	cmd := exec.Command(path)
	return runWithInput(cmd, input)
}

func runCandidate(path, input string) (string, error) {
	cmd := commandFor(path)
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
	return out.String(), err
}

func parseOutputs(out string, cases int, inputs []testCase) ([]verdict, error) {
	reader := bufio.NewReader(strings.NewReader(out))
	res := make([]verdict, cases)
	for i := 0; i < cases; i++ {
		var ans string
		if _, err := fmt.Fscan(reader, &ans); err != nil {
			return nil, fmt.Errorf("expected verdict for case %d: %v", i+1, err)
		}
		up := strings.ToUpper(ans)
		if up == "NO" {
			res[i] = verdict{ok: false}
		} else if up == "YES" {
			var line string
			if _, err := fmt.Fscan(reader, &line); err != nil {
				return nil, fmt.Errorf("case %d: expected string after YES", i+1)
			}
			res[i] = verdict{ok: true, str: line}
		} else {
			return nil, fmt.Errorf("case %d: invalid verdict %q", i+1, ans)
		}
	}
	var extra string
	if _, err := fmt.Fscan(reader, &extra); err == nil {
		return nil, fmt.Errorf("extra output detected (starts with %s)", extra)
	}
	return res, nil
}

func compareCases(exp, act verdict) bool {
	return exp.ok == act.ok
}

func validateCandidate(tc testCase, ans string) error {
	if len(ans) != len(tc.s) {
		return fmt.Errorf("expected string of length %d, got %d", len(tc.s), len(ans))
	}
	idx := 0
	for i := 0; i < len(ans); i++ {
		if ans[i] < 'a' || ans[i] > 'z' {
			return fmt.Errorf("invalid character %q at position %d", ans[i], i+1)
		}
		if tc.s[i] != '?' && tc.s[i] != ans[i] {
			return fmt.Errorf("character at position %d must be %q but got %q", i+1, tc.s[i], ans[i])
		}
		if idx < len(tc.t) && ans[i] == tc.t[idx] {
			idx++
		}
	}
	if idx != len(tc.t) {
		return fmt.Errorf("t is not a subsequence")
	}
	return nil
}

func formatCase(tc testCase) string {
	return fmt.Sprintf("1\n%s\n%s\n", tc.s, tc.t)
}

func generateTests() []testInput {
	var tests []testInput
	tests = append(tests, sampleTests())
	tests = append(tests, edgeTests())
	rng := rand.New(rand.NewSource(1999))
	tests = append(tests, randomTests(rng, []int{5, 7, 9}))
	tests = append(tests, randomTests(rng, []int{100, 200}))
	tests = append(tests, randomTests(rng, []int{200000}))
	return tests
}

func sampleTests() testInput {
	cases := []testCase{
		{s: "?????", t: "xbxab"},
		{s: "??eabc", t: "de"},
		{s: "ayy?x", t: "aab"},
		{s: "?edac", t: "paium"},
		{s: "om", t: "om"},
	}
	return buildTestInput(cases)
}

func edgeTests() testInput {
	cases := []testCase{
		{s: "a", t: "a"},
		{s: "?", t: "a"},
		{s: "b", t: "a"},
		{s: "??", t: "ab"},
		{s: "??", t: "ba"},
		{s: "abc", t: "abc"},
		{s: "aaa", t: "aaaa"},
	}
	return buildTestInput(cases)
}

func randomTests(rng *rand.Rand, sizes []int) testInput {
	var cases []testCase
	total := 0
	for _, t := range sizes {
		total += t
		cases = append(cases, randomCase(rng, t))
	}
	return buildTestInput(cases)
}

func randomCase(rng *rand.Rand, length int) testCase {
	sLen := rng.Intn(length) + 1
	tLen := rng.Intn(sLen) + 1
	s := make([]byte, sLen)
	for i := 0; i < sLen; i++ {
		if rng.Intn(5) == 0 {
			s[i] = '?'
		} else {
			s[i] = byte('a' + rng.Intn(26))
		}
	}
	t := make([]byte, tLen)
	for i := 0; i < tLen; i++ {
		t[i] = byte('a' + rng.Intn(26))
	}
	return testCase{s: string(s), t: string(t)}
}

func buildTestInput(cases []testCase) testInput {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", len(cases)))
	for _, tc := range cases {
		sb.WriteString(tc.s)
		sb.WriteByte('\n')
		sb.WriteString(tc.t)
		sb.WriteByte('\n')
	}
	return testInput{input: sb.String(), cases: cases}
}
