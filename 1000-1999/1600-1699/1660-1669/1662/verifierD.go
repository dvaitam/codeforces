package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
)

type testCase struct {
	input  string
	expect string
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	tests := genTests()
	for idx, tc := range tests {
		out, err := run(bin, tc.input)
		if err != nil {
			fmt.Printf("test %d runtime error: %v\ninput:\n%s\n", idx+1, err, tc.input)
			os.Exit(1)
		}
		if err := check(strings.TrimSpace(out), tc.expect); err != nil {
			fmt.Printf("test %d failed: %v\ninput:\n%s\nexpected:\n%s\nactual:\n%s\n", idx+1, err, tc.input, tc.expect, out)
			os.Exit(1)
		}
	}
	fmt.Println("all tests passed")
}

func check(actual, expect string) error {
	lines := strings.Fields(actual)
	if len(lines) == 0 {
		return fmt.Errorf("empty output")
	}
	exp := strings.Fields(expect)
	if len(lines) != len(exp) {
		return fmt.Errorf("expected %d answers but got %d", len(exp), len(lines))
	}
	for i := range exp {
		if strings.ToUpper(lines[i]) != exp[i] {
			return fmt.Errorf("mismatch at test %d: expected %s got %s", i+1, exp[i], lines[i])
		}
	}
	return nil
}

func genTests() []testCase {
	rand.Seed(42)
	tests := []testCase{
		makeTestCase([]string{"AB", "AB"}),
		makeTestCase([]string{"AA", "AA"}),
		makeTestCase([]string{"ABAB", "AB"}),
	}
	for i := 0; i < 100; i++ {
		tests = append(tests, randomTestCase())
	}
	return tests
}

func makeTestCase(pairs []string) testCase {
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d\n", len(pairs)/2)
	expect := make([]string, len(pairs)/2)
	for i := 0; i < len(pairs); i += 2 {
		u := pairs[i]
		v := pairs[i+1]
		sb.WriteString(u)
		sb.WriteByte('\n')
		sb.WriteString(v)
		sb.WriteByte('\n')
		if equivalent(u, v) {
			expect[i/2] = "YES"
		} else {
			expect[i/2] = "NO"
		}
	}
	return testCase{
		input:  sb.String(),
		expect: strings.Join(expect, "\n"),
	}
}

func randomTestCase() testCase {
	t := rand.Intn(5) + 1
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d\n", t)
	expect := make([]string, t)
	for i := 0; i < t; i++ {
		u := randomDNA(rand.Intn(10) + 1)
		v := randomDNA(rand.Intn(10) + 1)
		sb.WriteString(u)
		sb.WriteByte('\n')
		sb.WriteString(v)
		sb.WriteByte('\n')
		if equivalent(u, v) {
			expect[i] = "YES"
		} else {
			expect[i] = "NO"
		}
	}
	return testCase{
		input:  sb.String(),
		expect: strings.Join(expect, "\n"),
	}
}

func randomDNA(length int) string {
	letters := []byte{'A', 'B', 'C'}
	b := make([]byte, length)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

func equivalent(u, v string) bool {
	bu, su := canonical(u)
	bv, sv := canonical(v)
	return bu == bv && su == sv
}

func canonical(s string) (int, string) {
	b := 0
	var stack []byte
	for i := 0; i < len(s); i++ {
		ch := s[i]
		if ch == 'B' {
			b ^= 1
		} else {
			if len(stack) > 0 && stack[len(stack)-1] == ch {
				stack = stack[:len(stack)-1]
			} else {
				stack = append(stack, ch)
			}
		}
	}
	return b, string(stack)
}

func run(bin, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return strings.TrimSpace(out.String()), err
}
