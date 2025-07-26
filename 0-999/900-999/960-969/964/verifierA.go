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

// runCandidate executes the given binary or go source with the provided input
// and returns its output or an error.
func runCandidate(path, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(path, ".go") {
		cmd = exec.Command("go", "run", path)
	} else {
		cmd = exec.Command(path)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return out.String(), err
}

// solve computes the expected answer for problem 964A.
func solve(n int64) int64 {
	return n/2 + 1
}

// generateTests creates at least 100 test values for n.
func generateTests() []int64 {
	rand.Seed(1)
	tests := make([]int64, 0, 100)
	for i := 0; i < 97; i++ {
		val := rand.Int63n(1_000_000_000) + 1
		tests = append(tests, val)
	}
	// edge cases
	tests = append(tests, 1, 2, 1_000_000_000)
	return tests
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	cases := generateTests()
	for i, n := range cases {
		input := fmt.Sprintf("%d\n", n)
		expected := strconv.FormatInt(solve(n), 10)
		got, err := runCandidate(bin, input)
		if err != nil {
			fmt.Printf("runtime error on test %d: %v\n%s", i+1, err, got)
			fmt.Printf("input:\n%s", input)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != expected {
			fmt.Printf("test %d failed\ninput:\n%sexpected: %s\ngot: %s\n", i+1, input, expected, strings.TrimSpace(got))
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
