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

// solveB computes expected output for given input.
func solveB(input string) string {
	in := bufio.NewReader(strings.NewReader(input))
	var supply, demand string
	if _, err := fmt.Fscan(in, &supply); err != nil {
		return ""
	}
	if _, err := fmt.Fscan(in, &demand); err != nil {
		return ""
	}
	count := make([]int, 26)
	for _, c := range supply {
		count[c-'a']++
	}
	need := make([]bool, 26)
	for _, c := range demand {
		need[c-'a'] = true
	}
	total := 0
	for i, b := range need {
		if b {
			if count[i] == 0 {
				return "-1\n"
			}
			total += count[i]
		}
	}
	return fmt.Sprintf("%d\n", total)
}

type testCase struct {
	input    string
	expected string
}

func generateTests() []testCase {
	rng := rand.New(rand.NewSource(2))
	var tests []testCase
	fixed := []string{
		"abc\nabc\n",
		"aabbcc\nabc\n",
		"abc\nabz\n",
		"zzz\nz\n",
	}
	for _, f := range fixed {
		tests = append(tests, testCase{f, solveB(f)})
	}
	for len(tests) < 100 {
		l1 := rng.Intn(10) + 1
		l2 := rng.Intn(10) + 1
		var sb strings.Builder
		for i := 0; i < l1; i++ {
			sb.WriteByte(byte('a' + rng.Intn(26)))
		}
		sb.WriteByte('\n')
		for i := 0; i < l2; i++ {
			sb.WriteByte(byte('a' + rng.Intn(26)))
		}
		sb.WriteByte('\n')
		inp := sb.String()
		tests = append(tests, testCase{inp, solveB(inp)})
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
		fmt.Fprintln(os.Stderr, "usage: go run verifierB.go /path/to/binary")
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
