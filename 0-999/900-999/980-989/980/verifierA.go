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
	input  string
	output string
}

func solveA(s string) string {
	pearls := 0
	links := 0
	for _, ch := range s {
		if ch == 'o' {
			pearls++
		} else if ch == '-' {
			links++
		}
	}
	if pearls <= 1 || links%pearls == 0 {
		return "YES\n"
	}
	return "NO\n"
}

func generateTests() []Test {
	rand.Seed(42)
	tests := make([]Test, 0, 100)
	fixed := []string{"o--", "-o-", "ooo", "---", "o-o", "-oo", "--o"}
	for _, s := range fixed {
		tests = append(tests, Test{input: s + "\n", output: solveA(s)})
	}
	for len(tests) < 100 {
		n := rand.Intn(98) + 3
		b := make([]byte, n)
		for i := 0; i < n; i++ {
			if rand.Intn(2) == 0 {
				b[i] = 'o'
			} else {
				b[i] = '-'
			}
		}
		s := string(b)
		tests = append(tests, Test{input: s + "\n", output: solveA(s)})
	}
	return tests
}

func run(binary string, input string) (string, error) {
	cmd := exec.Command(binary)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	err := cmd.Run()
	if err != nil {
		return "", fmt.Errorf("%v: %s", err, stderr.String())
	}
	return out.String(), nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	binary := os.Args[1]
	tests := generateTests()
	for i, t := range tests {
		got, err := run(binary, t.input)
		if err != nil {
			fmt.Printf("Test %d runtime error: %v\n", i+1, err)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != strings.TrimSpace(t.output) {
			fmt.Printf("Test %d failed. Input: %q\nExpected: %qGot: %q\n", i+1, t.input, t.output, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed.\n", len(tests))
}
