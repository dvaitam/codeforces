package main

import (
	"bytes"
	"fmt"
	"math/big"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

const (
	refSource2005A = "2005A.go"
	refBinary2005A = "ref2005A.bin"
	totalTests     = 80
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

	refPath, err := buildReference()
	if err != nil {
		fmt.Println("failed to build reference:", err)
		return
	}
	defer os.Remove(refPath)

	tests := generateTests()
	input := formatInput(tests)

	refOut, err := runProgram(refPath, input)
	if err != nil {
		fmt.Printf("reference runtime error: %v\n", err)
		return
	}
	refStrings, err := parseStrings(refOut, len(tests))
	if err != nil {
		fmt.Printf("reference output parse error: %v\noutput:\n%s", err, refOut)
		return
	}

	minCounts := make([]*big.Int, len(tests))
	for i, tc := range tests {
		if len(refStrings[i]) != tc.n {
			fmt.Printf("reference produced invalid length for test %d\n", i+1)
			return
		}
		minCounts[i] = countPalindromicSubseq(refStrings[i])
	}

	candOut, err := runProgram(candidate, input)
	if err != nil {
		fmt.Printf("candidate runtime error: %v\n", err)
		return
	}
	candStrings, err := parseStrings(candOut, len(tests))
	if err != nil {
		fmt.Printf("candidate output parse error: %v\noutput:\n%s", err, candOut)
		return
	}

	vowels := map[rune]bool{'a': true, 'e': true, 'i': true, 'o': true, 'u': true}

	for i, tc := range tests {
		s := candStrings[i]
		if len(s) != tc.n {
			reportFailure(i, "string has incorrect length", input)
			return
		}
		for _, ch := range s {
			if !vowels[ch] {
				reportFailure(i, fmt.Sprintf("string contains non-vowel character %q", ch), input)
				return
			}
		}
		cnt := countPalindromicSubseq(s)
		if cnt.Cmp(minCounts[i]) != 0 {
			reportFailure(i, "string does not achieve minimal palindromic subsequences", input)
			return
		}
	}

	fmt.Printf("All %d tests passed\n", len(tests))
}

func buildReference() (string, error) {
	cmd := exec.Command("go", "build", "-o", refBinary2005A, refSource2005A)
	if out, err := cmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("%v\n%s", err, string(out))
	}
	return filepath.Join(".", refBinary2005A), nil
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

func parseStrings(out string, expected int) ([]string, error) {
	tokens := strings.Fields(out)
	if len(tokens) != expected {
		return nil, fmt.Errorf("expected %d strings, got %d", expected, len(tokens))
	}
	return tokens, nil
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
	rnd := rand.New(rand.NewSource(2005))
	tests := []testCase{{n: 1}, {n: 2}, {n: 3}, {n: 5}, {n: 100}}
	for len(tests) < totalTests {
		tests = append(tests, testCase{n: rnd.Intn(100) + 1})
	}
	return tests
}

func countPalindromicSubseq(s string) *big.Int {
	n := len(s)
	if n == 0 {
		return big.NewInt(1)
	}
	dp := make([][]*big.Int, n)
	for i := 0; i < n; i++ {
		dp[i] = make([]*big.Int, n)
		dp[i][i] = big.NewInt(1)
	}

	for length := 2; length <= n; length++ {
		for i := 0; i+length-1 < n; i++ {
			j := i + length - 1
			if s[i] == s[j] {
				val := new(big.Int).Set(dp[i+1][j])
				val.Add(val, dp[i][j-1])
				val.Add(val, big.NewInt(1))
				dp[i][j] = val
			} else {
				val := new(big.Int).Set(dp[i+1][j])
				val.Add(val, dp[i][j-1])
				val.Sub(val, dp[i+1][j-1])
				dp[i][j] = val
			}
		}
	}

	result := new(big.Int).Set(dp[0][n-1])
	result.Add(result, big.NewInt(1)) // include empty subsequence
	return result
}

func reportFailure(idx int, reason string, input []byte) {
	fmt.Printf("Failure on test %d: %s\n", idx+1, reason)
	fmt.Println("Input used:")
	fmt.Println(string(input))
}
