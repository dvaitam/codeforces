package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
)

func producesOutput(p string) string {
	for _, c := range p {
		if c == 'H' || c == 'Q' || c == '9' {
			return "YES"
		}
	}
	return "NO"
}

type test struct {
	p string
}

func randProgram(n int) string {
	const letters = "HQ9+abcdxyz123!@#"
	b := make([]byte, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

func generateTests() []test {
	rand.Seed(42)
	var tests []test
	fixed := []string{"H", "Q", "9", "+", "abc"}
	for _, s := range fixed {
		tests = append(tests, test{p: s})
	}
	for len(tests) < 100 {
		l := rand.Intn(10) + 1
		tests = append(tests, test{p: randProgram(l)})
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
		fmt.Fprintln(os.Stderr, "usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	tests := generateTests()
	for i, t := range tests {
		input := fmt.Sprintf("%s\n", t.p)
		exp := producesOutput(t.p)
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
