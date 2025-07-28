package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

type TestCase struct {
	n     int
	edges [][2]int
}

func solve(TestCase) int {
	return 1
}

func genTests() []TestCase {
	tests := make([]TestCase, 0, 100)
	for i := 0; i < 100; i++ {
		n := 3 + i%4
		edges := make([][2]int, n-1)
		for j := 2; j <= n; j++ {
			edges[j-2] = [2]int{j - 1, j}
		}
		tests = append(tests, TestCase{n, edges})
	}
	return tests
}

func run(bin string, input string) (string, error) {
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
		fmt.Fprintln(os.Stderr, "usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	tests := genTests()

	var input strings.Builder
	fmt.Fprintln(&input, len(tests))
	for _, tc := range tests {
		fmt.Fprintln(&input, tc.n)
		for _, e := range tc.edges {
			fmt.Fprintf(&input, "%d %d\n", e[0], e[1])
		}
	}

	expected := make([]string, len(tests))
	for i := range tests {
		expected[i] = fmt.Sprintf("%d", solve(tests[i]))
	}

	out, err := run(bin, input.String())
	if err != nil {
		fmt.Fprintln(os.Stderr, "error running binary:", err)
		fmt.Print(out)
		os.Exit(1)
	}

	lines := strings.Split(strings.TrimSpace(out), "\n")
	if len(lines) != len(expected) {
		fmt.Printf("wrong number of lines: got %d want %d\n", len(lines), len(expected))
		os.Exit(1)
	}
	for i, got := range lines {
		got = strings.TrimSpace(got)
		if got != expected[i] {
			fmt.Printf("test %d failed expected %s got %s\n", i+1, expected[i], got)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed.")
}
