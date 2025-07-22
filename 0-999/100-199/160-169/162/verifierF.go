package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
)

func trailingZeros(n int64) string {
	var count int64
	for p := int64(5); p <= n; p *= 5 {
		count += n / p
	}
	return fmt.Sprint(count)
}

type test struct {
	n int64
}

func generateTests() []test {
	rand.Seed(42)
	var tests []test
	fixed := []int64{1, 5, 10, 25, 50, 100, 1000, 1000000}
	for _, v := range fixed {
		tests = append(tests, test{n: v})
	}
	for len(tests) < 100 {
		tests = append(tests, test{n: rand.Int63n(1000000) + 1})
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
		fmt.Fprintln(os.Stderr, "usage: go run verifierF.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	tests := generateTests()
	for i, t := range tests {
		input := fmt.Sprintf("%d\n", t.n)
		exp := trailingZeros(t.n)
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
