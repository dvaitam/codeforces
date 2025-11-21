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
	refSourceE = "2086E.go"
	refBinaryE = "refE.bin"
)

type testCaseE struct {
	l int64
	r int64
	k int64
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierE.go /path/to/binary")
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
	exp, err := parseOutput(refOut, len(tests))
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
		if exp[i] != got[i] {
			fmt.Printf("Mismatch on test %d: expected %d, got %d\n", i+1, exp[i], got[i])
			fmt.Println("Input used:")
			fmt.Println(string(input))
			return
		}
	}
	fmt.Printf("All %d tests passed\n", len(tests))
}

func buildReference() (string, error) {
	cmd := exec.Command("go", "build", "-o", refBinaryE, refSourceE)
	if out, err := cmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("%v\n%s", err, string(out))
	}
	return filepath.Join(".", refBinaryE), nil
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
	res := make([]int64, t)
	for i, tok := range fields {
		val, err := strconv.ParseInt(tok, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("token %d is not an integer: %v", i+1, err)
		}
		res[i] = val
	}
	return res, nil
}

func formatInput(tests []testCaseE) []byte {
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d\n", len(tests))
	for _, tc := range tests {
		fmt.Fprintf(&sb, "%d %d %d\n", tc.l, tc.r, tc.k)
	}
	return []byte(sb.String())
}

func generateTests() []testCaseE {
	rnd := rand.New(rand.NewSource(2086))
	var tests []testCaseE

	seedCases := []testCaseE{
		{1, 100, 3},
		{1, 1, 1},
		{1, 1, 2},
		{1, 1_000_000_000_000_000_000, 1},
		{1, 1_000_000_000_000_000_000, 20},
		{1234567, 123456789101112131, 12},
		{42, 42, 10},
		{999999999999999999, 1_000_000_000_000_000_000, 1},
	}
	tests = append(tests, seedCases...)

	for len(tests) < 250 {
		l := rnd.Int63n(1_000_000_000_000_000_000) + 1
		r := l + rnd.Int63n(1_000_000_000_000_000_000-l+1)
		k := rnd.Int63n(1_000_000_000_000_000_000) + 1
		tests = append(tests, testCaseE{
			l: l,
			r: r,
			k: k,
		})
	}
	return tests
}
