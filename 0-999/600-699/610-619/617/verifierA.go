package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
)

type test struct {
	input    string
	expected string
}

func solve(x int) int {
	return (x + 4) / 5
}

func generateTests() []test {
	rng := rand.New(rand.NewSource(42))
	tests := make([]test, 100)
	for i := 0; i < 100; i++ {
		x := rng.Intn(1_000_000) + 1
		exp := solve(x)
		tests[i] = test{
			input:    fmt.Sprintf("%d\n", x),
			expected: fmt.Sprintf("%d", exp),
		}
	}
	return tests
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
	err := cmd.Run()
	return strings.TrimSpace(out.String()), err
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	tests := generateTests()
	for i, t := range tests {
		got, err := run(bin, t.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d runtime error: %v\n", i+1, err)
			os.Exit(1)
		}
		if got != strings.TrimSpace(t.expected) {
			fmt.Fprintf(os.Stderr, "case %d failed\ninput:%sexpected:%s\ngot:%s\n", i+1, t.input, t.expected, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed.\n", len(tests))
}
