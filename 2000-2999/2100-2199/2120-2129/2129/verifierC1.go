package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

const refSource = "./2129C1.go"

type testCase struct {
	input string
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierC1.go /path/to/binary")
		os.Exit(1)
	}
	candidate := os.Args[1]

	refBin, err := buildReference()
	if err != nil {
		fmt.Fprintln(os.Stderr, "failed to build reference:", err)
		os.Exit(1)
	}
	defer os.Remove(refBin)

	tests := generateTests()
	for i, tc := range tests {
		exp, err := runProgram(refBin, tc.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "reference failed on test %d: %v\ninput:\n%s", i+1, err, tc.input)
			os.Exit(1)
		}

		got, err := runCandidate(candidate, tc.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "runtime error on test %d: %v\ninput:\n%soutput:\n%s", i+1, err, tc.input, got)
			os.Exit(1)
		}

		if err := compare(tc.input, exp, got); err != nil {
			fmt.Fprintf(os.Stderr, "wrong answer on test %d: %v\ninput:\n%s\nexpected:\n%s\ngot:\n%s\n", i+1, err, tc.input, exp, got)
			os.Exit(1)
		}
	}

	fmt.Printf("All %d tests passed.\n", len(tests))
}

func buildReference() (string, error) {
	tmp, err := os.CreateTemp("", "2129C1-ref-*")
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

func compare(input, expected, got string) error {
	expLines := strings.Fields(expected)
	gotLines := strings.Fields(got)
	if len(expLines) != len(gotLines) {
		return fmt.Errorf("line count mismatch: expected %d got %d", len(expLines), len(gotLines))
	}

	// Parse n and s values to validate candidate produces correct length and characters.
	inFields := strings.Fields(input)
	ptr := 0
	readInt := func() (int, error) {
		if ptr >= len(inFields) {
			return 0, fmt.Errorf("unexpected end of input")
		}
		var x int
		_, err := fmt.Sscan(inFields[ptr], &x)
		ptr++
		return x, err
	}
	t, err := readInt()
	if err != nil {
		return fmt.Errorf("failed to read t: %v", err)
	}
	for caseIdx := 0; caseIdx < t; caseIdx++ {
		n, err := readInt()
		if err != nil {
			return fmt.Errorf("failed to read n for case %d: %v", caseIdx+1, err)
		}
		if ptr >= len(inFields) {
			return fmt.Errorf("missing hidden string for case %d", caseIdx+1)
		}
		hidden := inFields[ptr]
		ptr++
		caseOut := gotLines[caseIdx]
		if len(caseOut) != n {
			return fmt.Errorf("case %d: length mismatch, expected %d got %d", caseIdx+1, n, len(caseOut))
		}
		for i, ch := range caseOut {
			if ch != '(' && ch != ')' {
				return fmt.Errorf("case %d: invalid character %q at position %d", caseIdx+1, ch, i+1)
			}
		}
		if caseOut != hidden {
			return fmt.Errorf("case %d: output does not match hidden sequence", caseIdx+1)
		}
		if caseOut != expLines[caseIdx] {
			return fmt.Errorf("case %d: output differs from reference", caseIdx+1)
		}
	}
	return nil
}

func generateTests() []testCase {
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	var tests []testCase

	// Small fixed cases.
	tests = append(tests, buildCase([]string{"()"}))
	tests = append(tests, buildCase([]string{"())(", "()(())", "(())"}))

	// Random cases.
	for i := 0; i < 50; i++ {
		caseCnt := rng.Intn(5) + 1
		var seqs []string
		for j := 0; j < caseCnt; j++ {
			n := rng.Intn(20) + 2
			seqs = append(seqs, randomSeq(rng, n))
		}
		tests = append(tests, buildCase(seqs))
	}

	// Larger stress cases.
	tests = append(tests, buildCase([]string{randomSeq(rng, 200)}))
	tests = append(tests, buildCase([]string{randomSeq(rng, 500), randomSeq(rng, 1000)}))

	return tests
}

func buildCase(sequences []string) testCase {
	var b strings.Builder
	fmt.Fprintf(&b, "%d\n", len(sequences))
	for _, seq := range sequences {
		fmt.Fprintf(&b, "%d\n%s\n", len(seq), seq)
	}
	return testCase{input: b.String()}
}

func randomSeq(rng *rand.Rand, n int) string {
	// Ensure at least one '(' and one ')'.
	if n < 2 {
		n = 2
	}
	bytesSeq := make([]byte, n)
	countOpen := 0
	countClose := 0
	for i := 0; i < n; i++ {
		if rng.Intn(2) == 0 {
			bytesSeq[i] = '('
			countOpen++
		} else {
			bytesSeq[i] = ')'
			countClose++
		}
	}
	if countOpen == 0 {
		bytesSeq[0] = '('
	} else if countClose == 0 {
		bytesSeq[0] = ')'
	}
	return string(bytesSeq)
}
