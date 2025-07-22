package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
)

type test struct {
	n int
}

func factorize(n int) string {
	var factors []int
	for i := 2; i*i <= n; i++ {
		for n%i == 0 {
			factors = append(factors, i)
			n /= i
		}
	}
	if n > 1 {
		factors = append(factors, n)
	}
	var sb strings.Builder
	for i, f := range factors {
		if i > 0 {
			sb.WriteByte('*')
		}
		fmt.Fprintf(&sb, "%d", f)
	}
	return sb.String()
}

func generateTests() []test {
	rand.Seed(42)
	var tests []test
	fixed := []int{2, 3, 4, 6, 8, 9, 10, 9973, 10000}
	for _, v := range fixed {
		tests = append(tests, test{n: v})
	}
	for len(tests) < 100 {
		tests = append(tests, test{n: rand.Intn(9999) + 2})
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
		fmt.Fprintln(os.Stderr, "usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	tests := generateTests()
	for i, t := range tests {
		input := fmt.Sprintf("%d\n", t.n)
		exp := factorize(t.n)
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
