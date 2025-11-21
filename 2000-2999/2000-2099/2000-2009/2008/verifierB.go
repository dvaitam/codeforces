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

const refSource = "2000-2999/2000-2099/2000-2009/2008/2008B.go"

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierB.go /path/to/candidate")
		os.Exit(1)
	}
	candidate := os.Args[len(os.Args)-1]

	refBin, err := buildReference()
	if err != nil {
		fmt.Fprintln(os.Stderr, "failed to build reference:", err)
		os.Exit(1)
	}
	defer os.Remove(refBin)

	tests := generateTests()
	input := buildInput(tests)

	refOut, err := runProgram(refBin, input)
	if err != nil {
		fmt.Fprintf(os.Stderr, "reference failed: %v\noutput:\n%s\n", err, refOut)
		os.Exit(1)
	}

	candOut, err := runCandidate(candidate, input)
	if err != nil {
		fmt.Fprintf(os.Stderr, "candidate failed: %v\noutput:\n%s\n", err, candOut)
		os.Exit(1)
	}

	exp, err := parseOutputs(refOut, len(tests))
	if err != nil {
		fmt.Fprintf(os.Stderr, "could not parse reference output: %v\n", err)
		os.Exit(1)
	}
	got, err := parseOutputs(candOut, len(tests))
	if err != nil {
		fmt.Fprintf(os.Stderr, "could not parse candidate output: %v\n", err)
		os.Exit(1)
	}

	for i := range tests {
		if exp[i] != got[i] {
			fmt.Fprintf(os.Stderr, "wrong answer on test %d: expected %q got %q\n", i+1, exp[i], got[i])
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed.\n", len(tests))
}

func buildReference() (string, error) {
	tmp, err := os.CreateTemp("", "2008B-ref-*")
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

func runCandidate(path, input string) (string, error) {
	cmd := commandFor(path)
	return runWithInput(cmd, input)
}

func runProgram(path, input string) (string, error) {
	cmd := exec.Command(path)
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

type testCase struct {
	n int
	s string
}

func generateTests() []testCase {
	rng := rand.New(rand.NewSource(2008))
	var tests []testCase
	total := 0
	add := func(s string) {
		if len(s) == 0 {
			return
		}
		if total+len(s) > 200000 {
			return
		}
		tests = append(tests, testCase{n: len(s), s: s})
		total += len(s)
	}

	samples := []string{
		"1111",
		"111101111",
		"1111101111",
		"1010",
		"1111111",
		"10001",
	}
	for _, s := range samples {
		add(s)
	}

	for total < 200000 {
		side := rng.Intn(40) + 2
		add(makeSquareString(side))
		if len(tests) > 1500 {
			break
		}
	}

	for total < 200000 {
		n := rng.Intn(200) + 2
		str := randomString(rng, n)
		add(str)
	}

	return tests
}

func makeSquareString(m int) string {
	var b strings.Builder
	for i := 0; i < m; i++ {
		for j := 0; j < m; j++ {
			if i == 0 || i == m-1 || j == 0 || j == m-1 {
				b.WriteByte('1')
			} else {
				b.WriteByte('0')
			}
		}
	}
	return b.String()
}

func randomString(rng *rand.Rand, n int) string {
	var b strings.Builder
	for i := 0; i < n; i++ {
		if rng.Intn(2) == 0 {
			b.WriteByte('0')
		} else {
			b.WriteByte('1')
		}
	}
	return b.String()
}

func buildInput(tests []testCase) string {
	var b strings.Builder
	fmt.Fprintf(&b, "%d\n", len(tests))
	for _, tc := range tests {
		fmt.Fprintf(&b, "%d\n%s\n", tc.n, tc.s)
	}
	return b.String()
}

func parseOutputs(output string, t int) ([]string, error) {
	tokens := strings.Fields(output)
	if len(tokens) < t {
		return nil, fmt.Errorf("expected %d outputs, got %d", t, len(tokens))
	}
	if len(tokens) > t {
		return nil, fmt.Errorf("extra output detected starting at token %q", tokens[t])
	}
	return tokens, nil
}
