package main

import (
	"bytes"
	"errors"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
)

const (
	refSource2039A = "2039A.go"
	refBinary2039A = "ref2039A.bin"
	totalTests     = 50
)

type testCase struct {
	n int
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierA.go /path/to/binary")
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

	if _, err := runProgram(ref, input); err != nil {
		fmt.Printf("reference runtime error: %v\n", err)
		return
	}

	output, err := runProgram(candidate, input)
	if err != nil {
		fmt.Printf("candidate runtime error: %v\n", err)
		return
	}

	if err := validateOutput(tests, output); err != nil {
		fmt.Printf("candidate failed validation: %v\n", err)
		fmt.Println("Input used:")
		fmt.Println(string(input))
		return
	}

	fmt.Printf("All %d tests passed\n", len(tests))
}

func buildReference() (string, error) {
	cmd := exec.Command("go", "build", "-o", refBinary2039A, refSource2039A)
	if out, err := cmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("%v\n%s", err, string(out))
	}
	return filepath.Join(".", refBinary2039A), nil
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

func validateOutput(tests []testCase, out string) error {
	tokens := strings.Fields(out)
	pos := 0
	for idx, tc := range tests {
		if pos+tc.n > len(tokens) {
			return fmt.Errorf("test %d: expected %d integers, got %d", idx+1, tc.n, len(tokens)-pos)
		}
		prev := 0
		modSeen := make(map[int]bool)
		for i := 1; i <= tc.n; i++ {
			val, err := strconv.Atoi(tokens[pos])
			if err != nil {
				return fmt.Errorf("test %d: invalid integer %q", idx+1, tokens[pos])
			}
			pos++
			if val < 1 || val > 100 {
				return fmt.Errorf("test %d: value %d out of range [1,100]", idx+1, val)
			}
			if i > 1 && val <= prev {
				return fmt.Errorf("test %d: sequence not strictly increasing at position %d", idx+1, i)
			}
			prev = val
			remainder := val % i
			if modSeen[remainder] {
				return fmt.Errorf("test %d: duplicate remainder %d at position %d", idx+1, remainder, i)
			}
			modSeen[remainder] = true
		}
	}
	if pos != len(tokens) {
		return errors.New("extra tokens in output")
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
	rnd := rand.New(rand.NewSource(2039))
	var tests []testCase
	tests = append(tests, testCase{n: 2}, testCase{n: 50})
	for len(tests) < totalTests {
		tests = append(tests, testCase{n: rnd.Intn(49) + 2})
	}
	return tests
}
