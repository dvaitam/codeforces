package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

const refSource = "2000-2999/2000-2099/2000-2009/2002/2002F2.go"

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierF2.go /path/to/candidate")
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
		refOut, err := runBinary(refBin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "reference failed on test %d: %v\n", idx+1, err)
			os.Exit(1)
		}
		candOut, err := runBinary(candidate, input)
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
	tmp, err := os.CreateTemp("", "2002F2-ref-*")
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

func runBinary(path, input string) (string, error) {
	cmd := exec.Command(path)
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
		return fmt.Errorf("expected %d lines, got %d", len(expFields), len(actFields))
	}
	for i := range expFields {
		if expFields[i] != actFields[i] {
			return fmt.Errorf("mismatch on line %d: expected %s got %s", i+1, expFields[i], actFields[i])
		}
	}
	return nil
}

func buildTests() []string {
	return []string{
		"8\n3 4 2 5\n4 4 1 4\n6 6 2 2\n7 9 2 3\n8 9 9 1\n2 7 1 4\n5 9 1 4\n5 6 6 7\n",
		"2\n3082823 20000000 1341 331\n20000000 20000000 3 5\n",
		"1\n139 1293 193 412\n",
		"3\n2 3 1 1\n4 6 2 3\n10 15 5 6\n",
		"2\n100 200 1000000000 1\n200 1000 1 1000000000\n",
	}
}
