package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// refSource points to the local reference solution to avoid GOPATH resolution.
const refSource = "2000E.go"

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierE.go /path/to/candidate")
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
	tmp, err := os.CreateTemp("", "2000E-ref-*")
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
	expFields := strings.Fields(expected)
	actFields := strings.Fields(actual)
	if len(expFields) != len(actFields) {
		return fmt.Errorf("expected %d numbers, got %d", len(expFields), len(actFields))
	}
	for i := range expFields {
		if expFields[i] != actFields[i] {
			return fmt.Errorf("mismatch at position %d: expected %s got %s", i+1, expFields[i], actFields[i])
		}
	}
	return nil
}

func buildTests() []string {
	return []string{
		"5\n3 4 2\n9\n1 1 1 1 1 1 1 1 1\n2 1 1\n2 5\n7 20 1 5 7\n9 4 1\n4\n1 4 5 6\n1 1000000000 898 777\n19 84 1\n1\n45\n4 1 4\n9\n9 5 56 6 7 14 16 16 6\n",
		"1\n1 1 1\n1\n5\n",
		"1\n2 2 2\n3\n5 1 2\n",
		"1\n3 3 2\n5\n3 1 4 1 5\n",
		"2\n2 3 1\n6\n1 2 3 4 5 6\n1 4 1\n3\n10 20 30\n",
	}
}
