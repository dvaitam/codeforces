package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

const (
	refSource2091G = "2091G.go"
	refBinary2091G = "ref2091G.bin"
	maxTests       = 180
	maxTotalK      = 2000
)

type testCase struct {
	s int64
	k int
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierG.go /path/to/binary")
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

	for i := range expected {
		if expected[i] != got[i] {
			fmt.Printf("Mismatch on test %d: expected %d, got %d\n", i+1, expected[i], got[i])
			fmt.Println("Input used:")
			fmt.Println(string(input))
			return
		}
	}
	fmt.Printf("All %d tests passed\n", len(tests))
}

func buildReference() (string, error) {
	cmd := exec.Command("go", "build", "-o", refBinary2091G, refSource2091G)
	if out, err := cmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("%v\n%s", err, string(out))
	}
	return filepath.Join(".", refBinary2091G), nil
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

func parseOutput(out string, t int) ([]int, error) {
	fields := strings.Fields(out)
	if len(fields) != t {
		return nil, fmt.Errorf("expected %d outputs, got %d", t, len(fields))
	}
	res := make([]int, t)
	for i, tok := range fields {
		if _, err := fmt.Sscanf(tok, "%d", &res[i]); err != nil {
			return nil, fmt.Errorf("token %d: %v", i+1, err)
		}
	}
	return res, nil
}

func formatInput(tests []testCase) []byte {
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d\n", len(tests))
	for _, tc := range tests {
		fmt.Fprintf(&sb, "%d %d\n", tc.s, tc.k)
	}
	return []byte(sb.String())
}

func generateTests() []testCase {
	rnd := rand.New(rand.NewSource(time.Now().UnixNano()))
	var tests []testCase
	totalK := 0

	add := func(tc testCase) {
		tests = append(tests, tc)
		totalK += tc.k
	}

	// Deterministic small cases
	add(testCase{s: 1, k: 1})
	add(testCase{s: 2, k: 1})
	add(testCase{s: 5, k: 5})
	add(testCase{s: 9, k: 6})
	add(testCase{s: 10, k: 7})
	add(testCase{s: 24, k: 21})
	add(testCase{s: 123456, k: 777})
	add(testCase{s: 776, k: 499})

	for len(tests) < maxTests && totalK < maxTotalK {
		remaining := maxTotalK - totalK
		k := rnd.Intn(minInt(1000, remaining)) + 1
		if k > remaining {
			k = remaining
		}
		sOptions := rnd.Intn(4)
		var s int64
		switch sOptions {
		case 0:
			s = int64(k)
		case 1:
			s = int64(k + rnd.Intn(10))
		case 2:
			s = int64(k + rnd.Intn(1000))
		default:
			s = int64(rnd.Intn(1_000_000_000-k) + k)
		}
		add(testCase{s: s, k: k})
	}
	return tests
}

func minInt(a, b int) int {
	if a < b {
		return a
	}
	return b
}
