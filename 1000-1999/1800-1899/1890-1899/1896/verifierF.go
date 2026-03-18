package main

import (
	"bufio"
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

func generateInput() []byte {
	r := rand.New(rand.NewSource(42))
	t := 100
	var buf bytes.Buffer
	fmt.Fprintln(&buf, t)
	for i := 0; i < t; i++ {
		n := r.Intn(4) + 1
		fmt.Fprintln(&buf, n)
		var sb strings.Builder
		for j := 0; j < 2*n; j++ {
			if r.Intn(2) == 0 {
				sb.WriteByte('0')
			} else {
				sb.WriteByte('1')
			}
		}
		fmt.Fprintln(&buf, sb.String())
	}
	return buf.Bytes()
}

// isBalanced checks if b is a balanced bracket sequence.
func isBalanced(b string) bool {
	depth := 0
	for _, c := range b {
		if c == '(' {
			depth++
		} else if c == ')' {
			depth--
			if depth < 0 {
				return false
			}
		} else {
			return false
		}
	}
	return depth == 0
}

// applyOperation applies one bracket-sequence operation to binary string s (as []byte of '0'/'1').
func applyOperation(s []byte, b string) {
	n := len(s)
	// Find matching parentheses
	match := make([]int, n)
	stack := make([]int, 0, n)
	for i := 0; i < n; i++ {
		if b[i] == '(' {
			stack = append(stack, i)
		} else {
			j := stack[len(stack)-1]
			stack = stack[:len(stack)-1]
			match[j] = i
		}
	}
	// For each open bracket at position i, toggle s[i..match[i]]
	for i := 0; i < n; i++ {
		if b[i] == '(' {
			pi := match[i]
			for j := i; j <= pi; j++ {
				if s[j] == '0' {
					s[j] = '1'
				} else {
					s[j] = '0'
				}
			}
		}
	}
}

// isImpossible checks if it's impossible to zero out s.
// Impossible when s[0] != s[2n-1] or count of '1's is odd.
func isImpossible(s string) bool {
	if s[0] != s[len(s)-1] {
		return true
	}
	cnt := 0
	for _, c := range s {
		if c == '1' {
			cnt++
		}
	}
	return cnt%2 != 0
}

func run(cmd *exec.Cmd, input []byte) ([]byte, error) {
	cmd.Stdin = bytes.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return out.Bytes(), err
}

// parseTestCases parses the generated input into individual test cases (n, s).
func parseTestCases(input []byte) ([]struct {
	n int
	s string
}, error) {
	scanner := bufio.NewScanner(bytes.NewReader(input))
	scanner.Buffer(make([]byte, 1<<20), 1<<20)
	scanner.Scan()
	t, _ := strconv.Atoi(strings.TrimSpace(scanner.Text()))
	cases := make([]struct {
		n int
		s string
	}, t)
	for i := 0; i < t; i++ {
		scanner.Scan()
		n, _ := strconv.Atoi(strings.TrimSpace(scanner.Text()))
		scanner.Scan()
		s := strings.TrimSpace(scanner.Text())
		cases[i].n = n
		cases[i].s = s
	}
	return cases, nil
}

func fail(format string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, format+"\n", args...)
	os.Exit(1)
}

func main() {
	if len(os.Args) < 2 {
		fail("usage: verifierF /path/to/binary")
	}

	input := generateInput()

	out, err := run(exec.Command(os.Args[1]), input)
	if err != nil {
		fail("solution runtime error: %v\n%s", err, string(out))
	}

	cases, err := parseTestCases(input)
	if err != nil {
		fail("failed to parse input: %v", err)
	}

	outScanner := bufio.NewScanner(bytes.NewReader(out))
	outScanner.Buffer(make([]byte, 1<<20), 1<<20)

	for tc, c := range cases {
		if !outScanner.Scan() {
			fail("case %d: unexpected end of output", tc+1)
		}
		firstLine := strings.TrimSpace(outScanner.Text())
		k, err := strconv.Atoi(firstLine)
		if err != nil {
			fail("case %d: expected integer k, got %q", tc+1, firstLine)
		}

		if k == -1 {
			// Candidate says impossible
			if !isImpossible(c.s) {
				fail("case %d: candidate says -1 but answer is possible for s=%s", tc+1, c.s)
			}
			continue
		}

		// Candidate says possible
		if isImpossible(c.s) {
			fail("case %d: candidate says k=%d but answer is impossible for s=%s", tc+1, k, c.s)
		}

		if k < 0 || k > 10 {
			fail("case %d: k=%d out of range [0,10]", tc+1, k)
		}

		s := []byte(c.s)
		for i := 0; i < k; i++ {
			if !outScanner.Scan() {
				fail("case %d: expected %d bracket sequences but got only %d", tc+1, k, i)
			}
			b := strings.TrimSpace(outScanner.Text())
			if len(b) != 2*c.n {
				fail("case %d op %d: bracket sequence length %d != %d", tc+1, i+1, len(b), 2*c.n)
			}
			if !isBalanced(b) {
				fail("case %d op %d: bracket sequence is not balanced: %s", tc+1, i+1, b)
			}
			applyOperation(s, b)
		}

		// Check all zeros
		for i, ch := range s {
			if ch != '0' {
				fail("case %d: after %d operations, s[%d]=%c != '0'", tc+1, k, i, ch)
			}
		}
	}

	fmt.Println("all tests passed")
}
