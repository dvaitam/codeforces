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
	refSource2069A = "2069A.go"
	refBinary2069A = "ref2069A.bin"
	maxTests       = 400
)

type testCase struct {
	n int
	b []int
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
			fmt.Printf("Mismatch on case %d: expected %s, got %s\n", i+1, expected[i], got[i])
			fmt.Println("Input used:")
			fmt.Println(string(input))
			return
		}
	}
	fmt.Printf("All %d tests passed\n", len(tests))
}

func buildReference() (string, error) {
	cmd := exec.Command("go", "build", "-o", refBinary2069A, refSource2069A)
	if out, err := cmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("%v\n%s", err, string(out))
	}
	return filepath.Join(".", refBinary2069A), nil
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
	words := strings.Fields(out)
	if len(words) != t {
		return nil, fmt.Errorf("expected %d answers, got %d", t, len(words))
	}
	for i := range words {
		words[i] = strings.ToUpper(words[i])
		if words[i] != "YES" && words[i] != "NO" {
			return nil, fmt.Errorf("invalid answer at case %d: %s", i+1, words[i])
		}
	}
	return words, nil
}

func formatInput(tests []testCase) []byte {
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d\n", len(tests))
	for _, tc := range tests {
		fmt.Fprintf(&sb, "%d\n", tc.n)
		if len(tc.b) == 0 {
			sb.WriteByte('\n')
			continue
		}
		for i, v := range tc.b {
			if i > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(strconv.Itoa(v))
		}
		sb.WriteByte('\n')
	}
	return []byte(sb.String())
}

func generateTests() []testCase {
	rnd := rand.New(rand.NewSource(2069))
	var tests []testCase

	addCase := func(n int, b []int) {
		cpy := make([]int, len(b))
		copy(cpy, b)
		tests = append(tests, testCase{n: n, b: cpy})
	}

	addCase(3, nil)
	addCase(4, []int{0, 0})
	addCase(4, []int{1, 1})
	addCase(5, []int{0, 1, 0})
	addCase(6, []int{1, 0, 1, 0})

	for len(tests) < maxTests {
		n := rnd.Intn(98) + 3
		m := n - 2
		b := make([]int, m)
		for i := 0; i < m; i++ {
			if rnd.Intn(100) < 30 {
				b[i] = 1
			}
		}
		addCase(n, b)
	}
	return tests
}
