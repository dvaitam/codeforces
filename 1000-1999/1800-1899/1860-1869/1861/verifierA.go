package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
)

type test struct{ input, expected string }

var primes = []struct{ a, b, val int }{
	{1, 3, 13}, {1, 7, 17}, {1, 9, 19}, {2, 3, 23}, {2, 9, 29}, {3, 1, 31}, {3, 7, 37},
	{4, 1, 41}, {4, 3, 43}, {4, 7, 47}, {5, 3, 53}, {5, 9, 59}, {6, 1, 61}, {6, 7, 67},
	{7, 1, 71}, {7, 3, 73}, {8, 3, 83}, {8, 9, 89}, {9, 7, 97},
}

func solveCase(s string) string {
	var pos [10]int
	for i, ch := range s {
		pos[ch-'0'] = i
	}
	for _, p := range primes {
		if pos[p.a] < pos[p.b] {
			return fmt.Sprintf("%d", p.val)
		}
	}
	return "-1"
}

func generateTests() []test {
	rng := rand.New(rand.NewSource(42))
	digits := []byte{'1', '2', '3', '4', '5', '6', '7', '8', '9'}
	var tests []test
	for len(tests) < 100 {
		rng.Shuffle(len(digits), func(i, j int) { digits[i], digits[j] = digits[j], digits[i] })
		s := string(digits)
		input := fmt.Sprintf("1\n%s\n", s)
		tests = append(tests, test{input, solveCase(s)})
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
		if strings.TrimSpace(got) != strings.TrimSpace(t.expected) {
			fmt.Printf("Wrong answer on test %d\nInput:%sExpected:%s\nGot:%s\n", i+1, t.input, t.expected, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed.\n", len(tests))
}
