package main

import (
	"bufio"
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
)

func generateTests() []string {
	r := rand.New(rand.NewSource(1))
	tests := make([]string, 100)
	for i := range tests {
		n := r.Intn(100) + 1
		b := make([]byte, n)
		for j := range b {
			b[j] = byte('a' + r.Intn(26))
		}
		tests[i] = string(b)
	}
	return tests
}

func solve(s string) string {
	runes := []rune(s)
	for l, r := 0, len(runes)-1; l < r; l, r = l+1, r-1 {
		runes[l], runes[r] = runes[r], runes[l]
	}
	return string(runes) + s
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	tests := generateTests()
	var input bytes.Buffer
	fmt.Fprintln(&input, len(tests))
	for _, s := range tests {
		fmt.Fprintln(&input, s)
	}
	cmd := exec.Command(bin)
	cmd.Stdin = &input
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "execution failed: %v\n", err)
		os.Exit(1)
	}
	scanner := bufio.NewScanner(&out)
	for i, s := range tests {
		if !scanner.Scan() {
			fmt.Fprintf(os.Stderr, "missing output for test case %d\n", i+1)
			os.Exit(1)
		}
		got := strings.TrimSpace(scanner.Text())
		want := solve(s)
		if got != want {
			fmt.Fprintf(os.Stderr, "wrong answer on test %d: input=%s expected=%s got=%s\n", i+1, s, want, got)
			os.Exit(1)
		}
	}
	if scanner.Scan() {
		fmt.Fprintln(os.Stderr, "extra output detected")
		os.Exit(1)
	}
	fmt.Println("All test cases passed.")
}
