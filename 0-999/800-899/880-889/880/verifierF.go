package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

type testCase struct {
	input  string
	expect string
}

func isPrime(n int) bool {
	if n < 2 {
		return false
	}
	for i := 2; i*i <= n; i++ {
		if n%i == 0 {
			return false
		}
	}
	return true
}

func buildTests() []testCase {
	tests := make([]testCase, 100)
	for i := 0; i < 100; i++ {
		v := i + 1
		in := fmt.Sprintf("%d\n", v)
		expect := "NO"
		if isPrime(v) {
			expect = "YES"
		}
		tests[i] = testCase{input: in, expect: expect}
	}
	return tests
}

func run(binary, input string) (string, error) {
	cmd := exec.Command(binary)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	err := cmd.Run()
	if err != nil {
		return "", fmt.Errorf("run error: %v, stderr: %s", err, stderr.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierF.go /path/to/binary")
		os.Exit(1)
	}
	binary := os.Args[1]
	tests := buildTests()
	for i, t := range tests {
		got, err := run(binary, t.input)
		if err != nil {
			fmt.Printf("test %d failed: %v\n", i+1, err)
			os.Exit(1)
		}
		if got != t.expect {
			fmt.Printf("test %d failed: input %s expected %s got %s\n", i+1, strings.TrimSpace(t.input), t.expect, got)
			os.Exit(1)
		}
	}
	fmt.Printf("ok %d\n", len(tests))
}
