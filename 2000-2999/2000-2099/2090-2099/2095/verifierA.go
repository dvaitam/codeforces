package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
)

const expected = "printf(\"puzzling\");"

type test struct {
	input string
}

func generateTests() []test {
	rng := rand.New(rand.NewSource(2095))
	var tests []test

	// Empty input case.
	tests = append(tests, test{input: ""})

	// Random garbage inputs to ensure solution ignores them.
	for len(tests) < 50 {
		l := rng.Intn(50)
		var sb strings.Builder
		for i := 0; i < l; i++ {
			sb.WriteByte(byte(rng.Intn(26) + 'a'))
			if rng.Intn(5) == 0 {
				sb.WriteByte('\n')
			}
		}
		tests = append(tests, test{input: sb.String()})
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
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	tests := generateTests()

	for i, t := range tests {
		got, err := runBinary(bin, t.input)
		if err != nil {
			fmt.Printf("Runtime error on test %d: %v\n", i+1, err)
			os.Exit(1)
		}
		if got != expected {
			fmt.Printf("Wrong answer on test %d\nInput:\n%s\nExpected:\n%s\nGot:\n%s\n", i+1, t.input, expected, got)
			os.Exit(1)
		}
	}

	fmt.Printf("All %d tests passed.\n", len(tests))
}
