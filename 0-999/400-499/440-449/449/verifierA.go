package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
)

type test struct {
	input, expected string
}

func solve(input string) string {
	var n, m, k int64
	fmt.Sscan(strings.TrimSpace(input), &n, &m, &k)
	if k > n+m-2 {
		return "-1"
	}
	ans := int64(0)
	if k <= n-1 {
		h := n / (k + 1)
		area := h * m
		if area > ans {
			ans = area
		}
	}
	if k <= m-1 {
		w := m / (k + 1)
		area := w * n
		if area > ans {
			ans = area
		}
	}
	if k > n-1 {
		y := k - (n - 1)
		w := m / (y + 1)
		if w > ans {
			ans = w
		}
	}
	if k > m-1 {
		x := k - (m - 1)
		h := n / (x + 1)
		if h > ans {
			ans = h
		}
	}
	return fmt.Sprintf("%d", ans)
}

func generateTests() []test {
	rand.Seed(449)
	tests := []test{
		{"1 1 0\n", solve("1 1 0\n")},
		{"2 3 4\n", solve("2 3 4\n")},
		{"5 6 5\n", solve("5 6 5\n")},
		{"1000000000 1000000000 0\n", solve("1000000000 1000000000 0\n")},
		{"1 1000000000 999999999\n", solve("1 1000000000 999999999\n")},
	}
	for len(tests) < 100 {
		n := rand.Int63n(1e9) + 1
		m := rand.Int63n(1e9) + 1
		k := rand.Int63n(2e9) + 1
		inp := fmt.Sprintf("%d %d %d\n", n, m, k)
		tests = append(tests, test{inp, solve(inp)})
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
		if got != strings.TrimSpace(t.expected) {
			fmt.Printf("Wrong answer on test %d\nInput:%sExpected:%s\nGot:%s\n", i+1, t.input, strings.TrimSpace(t.expected), got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed.\n", len(tests))
}
