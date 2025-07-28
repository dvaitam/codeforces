package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
)

func solve(s string) int {
	allSame := true
	for i := 1; i < len(s); i++ {
		if s[i] != s[0] {
			allSame = false
			break
		}
	}
	if allSame {
		return -1
	}
	return len(s) - 1
}

func run(bin, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func generateTests() []string {
	rng := rand.New(rand.NewSource(42))
	tests := []string{"a", "aa", "aba", "abba", "aaaa", "abcba"}
	for len(tests) < 100 {
		n := rng.Intn(50) + 1
		b := make([]byte, n)
		if rng.Intn(5) == 0 {
			ch := byte('a' + rng.Intn(26))
			for i := range b {
				b[i] = ch
			}
		} else {
			for i := 0; i < (n+1)/2; i++ {
				ch := byte('a' + rng.Intn(26))
				b[i] = ch
				b[n-1-i] = ch
			}
		}
		tests = append(tests, string(b))
	}
	return tests
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	tests := generateTests()
	for i, s := range tests {
		input := fmt.Sprintf("1\n%s\n", s)
		expected := fmt.Sprintf("%d", solve(s))
		output, err := run(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "test %d: %v\n", i+1, err)
			os.Exit(1)
		}
		if output != expected {
			fmt.Fprintf(os.Stderr, "test %d failed:\ninput:\n%s\nexpected: %s\n got: %s\n", i+1, s, expected, output)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(tests))
}
