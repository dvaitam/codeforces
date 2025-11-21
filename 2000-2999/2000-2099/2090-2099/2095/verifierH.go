package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
)

const refSource2095H = "2000-2999/2000-2099/2090-2099/2095/2095H.go"

type testCase struct {
	name  string
	input string
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierH.go /path/to/binary")
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
		refAns, err := parseOutput(refOut)
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
		candAns, err := parseOutput(candOut)
		if err != nil {
			fmt.Fprintf(os.Stderr, "candidate output invalid on test %d (%s): %v\ninput:\n%soutput:\n%s",
				idx+1, tc.name, err, tc.input, candOut)
			os.Exit(1)
		}

		if refAns != candAns {
			fmt.Fprintf(os.Stderr, "test %d (%s) failed: expected %q got %q\ninput:\n%sreference output:\n%s\ncandidate output:\n%s",
				idx+1, tc.name, refAns, candAns, tc.input, refOut, candOut)
			os.Exit(1)
		}
	}

	fmt.Printf("All %d tests passed\n", len(tests))
}

func buildReference() (string, func(), error) {
	dir, err := os.MkdirTemp("", "cf-2095H-ref-")
	if err != nil {
		return "", nil, fmt.Errorf("failed to create temp dir: %v", err)
	}
	binPath := filepath.Join(dir, "ref2095H.bin")
	cmd := exec.Command("go", "build", "-o", binPath, refSource2095H)
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		_ = os.RemoveAll(dir)
		return "", nil, fmt.Errorf("failed to build reference: %v\n%s", err, stderr.String())
	}
	cleanup := func() {
		_ = os.RemoveAll(dir)
	}
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

func parseOutput(output string) (string, error) {
	fields := strings.Fields(output)
	if len(fields) != 1 {
		return "", fmt.Errorf("expected single word, got %d tokens", len(fields))
	}
	return fields[0], nil
}

func buildTests() []testCase {
	tests := make([]testCase, 0, 12)
	for x := 1; x <= 11; x++ {
		tests = append(tests, testCase{
			name:  fmt.Sprintf("line_%d", x),
			input: fmt.Sprintf("%d\n", x),
		})
	}
	tests = append(tests, testCase{name: "leading_spaces", input: "   5\n"})
	tests = append(tests, testCase{name: "trailing_spaces", input: "7    \n"})
	tests = append(tests, testCase{name: "newline_only", input: "9"})
	return tests
}
