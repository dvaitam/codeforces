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

const refSource = "2089B2.go"

type testCase struct {
	input string
}

type caseData struct {
	n int
	k int64
	a []int64
	b []int64
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierB2.go /path/to/binary")
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
		tCount, err := countCases(tc.input)
		if err != nil {
			fail("malformed generated test %d: %v", idx+1, err)
		}

		expectOut, err := runProgram(exec.Command(refBin), tc.input)
		if err != nil {
			fail("reference failed on test %d: %v\ninput:\n%s", idx+1, err, tc.input)
		}
		expect, err := parseOutputs(expectOut, tCount)
		if err != nil {
			fail("unable to parse reference output on test %d: %v\noutput:\n%s", idx+1, err, expectOut)
		}

		gotOut, err := runProgram(commandFor(candidate), tc.input)
		if err != nil {
			fail("runtime error on test %d: %v\ninput:\n%soutput:\n%s", idx+1, err, tc.input, gotOut)
		}
		got, err := parseOutputs(gotOut, tCount)
		if err != nil {
			fail("invalid output on test %d: %v\ninput:\n%soutput:\n%s", idx+1, err, tc.input, gotOut)
		}

		if len(got) != len(expect) {
			fail("wrong answer count on test %d: expected %d, got %d", idx+1, len(expect), len(got))
		}
		for i := range expect {
			if got[i] != expect[i] {
				fail("wrong answer on test %d case %d: expected %d, got %d\ninput:\n%s", idx+1, i+1, expect[i], got[i], tc.input)
			}
		}
	}

	fmt.Printf("All %d tests passed.\n", len(tests))
}

func buildReference() (string, error) {
	tmp, err := os.CreateTemp("", "2089B2-ref-*")
	if err != nil {
		return "", err
	}
	tmp.Close()

	source := filepath.Join(".", refSource)
	cmd := exec.Command("go", "build", "-o", tmp.Name(), source)
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

func parseOutputs(out string, expected int) ([]int64, error) {
	tokens := strings.Fields(out)
	if len(tokens) != expected {
		return nil, fmt.Errorf("expected %d numbers, got %d", expected, len(tokens))
	}
	res := make([]int64, expected)
	for i, t := range tokens {
		v, err := strconv.ParseInt(t, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("failed to parse token %q at position %d: %w", t, i+1, err)
		}
		if v < 0 {
			return nil, fmt.Errorf("negative answer at position %d", i+1)
		}
		res[i] = v
	}
	return res, nil
}

func countCases(input string) (int, error) {
	fields := strings.Fields(input)
	if len(fields) == 0 {
		return 0, fmt.Errorf("empty input")
	}
	pos := 0
	readInt := func() (int64, error) {
		if pos >= len(fields) {
			return 0, fmt.Errorf("unexpected end of input")
		}
		v, err := strconv.ParseInt(fields[pos], 10, 64)
		pos++
		return v, err
	}
	tRaw, err := readInt()
	if err != nil {
		return 0, fmt.Errorf("failed to read t: %w", err)
	}
	t := int(tRaw)
	for i := 0; i < t; i++ {
		nRaw, err := readInt()
		if err != nil {
			return 0, fmt.Errorf("failed to read n of case %d: %w", i+1, err)
		}
		qRaw, err := readInt()
		if err != nil {
			return 0, fmt.Errorf("failed to read k of case %d: %w", i+1, err)
		}
		n := int(nRaw)
		if n < 0 {
			return 0, fmt.Errorf("negative n")
		}
		for j := 0; j < n; j++ {
			if _, err := readInt(); err != nil {
				return 0, fmt.Errorf("failed to read a[%d] of case %d: %w", j, i+1, err)
			}
		}
		for j := 0; j < n; j++ {
			if _, err := readInt(); err != nil {
				return 0, fmt.Errorf("failed to read b[%d] of case %d: %w", j, i+1, err)
			}
		}
		if qRaw < 0 {
			return 0, fmt.Errorf("negative k")
		}
	}
	if pos != len(fields) {
		return 0, fmt.Errorf("extra data at end of input")
	}
	return t, nil
}

func buildTests() []testCase {
	var tests []testCase

	tests = append(tests, makeTest([]caseData{
		{n: 1, k: 0, a: []int64{1}, b: []int64{1}},
	}))
	tests = append(tests, makeTest([]caseData{
		{n: 3, k: 0, a: []int64{1, 1, 4}, b: []int64{5, 1, 4}},
		{n: 4, k: 6, a: []int64{1, 2, 3, 4}, b: []int64{4, 3, 2, 1}},
	}))

	rng := rand.New(rand.NewSource(20892089))
	for i := 0; i < 25; i++ {
		cases := make([]caseData, 0, 4)
		tc := rng.Intn(3) + 1
		for j := 0; j < tc; j++ {
			n := rng.Intn(60) + 1
			cases = append(cases, randomCase(rng, n))
		}
		tests = append(tests, makeTest(cases))
	}

	// Some larger cases for stress.
	tests = append(tests, makeTest([]caseData{randomCase(rand.New(rand.NewSource(1)), 300)}))
	tests = append(tests, makeTest([]caseData{randomCase(rand.New(rand.NewSource(2)), 600)}))

	return tests
}

func randomCase(rng *rand.Rand, n int) caseData {
	a := make([]int64, n)
	b := make([]int64, n)
	var asum int64
	for i := 0; i < n; i++ {
		a[i] = rng.Int63n(1_000_000_000) + 1
		asum += a[i]
	}
	var bsum int64
	for i := 0; i < n; i++ {
		// ensure b sum >= a sum; bias larger
		extra := rng.Int63n(1_000_000_000)
		b[i] = a[i]/2 + extra + 1
		bsum += b[i]
	}
	if bsum < asum {
		diff := asum - bsum
		b[0] += diff
		bsum += diff
	}
	var kLimit int64 = asum
	if kLimit < 0 {
		kLimit = 0
	}
	var k int64
	if kLimit > 0 {
		k = rng.Int63n(kLimit + 1)
	}
	return caseData{n: n, k: k, a: a, b: b}
}

func makeTest(cases []caseData) testCase {
	var b strings.Builder
	fmt.Fprintln(&b, len(cases))
	for _, c := range cases {
		fmt.Fprintf(&b, "%d %d\n", c.n, c.k)
		for i, v := range c.a {
			if i > 0 {
				fmt.Fprint(&b, " ")
			}
			fmt.Fprint(&b, v)
		}
		fmt.Fprintln(&b)
		for i, v := range c.b {
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
