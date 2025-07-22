package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
)

func balanced(s string) string {
	depth := 0
	for _, c := range s {
		if c == '(' {
			depth++
		} else {
			depth--
		}
		if depth < 0 {
			return "NO"
		}
	}
	if depth == 0 {
		return "YES"
	}
	return "NO"
}

type test struct {
	s string
}

func randSeq(n int) string {
	b := make([]byte, n)
	for i := range b {
		if rand.Intn(2) == 0 {
			b[i] = '('
		} else {
			b[i] = ')'
		}
	}
	return string(b)
}

func generateTests() []test {
	rand.Seed(42)
	var tests []test
	fixed := []string{"()", "(())", "(()())", "(", ")(", "(()"}
	for _, s := range fixed {
		tests = append(tests, test{s: s})
	}
	for len(tests) < 100 {
		l := rand.Intn(20) + 1
		tests = append(tests, test{s: randSeq(l)})
	}
	return tests
}

func runBinary(bin, input string) (string, error) {
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
		fmt.Fprintln(os.Stderr, "usage: go run verifierJ.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	tests := generateTests()
	for i, t := range tests {
		input := fmt.Sprintf("%s\n", t.s)
		exp := balanced(t.s)
		out, err := runBinary(bin, input)
		if err != nil {
			fmt.Printf("Runtime error on test %d: %v\n", i+1, err)
			os.Exit(1)
		}
		if out != exp {
			fmt.Printf("Wrong answer on test %d\nInput:\n%sExpected:%s\nGot:%s\n", i+1, input, exp, out)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed.\n", len(tests))
}
