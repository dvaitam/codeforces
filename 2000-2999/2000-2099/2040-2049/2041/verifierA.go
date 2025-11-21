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

const refSourceA = "2000-2999/2000-2099/2040-2049/2041/2041A.go"

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierA.go /path/to/candidate")
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
		fmt.Fprintf(os.Stderr, "reference failed: %v\n%s\n", err, refOut)
		os.Exit(1)
	}

	candOut, err := runCandidate(candidate, input)
	if err != nil {
		fmt.Fprintf(os.Stderr, "candidate failed: %v\n%s\n", err, candOut)
		os.Exit(1)
	}

	expected, err := parseOutputs(refOut, len(tests))
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
		if got[i] != expected[i] {
			fmt.Fprintf(os.Stderr, "wrong answer on test %d: expected %d got %d (input=%v)\n", i+1, expected[i], got[i], tests[i])
			os.Exit(1)
		}
	}

	fmt.Printf("All %d tests passed.\n", len(tests))
}

func buildReference() (string, error) {
	tmp, err := os.CreateTemp("", "2041A-ref-*")
	if err != nil {
		return "", err
	}
	tmp.Close()
	cmd := exec.Command("go", "build", "-o", tmp.Name(), filepath.Clean(refSourceA))
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

func parseOutputs(output string, t int) ([]int, error) {
	tokens := strings.Fields(output)
	if len(tokens) < t {
		return nil, fmt.Errorf("expected %d outputs, got %d", t, len(tokens))
	}
	if len(tokens) > t {
		return nil, fmt.Errorf("extra output detected starting with %q", tokens[t])
	}
	res := make([]int, t)
	for i := 0; i < t; i++ {
		val, err := strconv.Atoi(tokens[i])
		if err != nil {
			return nil, fmt.Errorf("token %q is not an integer", tokens[i])
		}
		if val < 1 || val > 5 {
			return nil, fmt.Errorf("output %d is not in [1,5]", val)
		}
		res[i] = val
	}
	return res, nil
}

func generateTests() [][]int {
	var tests [][]int
	numbers := []int{1, 2, 3, 4, 5}
	var dfs func(path []int, used []bool)
	dfs = func(path []int, used []bool) {
		if len(path) == 4 {
			tmp := make([]int, 4)
			copy(tmp, path)
			tests = append(tests, tmp)
			return
		}
		for i, v := range numbers {
			if !used[i] {
				used[i] = true
				dfs(append(path, v), used)
				used[i] = false
			}
		}
	}
	dfs([]int{}, make([]bool, len(numbers)))
	return tests
}

func buildInput(tests [][]int) string {
	var b strings.Builder
	fmt.Fprintf(&b, "%d\n", len(tests))
	for _, tc := range tests {
		fmt.Fprintf(&b, "%d %d %d %d\n", tc[0], tc[1], tc[2], tc[3])
	}
	return b.String()
}
