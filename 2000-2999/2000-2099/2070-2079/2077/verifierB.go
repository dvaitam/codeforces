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
	refSourceB = "2077B.go"
	refBinaryB = "refB.bin"
)

type testCaseB struct {
	x uint32
	y uint32
	m uint32
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierB.go /path/to/binary")
		return
	}
	candidate := os.Args[1]

	refBin, err := buildReference()
	if err != nil {
		fmt.Println("failed to build reference:", err)
		return
	}
	defer os.Remove(refBin)

	tests := generateTests()
	input := formatInput(tests)

	refOut, err := runProgram(refBin, input)
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
	cmd := exec.Command("go", "build", "-o", refBinaryB, refSourceB)
	if out, err := cmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("%v\n%s", err, string(out))
	}
	return filepath.Join(".", refBinaryB), nil
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

func parseOutput(out string, t int) ([]uint64, error) {
	fields := strings.Fields(out)
	if len(fields) != t {
		return nil, fmt.Errorf("expected %d numbers, got %d", t, len(fields))
	}
	res := make([]uint64, t)
	for i, tok := range fields {
		val, err := strconv.ParseUint(tok, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("token %d is not an integer: %v", i+1, err)
		}
		res[i] = val
	}
	return res, nil
}

func formatInput(tests []testCaseB) []byte {
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d\n", len(tests))
	for _, tc := range tests {
		fmt.Fprintf(&sb, "%d %d %d\n", tc.x, tc.y, tc.m)
	}
	return []byte(sb.String())
}

func generateTests() []testCaseB {
	rnd := rand.New(rand.NewSource(2077))
	var tests []testCaseB

	edge := []testCaseB{
		{0, 0, 0},
		{0, 0, (1 << 30) - 1},
		{(1 << 30) - 1, 0, 0},
		{(1 << 30) - 1, (1 << 30) - 1, (1 << 30) - 1},
		{0, (1 << 30) - 1, 123456789},
		{1, 2, 3},
	}
	tests = append(tests, edge...)

	for len(tests) < 300 {
		tc := testCaseB{
			x: rnd.Uint32() & ((1 << 30) - 1),
			y: rnd.Uint32() & ((1 << 30) - 1),
			m: rnd.Uint32() & ((1 << 30) - 1),
		}
		tests = append(tests, tc)
	}
	return tests
}
