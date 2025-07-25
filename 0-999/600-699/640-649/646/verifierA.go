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

func solve(a, b int) string {
	return fmt.Sprintf("%d", 6-a-b)
}

func generateTests() []test {
	rng := rand.New(rand.NewSource(42))
	var tests []test
	for len(tests) < 100 {
		a := rng.Intn(3) + 1
		b := rng.Intn(3) + 1
		for b == a {
			b = rng.Intn(3) + 1
		}
		input := fmt.Sprintf("%d %d\n", a, b)
		expected := solve(a, b)
		tests = append(tests, test{input, expected})
	}
	return tests
}

func runBinary(bin, input string) (string, error) {
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
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	if bin == "--" {
		if len(os.Args) < 3 {
			fmt.Println("usage: go run verifierA.go /path/to/binary")
			os.Exit(1)
		}
		bin = os.Args[2]
	}
	tests := generateTests()
	for i, t := range tests {
		got, err := runBinary(bin, t.input)
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
