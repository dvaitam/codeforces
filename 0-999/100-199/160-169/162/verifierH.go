package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
)

func altCase(s string) string {
	b := []byte(s)
	for i := range b {
		if i%2 == 0 {
			if b[i] >= 'a' && b[i] <= 'z' {
				b[i] = b[i] - 'a' + 'A'
			}
		} else {
			if b[i] >= 'A' && b[i] <= 'Z' {
				b[i] = b[i] - 'A' + 'a'
			}
		}
	}
	return string(b)
}

type test struct {
	s string
}

func randLetters(n int) string {
	letters := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	b := make([]byte, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

func generateTests() []test {
	rand.Seed(42)
	var tests []test
	fixed := []string{"a", "A", "abcDEF", "HelloWorld", "zzzz"}
	for _, s := range fixed {
		tests = append(tests, test{s: s})
	}
	for len(tests) < 100 {
		l := rand.Intn(20) + 1
		tests = append(tests, test{s: randLetters(l)})
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
		fmt.Fprintln(os.Stderr, "usage: go run verifierH.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	tests := generateTests()
	for i, t := range tests {
		input := fmt.Sprintf("%s\n", t.s)
		exp := altCase(t.s)
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
