package main

import (
	"bufio"
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"sort"
	"strings"
)

type testCase struct {
	input  string
	output string
}

func solve(input string) string {
	in := bufio.NewReader(strings.NewReader(input))
	var n, m, k int
	fmt.Fscan(in, &n, &m, &k)
	a := make([]int, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &a[i])
	}
	if k >= m {
		return "0\n"
	}
	gains := make([]int, n)
	for i, v := range a {
		gains[i] = v - 1
	}
	sort.Sort(sort.Reverse(sort.IntSlice(gains)))
	cur := k
	used := 0
	for _, g := range gains {
		if g <= 0 {
			break
		}
		cur += g
		used++
		if cur >= m {
			return fmt.Sprintf("%d\n", used)
		}
	}
	return "-1\n"
}

func generateTests() []testCase {
	rand.Seed(42)
	var tests []testCase
	// some fixed simple tests
	fixed := []string{
		"1 1 1\n1\n",
		"1 5 2\n3\n",
		"3 5 2\n3 1 3\n",
		"2 3 1\n2 2\n",
	}
	for _, f := range fixed {
		tests = append(tests, testCase{f, solve(f)})
	}
	for len(tests) < 100 {
		n := rand.Intn(50) + 1
		m := rand.Intn(50) + 1
		k := rand.Intn(50) + 1
		a := make([]int, n)
		var sb strings.Builder
		fmt.Fprintf(&sb, "%d %d %d\n", n, m, k)
		for i := 0; i < n; i++ {
			a[i] = rand.Intn(50) + 1
			fmt.Fprintf(&sb, "%d", a[i])
			if i+1 < n {
				sb.WriteByte(' ')
			}
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
		if strings.TrimSpace(got) != strings.TrimSpace(t.output) {
			fmt.Printf("Wrong answer on test %d\nInput:\n%sExpected: %sGot: %s\n", i+1, t.input, strings.TrimSpace(t.output), strings.TrimSpace(got))
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed.\n", len(tests))
}
