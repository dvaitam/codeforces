package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
)

const (
	refSource2072A = "2072A.go"
	refBinary2072A = "ref2072A.bin"
	maxTestsA      = 400
)

type testCaseA struct {
	n int
	k int
	p int
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

	refOut, err := runProgram(ref, input)
	if err != nil {
		fmt.Printf("reference runtime error: %v\n", err)
		return
	}
	expected, err := parseOutput(refOut, len(tests))
	if err != nil {
		fmt.Printf("reference output parse error: %v\noutput:\n%s", err, refOut)
		return
	}

	candOut, err := runProgram(candidate, input)
	if err != nil {
		fmt.Printf("candidate runtime error: %v\n", err)
		return
	}
	got, err := parseOutput(candOut, len(tests))
	if err != nil {
		fmt.Printf("candidate output parse error: %v\noutput:\n%s", err, candOut)
		return
	}

	for i := range tests {
		if expected[i] != got[i] {
			fmt.Printf("Mismatch on case %d: expected %d, got %d\n", i+1, expected[i], got[i])
			fmt.Println("Input used:")
			fmt.Println(string(input))
			return
		}
	}
	fmt.Printf("All %d tests passed\n", len(tests))
}

func buildReference() (string, error) {
	cmd := exec.Command("go", "build", "-o", refBinary2072A, refSource2072A)
	if out, err := cmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("%v\n%s", err, string(out))
	}
	return filepath.Join(".", refBinary2072A), nil
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

func parseOutput(out string, t int) ([]int64, error) {
	fields := strings.Fields(out)
	if len(fields) != t {
		return nil, fmt.Errorf("expected %d numbers, got %d", t, len(fields))
	}
	ans := make([]int64, t)
	for i, tok := range fields {
		val, err := strconv.ParseInt(tok, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("token %d not an integer: %v", i+1, err)
		}
		ans[i] = val
	}
	return ans, nil
}

func formatInput(tests []testCaseA) []byte {
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d\n", len(tests))
	for _, tc := range tests {
		fmt.Fprintf(&sb, "%d %d %d\n", tc.n, tc.k, tc.p)
	}
	return []byte(sb.String())
}

func generateTests() []testCaseA {
	rnd := rand.New(rand.NewSource(2072))
	var tests []testCaseA

	seed := []testCaseA{
		{n: 1, k: 0, p: 1},
		{n: 1, k: 1, p: 1},
		{n: 1, k: 2, p: 1},
		{n: 50, k: 2500, p: 50},
		{n: 50, k: -2500, p: 50},
		{n: 5, k: 7, p: 2},
		{n: 3, k: -7, p: 7},
	}
	tests = append(tests, seed...)

	for len(tests) < maxTestsA {
		n := rnd.Intn(50) + 1
		k := rnd.Intn(5001) - 2500
		p := rnd.Intn(50) + 1
		tests = append(tests, testCaseA{n: n, k: k, p: p})
	}
	return tests
}
