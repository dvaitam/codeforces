package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

type testCase struct {
	row int
	col int
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	target := os.Args[1]

	refBin, err := buildReference()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to build reference: %v\n", err)
		os.Exit(1)
	}
	defer os.Remove(refBin)

	tests := generateTests()
	for i, tc := range tests {
		input := fmt.Sprintf("%d %d\n", tc.row, tc.col)
		refOut, err := runProgram(refBin, input)
		if err != nil {
			fmt.Printf("Reference runtime error on test %d (%d,%d): %v\nOutput:\n%s\n", i+1, tc.row, tc.col, err, refOut)
			os.Exit(1)
		}
		exp, err := parseAnswer(refOut)
		if err != nil {
			fmt.Printf("Failed to parse reference output on test %d (%d,%d): %v\nOutput:\n%s\n", i+1, tc.row, tc.col, err, refOut)
			os.Exit(1)
		}

		out, err := runProgram(target, input)
		if err != nil {
			fmt.Printf("Runtime error on test %d (%d,%d): %v\nInput:\n%sOutput:\n%s\n", i+1, tc.row, tc.col, err, input, out)
			os.Exit(1)
		}
		got, err := parseAnswer(out)
		if err != nil {
			fmt.Printf("Failed to parse output on test %d (%d,%d): %v\nInput:\n%sOutput:\n%s\n", i+1, tc.row, tc.col, err, input, out)
			os.Exit(1)
		}
		if got != exp {
			fmt.Printf("Wrong answer on test %d (%d,%d): expected %s got %s\nInput:\n%s", i+1, tc.row, tc.col, exp, got, input)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed.\n", len(tests))
}

func buildReference() (string, error) {
	path := "./ref1331E.bin"
	cmd := exec.Command("go", "build", "-o", path, "1331E.go")
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("build failed: %v\n%s", err, stderr.String())
	}
	return path, nil
}

func runProgram(bin, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return out.String(), fmt.Errorf("runtime error: %v", err)
	}
	return out.String(), nil
}

func parseAnswer(out string) (string, error) {
	tokens := strings.Fields(out)
	if len(tokens) != 1 {
		return "", fmt.Errorf("expected single token, got %d", len(tokens))
	}
	ans := strings.ToUpper(tokens[0])
	if ans != "IN" && ans != "OUT" {
		return "", fmt.Errorf("invalid answer %q", tokens[0])
	}
	return ans, nil
}

func generateTests() []testCase {
	var tests []testCase
	for r := 0; r < 64; r++ {
		for c := 0; c < 64; c++ {
			tests = append(tests, testCase{row: r, col: c})
		}
	}
	return tests
}
