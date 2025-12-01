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

const refSourceC = "./2030C.go"

type caseData struct {
	s string
}

type testBundle struct {
	cases []caseData
}

func (tb testBundle) input() string {
	var b strings.Builder
	fmt.Fprintln(&b, len(tb.cases))
	for _, cs := range tb.cases {
		fmt.Fprintf(&b, "%d\n%s\n", len(cs.s), cs.s)
	}
	return b.String()
}

func main() {
	var candidate string
	if len(os.Args) == 2 {
		candidate = os.Args[1]
	} else if len(os.Args) == 3 && os.Args[1] == "--" {
		candidate = os.Args[2]
	} else {
		fmt.Fprintln(os.Stderr, "usage: go run verifierC.go /path/to/candidate")
		os.Exit(1)
	}

	refBin, err := buildReference()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to build reference: %v\n", err)
		os.Exit(1)
	}
	defer os.Remove(refBin)

	tests := buildTests()
	for idx, tb := range tests {
		input := tb.input()

		refOut, err := runProgram(refBin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "reference failed on test %d: %v\ninput:\n%s\n", idx+1, err, input)
			os.Exit(1)
		}
		refVals, err := parseOutputs(refOut, len(tb.cases))
		if err != nil {
			fmt.Fprintf(os.Stderr, "could not parse reference output on test %d: %v\noutput:\n%s\n", idx+1, err, refOut)
			os.Exit(1)
		}

		candOut, err := runCandidate(candidate, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "candidate runtime error on test %d: %v\ninput:\n%s\nstdout/stderr:\n%s\n", idx+1, err, input, candOut)
			os.Exit(1)
		}
		candVals, err := parseOutputs(candOut, len(tb.cases))
		if err != nil {
			fmt.Fprintf(os.Stderr, "invalid candidate output on test %d: %v\noutput:\n%s\n", idx+1, err, candOut)
			os.Exit(1)
		}

		for i := range refVals {
			if refVals[i] != candVals[i] {
				fmt.Fprintf(os.Stderr, "wrong answer on test %d case %d\ninput:\n%sreference: %s\ncandidate: %s\n", idx+1, i+1, input, refVals[i], candVals[i])
				os.Exit(1)
			}
		}
	}
	fmt.Printf("All %d tests passed.\n", len(tests))
}

func buildReference() (string, error) {
	tmp, err := os.CreateTemp("", "2030C-ref-*")
	if err != nil {
		return "", err
	}
	tmp.Close()

	cmd := exec.Command("go", "build", "-o", tmp.Name(), filepath.Clean(refSourceC))
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
	switch filepath.Ext(path) {
	case ".go":
		return exec.Command("go", "run", path)
	case ".py":
		return exec.Command("python3", path)
	default:
		return exec.Command(path)
	}
}

func runWithInput(cmd *exec.Cmd, input string) (string, error) {
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
	res := make([]string, expected)
	for i, f := range fields {
		up := strings.ToUpper(f)
		if up != "YES" && up != "NO" {
			return nil, fmt.Errorf("invalid token %q (expected YES/NO)", f)
		}
		res[i] = up
	}
	return res, nil
}

func buildTests() []testBundle {
	var tests []testBundle

	// Statement sample.
	tests = append(tests, testBundle{cases: []caseData{
		{s: "11"},
		{s: "010"},
		{s: "1011111111"},
		{s: "0010011111101"},
		{s: "01000010"},
	}})

	// Small fixed shapes.
	tests = append(tests, testBundle{cases: []caseData{
		{s: "10"},
		{s: "01"},
		{s: "101"},
		{s: "0101"},
		{s: "0010"},
		{s: "1111"},
	}})

	// Alternating and long runs.
	tests = append(tests, testBundle{cases: []caseData{
		{s: strings.Repeat("10", 25)},
		{s: strings.Repeat("01", 25)},
		{s: "0" + strings.Repeat("1", 48) + "0"},
		{s: "1" + strings.Repeat("0", 48) + "1"},
	}})

	rng := rand.New(rand.NewSource(2030))
	// Random bundles with modest total length.
	for i := 0; i < 4; i++ {
		var cs []caseData
		total := 0
		for total < 800 {
			n := rng.Intn(40) + 2
			cs = append(cs, caseData{s: randomBits(rng, n)})
			total += n
		}
		tests = append(tests, testBundle{cases: cs})
	}

	// Stress bundle approaching the total length limit.
	tests = append(tests, stressBundle(rng, 200000))

	return tests
}

func randomBits(rng *rand.Rand, n int) string {
	var b strings.Builder
	b.Grow(n)
	for i := 0; i < n; i++ {
		if rng.Intn(2) == 0 {
			b.WriteByte('0')
		} else {
			b.WriteByte('1')
		}
	}
	return b.String()
}

func stressBundle(rng *rand.Rand, total int) testBundle {
	var cs []caseData
	remaining := total

	// One large case plus several medium ones.
	mainLen := remaining / 2
	cs = append(cs, caseData{s: randomBits(rng, mainLen)})
	remaining -= mainLen

	for remaining > 0 {
		n := rng.Intn(4000) + 2
		if n > remaining {
			n = remaining
		}
		cs = append(cs, caseData{s: randomBits(rng, n)})
		remaining -= n
	}
	return testBundle{cases: cs}
}
