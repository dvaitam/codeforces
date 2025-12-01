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

const refSource = "./2052F.go"

type grid struct {
	n   int
	row [2]string
}

type testCase struct {
	input string
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierF.go /path/to/binary")
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
		tCount, err := readTestCount(tc.input)
		if err != nil {
			fail("invalid generated test %d: %v", idx+1, err)
		}

		expectOut, err := runProgram(exec.Command(refBin), tc.input)
		if err != nil {
			fail("reference failed on test %d: %v\ninput:\n%s", idx+1, err, tc.input)
		}
		expect, err := parseVerdicts(expectOut, tCount)
		if err != nil {
			fail("could not parse reference output on test %d: %v\noutput:\n%s", idx+1, err, expectOut)
		}

		gotOut, err := runProgram(commandFor(candidate), tc.input)
		if err != nil {
			fail("runtime error on test %d: %v\ninput:\n%soutput:\n%s", idx+1, err, tc.input, gotOut)
		}
		got, err := parseVerdicts(gotOut, tCount)
		if err != nil {
			fail("invalid output on test %d: %v\ninput:\n%soutput:\n%s", idx+1, err, tc.input, gotOut)
		}

		for i := 0; i < tCount; i++ {
			if got[i] != expect[i] {
				fail("wrong verdict on test %d case %d: expected %s, got %s\ninput:\n%s", idx+1, i+1, expect[i], got[i], tc.input)
			}
		}
	}

	fmt.Printf("All %d tests passed.\n", len(tests))
}

func buildReference() (string, error) {
	tmp, err := os.CreateTemp("", "2052F-ref-*")
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

func readTestCount(input string) (int, error) {
	fields := strings.Fields(input)
	if len(fields) == 0 {
		return 0, fmt.Errorf("empty input")
	}
	t, err := strconv.Atoi(fields[0])
	if err != nil {
		return 0, fmt.Errorf("failed to parse T: %w", err)
	}
	if t < 0 {
		return 0, fmt.Errorf("negative T")
	}
	return t, nil
}

func parseVerdicts(out string, t int) ([]string, error) {
	tokens := strings.Fields(out)
	if len(tokens) != t {
		return nil, fmt.Errorf("expected %d verdicts, got %d", t, len(tokens))
	}
	for i, v := range tokens {
		if v != "Unique" && v != "Multiple" && v != "None" {
			return nil, fmt.Errorf("invalid verdict %q at position %d", v, i+1)
		}
	}
	return tokens, nil
}

func buildTests() []testCase {
	var tests []testCase

	// Deterministic edge cases.
	tests = append(tests, makeTest([]grid{
		{n: 1, row: [2]string{".", "."}},
	}))
	tests = append(tests, makeTest([]grid{
		{n: 1, row: [2]string{".", "#"}},
	}))
	tests = append(tests, makeTest([]grid{
		{n: 2, row: [2]string{"##", "##"}},
		{n: 2, row: [2]string{"..", ".."}},
	}))
	tests = append(tests, makeTest([]grid{
		{n: 3, row: [2]string{"#.#", ".#."}},
		{n: 4, row: [2]string{"....", "####"}},
	}))

	rng := rand.New(rand.NewSource(20522052))
	for i := 0; i < 25; i++ {
		cases := make([]grid, 0, 6)
		cnt := rng.Intn(5) + 1
		for j := 0; j < cnt; j++ {
			n := rng.Intn(40) + 1
			cases = append(cases, randomGrid(rng, n))
		}
		tests = append(tests, makeTest(cases))
	}

	// A couple of larger cases.
	tests = append(tests, makeTest([]grid{randomGrid(rand.New(rand.NewSource(1)), 150)}))
	tests = append(tests, makeTest([]grid{randomGrid(rand.New(rand.NewSource(2)), 300)}))

	return tests
}

func randomGrid(rng *rand.Rand, n int) grid {
	var b1, b2 strings.Builder
	for i := 0; i < n; i++ {
		if rng.Intn(4) == 0 {
			b1.WriteByte('#')
		} else {
			b1.WriteByte('.')
		}
		if rng.Intn(4) == 0 {
			b2.WriteByte('#')
		} else {
			b2.WriteByte('.')
		}
	}
	return grid{n: n, row: [2]string{b1.String(), b2.String()}}
}

func makeTest(cases []grid) testCase {
	var b strings.Builder
	fmt.Fprintln(&b, len(cases))
	for _, c := range cases {
		fmt.Fprintln(&b, c.n)
		fmt.Fprintln(&b, c.row[0])
		fmt.Fprintln(&b, c.row[1])
	}
	return testCase{input: b.String()}
}

func fail(format string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, format+"\n", args...)
	os.Exit(1)
}
