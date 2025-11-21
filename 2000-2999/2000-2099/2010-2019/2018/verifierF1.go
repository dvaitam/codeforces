package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

const refSource = "2000-2999/2000-2099/2010-2019/2018/2018F1.go"

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierF1.go /path/to/candidate")
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
	tmp, err := os.CreateTemp("", "2018F1-ref-*")
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
	exp := normalizeLines(expected)
	act := normalizeLines(actual)
	if len(exp) != len(act) {
		return fmt.Errorf("expected %d lines, got %d", len(exp), len(act))
	}
	for i := range exp {
		expFields := strings.Fields(exp[i])
		actFields := strings.Fields(act[i])
		if len(expFields) != len(actFields) {
			return fmt.Errorf("line %d: expected %d numbers, got %d", i+1, len(expFields), len(actFields))
		}
		for j := range expFields {
			if expFields[j] != actFields[j] {
				return fmt.Errorf("line %d position %d mismatch: expected %s got %s", i+1, j+1, expFields[j], actFields[j])
			}
		}
	}
	return nil
}

func normalizeLines(out string) []string {
	out = strings.TrimSpace(out)
	if out == "" {
		return []string{}
	}
	return strings.Split(out, "\n")
}

func buildTests() []string {
	return []string{
		"1\n1 998244353\n",
		"1\n2 998244353\n",
		"1\n3 998244353\n",
		"2\n4 998244353\n5 998244353\n",
		"3\n6 102275857\n8 999999937\n10 100000007\n",
		"1\n12 100000007\n",
	}
}
