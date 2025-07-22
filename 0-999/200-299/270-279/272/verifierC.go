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
	input    string
	expected string
}

func solve(input string) string {
	r := strings.NewReader(strings.TrimSpace(input))
	var n int
	fmt.Fscan(r, &n)
	a := make([]int64, n+1)
	for i := 1; i <= n; i++ {
		fmt.Fscan(r, &a[i])
	}
	var m int
	fmt.Fscan(r, &m)
	cur := int64(0)
	var out strings.Builder
	for i := 0; i < m; i++ {
		var w int
		var h int64
		fmt.Fscan(r, &w, &h)
		if a[w] > cur {
			cur = a[w]
		}
		out.WriteString(fmt.Sprintf("%d\n", cur))
		cur += h
	}
	return out.String()
}

func generateTests() []test {
	rand.Seed(44)
	var tests []test
	fixed := []struct {
		n       int
		a       []int64
		queries [][2]int64
	}{
		{1, []int64{5}, [][2]int64{{1, 2}, {1, 3}}},
		{2, []int64{1, 2}, [][2]int64{{1, 1}, {2, 2}, {1, 3}}},
	}
	for _, f := range fixed {
		var sb strings.Builder
		sb.WriteString(fmt.Sprintf("%d\n", f.n))
		for i, v := range f.a {
			if i > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(fmt.Sprintf("%d", v))
		}
		sb.WriteString("\n")
		sb.WriteString(fmt.Sprintf("%d\n", len(f.queries)))
		for _, q := range f.queries {
			sb.WriteString(fmt.Sprintf("%d %d\n", q[0], q[1]))
		}
		inp := sb.String()
		tests = append(tests, test{inp, solve(inp)})
	}
	for len(tests) < 100 {
		n := rand.Intn(5) + 1
		var sb strings.Builder
		sb.WriteString(fmt.Sprintf("%d\n", n))
		a := make([]int64, n+1)
		cur := int64(0)
		for i := 1; i <= n; i++ {
			cur += int64(rand.Intn(5) + 1)
			a[i] = cur
			if i > 1 {
				sb.WriteByte(' ')
			}
			sb.WriteString(fmt.Sprintf("%d", a[i]))
		}
		sb.WriteString("\n")
		m := rand.Intn(5) + 1
		sb.WriteString(fmt.Sprintf("%d\n", m))
		for i := 0; i < m; i++ {
			w := rand.Intn(n) + 1
			h := rand.Intn(5) + 1
			sb.WriteString(fmt.Sprintf("%d %d\n", w, h))
		}
		inp := sb.String()
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
		fmt.Fprintln(os.Stderr, "usage: go run verifierC.go /path/to/binary")
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
			fmt.Printf("Wrong answer on test %d\nInput:\n%sExpected:%sGot:%s\n", i+1, t.input, t.expected, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed.\n", len(tests))
}
