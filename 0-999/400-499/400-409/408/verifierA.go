package main

import (
	"bufio"
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
)

// solveA computes expected output for given input.
func solveA(input string) string {
	in := bufio.NewReader(strings.NewReader(input))
	var n int
	if _, err := fmt.Fscan(in, &n); err != nil {
		return ""
	}
	ks := make([]int, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &ks[i])
	}
	minTime := -1
	for i := 0; i < n; i++ {
		total := 0
		for j := 0; j < ks[i]; j++ {
			var m int
			fmt.Fscan(in, &m)
			total += m*5 + 15
		}
		if minTime == -1 || total < minTime {
			minTime = total
		}
	}
	return fmt.Sprintf("%d\n", minTime)
}

type testCase struct {
	input    string
	expected string
}

func generateTests() []testCase {
	rng := rand.New(rand.NewSource(1))
	var tests []testCase
	// some fixed edge cases
	fixed := []string{
		"1\n1\n1\n",
		"1\n3\n1 1 1\n",
		"2\n1 1\n1\n1\n",
		"3\n2 2 2\n1 1\n1 1\n1 1\n",
	}
	for _, f := range fixed {
		tests = append(tests, testCase{f, solveA(f)})
	}
	for len(tests) < 100 {
		n := rng.Intn(10) + 1
		ks := make([]int, n)
		var sb strings.Builder
		sb.WriteString(fmt.Sprintf("%d\n", n))
		for i := 0; i < n; i++ {
			ks[i] = rng.Intn(5) + 1
			sb.WriteString(fmt.Sprintf("%d ", ks[i]))
		}
		sb.WriteString("\n")
		for i := 0; i < n; i++ {
			for j := 0; j < ks[i]; j++ {
				m := rng.Intn(10) + 1
				sb.WriteString(fmt.Sprintf("%d ", m))
			}
			sb.WriteString("\n")
		}
		inp := sb.String()
		tests = append(tests, testCase{inp, solveA(inp)})
	}
	return tests
}

func runBinary(path, input string) (string, error) {
	cmd := exec.Command(path)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return out.String(), err
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	tests := generateTests()
	for i, t := range tests {
		out, err := runBinary(bin, t.input)
		if err != nil {
			fmt.Printf("Runtime error on test %d: %v\n%s", i+1, err, out)
			os.Exit(1)
		}
		if strings.TrimSpace(out) != strings.TrimSpace(t.expected) {
			fmt.Printf("Wrong answer on test %d\nInput:\n%sExpected: %sGot: %s\n", i+1, t.input, strings.TrimSpace(t.expected), strings.TrimSpace(out))
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed.\n", len(tests))
}
