package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
)

const refSource = "./2064D.go"

type caseData struct {
	n       int
	w       []int
	queries []int
}

type testCase struct {
	input string
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	candidate := os.Args[1]

	refBin, err := buildReference()
	if err != nil {
		fail("failed to build reference: %v", err)
	}
	defer os.Remove(refBin)

	tests := buildTests()
	for idx, tc := range tests {
		tCount, totalQueries, err := countCasesAndQueries(tc.input)
		if err != nil {
			fail("invalid generated test %d: %v", idx+1, err)
		}

		expectOut, err := runProgram(exec.Command(refBin), tc.input)
		if err != nil {
			fail("reference failed on test %d: %v\ninput:\n%s", idx+1, err, tc.input)
		}
		expect, err := parseOutputs(expectOut, totalQueries)
		if err != nil {
			fail("could not parse reference output on test %d: %v\noutput:\n%s", idx+1, err, expectOut)
		}

		gotOut, err := runProgram(commandFor(candidate), tc.input)
		if err != nil {
			fail("runtime error on test %d: %v\ninput:\n%soutput:\n%s", idx+1, err, tc.input, gotOut)
		}
		got, err := parseOutputs(gotOut, totalQueries)
		if err != nil {
			fail("invalid output on test %d: %v\ninput:\n%soutput:\n%s", idx+1, err, tc.input, gotOut)
		}

		if len(got) != len(expect) {
			fail("wrong number of answers on test %d (expected %d, got %d)", idx+1, len(expect), len(got))
		}
		for i := range expect {
			if got[i] != expect[i] {
				fail("wrong answer on test %d query %d: expected %d, got %d\ninput:\n%s", idx+1, i+1, expect[i], got[i], tc.input)
			}
		}
		if tCount == 0 {
			fail("generated test %d has zero cases", idx+1)
		}
	}

	fmt.Printf("All %d tests passed.\n", len(tests))
}

func buildReference() (string, error) {
	tmp, err := os.CreateTemp("", "2064D-ref-*")
	if err != nil {
		return "", err
	}
	tmp.Close()

	cmd := exec.Command("go", "build", "-o", tmp.Name(), filepath.Clean(refSource))
	var buf bytes.Buffer
	cmd.Stdout = &buf
	cmd.Stderr = &buf
	if err := cmd.Run(); err != nil {
		os.Remove(tmp.Name())
		return "", fmt.Errorf("%v\n%s", err, buf.String())
	}
	return tmp.Name(), nil
}

func commandFor(path string) *exec.Cmd {
	switch filepath.Ext(path) {
	case ".go":
		return exec.Command("go", "run", path)
	case ".py":
		return exec.Command("python3", path)
	default:
		return exec.Command(path)
	}
}

func runProgram(cmd *exec.Cmd, input string) (string, error) {
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return out.String(), err
}

func parseOutputs(out string, expected int) ([]int, error) {
	tokens := strings.Fields(out)
	if len(tokens) != expected {
		return nil, fmt.Errorf("expected %d numbers, got %d", expected, len(tokens))
	}
	ans := make([]int, expected)
	for i, t := range tokens {
		v, err := strconv.Atoi(t)
		if err != nil {
			return nil, fmt.Errorf("failed to parse token %q at position %d: %w", t, i+1, err)
		}
		if v < 0 {
			return nil, fmt.Errorf("negative answer at position %d", i+1)
		}
		ans[i] = v
	}
	return ans, nil
}

func countCasesAndQueries(input string) (int, int, error) {
	fields := strings.Fields(input)
	if len(fields) == 0 {
		return 0, 0, fmt.Errorf("empty input")
	}
	pos := 0
	readInt := func() (int, error) {
		if pos >= len(fields) {
			return 0, fmt.Errorf("unexpected end of input")
		}
		v, err := strconv.Atoi(fields[pos])
		pos++
		return v, err
	}
	t, err := readInt()
	if err != nil {
		return 0, 0, fmt.Errorf("failed to read t: %w", err)
	}
	totalQ := 0
	for i := 0; i < t; i++ {
		n, err := readInt()
		if err != nil {
			return 0, 0, fmt.Errorf("failed to read n of case %d: %w", i+1, err)
		}
		q, err := readInt()
		if err != nil {
			return 0, 0, fmt.Errorf("failed to read q of case %d: %w", i+1, err)
		}
		for j := 0; j < n; j++ {
			if _, err := readInt(); err != nil {
				return 0, 0, fmt.Errorf("failed to read weight %d of case %d: %w", j+1, i+1, err)
			}
		}
		for j := 0; j < q; j++ {
			if _, err := readInt(); err != nil {
				return 0, 0, fmt.Errorf("failed to read query %d of case %d: %w", j+1, i+1, err)
			}
		}
		totalQ += q
	}
	if pos != len(fields) {
		return 0, 0, fmt.Errorf("extra data after parsing input")
	}
	return t, totalQ, nil
}

func buildTests() []testCase {
	var tests []testCase

	tests = append(tests, makeTest([]caseData{
		{n: 1, w: []int{1}, queries: []int{1}},
	}))
	tests = append(tests, makeTest([]caseData{
		{n: 3, w: []int{1, 5, 4}, queries: []int{8, 13, 16}},
		{n: 2, w: []int{10, 9}, queries: []int{10, 4, 3, 9, 7}},
	}))

	rng := rand.New(rand.NewSource(20642064))
	for i := 0; i < 20; i++ {
		cases := make([]caseData, 0, 4)
		tc := rng.Intn(3) + 1
		for j := 0; j < tc; j++ {
			n := rng.Intn(60) + 1
			q := rng.Intn(60) + 1
			cases = append(cases, randomCase(rng, n, q))
		}
		tests = append(tests, makeTest(cases))
	}

	// Larger stress-like single cases.
	tests = append(tests, makeTest([]caseData{randomCase(rand.New(rand.NewSource(1)), 200, 200)}))
	tests = append(tests, makeTest([]caseData{randomCase(rand.New(rand.NewSource(2)), 500, 500)}))

	return tests
}

func randomCase(rng *rand.Rand, n, q int) caseData {
	w := make([]int, n)
	for i := 0; i < n; i++ {
		// keep weights within range
		w[i] = rng.Intn(1<<30-1) + 1
	}
	queries := make([]int, q)
	for i := 0; i < q; i++ {
		queries[i] = rng.Intn(1<<30-1) + 1
	}
	return caseData{n: n, w: w, queries: queries}
}

func makeTest(cases []caseData) testCase {
	var b strings.Builder
	fmt.Fprintln(&b, len(cases))
	for _, c := range cases {
		fmt.Fprintf(&b, "%d %d\n", c.n, len(c.queries))
		for i, v := range c.w {
			if i > 0 {
				fmt.Fprint(&b, " ")
			}
			fmt.Fprint(&b, v)
		}
		fmt.Fprintln(&b)
		for i, v := range c.queries {
			if i > 0 {
				fmt.Fprint(&b, " ")
			}
			fmt.Fprint(&b, v)
		}
		fmt.Fprintln(&b)
	}
	return testCase{input: b.String()}
}

func fail(format string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, format+"\n", args...)
	os.Exit(1)
}
