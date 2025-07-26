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

type test struct {
	input    string
	expected string
}

func solveCase(input string) string {
	in := bufio.NewReader(strings.NewReader(strings.TrimSpace(input)))
	var n, h, m int
	if _, err := fmt.Fscan(in, &n, &h, &m); err != nil {
		return ""
	}
	heights := make([]int, n)
	for i := range heights {
		heights[i] = h
	}
	for i := 0; i < m; i++ {
		var l, r, x int
		fmt.Fscan(in, &l, &r, &x)
		for j := l - 1; j < r; j++ {
			if heights[j] > x {
				heights[j] = x
			}
		}
	}
	total := 0
	for _, v := range heights {
		total += v * v
	}
	return fmt.Sprintf("%d", total)
}

func generateTests() []test {
	rng := rand.New(rand.NewSource(42))
	var tests []test
	// simple deterministic cases
	tests = append(tests, test{input: "1 1 0\n", expected: solveCase("1 1 0\n")})
	tests = append(tests, test{input: "3 3 3\n1 1 1\n2 2 3\n3 3 2\n", expected: solveCase("3 3 3\n1 1 1\n2 2 3\n3 3 2\n")})
	tests = append(tests, test{input: "4 10 2\n2 3 8\n3 4 7\n", expected: solveCase("4 10 2\n2 3 8\n3 4 7\n")})
	for len(tests) < 100 {
		n := rng.Intn(50) + 1
		h := rng.Intn(50) + 1
		m := rng.Intn(50) + 1
		var sb strings.Builder
		sb.WriteString(fmt.Sprintf("%d %d %d\n", n, h, m))
		for i := 0; i < m; i++ {
			l := rng.Intn(n) + 1
			r := l + rng.Intn(n-l+1)
			x := rng.Intn(h + 1)
			sb.WriteString(fmt.Sprintf("%d %d %d\n", l, r, x))
		}
		inp := sb.String()
		tests = append(tests, test{inp, solveCase(inp)})
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
			fmt.Printf("Wrong answer on test %d\nInput:\n%sExpected:%s\nGot:%s\n", i+1, t.input, t.expected, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed.\n", len(tests))
}
