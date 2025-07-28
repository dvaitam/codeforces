package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
)

type Test struct {
	input    string
	expected string
}

func solve(n, s int64) string {
	return fmt.Sprint(s / (n * n))
}

func generateTests() []Test {
	rng := rand.New(rand.NewSource(42))
	tests := []Test{
		{"1\n1 0\n", solve(1, 0)},
		{"1\n2 5\n", solve(2, 5)},
		{"1\n3 27\n", solve(3, 27)},
	}
	for len(tests) < 100 {
		n := int64(rng.Intn(1000) + 1)
		k := int64(rng.Intn(int(n) + 2))
		if k > n+1 {
			k = n + 1
		}
		maxRem := (n + 1 - k) * (n - 1)
		var rem int64
		if maxRem > 0 {
			rem = rng.Int63n(maxRem + 1)
		}
		s := k*n*n + rem
		input := fmt.Sprintf("1\n%d %d\n", n, s)
		tests = append(tests, Test{input: input, expected: solve(n, s)})
	}
	return tests
}

func runBinary(bin, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierA.go /path/to/binary")
		return
	}
	bin := os.Args[1]
	tests := generateTests()
	for i, t := range tests {
		got, err := runBinary(bin, t.input)
		if err != nil {
			fmt.Printf("test %d failed: %v\ninput:%s", i+1, err, t.input)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != strings.TrimSpace(t.expected) {
			fmt.Printf("wrong answer on test %d\ninput:%sexpected:%s\ngot:%s\n", i+1, t.input, t.expected, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed.\n", len(tests))
}
