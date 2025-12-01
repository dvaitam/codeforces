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

const refSource = "./987B.go"

type testCase struct {
	input string
	desc  string
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}

	refBin, err := buildReference()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to build reference: %v\n", err)
		os.Exit(1)
	}
	defer os.Remove(refBin)

	tests := generateTests()
	candidate := os.Args[1]
	for i, tc := range tests {
		expect, err := runBinary(refBin, tc.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "reference failed on test %d (%s): %v\ninput:\n%s", i+1, tc.desc, err, tc.input)
			os.Exit(1)
		}
		got, err := runProgram(candidate, tc.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "runtime error on test %d (%s): %v\ninput:\n%soutput:\n%s\n", i+1, tc.desc, err, tc.input, got)
			os.Exit(1)
		}
		if !equalVerdict(expect, got) {
			fmt.Fprintf(os.Stderr, "wrong answer on test %d (%s)\ninput:\n%s\nexpected:\n%s\ngot:\n%s\n", i+1, tc.desc, tc.input, strings.TrimSpace(expect), strings.TrimSpace(got))
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed.\n", len(tests))
}

func equalVerdict(a, b string) bool {
	trimA := strings.TrimSpace(a)
	trimB := strings.TrimSpace(b)
	if trimA == "" || trimB == "" {
		return false
	}
	va := trimA[0]
	vb := trimB[0]
	return (va == '<' || va == '>' || va == '=') && va == vb
}

func buildReference() (string, error) {
	tmp, err := os.CreateTemp("", "987B-ref-*")
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

func runProgram(path, input string) (string, error) {
	cmd := commandFor(path)
	return runWithInput(cmd, input)
}

func runBinary(path, input string) (string, error) {
	cmd := exec.Command(path)
	return runWithInput(cmd, input)
}

func runWithInput(cmd *exec.Cmd, input string) (string, error) {
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return out.String(), err
}

func generateTests() []testCase {
	var tests []testCase
	tests = append(tests, testCase{"5 8\n", "sample1"})
	tests = append(tests, testCase{"10 3\n", "sample2"})
	tests = append(tests, testCase{"6 6\n", "sample3"})
	tests = append(tests, testCase{"1 1\n", "both ones"})
	tests = append(tests, testCase{"1 1000000000\n", "x=1 small"})
	tests = append(tests, testCase{"1000000000 1\n", "y=1 small"})
	tests = append(tests, testCase{"2 4\n", "classic"})
	tests = append(tests, testCase{"4 2\n", "classic swap"})
	tests = append(tests, testCase{"2 3\n", "small pair"})
	tests = append(tests, testCase{"3 2\n", "small pair swap"})
	tests = append(tests, testCase{"2 2\n", "equal small"})
	tests = append(tests, testCase{"999999937 999999937\n", "large equal"})

	rng := rand.New(rand.NewSource(98709987))
	for i := 0; i < 200; i++ {
		x := rng.Int63n(1_000_000_000) + 1
		y := rng.Int63n(1_000_000_000) + 1
		tests = append(tests, testCase{
			input: fmt.Sprintf("%d %d\n", x, y),
			desc:  fmt.Sprintf("rand-%d", i+1),
		})
	}
	return tests
}
