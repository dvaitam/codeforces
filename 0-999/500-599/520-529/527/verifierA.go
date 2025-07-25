package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
)

// solve computes expected result for problem A
func solve(a, b int64) int64 {
	var ships int64
	for b != 0 {
		ships += a / b
		a, b = b, a%b
	}
	return ships
}

// generateTests returns at least 100 deterministic test cases
func generateTests() [][2]int64 {
	rnd := rand.New(rand.NewSource(1))
	tests := make([][2]int64, 0, 110)
	// some boundary cases
	tests = append(tests, [2]int64{2, 1})
	tests = append(tests, [2]int64{1000000000000, 1})
	tests = append(tests, [2]int64{1000000000000, 999999999999})
	// random cases
	for len(tests) < 100 {
		a := rnd.Int63n(1_000_000_000_000-2) + 2
		b := rnd.Int63n(a-1) + 1
		tests = append(tests, [2]int64{a, b})
	}
	return tests
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierA.go /path/to/binary")
		return
	}
	binary := os.Args[1]
	tests := generateTests()
	for i, t := range tests {
		input := fmt.Sprintf("%d %d\n", t[0], t[1])
		expected := fmt.Sprintf("%d\n", solve(t[0], t[1]))
		cmd := exec.Command(binary)
		cmd.Stdin = strings.NewReader(input)
		var out bytes.Buffer
		cmd.Stdout = &out
		err := cmd.Run()
		if err != nil {
			fmt.Printf("test %d: runtime error: %v\n", i+1, err)
			os.Exit(1)
		}
		got := strings.TrimSpace(out.String())
		want := strings.TrimSpace(expected)
		if got != want {
			fmt.Printf("test %d failed\ninput: %sexpected: %s\n got: %s\n", i+1, input, want, got)
			os.Exit(1)
		}
	}
	fmt.Printf("all %d tests passed\n", len(tests))
}
