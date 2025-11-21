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
	"time"
)

const refSource = "0-999/700-799/720-729/720/720E.go"

type testCase struct {
	name  string
	input string
	cases int
}

type singleCase struct {
	n      int
	number string
	codes  []string
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	candidate := os.Args[1]

	refBin, cleanup, err := buildReference()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer cleanup()

	tests := buildTests()
	for idx, tc := range tests {
		refOut, err := runProgram(refBin, tc.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "reference runtime error on test %d (%s): %v\ninput:\n%soutput:\n%s",
				idx+1, tc.name, err, tc.input, refOut)
			os.Exit(1)
		}
		expVals, err := parseOutputs(refOut, tc.cases)
		if err != nil {
			fmt.Fprintf(os.Stderr, "reference output invalid on test %d (%s): %v\ninput:\n%soutput:\n%s",
				idx+1, tc.name, err, tc.input, refOut)
			os.Exit(1)
		}

		candOut, err := runProgram(candidate, tc.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "candidate runtime error on test %d (%s): %v\ninput:\n%soutput:\n%s",
				idx+1, tc.name, err, tc.input, candOut)
			os.Exit(1)
		}
		gotVals, err := parseOutputs(candOut, tc.cases)
		if err != nil {
			fmt.Fprintf(os.Stderr, "candidate output invalid on test %d (%s): %v\ninput:\n%soutput:\n%s",
				idx+1, tc.name, err, tc.input, candOut)
			os.Exit(1)
		}
		for i := 0; i < tc.cases; i++ {
			if gotVals[i] != expVals[i] {
				fmt.Fprintf(os.Stderr, "test %d (%s) failed on case %d: expected %s got %s\ninput:\n%sreference output:\n%s\ncandidate output:\n%s",
					idx+1, tc.name, i+1, expVals[i], gotVals[i], tc.input, refOut, candOut)
				os.Exit(1)
			}
		}
	}
	fmt.Printf("All %d tests passed\n", len(tests))
}

func buildReference() (string, func(), error) {
	dir, err := os.MkdirTemp("", "cf-720E-ref-")
	if err != nil {
		return "", nil, fmt.Errorf("failed to create temp dir: %v", err)
	}
	binPath := filepath.Join(dir, "ref720E.bin")
	cmd := exec.Command("go", "build", "-o", binPath, refSource)
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		os.RemoveAll(dir)
		return "", nil, fmt.Errorf("failed to build reference: %v\n%s", err, stderr.String())
	}
	cleanup := func() { _ = os.RemoveAll(dir) }
	return binPath, cleanup, nil
}

func runProgram(bin, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return out.String(), err
}

func parseOutputs(out string, expected int) ([]string, error) {
	fields := strings.Fields(out)
	if len(fields) != expected {
		return nil, fmt.Errorf("expected %d answers, got %d", expected, len(fields))
	}
	for _, f := range fields {
		if _, err := strconv.ParseUint(f, 10, 64); err != nil {
			return nil, fmt.Errorf("invalid integer %q", f)
		}
	}
	return fields, nil
}

func buildTests() []testCase {
	tests := []testCase{
		makeTestCase("single_digit_unique", []singleCase{
			{
				n:      1,
				number: "5",
				codes:  []string{"abcdefghij"},
			},
		}),
		makeTestCase("two_digits_same_code", []singleCase{
			{
				n:      2,
				number: "09",
				codes: []string{
					strings.Repeat("a", 10),
					"bbbbbbbbbb",
				},
			},
		}),
		makeTestCase("multi_cases", []singleCase{
			{
				n:      3,
				number: "123",
				codes: []string{
					"abcdefghij",
					"jihgfedcba",
					"klmnopqrst",
				},
			},
			{
				n:      1,
				number: "0",
				codes:  []string{"zzzzzzzzzz"},
			},
		}),
	}

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 200; i++ {
		tests = append(tests, randomCase(rng, i))
	}
	tests = append(tests, maxCase())
	return tests
}

func makeTestCase(name string, cases []singleCase) testCase {
	input, cnt := formatInput(cases)
	return testCase{name: name, input: input, cases: cnt}
}

func formatInput(cases []singleCase) (string, int) {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", len(cases)))
	for _, cs := range cases {
		sb.WriteString(fmt.Sprintf("%d\n", cs.n))
		sb.WriteString(cs.number)
		sb.WriteByte('\n')
		for _, line := range cs.codes {
			sb.WriteString(line)
			sb.WriteByte('\n')
		}
	}
	return sb.String(), len(cases)
}

func randomCase(rng *rand.Rand, idx int) testCase {
	tcCount := rng.Intn(4) + 1
	cases := make([]singleCase, tcCount)
	for i := 0; i < tcCount; i++ {
		cases[i] = randomSingleCase(rng)
	}
	return makeTestCase(fmt.Sprintf("random_%d", idx+1), cases)
}

func randomSingleCase(rng *rand.Rand) singleCase {
	n := rng.Intn(10) + 1
	if rng.Intn(5) == 0 {
		n = rng.Intn(18) + 1
	}
	number := randomDigits(rng, n)
	codes := make([]string, n)
	for i := 0; i < n; i++ {
		var sb strings.Builder
		for d := 0; d < 10; d++ {
			sb.WriteByte(byte('a' + rng.Intn(26)))
		}
		codes[i] = sb.String()
	}
	return singleCase{n: n, number: number, codes: codes}
}

func randomDigits(rng *rand.Rand, n int) string {
	var sb strings.Builder
	for i := 0; i < n; i++ {
		sb.WriteByte(byte('0' + rng.Intn(10)))
	}
	return sb.String()
}

func maxCase() testCase {
	cases := []singleCase{
		{
			n:      18,
			number: strings.Repeat("9", 18),
			codes:  makeIdenticalCodes(18, 'm'),
		},
	}
	return makeTestCase("max_case", cases)
}

func makeIdenticalCodes(n int, ch byte) []string {
	codes := make([]string, n)
	line := strings.Repeat(string(ch), 10)
	for i := 0; i < n; i++ {
		codes[i] = line
	}
	return codes
}
