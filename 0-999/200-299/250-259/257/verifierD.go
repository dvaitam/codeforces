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

type testCase struct {
	input  string
	output string
}

func solve(input string) string {
	in := bufio.NewReader(strings.NewReader(input))
	var n int
	fmt.Fscan(in, &n)
	a := make([]int, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &a[i])
	}
	signs := make([]bool, n)
	var sum int64
	for i := n - 1; i >= 0; i-- {
		if sum > 0 {
			sum -= int64(a[i])
			signs[i] = true
		} else {
			sum += int64(a[i])
		}
	}
	if sum > 0 {
		for i := range signs {
			signs[i] = !signs[i]
		}
	}
	var sb strings.Builder
	for i := 0; i < n; i++ {
		if signs[i] {
			sb.WriteByte('+')
		} else {
			sb.WriteByte('-')
		}
	}
	sb.WriteByte('\n')
	return sb.String()
}

func generateTests() []testCase {
	rand.Seed(45)
	var tests []testCase
	fixed := []string{
		"1\n0\n",
		"1\n5\n",
		"2\n1 2\n",
		"3\n1 1 1\n",
	}
	for _, f := range fixed {
		tests = append(tests, testCase{f, solve(f)})
	}
	for len(tests) < 100 {
		n := rand.Intn(20) + 1
		var sb strings.Builder
		fmt.Fprintf(&sb, "%d\n", n)
		a := make([]int, n)
		a[0] = rand.Intn(1000) + 1
		fmt.Fprintf(&sb, "%d", a[0])
		for i := 1; i < n; i++ {
			a[i] = a[i-1] + rand.Intn(a[i-1]+1)
			if a[i] > 1000000000 {
				a[i] = 1000000000
			}
			if a[i] < a[i-1] {
				a[i] = a[i-1]
			}
			if a[i] > 2*a[i-1] {
				a[i] = 2 * a[i-1]
			}
			fmt.Fprintf(&sb, " %d", a[i])
		}
		sb.WriteByte('\n')
		inp := sb.String()
		tests = append(tests, testCase{inp, solve(inp)})
	}
	return tests
}

func runBinary(bin, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return out.String(), err
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierD.go /path/to/binary")
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
		if strings.TrimSpace(got) != strings.TrimSpace(t.output) {
			fmt.Printf("Wrong answer on test %d\nInput:\n%sExpected: %sGot: %s\n", i+1, t.input, strings.TrimSpace(t.output), strings.TrimSpace(got))
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed.\n", len(tests))
}
