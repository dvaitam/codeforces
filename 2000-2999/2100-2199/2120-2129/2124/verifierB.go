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

const refSource2124B = "2000-2999/2100-2199/2120-2129/2124/2124B.go"

type testCase struct {
	name  string
	input string
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierB.go /path/to/binary")
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
		refVals, err := parseOutput(tc.input, refOut)
		if err != nil {
			fmt.Fprintf(os.Stderr, "reference produced invalid output on test %d (%s): %v\ninput:\n%soutput:\n%s",
				idx+1, tc.name, err, tc.input, refOut)
			os.Exit(1)
		}

		candOut, err := runProgram(candidate, tc.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "candidate runtime error on test %d (%s): %v\ninput:\n%soutput:\n%s",
				idx+1, tc.name, err, tc.input, candOut)
			os.Exit(1)
		}
		candVals, err := parseOutput(tc.input, candOut)
		if err != nil {
			fmt.Fprintf(os.Stderr, "candidate output invalid on test %d (%s): %v\ninput:\n%soutput:\n%s",
				idx+1, tc.name, err, tc.input, candOut)
			os.Exit(1)
		}

		if len(refVals) != len(candVals) {
			fmt.Fprintf(os.Stderr, "test %d (%s) failed: expected %d answers got %d\ninput:\n%sreference output:\n%s\ncandidate output:\n%s",
				idx+1, tc.name, len(refVals), len(candVals), tc.input, refOut, candOut)
			os.Exit(1)
		}
		for i := range refVals {
			if refVals[i] != candVals[i] {
				fmt.Fprintf(os.Stderr, "test %d (%s) failed on case %d: expected %d got %d\ninput:\n%sreference output:\n%s\ncandidate output:\n%s",
					idx+1, tc.name, i+1, refVals[i], candVals[i], tc.input, refOut, candOut)
				os.Exit(1)
			}
		}
	}

	fmt.Printf("All %d tests passed\n", len(tests))
}

func buildReference() (string, func(), error) {
	dir, err := os.MkdirTemp("", "cf-2124B-ref-")
	if err != nil {
		return "", nil, fmt.Errorf("failed to create temp dir: %v", err)
	}
	binPath := filepath.Join(dir, "ref2124B.bin")
	cmd := exec.Command("go", "build", "-o", binPath, refSource2124B)
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		_ = os.RemoveAll(dir)
		return "", nil, fmt.Errorf("failed to build reference: %v\n%s", err, stderr.String())
	}
	cleanup := func() { _ = os.RemoveAll(dir) }
	return binPath, cleanup, nil
}

func runProgram(bin string, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return out.String(), err
}

func parseOutput(input, output string) ([]int64, error) {
	inFields := strings.Fields(input)
	if len(inFields) == 0 {
		return nil, fmt.Errorf("empty input")
	}
	t, err := strconv.Atoi(inFields[0])
	if err != nil || t < 1 || t > 10000 {
		return nil, fmt.Errorf("invalid test count %q", inFields[0])
	}
	outFields := strings.Fields(output)
	if len(outFields) != t {
		return nil, fmt.Errorf("expected %d answers, got %d", t, len(outFields))
	}
	res := make([]int64, t)
	for i, tok := range outFields {
		v, err := strconv.ParseInt(tok, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("invalid integer %q", tok)
		}
		res[i] = v
	}
	return res, nil
}

func buildTests() []testCase {
	tests := []testCase{
		makeManual("samples", []manualCase{
			{a: []int{2, 1}, expect: 2},
			{a: []int{1, 2, 3}, expect: 2},
			{a: []int{3, 0, 2, 3}, expect: 3},
		}),
		makeManual("no_op_needed", []manualCase{
			{a: []int{0, 0, 0, 0}, expect: 0},
			{a: []int{5, 4, 3, 2, 1}, expect: 15},
		}),
		makeManual("operation_helpful", []manualCase{
			{a: []int{1, 0, 5}, expect: 1},
			{a: []int{2, 5, 0}, expect: 2},
		}),
	}

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 180; i++ {
		tests = append(tests, randomTest(rng, i))
	}
	tests = append(tests, largeTest())
	return tests
}

type manualCase struct {
	a      []int
	expect int64 // for readability, not used
}

func makeManual(name string, cases []manualCase) testCase {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", len(cases)))
	for _, cs := range cases {
		n := len(cs.a)
		sb.WriteString(fmt.Sprintf("%d\n", n))
		writeArray(&sb, cs.a)
		sb.WriteByte('\n')
	}
	return testCase{name: name, input: sb.String()}
}

func randomTest(rng *rand.Rand, idx int) testCase {
	t := rng.Intn(5) + 1
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", t))
	for i := 0; i < t; i++ {
		n := rng.Intn(50) + 2
		sb.WriteString(fmt.Sprintf("%d\n", n))
		for j := 0; j < n; j++ {
			if j > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(strconv.Itoa(rng.Intn(n + 1)))
		}
		sb.WriteByte('\n')
	}
	return testCase{name: fmt.Sprintf("random_%d", idx+1), input: sb.String()}
}

func largeTest() testCase {
	n := 200000
	var sb strings.Builder
	sb.WriteString("1\n")
	sb.WriteString(fmt.Sprintf("%d\n", n))
	for i := 0; i < n; i++ {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString("0")
	}
	sb.WriteByte('\n')
	return testCase{name: "large_zeroes", input: sb.String()}
}

func writeArray(sb *strings.Builder, arr []int) {
	for i, v := range arr {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(strconv.Itoa(v))
	}
}
