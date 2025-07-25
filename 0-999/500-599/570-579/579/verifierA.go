package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

// expected computes the number of set bits in x.
func expected(x int) int {
	count := 0
	for x > 0 {
		if x&1 == 1 {
			count++
		}
		x >>= 1
	}
	return count
}

// runBinary runs the provided binary with given input and returns trimmed stdout.
func runBinary(binary string, input string) (string, error) {
	cmd := exec.Command(binary)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return strings.TrimSpace(out.String()), err
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	binary := os.Args[1]

	tests := []int{}
	for i := 1; i <= 100; i++ {
		tests = append(tests, i)
	}
	extra := []int{101, 123, 234, 345, 456, 567, 678, 789, 890, 901,
		1000, 1234, 2048, 4096, 8192, 16384, 32768, 65536,
		99999, 100000, 500000, 999999, 1000000, 1234567, 7654321,
		9999999, 12345678, 50000000, 100000000, 200000000, 500000000, 1000000000}
	tests = append(tests, extra...)

	for idx, x := range tests {
		exp := expected(x)
		output, err := runBinary(binary, fmt.Sprintf("%d\n", x))
		if err != nil {
			fmt.Printf("Test %d (x=%d) runtime error: %v\n", idx+1, x, err)
			os.Exit(1)
		}
		if output != fmt.Sprintf("%d", exp) {
			fmt.Printf("Test %d (x=%d) failed: expected %d, got %s\n", idx+1, x, exp, output)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed!\n", len(tests))
}
