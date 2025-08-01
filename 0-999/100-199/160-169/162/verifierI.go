package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
)

func isPrime(n int) bool {
	if n < 2 {
		return false
	}
	if n%2 == 0 {
		return n == 2
	}
	for i := 3; i*i <= n; i += 2 {
		if n%i == 0 {
			return false
		}
	}
	return true
}

func truncatable(n int) string {
	if strings.ContainsRune(fmt.Sprint(n), '0') {
		return "NO"
	}
	if !isPrime(n) {
		return "NO"
	}
	p := 10
	for p <= n {
		suffix := n % p
		if !isPrime(suffix) {
			return "NO"
		}
		p *= 10
	}
	return "YES"
}

type test struct {
	n int
}

func generateTests() []test {
	rand.Seed(42)
	var tests []test
	fixed := []int{2, 3, 13, 23, 233, 2333, 2339, 53, 101}
	for _, v := range fixed {
		tests = append(tests, test{n: v})
	}
	for len(tests) < 100 {
		tests = append(tests, test{n: rand.Intn(10000000-2) + 2})
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
		fmt.Fprintln(os.Stderr, "usage: go run verifierI.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	tests := generateTests()
	for i, t := range tests {
		input := fmt.Sprintf("%d\n", t.n)
		exp := truncatable(t.n)
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
