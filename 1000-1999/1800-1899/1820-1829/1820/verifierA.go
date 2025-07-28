package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

func solveA(s string) int {
	hasCaret := strings.Contains(s, "^")
	if !hasCaret {
		return len(s) + 1
	}
	if len(s) == 1 {
		return 1
	}
	ans := 0
	if s[0] == '_' {
		ans++
	}
	if s[len(s)-1] == '_' {
		ans++
	}
	for i := 1; i < len(s); i++ {
		if s[i] == '_' && s[i-1] == '_' {
			ans++
		}
	}
	return ans
}

func generateTests() []string {
	rand.Seed(1)
	tests := []string{"^", "_", "^^", "__", "^_", "_^"}
	for len(tests) < 100 {
		n := rand.Intn(100) + 1
		b := make([]byte, n)
		for i := range b {
			if rand.Intn(2) == 0 {
				b[i] = '^'
			} else {
				b[i] = '_'
			}
		}
		tests = append(tests, string(b))
	}
	return tests
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintf(os.Stderr, "usage: go run verifierA.go /path/to/binary\n")
		os.Exit(1)
	}
	binary := os.Args[1]
	tests := generateTests()

	var input bytes.Buffer
	fmt.Fprintln(&input, len(tests))
	for _, t := range tests {
		fmt.Fprintln(&input, t)
	}

	cmd := exec.Command(binary)
	cmd.Stdin = bytes.NewReader(input.Bytes())
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		fmt.Fprintln(os.Stderr, "program failed:", err)
		os.Exit(1)
	}

	outputs := strings.Fields(out.String())
	if len(outputs) != len(tests) {
		fmt.Fprintf(os.Stderr, "expected %d outputs, got %d\n", len(tests), len(outputs))
		os.Exit(1)
	}

	for i, t := range tests {
		got, err := strconv.Atoi(outputs[i])
		if err != nil {
			fmt.Fprintf(os.Stderr, "invalid output on test %d: %v\n", i+1, err)
			os.Exit(1)
		}
		want := solveA(t)
		if got != want {
			fmt.Fprintf(os.Stderr, "test %d failed: input %q expected %d got %d\n", i+1, t, want, got)
			os.Exit(1)
		}
	}

	fmt.Println("all tests passed")
}
