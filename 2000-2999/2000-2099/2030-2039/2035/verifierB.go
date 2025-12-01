package main

import (
	"bytes"
	"errors"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

const (
	refSource2035B = "2035B.go"
	refBinary2035B = "ref2035B.bin"
	totalTests     = 60
)

type testCase struct {
	n int
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierB.go /path/to/binary")
		return
	}
	candidate := os.Args[1]

	ref, err := buildReference()
	if err != nil {
		fmt.Println("failed to build reference:", err)
		return
	}
	defer os.Remove(ref)

	tests := generateTests()
	input := formatInput(tests)

	refOut, err := runProgram(ref, input)
	if err != nil {
		fmt.Printf("reference runtime error: %v\n", err)
		return
	}
	refSeq, err := parseOutput(refOut, len(tests))
	if err != nil {
		fmt.Printf("reference output parse error: %v\noutput:\n%s", err, refOut)
		return
	}

	candOut, err := runProgram(candidate, input)
	if err != nil {
		fmt.Printf("candidate runtime error: %v\n", err)
		return
	}
	candSeq, err := parseOutput(candOut, len(tests))
	if err != nil {
		fmt.Printf("candidate output parse error: %v\noutput:\n%s", err, candOut)
		return
	}

	for i, tc := range tests {
		refAns := refSeq[i]
		candAns := candSeq[i]
		if refAns == "-1" {
			if candAns != "-1" {
				reportMismatch(i, refAns, candAns, input)
				return
			}
			continue
		}
		if candAns == "-1" {
			reportMismatch(i, refAns, candAns, input)
			return
		}
		if err := validateString(tc.n, candAns); err != nil {
			fmt.Printf("test %d failed validation: %v\n", i+1, err)
			fmt.Println("Input used:")
			fmt.Println(string(input))
			return
		}
		if candAns != refAns {
			reportMismatch(i, refAns, candAns, input)
			return
		}
	}

	fmt.Printf("All %d tests passed\n", len(tests))
}

func buildReference() (string, error) {
	cmd := exec.Command("go", "build", "-o", refBinary2035B, refSource2035B)
	if out, err := cmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("%v\n%s", err, string(out))
	}
	return filepath.Join(".", refBinary2035B), nil
}

func runProgram(path string, input []byte) (string, error) {
	cmd := exec.Command(path)
	cmd.Stdin = bytes.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	return out.String(), err
}

func parseOutput(out string, t int) ([]string, error) {
	lines := strings.Fields(out)
	if len(lines) != t {
		return nil, fmt.Errorf("expected %d outputs, got %d", t, len(lines))
	}
	return lines, nil
}

func validateString(n int, s string) error {
	if len(s) != n {
		return fmt.Errorf("length mismatch: expected %d, got %d", n, len(s))
	}
	mod := 0
	for i := 0; i < len(s); i++ {
		ch := s[i]
		if ch != '3' && ch != '6' {
			return fmt.Errorf("invalid digit %q", ch)
		}
		val := int(ch - '0')
		mod = (mod*10 + val) % 66
	}
	if s[len(s)-1] != '6' { // divisibility by 2
		return errors.New("number not divisible by 2")
	}
	if mod != 0 {
		return errors.New("number not divisible by 66")
	}
	return nil
}

func formatInput(tests []testCase) []byte {
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d\n", len(tests))
	for _, tc := range tests {
		fmt.Fprintf(&sb, "%d\n", tc.n)
	}
	return []byte(sb.String())
}

func generateTests() []testCase {
	rnd := rand.New(rand.NewSource(2035))
	var tests []testCase
	tests = append(tests, testCase{n: 1}, testCase{n: 2}, testCase{n: 3}, testCase{n: 500})
	for len(tests) < totalTests {
		tests = append(tests, testCase{n: rnd.Intn(500) + 1})
	}
	return tests
}

func reportMismatch(idx int, refAns, candAns string, input []byte) {
	fmt.Printf("Mismatch on test %d: expected %s, got %s\n", idx+1, refAns, candAns)
	fmt.Println("Input used:")
	fmt.Println(string(input))
}
