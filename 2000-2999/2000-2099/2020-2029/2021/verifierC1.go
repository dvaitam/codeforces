package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

const refSource = "./2021C1.go"

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierC1.go /path/to/candidate")
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
	tmp, err := os.CreateTemp("", "2021C1-ref-*")
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
		e := normalizeAnswer(exp[i])
		a := normalizeAnswer(act[i])
		if e != a {
			return fmt.Errorf("mismatch at answer %d: expected %s got %s", i+1, exp[i], act[i])
		}
	}
	return nil
}

func normalizeAnswer(s string) string {
	s = strings.ToUpper(s)
	if s == "YA" {
		return "YA"
	}
	return "TIDAK"
}

func buildTests() []string {
	return []string{
		"3\n4 2 0\n1 2 3 4\n1 1\n3 6 0\n1 2 3\n1 1 2 3 3 2\n4 6 0\n3 1 4 2\n3 1 1 2 3 4\n",
		"2\n1 1 0\n1\n1\n3 3 0\n2 1 3\n2 3 1\n",
		"1\n5 7 0\n5 4 3 2 1\n1 2 3 4 5 1 2\n",
		"1\n6 3 0\n6 5 4 3 2 1\n6 5 4\n",
	}
}
