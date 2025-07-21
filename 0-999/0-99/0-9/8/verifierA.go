package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
)

// can checks if pattern a occurs in s and pattern b occurs after it
func can(s, a, b string) bool {
	i := strings.Index(s, a)
	if i == -1 {
		return false
	}
	j := strings.Index(s[i+len(a):], b)
	return j != -1
}

// reverse returns the reversed string
func reverse(s string) string {
	r := []rune(s)
	for i, j := 0, len(r)-1; i < j; i, j = i+1, j-1 {
		r[i], r[j] = r[j], r[i]
	}
	return string(r)
}

// solve computes the expected answer for a single test case
func solve(input string) string {
	lines := strings.Split(strings.TrimSpace(input), "\n")
	s := strings.TrimSpace(lines[0])
	a := strings.TrimSpace(lines[1])
	b := strings.TrimSpace(lines[2])
	forward := can(s, a, b)
	backward := can(reverse(s), a, b)
	switch {
	case forward && backward:
		return "both"
	case forward:
		return "forward"
	case backward:
		return "backward"
	default:
		return "fantasy"
	}
}

type test struct {
	input    string
	expected string
}

// generateTests creates at least 100 deterministic test cases
func generateTests() []test {
	rand.Seed(42)
	var tests []test
	// some fixed edge cases
	fixed := []string{
		"aaaaa\na\na",
		"abcde\nab\ncd",
		"ababa\naba\nba",
		"xyz\nxy\nz",
		"abcdef\nabc\ndef",
	}
	for _, f := range fixed {
		tests = append(tests, test{f + "\n", solve(f)})
	}
	for len(tests) < 100 {
		n := rand.Intn(20) + 5
		var sb strings.Builder
		for i := 0; i < n; i++ {
			sb.WriteByte(byte('a' + rand.Intn(3)))
		}
		la := rand.Intn(3) + 1
		var sa strings.Builder
		for i := 0; i < la; i++ {
			sa.WriteByte(byte('a' + rand.Intn(3)))
		}
		lb := rand.Intn(3) + 1
		var sb2 strings.Builder
		for i := 0; i < lb; i++ {
			sb2.WriteByte(byte('a' + rand.Intn(3)))
		}
		inp := fmt.Sprintf("%s\n%s\n%s\n", sb.String(), sa.String(), sb2.String())
		tests = append(tests, test{inp, solve(inp)})
	}
	return tests
}

func runBinary(bin string, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return strings.TrimSpace(out.String()), err
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	tests := generateTests()
	for i, t := range tests {
		got, err := runBinary(bin, t.input)
		if err != nil {
			fmt.Printf("Runtime error on test %d: %v\n", i+1, err)
			os.Exit(1)
		}
		if got != strings.TrimSpace(t.expected) {
			fmt.Printf("Wrong answer on test %d\nInput:\n%sExpected: %s\nGot: %s\n", i+1, t.input, strings.TrimSpace(t.expected), got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed.\n", len(tests))
}
