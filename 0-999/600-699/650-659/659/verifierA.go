package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
)

func solve(input string) string {
	var n, a, b int
	fmt.Sscan(input, &n, &a, &b)
	pos := (a - 1 + b) % n
	if pos < 0 {
		pos += n
	}
	return fmt.Sprintln(pos + 1)
}

func generateTests() []string {
	rand.Seed(42)
	tests := make([]string, 100)
	for i := 0; i < 100; i++ {
		n := rand.Intn(100) + 1
		a := rand.Intn(n) + 1
		b := rand.Intn(201) - 100
		tests[i] = fmt.Sprintf("%d %d %d\n", n, a, b)
	}
	return tests
}

func runBinary(bin string, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	return out.String(), err
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	tests := generateTests()
	for idx, t := range tests {
		expect := strings.TrimSpace(solve(t))
		got, err := runBinary(bin, t)
		if err != nil {
			fmt.Fprintf(os.Stderr, "test %d: runtime error: %v\n", idx+1, err)
			os.Exit(1)
		}
		got = strings.TrimSpace(got)
		if expect != got {
			fmt.Printf("test %d failed\ninput:\n%sexpected:%s\ngot:%s\n", idx+1, t, expect, got)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed.")
}
