package main

import (
	"bytes"
	"fmt"
	"math"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

const refSource = "./1666B.go"

type testCase struct {
	input       string
	expectLines int
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierB.go /path/to/candidate")
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
	for idx, tc := range tests {
		refOut, err := runBinary(refBin, tc.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "reference failed on test %d: %v\n", idx+1, err)
			os.Exit(1)
		}
		candOut, err := runCandidate(candidate, tc.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "candidate failed on test %d: %v\n", idx+1, err)
			os.Exit(1)
		}
		if err := compareOutputs(refOut, candOut, tc.expectLines); err != nil {
			fmt.Fprintf(os.Stderr, "test %d failed: %v\nInput:\n%s", idx+1, err, tc.input)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed.\n", len(tests))
}

func buildReference() (string, error) {
	tmp, err := os.CreateTemp("", "1666B-ref-*")
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

func runCandidate(path, input string) (string, error) {
	cmd := commandFor(path)
	return runWithInput(cmd, input)
}

func runBinary(path, input string) (string, error) {
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

func parseFloats(out string, expected int) ([]float64, error) {
	fields := strings.Fields(out)
	if len(fields) < expected {
		return nil, fmt.Errorf("expected %d numbers, got %d", expected, len(fields))
	}
	values := make([]float64, expected)
	for i := 0; i < expected; i++ {
		var v float64
		if _, err := fmt.Sscan(fields[i], &v); err != nil {
			return nil, fmt.Errorf("invalid float %q: %v", fields[i], err)
		}
		values[i] = v
	}
	return values, nil
}

func compareOutputs(refOut, candOut string, expected int) error {
	refVals, err := parseFloats(refOut, expected)
	if err != nil {
		return fmt.Errorf("reference parsing error: %v", err)
	}
	candVals, err := parseFloats(candOut, expected)
	if err != nil {
		return fmt.Errorf("candidate parsing error: %v", err)
	}
	for i := 0; i < expected; i++ {
		if !closeEnough(refVals[i], candVals[i]) {
			return fmt.Errorf("mismatch on line %d: expected %.10f got %.10f", i+1, refVals[i], candVals[i])
		}
	}
	return nil
}

func closeEnough(a, b float64) bool {
	diff := math.Abs(a - b)
	limit := 1e-6 + 1e-6*math.Max(1, math.Abs(a))
	return diff <= limit
}

func buildTests() []testCase {
	return []testCase{
		{
			input:       "1 5\n3 1 7 10 700 400 100\n0 2 10 50 102\n",
			expectLines: 5,
		},
		{
			input:       "2 5\n3 10 70 100 700 400 100\n3 10 30 100 700 400 100\n2 10 50 70 110\n",
			expectLines: 5,
		},
		{
			input:       "1 3\n2 5 5 1 1\n0 10 100\n",
			expectLines: 3,
		},
	}
}
