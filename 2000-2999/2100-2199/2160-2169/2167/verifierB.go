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
	refSource2167B = "2167B.go"
	refBinary2167B = "ref2167B.bin"
	maxTests       = 400
)

type testCase struct {
	s string
	t string
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
			fmt.Printf("Mismatch on case %d: expected %s, got %s\n", i+1, expected[i], got[i])
			fmt.Println("Input used:")
			fmt.Println(string(input))
			return
		}
	}
	fmt.Printf("All %d tests passed\n", len(tests))
}

func buildReference() (string, error) {
	cmd := exec.Command("go", "build", "-o", refBinary2167B, refSource2167B)
	if out, err := cmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("%v\n%s", err, string(out))
	}
	return filepath.Join(".", refBinary2167B), nil
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
		return nil, fmt.Errorf("expected %d answers, got %d", t, len(lines))
	}
	for i := range lines {
		lines[i] = strings.ToUpper(lines[i])
		if lines[i] != "YES" && lines[i] != "NO" {
			return nil, fmt.Errorf("invalid answer at case %d: %s", i+1, lines[i])
		}
	}
	return lines, nil
}

func formatInput(tests []testCase) []byte {
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d\n", len(tests))
	for _, tc := range tests {
		fmt.Fprintf(&sb, "%d\n%s %s\n", len(tc.s), tc.s, tc.t)
	}
	return []byte(sb.String())
}

func generateTests() []testCase {
	rnd := rand.New(rand.NewSource(2167))
	var tests []testCase

	addCase := func(s, t string) {
		tests = append(tests, testCase{s: s, t: t})
	}

	addCase("a", "a")
	addCase("a", "b")
	addCase("abc", "bca")
	addCase("abc", "abd")
	addCase("zzzz", "zzzz")

	for len(tests) < maxTests {
		n := rnd.Intn(20) + 1
		var sb1, sb2 strings.Builder
		for i := 0; i < n; i++ {
			sb1.WriteByte(byte('a' + rnd.Intn(26)))
		}
		bytes1 := []byte(sb1.String())
		shuffled := append([]byte(nil), bytes1...)
		rnd.Shuffle(len(shuffled), func(i, j int) {
			shuffled[i], shuffled[j] = shuffled[j], shuffled[i]
		})
		if rnd.Intn(2) == 1 {
			shuffled[rnd.Intn(len(shuffled))] = byte('a' + rnd.Intn(26))
		}
		sb2.Write(shuffled)
		addCase(string(bytes1), sb2.String())
	}
	return tests
}
