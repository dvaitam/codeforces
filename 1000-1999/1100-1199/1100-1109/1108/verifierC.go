package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

// Test represents a single generated test case input.
type Test struct {
	input string
}

// runExe runs the binary at path with the given input and returns its output.
func runExe(path, input string) (string, error) {
	cmd := exec.Command(path)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return out.String(), err
}

// genTests generates random test cases.
func genTests() []Test {
	rand.Seed(2)
	tests := make([]Test, 0, 100)
	colors := []byte{'R', 'G', 'B'}
	for i := 0; i < 100; i++ {
		n := rand.Intn(20) + 1
		var sb strings.Builder
		sb.WriteString(strconv.Itoa(n))
		sb.WriteByte('\n')
		for j := 0; j < n; j++ {
			sb.WriteByte(colors[rand.Intn(3)])
		}
		sb.WriteByte('\n')
		tests = append(tests, Test{sb.String()})
	}
	return tests
}

// parseInput extracts n and s from the generated input.
func parseInput(in string) (int, string, error) {
	fields := strings.Fields(in)
	if len(fields) != 2 {
		return 0, "", fmt.Errorf("unexpected input: %q", in)
	}
	n, err := strconv.Atoi(fields[0])
	if err != nil {
		return 0, "", err
	}
	return n, fields[1], nil
}

// bestPatterns returns the minimal number of changes and all valid patterns.
func bestPatterns(s string) (int, []string) {
	perms := []string{"RGB", "RBG", "GRB", "GBR", "BRG", "BGR"}
	n := len(s)
	min := n + 1
	var patterns []string
	for _, p := range perms {
		diff := 0
		var b strings.Builder
		b.Grow(n)
		for i := 0; i < n; i++ {
			ch := p[i%3]
			if s[i] != ch {
				diff++
			}
			b.WriteByte(ch)
		}
		pat := b.String()
		if diff < min {
			min = diff
			patterns = []string{pat}
		} else if diff == min {
			patterns = append(patterns, pat)
		}
	}
	return min, patterns
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: go run verifierC.go /path/to/binary")
		return
	}
	bin := os.Args[1]
	tests := genTests()
	for i, tc := range tests {
		n, s, err := parseInput(tc.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "bad generated test: %v\n", err)
			os.Exit(1)
		}
		expMin, patterns := bestPatterns(s)
		got, err := runExe(bin, tc.input)
		if err != nil {
			fmt.Printf("candidate runtime error on test %d: %v\n", i+1, err)
			os.Exit(1)
		}
		fields := strings.Fields(got)
		if len(fields) < 2 {
			fmt.Printf("Test %d invalid output format:\n%s\n", i+1, got)
			os.Exit(1)
		}
		r, err := strconv.Atoi(fields[0])
		if err != nil {
			fmt.Printf("Test %d invalid integer: %v\n", i+1, err)
			os.Exit(1)
		}
		t := fields[1]
		if len(t) != n {
			fmt.Printf("Test %d wrong length. Expected %d got %d\n", i+1, n, len(t))
			os.Exit(1)
		}
		diff := 0
		for j := 0; j < n; j++ {
			if s[j] != t[j] {
				diff++
			}
		}
		if diff != r {
			fmt.Printf("Test %d inconsistent changes count. got %d expected %d\n", i+1, r, diff)
			os.Exit(1)
		}
		if r != expMin {
			fmt.Printf("Test %d non optimal answer. expected %d got %d\n", i+1, expMin, r)
			os.Exit(1)
		}
		valid := false
		for _, pat := range patterns {
			if t == pat {
				valid = true
				break
			}
		}
		if !valid {
			fmt.Printf("Test %d produced invalid pattern\nInput:\n%sGot:\n%s %s\n", i+1, tc.input, fields[0], t)
			os.Exit(1)
		}
	}
	fmt.Println("all tests passed")
}
