package main

import (
	"bytes"
	"fmt"
	"math/big"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

const refSource = "./2037G.go"

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierG.go /path/to/candidate")
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
	tmp, err := os.CreateTemp("", "2037G-ref-*")
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
		return fmt.Errorf("expected %d answers, got %d", len(expFields), len(actFields))
	}
	for i := range expFields {
		if expFields[i] == actFields[i] {
			continue
		}
		expVal, ok1 := new(big.Int).SetString(expFields[i], 10)
		actVal, ok2 := new(big.Int).SetString(actFields[i], 10)
		if !ok1 || !ok2 || expVal.Cmp(actVal) != 0 {
			return fmt.Errorf("mismatch at answer %d: expected %s got %s", i+1, expFields[i], actFields[i])
		}
	}
	return nil
}

func buildTests() []string {
	return []string{
		"5\n2 6 3 4 6\n",
		"5\n4 196 2662 2197 121\n",
		"7\n3 6 8 9 11 12 20\n",
		"2\n2 3\n",
		"10\n2 4 6 8 10 12 14 16 18 20\n",
		"10\n6 10 15 21 28 36 45 55 66 78\n",
		"12\n2 3 5 7 11 13 17 19 23 29 31 37\n",
	}
}
