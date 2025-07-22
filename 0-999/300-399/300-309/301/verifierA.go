package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
)

type TestCaseA struct {
	input    string
	expected string
}

func runBinary(path string, input string) (string, error) {
	cmd := exec.Command(path)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	return strings.TrimSpace(out.String()), err
}

func solveLocal(n int, arr []int) int {
	m := 2*n - 1
	sumAbs := 0
	neg := 0
	minAbs := int(1e9)
	for i := 0; i < m; i++ {
		x := arr[i]
		if x < 0 {
			neg++
		}
		ax := x
		if ax < 0 {
			ax = -ax
		}
		sumAbs += ax
		if ax < minAbs {
			minAbs = ax
		}
	}
	if n%2 == 1 || neg%2 == 0 {
		return sumAbs
	}
	return sumAbs - 2*minAbs
}

func genTests() []TestCaseA {
	rand.Seed(1)
	tests := make([]TestCaseA, 0, 100)
	for t := 0; t < 100; t++ {
		n := rand.Intn(99) + 2 // 2..100
		m := 2*n - 1
		arr := make([]int, m)
		for i := 0; i < m; i++ {
			arr[i] = rand.Intn(2001) - 1000
		}
		expected := solveLocal(n, arr)
		var b strings.Builder
		fmt.Fprintln(&b, n)
		for i, v := range arr {
			if i > 0 {
				b.WriteByte(' ')
			}
			fmt.Fprintf(&b, "%d", v)
		}
		b.WriteByte('\n')
		tests = append(tests, TestCaseA{input: b.String(), expected: fmt.Sprint(expected)})
	}
	return tests
}

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	tests := genTests()
	passed := 0
	for i, tc := range tests {
		out, err := runBinary(bin, tc.input)
		if err != nil {
			fmt.Printf("test %d: runtime error: %v\n", i+1, err)
			continue
		}
		if strings.TrimSpace(out) != tc.expected {
			fmt.Printf("test %d failed: expected %s got %s\n", i+1, tc.expected, out)
		} else {
			passed++
		}
	}
	fmt.Printf("passed %d/%d tests\n", passed, len(tests))
	if passed != len(tests) {
		os.Exit(1)
	}
}
