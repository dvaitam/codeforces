package main

import (
	"bytes"
	"context"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
	"time"
)

type TestCase struct {
	Input  string
	Output string
}

func runBinary(bin string, input string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	cmd := exec.CommandContext(ctx, bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	err := cmd.Run()
	if err != nil {
		if ctx.Err() == context.DeadlineExceeded {
			return "", fmt.Errorf("timeout")
		}
		return "", fmt.Errorf("%v: %s", err, stderr.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func prefixFunction(s string) []int {
	n := len(s)
	pi := make([]int, n)
	for i := 1; i < n; i++ {
		j := pi[i-1]
		for j > 0 && s[i] != s[j] {
			j = pi[j-1]
		}
		if s[i] == s[j] {
			j++
		}
		pi[i] = j
	}
	return pi
}

func solveA(n, k int, t string) string {
	pi := prefixFunction(t)
	l := pi[n-1]
	p := t[l:]
	b := strings.Builder{}
	b.WriteString(t)
	for i := 1; i < k; i++ {
		b.WriteString(p)
	}
	return b.String()
}

func randString(n int) string {
	letters := []rune("abcdefghijklmnopqrstuvwxyz")
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

func generateTests() []TestCase {
	rand.Seed(42)
	tests := make([]TestCase, 100)
	for i := 0; i < 100; i++ {
		n := rand.Intn(50) + 1
		k := rand.Intn(50) + 1
		t := randString(n)
		input := fmt.Sprintf("%d %d\n%s\n", n, k, t)
		output := solveA(n, k, t)
		tests[i] = TestCase{Input: input, Output: output}
	}
	return tests
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "Usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	tests := generateTests()
	for i, tc := range tests {
		got, err := runBinary(bin, tc.Input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Test %d: runtime error: %v\n", i+1, err)
			os.Exit(1)
		}
		if got != strings.TrimSpace(tc.Output) {
			fmt.Fprintf(os.Stderr, "Test %d failed:\ninput:\n%s\nexpected:%s\n got:%s\n", i+1, tc.Input, tc.Output, got)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed!")
}
