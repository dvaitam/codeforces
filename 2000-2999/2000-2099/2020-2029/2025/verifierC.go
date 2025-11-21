package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

const refSource = "2000-2999/2000-2099/2020-2029/2025/2025C.go"

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierC.go /path/to/candidate")
		os.Exit(1)
	}
	candidate := os.Args[1]

	refBin, err := buildReference()
	if err != nil {
		fmt.Fprintln(os.Stderr, "failed to build reference:", err)
		os.Exit(1)
	}
	defer os.Remove(refBin)

	tests := buildTests()
	for idx, input := range tests {
		refOut, err := runProgram(refBin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "reference failed on test %d: %v\n", idx+1, err)
			os.Exit(1)
		}
		candOut, err := runProgram(candidate, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "candidate failed on test %d: %v\n", idx+1, err)
			os.Exit(1)
		}
		if err := compareOutputs(refOut, candOut); err != nil {
			fmt.Fprintf(os.Stderr, "test %d failed: %v\nInput:\n%s", idx+1, err, input)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed.\n", len(tests))
}

func buildReference() (string, error) {
	tmp, err := os.CreateTemp("", "2025C-ref-*")
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

func runProgram(target, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(target, ".go") {
		cmd = exec.Command("go", "run", target)
	} else {
		cmd = exec.Command(target)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	return out.String(), nil
}

func compareOutputs(expected, actual string) error {
	exp := strings.Fields(expected)
	act := strings.Fields(actual)
	if len(exp) != len(act) {
		return fmt.Errorf("expected %d outputs, got %d", len(exp), len(act))
	}
	for i := range exp {
		if exp[i] != act[i] {
			return fmt.Errorf("mismatch at answer %d: expected %s got %s", i+1, exp[i], act[i])
		}
	}
	return nil
}

func buildTests() []string {
	return []string{
		"4\n10 2\n5 2 4 3 4 3 4 5 3 2\n5 1\n10 11 10 11 10\n9 3\n4 5 4 4 6 5 4 4 6\n3 2\n1 3 1\n",
		"3\n1 1\n100\n5 2\n1 2 3 4 5\n6 3\n1 2 2 2 3 3\n",
		"2\n7 3\n3 3 3 4 4 5 5\n8 4\n1 1 2 2 3 3 4 4\n",
		"1\n12 5\n5 5 5 5 6 6 6 7 7 8 8 8\n",
	}
}
