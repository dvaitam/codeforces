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
	refSource2072B = "2072B.go"
	refBinary2072B = "ref2072B.bin"
	maxTotalN      = 200000
	maxTests       = 250
)

type testCase struct {
	s string
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
	cmd := exec.Command("go", "build", "-o", refBinary2072B, refSource2072B)
	if out, err := cmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("%v\n%s", err, string(out))
	}
	return filepath.Join(".", refBinary2072B), nil
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
			return nil, fmt.Errorf("token %d not an integer: %v", i+1, err)
		}
		res[i] = val
	}
	return res, nil
}

func formatInput(tests []testCase) []byte {
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d\n", len(tests))
	for _, tc := range tests {
		fmt.Fprintf(&sb, "%d\n%s\n", len(tc.s), tc.s)
	}
	return []byte(sb.String())
}

func generateTests() []testCase {
	rnd := rand.New(rand.NewSource(2072))
	var tests []testCase
	totalN := 0

	addCase := func(s string) {
		if len(s) == 0 {
			s = "-"
		}
		tests = append(tests, testCase{s: s})
		totalN += len(s)
	}

	addCase("-")
	addCase("_")
	addCase("-_-")
	addCase("--___")
	addCase(strings.Repeat("-", 5))
	addCase(strings.Repeat("_", 5))
	addCase("-__-__-__-__")

	for len(tests) < maxTests && totalN < maxTotalN {
		maxLen := maxTotalN - totalN
		if maxLen <= 0 {
			break
		}
		limit := 1000
		if maxLen < limit {
			limit = maxLen
		}
		n := rnd.Intn(limit) + 1
		var sb strings.Builder
		for i := 0; i < n; i++ {
			if rnd.Intn(2) == 0 {
				sb.WriteByte('-')
			} else {
				sb.WriteByte('_')
			}
		}
		addCase(sb.String())
	}
	return tests
}
